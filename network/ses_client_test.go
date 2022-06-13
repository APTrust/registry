package network_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/require"
)

func TestSESClient(t *testing.T) {
	aptContext := common.Context()
	err := aptContext.SESClient.Send("nobody@example.com", "Whassup?", "Heart emojis.")
	require.Nil(t, err)
}
