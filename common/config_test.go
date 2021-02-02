package common_test

import (
	"os"
	"strings"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test assumes we're reading .env.test
func TestNewConfig(t *testing.T) {
	config := common.NewConfig()
	require.NotNil(t, config)
	assert.Equal(t, "postgres", config.DB.Driver)
	assert.Equal(t, "localhost", config.DB.Host)
	assert.Equal(t, 5432, config.DB.Port)
	assert.True(t, strings.HasSuffix(config.Logging.File, "registry_test.log"))
	assert.Equal(t, zerolog.WarnLevel, config.Logging.Level)

	// Local tests vs. Travis-CI tests.
	// Travis requires DB user 'postgres', which we don't do locally
	// because it's dangerous.
	if os.Getenv("APT_ENV") == "test" {
		assert.Equal(t, "apt_registry_test", config.DB.Name)
		assert.Equal(t, "dev_user", config.DB.User)
		assert.Equal(t, "password", config.DB.Password)
	} else if os.Getenv("APT_ENV") == "travis" {
		assert.Equal(t, "apt_registry_travis", config.DB.Name)
		assert.Equal(t, "postgres", config.DB.User)
		assert.Equal(t, "", config.DB.Password)
	} else {
		require.False(t, true, "Wrong APT_ENV environment for testing")
	}

	assert.Equal(t, "localhost", config.Cookies.Domain)
	assert.Equal(t, 43200, config.Cookies.MaxAge)
	assert.Equal(t, "aptrust_session", config.Cookies.SessionCookie)
	assert.False(t, config.Cookies.HTTPSOnly)
}
