package network_test

import (
	"fmt"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNSQClient(t *testing.T) {
	var topic = "throwaway_test_topic"
	var channel = "test_channel"
	aptContext := common.Context()
	require.NotNil(t, aptContext.NSQClient)

	// Start with a clean slate
	resetNSQ(aptContext.NSQClient)

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
	err = aptContext.NSQClient.CreateChannel(topic, channel)
	require.Nil(t, err)

	// Re-fetch the stats, so we can look into that channel
	stats, err = aptContext.NSQClient.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)
}

// Clear out any items remaining from prior tests.
func resetNSQ(client *network.NSQClient) {
	client.EmptyAllChannels()
	client.EmptyAllTopics()
	client.DeleteAllChannels()
	client.DeleteAllTopics()
}

func queueSomeItems(t *testing.T, client *network.NSQClient) {
	for i := 0; i < 5; i++ {
		topic := fmt.Sprintf("admin_topic_%d", i)
		err := client.CreateTopic(topic)
		require.Nil(t, err)
		for j := 0; j < 5; j++ {
			channel := fmt.Sprintf("channel_%d", j)
			err = client.CreateChannel(topic, channel)
			require.Nil(t, err)
		}
		for j := 0; j < 5; j++ {
			err = client.Enqueue(topic, int64(j+1000))
			require.Nil(t, err)
		}
	}
}

func TestNSQAdminFunctions(t *testing.T) {
	client := common.Context().NSQClient
	require.NotNil(t, client)

	// Start with a clean slate
	resetNSQ(client)

	queueSomeItems(t, client)
	testItemsAreQueued(t, client)
	testPauseOne(t, client)
	testPauseAll(t, client)
	testUnpauseAll(t, client)
	testEmptyOne(t, client)
	testEmptyAll(t, client)
}

func testItemsAreQueued(t *testing.T, client *network.NSQClient) {
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 5, len(stats.Topics))
	for _, topicStats := range stats.Topics {
		assert.False(t, topicStats.Paused)
		assert.Equal(t, 5, len(topicStats.Channels))
		assert.True(t, topicStats.MessageCount > 0)
		for _, channelStats := range topicStats.Channels {
			assert.False(t, channelStats.Paused)
			assert.True(t, channelStats.MessageCount > 0)
			assert.True(t, channelStats.Depth > 0)
		}
	}
}

func testPauseOne(t *testing.T, client *network.NSQClient) {
	err := client.PauseTopic("admin_topic_1")
	require.Nil(t, err)

	err = client.PauseChannel("admin_topic_2", "channel_2")
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	topic1 := stats.GetTopic("admin_topic_1")
	require.NotNil(t, topic1)
	assert.True(t, topic1.Paused)

	channel2 := stats.GetChannel("admin_topic_2", "channel_2")
	require.NotNil(t, channel2)
	assert.True(t, channel2.Paused)
}

func testUnpauseOne(t *testing.T, client *network.NSQClient) {
	err := client.UnpauseTopic("admin_topic_1")
	require.Nil(t, err)

	err = client.UnpauseChannel("admin_topic_2", "channel_2")
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	topic1 := stats.GetTopic("admin_topic_1")
	require.NotNil(t, topic1)
	assert.False(t, topic1.Paused)

	channel2 := stats.GetChannel("admin_topic_2", "channel_2")
	require.NotNil(t, channel2)
	assert.False(t, channel2.Paused)
}

func testEmptyOne(t *testing.T, client *network.NSQClient) {
	err := client.EmptyTopic("admin_topic_3")
	require.Nil(t, err)

	err = client.EmptyChannel("admin_topic_4", "channel_4")
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	topic3 := stats.GetTopic("admin_topic_3")
	require.NotNil(t, topic3)
	assert.EqualValues(t, 0, topic3.Depth)

	channel4 := stats.GetChannel("admin_topic_4", "channel_4")
	require.NotNil(t, channel4)
	assert.EqualValues(t, 0, channel4.Depth)
}

func testPauseAll(t *testing.T, client *network.NSQClient) {
	err := client.PauseAllTopics()
	require.Nil(t, err)

	err = client.PauseAllChannels()
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	for _, topic := range stats.Topics {
		assert.True(t, topic.Paused)
		for _, channel := range topic.Channels {
			assert.True(t, channel.Paused)
		}
	}
}

func testUnpauseAll(t *testing.T, client *network.NSQClient) {
	err := client.UnpauseAllTopics()
	require.Nil(t, err)

	err = client.UnpauseAllChannels()
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	for _, topic := range stats.Topics {
		assert.False(t, topic.Paused)
		for _, channel := range topic.Channels {
			assert.False(t, channel.Paused)
		}
	}
}

func testEmptyAll(t *testing.T, client *network.NSQClient) {
	err := client.EmptyAllTopics()
	require.Nil(t, err)

	err = client.EmptyAllChannels()
	require.Nil(t, err)

	// Check that the above operations worked
	stats, err := client.GetStats()
	require.Nil(t, err)
	require.NotNil(t, stats)

	for _, topic := range stats.Topics {
		assert.EqualValues(t, 0, topic.Depth)
		for _, channel := range topic.Channels {
			assert.EqualValues(t, 0, channel.Depth)
		}
	}
}
