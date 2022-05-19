package webui_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/web/testutil"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nsqOps = []string{
	"pause",
	"unpause",
	"empty",
	"delete",
}
var nsqTopics = []string{
	"admin_topic_0",
	"admin_topic_1",
	"admin_topic_2",
	"admin_topic_3",
	"admin_topic_4",
}
var nsqChannels = []string{
	"channel_0",
	"channel_1",
	"channel_2",
	"channel_3",
	"channel_4",
}

func TestNsqShow(t *testing.T) {
	items := []string{
		"NSQ",
		"Version",
		"running",
		"TCP",
		"Started",
		"Health: OK",
	}

	// Sysadmin can see NSQ stats
	html := testutil.SysAdminClient.GET("/nsq").
		Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	// Non-admin cannot see this page
	for _, client := range testutil.AllClients {
		if client != testutil.SysAdminClient {
			client.GET("/nsq").Expect().Status(http.StatusForbidden)
		}
	}
}

func TestNsqInit(t *testing.T) {
	// Sysadmin can hit this
	testutil.SysAdminClient.POST("/nsq/init").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
		Expect().
		Status(http.StatusOK)

	// Non-admin cannot see this page
	for _, client := range testutil.AllClients {
		if client != testutil.SysAdminClient {
			client.POST("/nsq/init").
				WithHeader("Referer", testutil.BaseURL).
				WithFormField(constants.CSRFTokenName, testutil.TokenFor[client]).
				Expect().Status(http.StatusForbidden)
		}
	}
}

func TestNsqAdmin(t *testing.T) {
	testAdminSingleOps(t)
	testAdminMultiOps(t)
	testNonAdminAnyOps(t)
}

func queueSomeItems(t *testing.T) {
	client := common.Context().NSQClient
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

func testAdminSingleOps(t *testing.T) {
	for _, op := range nsqOps {
		for _, topic := range nsqTopics {
			queueSomeItems(t)
			exp := testutil.SysAdminClient.POST("/nsq/admin").
				WithHeader("Referer", testutil.BaseURL).
				WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
				WithFormField("operation", op).
				WithFormField("targetType", "topic").
				WithFormField("topicName", topic).
				WithFormField("channelName", "").
				WithFormField("applyToAll", "false").
				Expect()
			exp.Status(http.StatusOK)
			exp.Body().Contains("Succeeded: ")

			for _, channel := range nsqChannels {
				queueSomeItems(t)
				exp = testutil.SysAdminClient.POST("/nsq/admin").
					WithHeader("Referer", testutil.BaseURL).
					WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
					WithFormField("operation", op).
					WithFormField("targetType", "channel").
					WithFormField("topicName", topic).
					WithFormField("channelName", channel).
					WithFormField("applyToAll", "false").
					Expect()
				exp.Status(http.StatusOK)
				exp.Body().Contains("Succeeded: ")
			}
		}
	}
}

func testAdminMultiOps(t *testing.T) {
	for _, op := range nsqOps {
		queueSomeItems(t)
		exp := testutil.SysAdminClient.POST("/nsq/admin").
			WithHeader("Referer", testutil.BaseURL).
			WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
			WithFormField("operation", op).
			WithFormField("targetType", "topic").
			WithFormField("topicName", "").
			WithFormField("channelName", "").
			WithFormField("applyToAll", "true").
			Expect()
		exp.Status(http.StatusOK)
		exp.Body().Contains("Succeeded: ")

		queueSomeItems(t)
		exp = testutil.SysAdminClient.POST("/nsq/admin").
			WithHeader("Referer", testutil.BaseURL).
			WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
			WithFormField("operation", op).
			WithFormField("targetType", "channel").
			WithFormField("topicName", "").
			WithFormField("channelName", "").
			WithFormField("applyToAll", "true").
			Expect()
		exp.Status(http.StatusOK)
		exp.Body().Contains("Succeeded: ")
	}
}

func testNonAdminAnyOps(t *testing.T) {
	for _, client := range testutil.AllClients {
		if client != testutil.SysAdminClient {
			client.POST("/nsq/admin").
				WithHeader("Referer", testutil.BaseURL).
				WithFormField(constants.CSRFTokenName, testutil.TokenFor[client]).
				Expect().Status(http.StatusForbidden)
		}
	}
}
