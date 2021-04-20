package network_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/require"
)

func TestNSQClientEnqueue(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(NsqOkHandler()))
	defer testServer.Close()

	aptContext := common.Context()
	aptContext.NSQClient.URL = testServer.URL

	err := aptContext.NSQClient.Enqueue("some_topic", 788)
	require.Nil(t, err)
}

// NsqOkhandler just returns Ok, which is the NSQ response when you queue
// an item in a topic.
func NsqOkHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
