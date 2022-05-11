package network_test

import (
	//	"fmt"
	//	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/network"
	"github.com/stretchr/testify/assert"
	//	"github.com/stretchr/testify/require"
)

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
