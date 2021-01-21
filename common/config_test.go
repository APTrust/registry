package common_test

import (
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
	assert.Equal(t, "apt_registry_test", config.DB.Name)
	assert.Equal(t, "dev_user", config.DB.User)
	assert.Equal(t, "password", config.DB.Password)
	assert.Equal(t, 5432, config.DB.Port)

	assert.True(t, strings.HasSuffix(config.Logging.File, "logs/registry_test.log"))
	assert.Equal(t, zerolog.WarnLevel, config.Logging.Level)

	assert.Equal(t, "localhost", config.Cookies.Domain)
	assert.Equal(t, 43200, config.Cookies.MaxAge)
	assert.Equal(t, "aptrust_session", config.Cookies.SessionCookie)
	assert.False(t, config.Cookies.HTTPSOnly)
}
