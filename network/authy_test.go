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

	// We can't test registration and push notifications
	// in unit tests, but we do want to make sure our
	// client is properly constructed. If so, it will return
	// the following errors with Authy disabled (instead of
	// panicking).
	response, err := client.RegisterUser("homer@example.com", 1, "302-555-1212")
	assert.Empty(t, response)
	assert.Equal(t, network.ErrAuthyDisabled, err)

	ok, err := client.AwaitOneTouch("homer@example.com", "no-id")
	assert.False(t, ok)
	assert.Equal(t, network.ErrAuthyDisabled, err)
}
