package network_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSmtpClient(t *testing.T) {
	client := network.NewSMTPClient(true, "us-east-1", "mail.example.com", "yoozer", "password", "sender@example.com", common.Context().Log)
	require.NotNil(t, client)
	assert.True(t, client.ServiceEnabled)
	assert.Equal(t, "sender@example.com", client.FromAddress)
}

func TestSend(t *testing.T) {
	err := common.Context().SMTPClient.Send("joe@example.com", "Hey Joe", "Where you going with that gun of yours?")
	require.NoError(t, err)
}

func TestFormatHeader(t *testing.T) {
	header := network.SMTPFormatHeader("To", "recipient@example.com")
	assert.Equal(t, "To: recipient@example.com\r\n", header)

	subject := network.SMTPFormatHeader("Subject", "Contains \r\n Newlines")
	assert.Equal(t, "Subject: Contains \n Newlines\r\n", subject)
}

func TestFormatMessage(t *testing.T) {
	message := network.SMTPFormatMessage("joe@example.com", "Hey Joe", "Where you going with that gun of yours?")
	assert.Equal(t, "To: joe@example.com\r\nSubject: Hey Joe\r\n\r\nWhere you going with that gun of yours?\r\n", string(message))
}
