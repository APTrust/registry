package common_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPTContext(t *testing.T) {
	ctx := common.Context()
	require.NotNil(t, ctx)
	require.NotNil(t, ctx.Config)
	require.NotNil(t, ctx.DB)
	require.NotNil(t, ctx.Log)
	require.NotNil(t, ctx.AuthyClient)
	require.NotNil(t, ctx.NSQClient)
	assert.NotEmpty(t, ctx.NSQClient.URL)
	require.NotNil(t, ctx.SNSClient)
}
