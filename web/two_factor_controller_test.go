package web_test

import (
	"net/http"
	"regexp"
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
	wasEnabled := smsUser.EnabledTwoFactor
	wasConfirmed := smsUser.ConfirmedTwoFactor
	oldMethod := smsUser.AuthyStatus
	defer func() {
		smsUser.EnabledTwoFactor = wasEnabled
		smsUser.ConfirmedTwoFactor = wasConfirmed
		smsUser.AuthyStatus = oldMethod
		smsUser.Save()
	}()

	smsUser.EnabledTwoFactor = true
	smsUser.ConfirmedTwoFactor = true
	smsUser.AuthyStatus = constants.TwoFactorSMS
	require.Nil(t, smsUser.Save())

	signInForm := map[string]string{
		"email":    smsUser.Email,
		"password": "password",
	}

	client := getAnonymousClient(t)
	resp := client.POST("/users/sign_in").WithForm(signInForm).Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)

	itemsOnChoosePage := []string{
		"csrf_token",
		"submitSecondFactor('authy')",
		"submitSecondFactor('sms')",
		"submitSecondFactor('backup')",
	}

	html := resp.Body().Raw()
	AssertMatchesAll(t, html, itemsOnChoosePage)

	// A user without 2FA turned on should see the dashboard,
	// not the 2fa_choose page. Not sure how to test the url
	// to which we're redirected using httpexpect.  ??
	signInForm["email"] = inst1User.Email
	client = getAnonymousClient(t)
	resp = client.POST("/users/sign_in").WithForm(signInForm).Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)

	itemsOnDashboard := []string{
		"Recent Work Items",
		"Notifications",
	}
	html = resp.Body().Raw()
	AssertMatchesAll(t, html, itemsOnDashboard)
}

func TestUserTwoFactorBackup(t *testing.T) {
	initHTTPTests(t)

	expected := []string{
		"two_factor_method",
		"Backup Code",
	}

	html := instUserClient.GET("/users/2fa_backup").
		Expect().Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, expected)
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

	// Generate backup codes
	html := instUserClient.POST("/users/backup_codes").
		WithHeader("Referer", baseURL).
		WithFormField(constants.CSRFTokenName, instUserToken).
		Expect().Status(http.StatusOK).Body().Raw()

	backupCodes := extractBackupCodes(t, html)

	// Make sure they were generated
	reloadedUser, err = pgmodels.UserByEmail(inst1User.Email)
	require.Nil(t, err)
	assert.Equal(t, 6, len(reloadedUser.OTPBackupCodes))

	// Make sure backup code can be verified
	instUserClient.POST("/users/2fa_verify").
		WithHeader("Referer", baseURL).
		WithFormField(constants.CSRFTokenName, instUserToken).
		WithFormField("two_factor_method", "Backup Code").
		WithFormField("otp", backupCodes[0]).
		Expect().Status(http.StatusOK)

	// Backup code should be deleted after use
	reloadedUser, err = pgmodels.UserByEmail(inst1User.Email)
	require.Nil(t, err)
	assert.Equal(t, 5, len(reloadedUser.OTPBackupCodes))
}

func extractBackupCodes(t *testing.T, html string) []string {
	re := regexp.MustCompile(`<span class="backup-code">(\w+)</span>`)
	m := re.FindAllStringSubmatch(html, -1)
	require.True(t, len(m) > 0)
	require.True(t, len(m[0]) > 0)
	codes := make([]string, 6)
	for i, match := range m {
		codes[i] = match[1]
	}
	return codes
}
