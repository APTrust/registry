package webui_test

import (
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/APTrust/registry/web/webui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOTPTokenIsExpired(t *testing.T) {
	ok := time.Now().Add(10 * time.Second)
	assert.False(t, webui.OTPTokenIsExpired(ok))

	notOK := time.Now().Add(-2 * common.Context().Config.TwoFactor.OTPExpiration)
	assert.True(t, webui.OTPTokenIsExpired(notOK))
}

func TestUserCompleteSMSSetup(t *testing.T) {
	defer func() {
		testutil.Inst1User.AwaitingSecondFactor = false
		testutil.Inst1User.ClearOTPSecret()
	}()

	req := &webui.Request{
		CurrentUser: testutil.Inst1User,
	}

	testutil.Inst1User.ClearOTPSecret()
	require.Empty(t, testutil.Inst1User.EncryptedOTPSecret)
	require.Empty(t, testutil.Inst1User.EncryptedOTPSentAt)

	webui.UserCompleteSMSSetup(req)

	require.NotEmpty(t, testutil.Inst1User.EncryptedOTPSecret)
	require.NotEmpty(t, testutil.Inst1User.EncryptedOTPSentAt)
}

func TestUserTwoFactorChoose(t *testing.T) {
	// Sign in as two-factor user and make sure we get the choice page.
	wasEnabled := testutil.SmsUser.EnabledTwoFactor
	wasConfirmed := testutil.SmsUser.ConfirmedTwoFactor
	oldMethod := testutil.SmsUser.AuthyStatus
	defer func() {
		testutil.SmsUser.EnabledTwoFactor = wasEnabled
		testutil.SmsUser.ConfirmedTwoFactor = wasConfirmed
		testutil.SmsUser.AuthyStatus = oldMethod
		testutil.SmsUser.Save()
	}()

	testutil.SmsUser.EnabledTwoFactor = true
	testutil.SmsUser.ConfirmedTwoFactor = true
	testutil.SmsUser.AuthyStatus = constants.TwoFactorSMS
	require.Nil(t, testutil.SmsUser.Save())

	signInForm := map[string]string{
		"email":    testutil.SmsUser.Email,
		"password": "password",
	}

	client := testutil.GetAnonymousClient(t)
	resp := client.POST("/users/sign_in").WithForm(signInForm).Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)

	itemsOnChoosePage := []string{
		"csrf_token",
		"submitSecondFactor('authy')",
		"submitSecondFactor('sms')",
		"submitSecondFactor('backup')",
	}

	html := resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, itemsOnChoosePage)

	// A user without 2FA turned on should see the dashboard,
	// not the 2fa_choose page. Not sure how to test the url
	// to which we're redirected using httpexpect.  ??
	signInForm["email"] = testutil.Inst1User.Email
	client = testutil.GetAnonymousClient(t)
	resp = client.POST("/users/sign_in").WithForm(signInForm).Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)

	itemsOnDashboard := []string{
		"Recent Work Items",
		"Notifications",
	}
	html = resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, itemsOnDashboard)
}

func TestUserTwoFactorBackup(t *testing.T) {
	testutil.InitHTTPTests(t)

	expected := []string{
		"two_factor_method",
		"Backup Code",
	}

	html := testutil.Inst1UserClient.GET("/users/2fa_backup").
		Expect().Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, expected)
}

func TestUserTwoFactorGenerateSMS(t *testing.T) {
	testutil.InitHTTPTests(t)
	oldToken := testutil.Inst1User.EncryptedOTPSecret
	oldTimestamp := testutil.Inst1User.EncryptedOTPSentAt
	defer func() {
		testutil.Inst1User.AwaitingSecondFactor = false
		testutil.Inst1User.ClearOTPSecret()
	}()

	html := testutil.Inst1UserClient.POST("/users/2fa_sms").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusOK).Body().Raw()

	expectedStrings := []string{
		`type="hidden" name="two_factor_method" value="sms"`,
		"SMS Code",
		"csrf_token",
	}
	testutil.AssertMatchesAll(t, html, expectedStrings)

	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.NotEqual(t, oldToken, reloadedUser.EncryptedOTPSecret)
	assert.NotEqual(t, oldTimestamp, reloadedUser.EncryptedOTPSentAt)
}

func TestUserTwoFactorPush(t *testing.T) {
	// We can't test this without an Authy API key
	// and Authy user id AND a user with a phone to
	// respond to the push.
}

func TestUserTwoFactorVerify(t *testing.T) {
	targetURL := "/users/2fa_verify"
	failureStrings := []string{
		"One-time password is incorrect",
	}
	// These appear on the dashboard, to which user
	// is redirected on successful login.
	successStrings := []string{
		"Recent Work Items",
		"Notifications",
	}
	testSMSVerify(t, targetURL, successStrings, failureStrings)
}

func TestUserTwoInit2FASetup(t *testing.T) {
	testutil.InitHTTPTests(t)
	expected := []string{
		`name="PhoneNumber"`,
		`name="AuthyStatus"`,
		"confirmChange()",
	}
	html := testutil.Inst1UserClient.GET("/users/2fa_setup").
		Expect().Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, expected)
}

func TestUserTwoComplete2FASetup(t *testing.T) {
	// There are many success/failure paths here to test.
	// Each leaves the user in an altered state, so best
	// to use a throwaway user, or ensure we can revert user
	// to known state expected by other tests.
	//
	// - Changing to Authy (can't test fully unless we mock Authy)
	// - Changing to SMS   (partially tested above in TestCompleteSMSSetup)
	// - Changing to None
	//
	// These have been manually tested. Need time to write the
	// automated tests.
}

func TestUserConfirmPhone(t *testing.T) {
	targetURL := "/users/confirm_phone"
	failureStrings := []string{
		"Confirm Your Phone Number",
	}
	// These appear on the My Account page, where
	// user goes after successful confirmation.
	successStrings := []string{
		testutil.Inst1User.Name,
		"getAPIKey()",
	}
	testSMSVerify(t, targetURL, successStrings, failureStrings)
}

func TestUserAuthyRegister(t *testing.T) {
	// We can't test this without an Authy API key
	// and Authy user id AND a user with a phone to
	// respond to the push.
}

// This tests both backup code generation and verification.
func TestUserBackupCodes(t *testing.T) {
	testutil.InitHTTPTests(t)
	testutil.Inst1User.OTPBackupCodes = []string{}
	require.Nil(t, testutil.Inst1User.Save())

	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.Empty(t, reloadedUser.OTPBackupCodes)

	// Generate backup codes
	html := testutil.Inst1UserClient.POST("/users/backup_codes").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusOK).Body().Raw()

	backupCodes := extractBackupCodes(t, html)

	// Make sure they were generated
	reloadedUser, err = pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.Equal(t, 6, len(reloadedUser.OTPBackupCodes))

	// Make sure backup code can be verified
	testutil.Inst1UserClient.POST("/users/2fa_verify").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("two_factor_method", "Backup Code").
		WithFormField("otp", backupCodes[0]).
		Expect().Status(http.StatusOK)

	// Backup code should be deleted after use
	reloadedUser, err = pgmodels.UserByEmail(testutil.Inst1User.Email)
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

func testSMSVerify(t *testing.T, targetURL string, successStrings, failureStrings []string) {
	testutil.InitHTTPTests(t)
	defer func() {
		testutil.Inst1User.AwaitingSecondFactor = false
		testutil.Inst1User.ClearOTPSecret()
	}()
	otp, err := testutil.Inst1User.CreateOTPToken()
	require.Nil(t, err)

	// Normally, the controller sets this after it knows
	// the OTP was sent successfully. For testing, we have
	// to set it manually.
	testutil.Inst1User.EncryptedOTPSentAt = time.Now()
	require.Nil(t, testutil.Inst1User.Save())

	// Bad token should be rejected. This user will be signed out,
	// so don't use one of our reusable clients.
	client, csrfToken := testutil.InitClient(t, "user@inst1.edu")
	html := client.POST(targetURL).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, csrfToken).
		WithFormField("otp", "this token is not valid").
		WithFormField("two_factor_method", constants.TwoFactorSMS).
		Expect().Status(http.StatusBadRequest).Body().Raw()
	testutil.AssertMatchesAll(t, html, failureStrings)

	// Try again with a good token.
	otp, err = testutil.Inst1User.CreateOTPToken()
	require.Nil(t, err)
	testutil.Inst1User.EncryptedOTPSentAt = time.Now()
	require.Nil(t, testutil.Inst1User.Save())

	// We need to log this user back in...
	// Good token should be accepted.
	client, csrfToken = testutil.InitClient(t, "user@inst1.edu")
	html = client.POST(targetURL).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, csrfToken).
		WithFormField("otp", otp).
		WithFormField("two_factor_method", constants.TwoFactorSMS).
		Expect().Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, successStrings)

	// These fields should be cleared after successful
	// two-factor auth.
	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.Empty(t, reloadedUser.EncryptedOTPSecret)
	assert.Empty(t, reloadedUser.EncryptedOTPSentAt)
	assert.False(t, reloadedUser.AwaitingSecondFactor)
}
