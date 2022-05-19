package webui

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/network"
	"github.com/gin-gonic/gin"
)

// NsqShow returns NSQ stats.
//
// GET /nsq
func NsqShow(c *gin.Context) {
	req := NewRequest(c)
	stats, err := common.Context().NSQClient.GetStats()
	if AbortIfError(c, err) {
		return
	}
	nsqInitTopics(stats)
	req.TemplateData["stats"] = stats
	c.HTML(http.StatusOK, "nsq/show.html", req.TemplateData)
}

// NsqAdmin executes administrative actions against nsqd.
//
// POST /nsq/admind
func NsqAdmin(c *gin.Context) {
	var err error
	operation := c.PostForm("operation")
	targetType := c.PostForm("targetType")
	topicName := c.PostForm("topicName")
	channelName := c.PostForm("channelName")
	applyToAll := c.PostForm("applyToAll")
	nsqClient := common.Context().NSQClient
	if applyToAll == "true" {
		err = nsqDoToAll(nsqClient, operation, targetType)
	} else {
		err = nsqDoToOne(nsqClient, operation, targetType, topicName, channelName)
	}
	if AbortIfError(c, err) {
		return
	}
	helpers.SetFlashCookie(c, nsqSuccessMessage(operation, topicName, channelName, targetType, applyToAll))
	c.Redirect(http.StatusSeeOther, "/nsq")
}

// TODO: Move this to it's own endpoint.
//
// initTopics creates the initial NSQ topics and channels, so we'll have
// somthing to look at. Outside of our dev and test machines, these topics
// and channels are created by the queue workers that push items into NSQ.
func nsqInitTopics(stats *network.NSQStatsData) {
	ctx := common.Context()
	client := ctx.NSQClient
	allTopics := make([]string, len(constants.NonIngestTopics))
	copy(allTopics, constants.NonIngestTopics)
	for _, topicName := range constants.NSQIngestTopicFor {
		allTopics = append(allTopics, topicName)
	}
	for _, topicName := range constants.NSQIngestTopicFor {
		if stats.GetTopic(topicName) == nil {
			err := client.CreateTopic(topicName)
			if err != nil {
				ctx.Log.Warn().Msgf("Could not create NSQ topic %s: %v", topicName, err)
			}
			channelName := fmt.Sprintf("%s_worker_chan", topicName)
			err = client.CreateChannel(topicName, channelName)
			if err != nil {
				ctx.Log.Warn().Msgf("Could not create NSQ channel %s: %v", channelName, err)
			}
		}
	}
}

// nsqDoToAll performs the requested operation on all topics or channels.
func nsqDoToAll(nsqClient *network.NSQClient, operation, targetType string) error {
	var err error
	if targetType == "topic" {
		switch operation {
		case "pause":
			err = nsqClient.PauseAllTopics()
		case "unpause":
			err = nsqClient.UnpauseAllTopics()
		case "delete":
			err = nsqClient.DeleteAllTopics()
		case "empty":
			err = nsqClient.EmptyAllTopics()
		default:
			err = fmt.Errorf("unknown operation: %s", operation)
		}
	} else {
		switch operation {
		case "pause":
			err = nsqClient.PauseAllChannels()
		case "unpause":
			err = nsqClient.UnpauseAllChannels()
		case "delete":
			err = nsqClient.DeleteAllChannels()
		case "empty":
			err = nsqClient.EmptyAllChannels()
		default:
			err = fmt.Errorf("unknown operation: %s", operation)
		}
	}
	return err
}

func nsqDoToOne(nsqClient *network.NSQClient, operation, targetType, topicName, channelName string) error {
	var err error
	if targetType == "topic" {
		switch operation {
		case "pause":
			err = nsqClient.PauseTopic(topicName)
		case "unpause":
			err = nsqClient.UnpauseTopic(topicName)
		case "delete":
			err = nsqClient.DeleteTopic(topicName)
		case "empty":
			err = nsqClient.EmptyTopic(topicName)
		default:
			err = fmt.Errorf("unknown operation: %s", operation)
		}
	} else if targetType == "channel" {
		switch operation {
		case "pause":
			err = nsqClient.PauseChannel(topicName, channelName)
		case "unpause":
			err = nsqClient.UnpauseChannel(topicName, channelName)
		case "delete":
			err = nsqClient.DeleteChannel(topicName, channelName)
		case "empty":
			err = nsqClient.EmptyChannel(topicName, channelName)
		default:
			err = fmt.Errorf("unknown operation: %s", operation)
		}
	} else {
		err = fmt.Errorf("unknown target type: %s", targetType)
	}
	return err
}

func nsqSuccessMessage(operation, topicName, channelName, targetType, applyToAll string) string {
	target := topicName
	if channelName != "" {
		target = channelName
	}
	if applyToAll == "true" {
		if targetType == "channel" {
			target = "all channels"
		} else {
			target = "all topics"
		}
	}
	return fmt.Sprintf("Succeeded: %s %s", operation, target)
}
