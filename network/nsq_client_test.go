package network_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var topic = "throwaway_test_topic"
var channel = "test_channel"

func TestNSQClient(t *testing.T) {
	aptContext := common.Context()
	require.NotNil(t, aptContext.NSQClient)

	for i := 0; i < 5; i++ {
		err := aptContext.NSQClient.Enqueue(topic, int64(i+1000))
		require.Nil(t, err)
	}
	for i := 0; i < 5; i++ {
		identifier := fmt.Sprintf("some/object/identifier_%d", i)
		err := aptContext.NSQClient.EnqueueString(topic, identifier)
		require.Nil(t, err)
	}

	stats, err := aptContext.NSQClient.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	topicStats := stats.GetTopic(topic)
	assert.NotNil(t, topicStats)
	assert.Equal(t, topic, topicStats.TopicName)
	assert.EqualValues(t, 10, topicStats.Depth)

	// Create a channel for tests below
	url := fmt.Sprintf("%s/channel/create?topic=%s&channel=%s", aptContext.Config.NsqUrl, topic, channel)
	fmt.Println(url)
	resp, err := http.Post(url, "text/html", nil)
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Re-fetch the stats, so we can look into that channel
	stats, err = aptContext.NSQClient.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)
}
