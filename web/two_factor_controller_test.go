package web_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/web"
	"github.com/stretchr/testify/assert"
)

func TestOTPTokenIsExpired(t *testing.T) {
	ok := time.Now().Add(10 * time.Second)
	assert.False(t, web.OTPTokenIsExpired(ok))

	notOK := time.Now().Add(-2 * common.Context().Config.TwoFactor.OTPExpiration)
	assert.True(t, web.OTPTokenIsExpired(notOK))
}
