package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/APTrust/registry/constants"
	"github.com/go-redis/redis/v7"
)

// RedisClient is a crude and deliberately limited implementation that
// returns JSON only. The JSON is intended for human consumption, primarily
// for APTrust admins to view when debugging problems.
//
// The JSON represents internal state information from the ingest and
// restoration workers. The structure of the data may change over time,
// so this client simply uses map[string]interface{} structures to
// accomodate arbitrary JSON structures.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient. Param address is the net address
// of the Redis server. Param password is the password required to connect.
// It may be blank, but shouldn't be in production. Param db is the id of the
// Redis database.
func NewRedisClient(address, password string, db int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}
}

// Ping pings the Redis server. It should return "PONG" if the server is
// running and we can connect.
func (c *RedisClient) Ping() (string, error) {
	return c.client.Ping().Result()
}

// KeyExists returns true if the specified key exists in our Redis DB.
func (c *RedisClient) KeyExists(workItemID int64) bool {
	key := strconv.FormatInt(workItemID, 10)
	count, err := c.client.Exists(key).Result()
	return count > 0 && err == nil
}

// IngestObjectGet returns a JSON string representing an ingest object
// and its associated work results.
func (c *RedisClient) IngestObjectGet(workItemID int64, objIdentifier string) (string, error) {
	obj := make(map[string]interface{})
	errMessages := make([]string, 0)
	ingestObj, err := c.ingestObjectGet(workItemID, objIdentifier)
	if err != nil {
		return "", err
	}
	obj["object"] = ingestObj

	// Ideally, we'd log errors here, but logging is in common
	// and common imports this package, and we can't do a circular
	// import. So we collect errors and return them as one.

	for _, operationName := range constants.NSQIngestTopicFor {
		op, err := c.workResultGet(workItemID, operationName)
		obj[operationName] = op
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	if len(errMessages) > 0 {
		err = errors.New(strings.Join(errMessages, "; "))
	}
	return string(data), err
}

func (c *RedisClient) ingestObjectGet(workItemID int64, objIdentifier string) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	key := strconv.FormatInt(workItemID, 10)
	field := fmt.Sprintf("object:%s", objIdentifier)
	data, err := c.client.HGet(key, field).Result()
	if err != nil {
		return nil, fmt.Errorf("IngestObjectGet from Redis (%d, %s): %s",
			workItemID, objIdentifier, err.Error())
	}
	err = json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return nil, fmt.Errorf("IngestObjectGet unmarshal JSON (%d, %s): %s",
			workItemID, objIdentifier, err.Error())
	}
	return obj, nil
}

func (c *RedisClient) workResultGet(workItemID int64, operationName string) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	key := strconv.FormatInt(workItemID, 10)
	field := fmt.Sprintf("workresult:%s", operationName)
	data, err := c.client.HGet(key, field).Result()
	if err != nil && err.Error() != "redis: nil" {
		return nil, fmt.Errorf("WorkResultGet (%d, %s): %s",
			workItemID, operationName, err.Error())
	}
	if strings.TrimSpace(data) == "" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return nil, fmt.Errorf("WorkResultGet unmarshal JSON (%d, %s): %s",
			workItemID, operationName, err.Error())
	}
	return obj, nil
}

// RestorationObjectGet returns a JSON string representing the specified
// restoration object.
func (c *RedisClient) RestorationObjectGet(workItemID int64, objIdentifier string) (string, error) {
	key := strconv.FormatInt(workItemID, 10)
	field := fmt.Sprintf("restoration:%s", objIdentifier)
	data, err := c.client.HGet(key, field).Result()
	if err != nil {
		return "", fmt.Errorf("RestorationObjectGet (%d, %s): %s",
			workItemID, objIdentifier, err.Error())
	}
	return string(data), nil
}

// WorkItemDelete deletes the Redis copy (NOT the Registry copy) of a WorkItem,
// along with its associated IngestObject and IngestFile records.
// This is dangerous and should be called only in two cases:
//
// 1. We want to delete old ingest data from Redis after a failed ingest
//    that we know we will never retry. This is essentially a cleanup operation.
//
// 2. We are forcing an item back to the very first step of ingest or
//    restoration, and we want the workers to redo all work from scratch
//    instead of relying on the already completed work recorded in Redis.
//    This is extremely rare, but it will come up a few times a year.
func (c *RedisClient) WorkItemDelete(workItemID int64) (int64, error) {
	key := strconv.FormatInt(workItemID, 10)
	return c.client.Del(key).Result()
}

// SaveItem saves value to Redis. This is used only for testing.
func (c *RedisClient) SaveItem(workItemID int64, field, value string) error {
	key := strconv.FormatInt(workItemID, 10)
	_, err := c.client.HSet(key, field, value).Result()
	return err
}

// List returns up to the first 500 keys in the Redis DB matching
// the specified pattern. The 500 limit is to prevent overload if
// our Redis DB fills up with lots of entries.
//
// Realistically, we will almost never have more than a few dozen keys
// at any given time, since Redis data  is deleted as soon as processing
// completes. Each key is a WorkItem.ID in string form.
func (c *RedisClient) List(pattern string) ([]string, error) {
	keys, _, err := c.client.Scan(0, pattern, 500).Result()
	return keys, err
}
