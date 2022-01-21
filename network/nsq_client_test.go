package network_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/require"
)

func TestNSQClient(t *testing.T) {
	aptContext := common.Context()

	err := aptContext.NSQClient.Enqueue("throwaway_test_topic", 788)
	require.Nil(t, err)

	err = aptContext.NSQClient.EnqueueString("throwaway_test_topic", "some/object/identifier")
	require.Nil(t, err)
}
