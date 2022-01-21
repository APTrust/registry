package network_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/require"
)

func TestSNSClient(t *testing.T) {
	aptContext := common.Context()

	// Note that SNSClient is disabled in testing, so this just
	// prints a message to the console.
	err := aptContext.SNSClient.SendSMS("867-5309", "Jenny, I got your number.")
	require.Nil(t, err)
}
