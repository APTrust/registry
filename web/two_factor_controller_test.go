package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
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

func TestUserTwoFactorChoose(t *testing.T) {
	// Sign in as two-factor user and make sure we get the choice page.
}

func TestUserTwoFactorBackup(t *testing.T) {

}

func TestUserTwoFactorGenerateSMS(t *testing.T) {

}

func TestUserTwoFactorPush(t *testing.T) {

}

func TestUserTwoFactorVerify(t *testing.T) {

}

func TestUserTwoInit2FASetup(t *testing.T) {

}

func TestUserTwoComplete2FASetup(t *testing.T) {

}

func TestUserConfirmPhone(t *testing.T) {

}

func TestUserAuthyRegister(t *testing.T) {

}

func TestUserGenerateBackupCodes(t *testing.T) {
	initHTTPTests(t)
	inst1User.OTPBackupCodes = []string{}
	require.Nil(t, inst1User.Save())

	reloadedUser, err := pgmodels.UserByEmail(inst1User.Email)
	require.Nil(t, err)
	assert.Empty(t, reloadedUser.OTPBackupCodes)

	instUserClient.POST("/users/backup_codes").
		WithHeader("Referer", baseURL).
		WithFormField(constants.CSRFTokenName, instUserToken).
		Expect().Status(http.StatusOK)

	reloadedUser, err = pgmodels.UserByEmail(inst1User.Email)
	require.Nil(t, err)
	assert.Equal(t, 6, len(reloadedUser.OTPBackupCodes))
}
