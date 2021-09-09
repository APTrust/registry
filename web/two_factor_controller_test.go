package web_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOTPTokenIsExpired(t *testing.T) {
	ok := time.Now().Add(10 * time.Second)
	assert.False(t, web.OTPTokenIsExpired(ok))

	notOK := time.Now().Add(-2 * common.Context().Config.TwoFactor.OTPExpiration)
	assert.True(t, web.OTPTokenIsExpired(notOK))
}

func TestUserCompleteSMSSetup(t *testing.T) {
	defer func() { inst1User.ClearOTPSecret() }()

	req := &web.Request{
		CurrentUser: inst1User,
	}

	inst1User.ClearOTPSecret()
	require.Empty(t, inst1User.EncryptedOTPSecret)
	require.Empty(t, inst1User.EncryptedOTPSentAt)

	web.UserCompleteSMSSetup(req)

	require.NotEmpty(t, inst1User.EncryptedOTPSecret)
	require.NotEmpty(t, inst1User.EncryptedOTPSentAt)
}
