package network_test

import (
	"io/ioutil"
	"testing"

	"github.com/APTrust/registry/network"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthyClient(t *testing.T) {
	client := network.NewAuthyClient(false, "SecretKey",
		zerolog.New(ioutil.Discard))
	assert.NotNil(t, client)

	// Test authy when disabled.
	response, err := client.RegisterUser("homer@example.com", 1, "302-555-1212")
	assert.Empty(t, response)
	assert.Equal(t, network.ErrAuthyDisabled, err)

	ok, err := client.AwaitOneTouch("homer@example.com", "no-id")
	assert.False(t, ok)
	assert.Equal(t, network.ErrAuthyDisabled, err)

	// Test our mock authy client. This seems superfluous, but
	// we want to make sure it works for tests in web/webui.
	client = network.NewMockAuthyClient()
	assert.NotNil(t, client)

	response, err = client.RegisterUser("homer@example.com", 1, "302-555-1212")
	assert.NotEmpty(t, response)
	assert.Nil(t, err)

	ok, err = client.AwaitOneTouch("homer@example.com", "no-id")
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = client.AwaitOneTouch("homer@example.com", "fail")
	assert.False(t, ok)
	assert.Nil(t, err)
}
