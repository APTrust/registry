package network

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/APTrust/registry/constants"
	"github.com/go-redis/redis/v7"
)

// RedisClient is a client that lets workers store and retrieve working
// data from a Redis server.
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

// Keys returns all keys in the Redis DB matching the specified pattern.
// Each key is a WorkItem.ID in string form. It's generally safe to call
// this with pattern "*" because we rarely have more than a few dozen items
// in Redis at any given time.
func (c *RedisClient) Keys(pattern string) ([]string, error) {
	return c.client.Keys(pattern).Result()
}

func (c *RedisClient) IngestObjectGet(workItemID int64, objIdentifier string) (string, error) {
	obj, err := c.ingestObjectGet(workItemID, objIdentifier)
	if err != nil {
		return "", err
	}
	for _, operationName := range constants.NSQIngestTopicFor {
		op, err := c.workResultGet(workItemID, operationName)
		if err == nil {
			obj[operationName] = op
		} else {
			obj[operationName] = err.Error()
		}
	}
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
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
	err = json.Unmarshal([]byte(data), obj)
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
	if err != nil {
		return nil, fmt.Errorf("WorkResultGet (%d, %s): %s",
			workItemID, operationName, err.Error())
	}
	err = json.Unmarshal([]byte(data), obj)
	if err != nil {
		return nil, fmt.Errorf("WorkResultGet unmarshal JSON (%d, %s): %s",
			workItemID, operationName, err.Error())
	}
	return obj, nil
}

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
