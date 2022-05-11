package network_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var redisItemID = int64(31337)
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
			err = c.SaveItem(redisItemID, rkey, string(jsonBytes))
			require.Nil(t, err, key)
		} else {
			rkey := fmt.Sprintf("workresult:%s", key)
			err = c.SaveItem(redisItemID, rkey, string(jsonBytes))
			require.Nil(t, err, key)
		}
	}
}

func TestRedisIngestObjectGet(t *testing.T) {
	client := getRedisClient()
	assert.NotNil(t, client)
	createRedisIngestObject(t, client)

	objJson, err := client.IngestObjectGet(redisItemID, redisObjIdentifier)
	require.Nil(t, err)

	expected, err := jsonToMap(ingestObjectJson)
	require.Nil(t, err)

	actual, err := jsonToMap(objJson)
	require.Nil(t, err)

	assert.Equal(t, expected, actual)
}

// TODO:
// Put Restoration json
// Test RestorationObjectGet
// Test WorkItemDelete
