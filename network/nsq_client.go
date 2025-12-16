package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nsqio/nsq/nsqd"
	"github.com/rs/zerolog"
)

type NSQClient struct {
	URL    string
	logger zerolog.Logger
}

// NSQStatsData contains the important info returned by a call
// to NSQ's /stats endpoint, including the number of items in each
// topic and queue.
type NSQStatsData struct {
	Version   string             `json:"version"`
	Health    string             `json:"health"`
	StartTime uint64             `json:"start_time"`
	Topics    []*nsqd.TopicStats `json:"topics"`
	Info      *NSQInfo           `json:"info"`
}

type NSQInfo struct {
	Version          string `json:"version"`
	BroadcastAddress string `json:"broadcast_address"`
	Hostname         string `json:"hostname"`
	HttpPort         int    `json:"http_port"`
	TcpPort          int    `json:"tcp_port"`
	StartTime        int64  `json:"start_time"`
}

func (data *NSQStatsData) GetTopic(name string) *nsqd.TopicStats {
	for _, topic := range data.Topics {
		if topic.TopicName == name {
			return topic
		}
	}
	return nil
}

func (data *NSQStatsData) GetChannel(topicName, channelName string) *nsqd.ChannelStats {
	topic := data.GetTopic(topicName)
	if topic != nil {
		for _, channel := range topic.Channels {
			if channel.ChannelName == channelName {
				return &channel
			}
		}
	}
	return nil
}

func (data *NSQStatsData) ClientIsRunning(hostname string) bool {
	for _, topic := range data.Topics {
		for _, channel := range topic.Channels {
			for _, client := range channel.Clients {
				if strings.Contains(client.String(), hostname) {
					return true
				}
			}
		}
	}
	return false
}

// NewNSQClient returns a new NSQ client that will connect to the NSQ
// server and the specified url. The URL is typically available through
// Config.NsqdHttpAddress, and usually ends with :4151. This is
// the URL to which we post items we want to queue, and from
// which our workers read.
//
// Note that this client provides write access to queue, so we can
// add things. It does not provide read access. The workers do the
// reading.
func NewNSQClient(url string, logger zerolog.Logger) *NSQClient {
	return &NSQClient{
		URL:    url,
		logger: logger,
	}
}

// Enqueue posts data to NSQ, which essentially means putting it into a work
// topic. Param topic is the topic under which you want to queue something.
// For example, prepare_topic, fixity_topic, etc.
// Param workItemId is the id of the WorkItem record in Pharos we want to queue.
func (client *NSQClient) Enqueue(topic string, workItemID int64) error {
	idAsString := strconv.FormatInt(workItemID, 10)
	return client.EnqueueString(topic, idAsString)
}

// EnqueueString posts string data to the specified NSQ topic
func (client *NSQClient) EnqueueString(topic string, data string) error {
	url := fmt.Sprintf("%s/pub?topic=%s", client.URL, topic)
	client.logger.Info().Msgf("Enqueue: %s", url)
	resp, err := http.Post(url, "text/html", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return fmt.Errorf("Nsqd returned an error when queuing data: %v", err)
	}
	if resp == nil {
		return fmt.Errorf("No response from nsqd at '%s'. Is it running?", url)
	}

	// nsqd sends a simple OK. We have to read the response body,
	// or the connection will hang open forever.
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// NSQ response body is short. "OK" on success,
	// or about 100-200 bytes on error.
	client.logger.Info().Msgf("NSQ response from %s: [%d]  %s",
		url, resp.StatusCode, body)

	if resp.StatusCode != 200 {
		bodyText := "[no response body]"
		if len(body) > 0 {
			bodyText = string(body)
		}
		return fmt.Errorf("nsqd returned status code %d when attempting to queue data. "+
			"Response body: %s", resp.StatusCode, bodyText)
	}
	return nil
}

// get performs an HTTP get request and returns the body as a byte slice.
func (client *NSQClient) get(_url string) ([]byte, error) {
	resp, err := http.Get(_url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nsq at %s: %v", client.URL, err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading nsq response body: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("nsq returned status code %d, body: %s",
			resp.StatusCode, body)
	}
	return body, err
}

// doEmptyPost posts an empty request body to the specifid url.
func (client *NSQClient) doEmptyPost(_url string) error {
	resp, err := http.Post(_url, "text/html", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 204 {
		return fmt.Errorf("post to %s got unexpected status %d", _url, resp.StatusCode)
	}
	return nil
}

// GetStats allows us to get some basic stats from NSQ. The NSQ /stats endpoint
// returns a richer set of stats than what this fuction returns, but we only
// need some basic data for integration tests, so that's all we're parsing.
// The return value is a map whose key is the topic name and whose value is
// an NSQTopicStats object. NSQ is supposed to support topic_name as a query
// param, but this doesn't seem to be working in NSQ 0.3.0, so we're just
// returning stats for all topics right now. Also note that requests to
// /stats/ (with trailing slash) produce a 404.
func (client *NSQClient) GetStats() (*NSQStatsData, error) {
	url := fmt.Sprintf("%s/stats?format=json", client.URL)
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}
	stats := &NSQStatsData{}
	err = json.Unmarshal(body, stats)
	if err != nil {
		return nil, fmt.Errorf("error parsing nsq response json: %v", err)
	}
	info, err := client.GetInfo()
	if err != nil {
		return nil, err
	}
	stats.Info = info
	return stats, nil
}

// GetInfo returns basic info about nsqd, including hostname and port numbers.
func (client *NSQClient) GetInfo() (*NSQInfo, error) {
	url := fmt.Sprintf("%s/info?format=json", client.URL)
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}
	info := &NSQInfo{}
	err = json.Unmarshal(body, info)
	if err != nil {
		return nil, fmt.Errorf("error parsing nsq response json: %v", err)
	}
	return info, nil
}

// applyToAll is a generic function for pausing/unpausing all topics and/or
// channels.
func (client *NSQClient) applyToAll(action string) error {
	stats, err := client.GetStats()
	if err != nil {
		return err
	}
	errMsg := ""
	for _, topic := range stats.Topics {
		switch action {
		case "PauseTopics":
			err = client.PauseTopic(topic.TopicName)
		case "UnpauseTopics":
			err = client.UnpauseTopic(topic.TopicName)
		case "DeleteTopics":
			err = client.DeleteTopic(topic.TopicName)
		case "EmptyTopics":
			err = client.EmptyTopic(topic.TopicName)
		}
		for _, channel := range topic.Channels {
			switch action {
			case "PauseChannels":
				err = client.PauseChannel(topic.TopicName, channel.ChannelName)
			case "UnpauseChannels":
				err = client.UnpauseChannel(topic.TopicName, channel.ChannelName)
			case "DeleteChannels":
				err = client.DeleteChannel(topic.TopicName, channel.ChannelName)
			case "EmptyChannels":
				err = client.EmptyChannel(topic.TopicName, channel.ChannelName)
			}
		}
		if err != nil {
			errMsg += fmt.Sprintf("%s; ", err.Error())
		}
	}
	if len(errMsg) > 0 {
		err = fmt.Errorf(errMsg)
	}
	return err
}

// PauseTopic pauses the topic with the specified name.
func (client *NSQClient) PauseTopic(topicName string) error {
	_url := fmt.Sprintf("%s/topic/pause?topic=%s", client.URL, url.QueryEscape(topicName))
	return client.doEmptyPost(_url)
}

// UnpauseTopic unpauses the specified topic.
func (client *NSQClient) UnpauseTopic(topicName string) error {
	_url := fmt.Sprintf("%s/topic/unpause?topic=%s", client.URL, url.QueryEscape(topicName))
	return client.doEmptyPost(_url)
}

// PauseAllTopics pauses all topics.
func (client *NSQClient) PauseAllTopics() error {
	return client.applyToAll("PauseTopics")
}

// UnpauseAllTopics unpauses all topics.
func (client *NSQClient) UnpauseAllTopics() error {
	return client.applyToAll("UnpauseTopics")
}

// EmptyTopic empties the specified topic.
func (client *NSQClient) EmptyTopic(topicName string) error {
	_url := fmt.Sprintf("%s/topic/empty?topic=%s", client.URL, url.QueryEscape(topicName))
	return client.doEmptyPost(_url)
}

// EmptyAllTopics empties all topics.
func (client *NSQClient) EmptyAllTopics() error {
	return client.applyToAll("EmptyTopics")
}

// PauseChannel pauses the specified channel.
func (client *NSQClient) PauseChannel(topicName, channelName string) error {
	_url := fmt.Sprintf("%s/channel/pause?topic=%s&channel=%s", client.URL, url.QueryEscape(topicName), url.QueryEscape(channelName))
	return client.doEmptyPost(_url)
}

// UnpauseChannel unpauses the specified channel.
func (client *NSQClient) UnpauseChannel(topicName, channelName string) error {
	_url := fmt.Sprintf("%s/channel/unpause?topic=%s&channel=%s", client.URL, url.QueryEscape(topicName), url.QueryEscape(channelName))
	return client.doEmptyPost(_url)
}

// PauseAllChannels pauses all channels.
func (client *NSQClient) PauseAllChannels() error {
	return client.applyToAll("PauseChannels")
}

// UnpauseAllChannels unpauses all channels.
func (client *NSQClient) UnpauseAllChannels() error {
	return client.applyToAll("UnpauseChannels")
}

// EmptyChannel empties the specified channel.
func (client *NSQClient) EmptyChannel(topicName, channelName string) error {
	_url := fmt.Sprintf("%s/channel/empty?topic=%s&channel=%s", client.URL, url.QueryEscape(topicName), url.QueryEscape(channelName))
	return client.doEmptyPost(_url)
}

// EmptyAllChannels empties all channels.
func (client *NSQClient) EmptyAllChannels() error {
	return client.applyToAll("EmptyChannels")
}

// CreateTopic creates a topic. This is used only in testing.
func (client *NSQClient) CreateTopic(topicName string) error {
	_url := fmt.Sprintf("%s/topic/create?topic=%s", client.URL, url.QueryEscape(topicName))
	return client.doEmptyPost(_url)
}

// CreateChannel creates a channel in the specified topic. This is used
// only in testing.
func (client *NSQClient) CreateChannel(topicName, channelName string) error {
	_url := fmt.Sprintf("%s/channel/create?topic=%s&channel=%s", client.URL, url.QueryEscape(topicName), url.QueryEscape(channelName))
	return client.doEmptyPost(_url)
}

// DeleteTopic deletes a topic. This is used only in testing.
func (client *NSQClient) DeleteTopic(topicName string) error {
	_url := fmt.Sprintf("%s/topic/delete?topic=%s", client.URL, url.QueryEscape(topicName))
	return client.doEmptyPost(_url)
}

// DeleteChannel deletes a channel in the specified topic. This is used
// only in testing.
func (client *NSQClient) DeleteChannel(topicName, channelName string) error {
	_url := fmt.Sprintf("%s/channel/delete?topic=%s&channel=%s", client.URL, url.QueryEscape(topicName), url.QueryEscape(channelName))
	return client.doEmptyPost(_url)
}

// DeleteAllTopics deletes all topics
func (client *NSQClient) DeleteAllTopics() error {
	return client.applyToAll("DeleteTopics")
}

// DeleteAllChannels deletes all channels
func (client *NSQClient) DeleteAllChannels() error {
	return client.applyToAll("DeleteChannels")
}
