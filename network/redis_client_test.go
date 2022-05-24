package network_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var redisIngestItemID = int64(867)
var redisRestoreItemID = int64(5309)
var redisObjIdentifier = "test.edu/MaySample"

var ingestObjectJson = `{
  "object": {
    "key1": "value1"
  },
  "ingest01_prefetch": {
    "result": "prefetch result"
  },
  "ingest02_bag_validation": {
    "result": "validation result"
  },
  "ingest03_reingest_check": {
    "result": "reingest result"
  },
  "ingest04_staging": {
    "result": "staging result"
  },
  "ingest05_format_identification": {
    "result": "format identifier result"
  },
  "ingest06_storage": {
    "result": "storage result"
  },
  "ingest07_storage_validation": {
    "result": "storage validation result"
  },
  "ingest08_record": {
    "result": "record result"
  },
  "ingest09_cleanup": {
    "result": "cleanup result"
  }
}`

var restoreItemJSON = `{
  "key1":"value1",
  "key2": 2.02
}`

func getRedisClient() *network.RedisClient {
	config := common.NewConfig()
	return network.NewRedisClient(
		config.Redis.URL,
		config.Redis.Password,
		config.Redis.DefaultDB,
	)
}

func TestNewRedisClient(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
}

func TestRedisPing(t *testing.T) {
	client := getRedisClient()
	response, err := client.Ping()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", response)
}

func jsonToMap(jsonData string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(ingestObjectJson), &data)
	return data, err
}

func createRedisIngestObject(t *testing.T, c *network.RedisClient) {
	data, err := jsonToMap(ingestObjectJson)
	require.Nil(t, err)

	for key, value := range data {
		jsonBytes, err := json.Marshal(value)
		require.Nil(t, err)
		if key == "object" {
			rkey := fmt.Sprintf("object:%s", redisObjIdentifier)
			err = c.SaveItem(redisIngestItemID, rkey, string(jsonBytes))
			require.Nil(t, err, key)
		} else {
			rkey := fmt.Sprintf("workresult:%s", key)
			err = c.SaveItem(redisIngestItemID, rkey, string(jsonBytes))
			require.Nil(t, err, key)
		}
	}
}

func createRedisRestoreObject(t *testing.T, c *network.RedisClient) {
	key := fmt.Sprintf("restoration:%s", redisObjIdentifier)
	err := c.SaveItem(redisRestoreItemID, key, string(restoreItemJSON))
	require.Nil(t, err, key)
}

func TestRedisIngestObjectGet(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
	createRedisIngestObject(t, client)

	objJson, err := client.IngestObjectGet(redisIngestItemID, redisObjIdentifier)
	require.Nil(t, err)

	expected, err := jsonToMap(ingestObjectJson)
	require.Nil(t, err)

	actual, err := jsonToMap(objJson)
	require.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func TestRedisRestoreObjectGet(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
	createRedisRestoreObject(t, client)

	str, err := client.RestorationObjectGet(redisRestoreItemID, redisObjIdentifier)
	require.Nil(t, err)
	assert.Equal(t, restoreItemJSON, str)
}

func TestRedisWorkItemExistsDelete(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
	createRedisIngestObject(t, client)
	createRedisRestoreObject(t, client)

	str, err := client.RestorationObjectGet(redisRestoreItemID, redisObjIdentifier)
	require.Nil(t, err)
	assert.NotEmpty(t, str)

	str, err = client.IngestObjectGet(redisIngestItemID, redisObjIdentifier)
	require.Nil(t, err)
	assert.NotEmpty(t, str)

	assert.True(t, client.KeyExists(redisIngestItemID))
	assert.True(t, client.KeyExists(redisRestoreItemID))

	count, err := client.WorkItemDelete(redisRestoreItemID)
	require.Nil(t, err)
	assert.Equal(t, int64(1), count)

	count, err = client.WorkItemDelete(redisIngestItemID)
	require.Nil(t, err)
	assert.Equal(t, int64(1), count)

	// If keys were successfully deleted, these calls should return error.
	_, err = client.RestorationObjectGet(redisRestoreItemID, redisObjIdentifier)
	assert.NotNil(t, err)

	_, err = client.IngestObjectGet(redisRestoreItemID, redisObjIdentifier)
	assert.NotNil(t, err)

	assert.False(t, client.KeyExists(redisIngestItemID))
	assert.False(t, client.KeyExists(redisRestoreItemID))
}

func TestRedisList(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
	createRedisIngestObject(t, client)
	createRedisRestoreObject(t, client)

	keys, err := client.List("*")
	require.Nil(t, err)
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, strconv.FormatInt(redisIngestItemID, 10))
	assert.Contains(t, keys, strconv.FormatInt(redisRestoreItemID, 10))
}
