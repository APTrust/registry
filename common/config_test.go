package common_test

import (
	"fmt"
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
	assert.True(t, strings.HasSuffix(config.Logging.File, fmt.Sprintf("registry_%s.log", os.Getenv("APT_ENV"))))
	assert.Equal(t, zerolog.DebugLevel, config.Logging.Level)

	// Local tests vs. Travis-CI tests.
	// Travis requires DB user 'postgres', which we don't do locally
	// because it's dangerous.
	if os.Getenv("APT_ENV") == "test" {
		assert.Equal(t, "apt_registry_test", config.DB.Name)
		assert.Equal(t, "dev_user", config.DB.User)
		assert.Equal(t, "password", config.DB.Password)
		assert.Equal(t, "test", config.EnvName)
	} else if os.Getenv("APT_ENV") == "dev" {
		assert.Equal(t, "apt_registry_development", config.DB.Name)
		assert.Equal(t, "dev_user", config.DB.User)
		assert.Equal(t, "password", config.DB.Password)
		assert.Equal(t, "dev", config.EnvName)
	} else if os.Getenv("APT_ENV") == "travis" {
		assert.Equal(t, "apt_registry_travis", config.DB.Name)
		assert.Equal(t, "postgres", config.DB.User)
		assert.Equal(t, "", config.DB.Password)
		assert.Equal(t, "travis", config.EnvName)
	} else {
		// TODO: Handle integration test env
		require.False(t, true, "Wrong APT_ENV environment for testing")
	}

	assert.False(t, config.Email.Enabled)
	assert.Equal(t, "help@aptrust.org", config.Email.FromAddress)

	assert.Equal(t, "localhost", config.Cookies.Domain)
	assert.Equal(t, 43200, config.Cookies.MaxAge)
	assert.Equal(t, "aptrust_session", config.Cookies.SessionCookie)
	assert.False(t, config.Cookies.HTTPSOnly)

	// Set sensitive data and make sure it's masked in JSON.
	config.BatchDeletionKey = "mask-me-batchdeletionkey"
	config.DB.Password = "mask-me-dbpassword"
	config.DB.User = "mask-me-dbuser"
	config.Email.SesUser = "mask-me-sesuser"
	config.Email.SesPassword = "mask-me-pwd-ses"
	config.Redis.Password = "mask-me-pwd-redis"
	config.TwoFactor.AuthyAPIKey = "mask-me-authy"

	jsonString, err := config.ToJSON()
	require.NoError(t, err)
	assert.Equal(t, expectedConfigJson, jsonString)
}

func TestConfigBucketQualifier(t *testing.T) {
	config := common.NewConfig()
	assert.Equal(t, ".test", config.BucketQualifier())

	config.Cookies.Domain = "demo.aptrust.org"
	assert.Equal(t, ".test", config.BucketQualifier())

	config.Cookies.Domain = "test.aptrust.org"
	assert.Equal(t, ".test", config.BucketQualifier())

	config.Cookies.Domain = "repo.aptrust.org"
	assert.Equal(t, "", config.BucketQualifier())

	config.Cookies.Domain = "staging.aptrust.org"
	assert.Equal(t, ".staging", config.BucketQualifier())

	config.Cookies.Domain = "localhost"
	assert.Equal(t, ".test", config.BucketQualifier())

	config.Cookies.Domain = "example.com"
	assert.Equal(t, ".test", config.BucketQualifier())
}

func TestIsTestEnv(t *testing.T) {
	config := common.NewConfig()

	config.EnvName = "test"
	assert.True(t, config.IsTestOrDevEnv())

	config.EnvName = "dev"
	assert.True(t, config.IsTestOrDevEnv())

	config.EnvName = "travis"
	assert.True(t, config.IsTestOrDevEnv())

	config.EnvName = "ci"
	assert.True(t, config.IsTestOrDevEnv())

	config.EnvName = "demo"
	assert.False(t, config.IsTestOrDevEnv())

	config.EnvName = "staging"
	assert.False(t, config.IsTestOrDevEnv())

	config.EnvName = "production"
	assert.False(t, config.IsTestOrDevEnv())
}

func TestHTTPScheme(t *testing.T) {
	config := common.NewConfig()

	config.EnvName = "test"
	assert.Equal(t, "http", config.HTTPScheme())

	config.EnvName = "dev"
	assert.Equal(t, "http", config.HTTPScheme())

	config.EnvName = "travis"
	assert.Equal(t, "http", config.HTTPScheme())

	config.EnvName = "ci"
	assert.Equal(t, "http", config.HTTPScheme())

	config.EnvName = "demo"
	assert.Equal(t, "https", config.HTTPScheme())

	config.EnvName = "staging"
	assert.Equal(t, "https", config.HTTPScheme())

	config.EnvName = "production"
	assert.Equal(t, "https", config.HTTPScheme())
}

var expectedConfigJson = `{
  "Cookies": {
    "Secure": {},
    "Domain": "localhost",
    "HTTPSOnly": false,
    "MaxAge": 43200,
    "SessionCookie": "aptrust_session",
    "FlashCookie": "aptrust_flash",
    "PrefsCookie": "aptrust_prefs"
  },
  "DB": {
    "Host": "localhost",
    "Name": "apt_registry_test",
    "User": "****ser",
    "Password": "****ord",
    "Port": 5432,
    "Driver": "postgres",
    "UseSSL": false
  },
  "EnvName": "test",
  "Logging": {
    "File": "/home/diamond/tmp/logs/registry_test.log",
    "Level": 0,
    "LogCaller": false,
    "LogToConsole": false,
    "LogSql": false
  },
  "NsqUrl": "http://localhost:4151",
  "TwoFactor": {
    "AuthyEnabled": false,
    "AuthyAPIKey": "****thy",
    "AWSRegion": "",
    "SMSEnabled": false,
    "OTPExpiration": 900000000000,
    "SNSEndpoint": ""
  },
  "Email": {
    "AWSRegion": "",
    "Enabled": false,
    "FromAddress": "help@aptrust.org",
    "SesUser": "****ser",
    "SesPassword": "****ses",
    "SesEndpoint": ""
  },
  "Redis": {
    "URL": "localhost:6379",
    "Password": "****dis",
    "DefaultDB": 0
  },
  "BatchDeletionKey": "****key"
}`
