package webui_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/network"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/APTrust/registry/web/webui"
	"github.com/pquerna/otp/totp"
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
	aptContext := common.Context()
	originalAuthyClient := aptContext.AuthyClient
	aptContext.AuthyClient = network.NewMockAuthyClient()

	origAuthyID := testutil.Inst1User.AuthyID
	testutil.Inst1User.AuthyID = "abc123"
	require.Nil(t, testutil.Inst1User.Save())

	defer func() {
		aptContext.AuthyClient = originalAuthyClient
		testutil.Inst1User.AuthyID = origAuthyID
		testutil.Inst1User.Save()
	}()

	testutil.Inst1UserClient.POST("/users/2fa_push").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusOK)
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
	// These have been manually tested. Automated tests cover only
	// a few cases.

	aptContext := common.Context()
	originalAuthyClient := aptContext.AuthyClient
	aptContext.AuthyClient = network.NewMockAuthyClient()

	defer func() {
		aptContext.AuthyClient = originalAuthyClient
	}()

	// Submit with no change
	expect := testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", testutil.Inst1User.PhoneNumber).
		WithFormField("AuthyStatus", testutil.Inst1User.AuthyStatus).
		Expect()
	html := expect.Body().Raw()
	assert.True(t, strings.Contains(html, "Your two-factor preferences remain unchanged."))

	// Submit with change to Phone number
	expect = testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", "+12223334444").
		WithFormField("AuthyStatus", "").
		Expect()
	html = expect.Body().Raw()
	assert.True(t, strings.Contains(html, "Your phone number has been updated."))

	// Submit with change to Phone number and AuthyStatus
	expect = testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", "+15556662627").
		WithFormField("AuthyStatus", constants.TwoFactorAuthy).
		Expect()
	html = expect.Body().Raw()
	assert.True(t, strings.Contains(html, "Your two-factor setup is complete."))
	assert.True(t, strings.Contains(html, "receive a push notification from Authy to complete the sign-in process."))

	// Submit with change to Phone number and SMS
	expect = testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", "+15556662888").
		WithFormField("AuthyStatus", constants.TwoFactorSMS).
		Expect()
	html = expect.Body().Raw()
	assert.True(t, strings.Contains(html, "Enter the code we just texted you into the box below."))

	// Turn off two-factor
	expect = testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", "+15556662888").
		WithFormField("AuthyStatus", constants.TwoFactorNone).
		Expect()
	html = expect.Body().Raw()
	assert.True(t, strings.Contains(html, "Two-factor authentication has been turned off for your account."))

	// Confirm a couple of settings after turning off two-factor auth.
	user, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.False(t, user.EnabledTwoFactor)
	assert.Empty(t, user.AuthyStatus)
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
	re := regexp.MustCompile(`<p class="backup-code">(\w+)</p>`)
	m := re.FindAllStringSubmatch(html, -1)
	require.True(t, len(m) > 0)
	require.True(t, len(m[0]) > 0)
	codes := make([]string, 6)
	for i, match := range m {
		codes[i] = match[1]
	}
	return codes
}

// TestUserGenerateTOTP hits GET /users/generate_totp. The first call should
// generate and persist a new authenticator-app secret for the user and
// display the setup page with the QR code and setup key. Subsequent calls
// should reuse the existing secret.
func TestUserGenerateTOTP(t *testing.T) {
	testutil.InitHTTPTests(t)

	origSecret := testutil.Inst1User.EncryptedAuthAppSecret
	defer func() {
		testutil.Inst1User.EncryptedAuthAppSecret = origSecret
		testutil.Inst1User.Save()
	}()

	// Clear the secret so the endpoint generates a new one.
	testutil.Inst1User.EncryptedAuthAppSecret = ""
	require.Nil(t, testutil.Inst1User.Save())

	html := testutil.Inst1UserClient.GET("/users/generate_totp").
		Expect().Status(http.StatusOK).Body().Raw()

	expected := []string{
		"Set Up Authenticator App",
		"Show Setup Key",
		`action="/users/validate_totp"`,
		"data:image/png;base64,",
	}
	testutil.AssertMatchesAll(t, html, expected)

	// Secret should have been generated and saved to the user.
	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.NotEmpty(t, reloadedUser.EncryptedAuthAppSecret)

	// The setup key shown in the page should match the saved secret.
	assert.True(t, strings.Contains(html, reloadedUser.EncryptedAuthAppSecret))

	// A second call should reuse the existing secret, not generate a new one.
	existingSecret := reloadedUser.EncryptedAuthAppSecret
	html = testutil.Inst1UserClient.GET("/users/generate_totp").
		Expect().Status(http.StatusOK).Body().Raw()
	assert.True(t, strings.Contains(html, existingSecret))

	reloadedUser, err = pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.Equal(t, existingSecret, reloadedUser.EncryptedAuthAppSecret)
}

// TestUserValidateTOTPView hits GET /users/validate_totp and verifies that
// the form for entering an authenticator-app code is rendered.
func TestUserValidateTOTPView(t *testing.T) {
	testutil.InitHTTPTests(t)

	expected := []string{
		"Enter One-Time Code From Your Authenticator App",
		`name="totpCode"`,
		`action="/users/validate_totp"`,
		"csrf_token",
	}
	html := testutil.Inst1UserClient.GET("/users/validate_totp").
		Expect().Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, expected)
}

// TestUserValidateTOTP hits POST /users/validate_totp. It verifies that an
// invalid code is rejected with an error message and that a valid code
// signs the user in and clears the awaiting-second-factor flag.
func TestUserValidateTOTP(t *testing.T) {
	testutil.InitHTTPTests(t)

	origSecret := testutil.Inst1User.EncryptedAuthAppSecret
	origEnabled := testutil.Inst1User.EnabledTwoFactor
	origConfirmed := testutil.Inst1User.ConfirmedTwoFactor
	origAwaiting := testutil.Inst1User.AwaitingSecondFactor
	defer func() {
		testutil.Inst1User.EncryptedAuthAppSecret = origSecret
		testutil.Inst1User.EnabledTwoFactor = origEnabled
		testutil.Inst1User.ConfirmedTwoFactor = origConfirmed
		testutil.Inst1User.AwaitingSecondFactor = origAwaiting
		testutil.Inst1User.Save()
	}()

	// Generate a TOTP secret and save it on the user so that the
	// controller can validate codes against it.
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.TOTPSecretIssuer,
		AccountName: testutil.Inst1User.Email,
	})
	require.Nil(t, err)
	testutil.Inst1User.EncryptedAuthAppSecret = key.Secret()
	require.Nil(t, testutil.Inst1User.Save())

	// An invalid code should be rejected with the expected error message.
	html := testutil.Inst1UserClient.POST("/users/validate_totp").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("totpCode", "000000").
		Expect().Status(http.StatusOK).Body().Raw()
	assert.True(t, strings.Contains(html, "Oops! That wasn't the right code. Please try again."))

	// A valid code should be accepted. The client follows the redirect
	// to /dashboard, so we expect a 200 with the dashboard content.
	validCode, err := totp.GenerateCode(key.Secret(), time.Now())
	require.Nil(t, err)

	successStrings := []string{
		"Recent Work Items",
		"Notifications",
	}
	html = testutil.Inst1UserClient.POST("/users/validate_totp").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("totpCode", validCode).
		Expect().Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, successStrings)

	// After a successful validation, the user should no longer be
	// awaiting a second factor and two-factor should be enabled and
	// confirmed.
	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.False(t, reloadedUser.AwaitingSecondFactor)
	assert.True(t, reloadedUser.EnabledTwoFactor)
	assert.True(t, reloadedUser.ConfirmedTwoFactor)
}

// TestUserValidateTOTPFirstConfirm hits POST /users/validate_totp with the
// firstConfirm form field set, which represents the user completing initial
// authenticator-app setup. On success the user should be redirected to the
// My Account page with a setup-complete flash message.
func TestUserValidateTOTPFirstConfirm(t *testing.T) {
	testutil.InitHTTPTests(t)

	origSecret := testutil.Inst1User.EncryptedAuthAppSecret
	origEnabled := testutil.Inst1User.EnabledTwoFactor
	origConfirmed := testutil.Inst1User.ConfirmedTwoFactor
	origAwaiting := testutil.Inst1User.AwaitingSecondFactor
	defer func() {
		testutil.Inst1User.EncryptedAuthAppSecret = origSecret
		testutil.Inst1User.EnabledTwoFactor = origEnabled
		testutil.Inst1User.ConfirmedTwoFactor = origConfirmed
		testutil.Inst1User.AwaitingSecondFactor = origAwaiting
		testutil.Inst1User.Save()
	}()

	// Set up a fresh TOTP secret for the user and make sure two-factor
	// starts out disabled so we can confirm the controller turns it on.
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.TOTPSecretIssuer,
		AccountName: testutil.Inst1User.Email,
	})
	require.Nil(t, err)
	testutil.Inst1User.EncryptedAuthAppSecret = key.Secret()
	testutil.Inst1User.EnabledTwoFactor = false
	testutil.Inst1User.ConfirmedTwoFactor = false
	require.Nil(t, testutil.Inst1User.Save())

	validCode, err := totp.GenerateCode(key.Secret(), time.Now())
	require.Nil(t, err)

	html := testutil.Inst1UserClient.POST("/users/validate_totp").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("totpCode", validCode).
		WithFormField("firstConfirm", "firstConfirm").
		Expect().Status(http.StatusOK).Body().Raw()
	assert.True(t, strings.Contains(html, "Your two-factor setup is complete."))

	reloadedUser, err := pgmodels.UserByEmail(testutil.Inst1User.Email)
	require.Nil(t, err)
	assert.True(t, reloadedUser.EnabledTwoFactor)
	assert.True(t, reloadedUser.ConfirmedTwoFactor)
	assert.False(t, reloadedUser.AwaitingSecondFactor)
}

// TestUserComplete2FASetupTOTP exercises the authenticator-app path through
// POST /users/2fa_setup. Selecting the TOTP method should redirect the user
// to the generate_totp page so they can scan the QR code.
func TestUserComplete2FASetupTOTP(t *testing.T) {
	testutil.InitHTTPTests(t)

	origStatus := testutil.Inst1User.AuthyStatus
	origPhone := testutil.Inst1User.PhoneNumber
	origSecret := testutil.Inst1User.EncryptedAuthAppSecret
	defer func() {
		testutil.Inst1User.AuthyStatus = origStatus
		testutil.Inst1User.PhoneNumber = origPhone
		testutil.Inst1User.EncryptedAuthAppSecret = origSecret
		testutil.Inst1User.Save()
	}()

	// Start with no existing authenticator-app secret so that the
	// generate_totp page will issue a new one.
	testutil.Inst1User.EncryptedAuthAppSecret = ""
	require.Nil(t, testutil.Inst1User.Save())

	// Submitting the 2FA setup form with AuthyStatus=totp should cause
	// the controller to redirect to /users/generate_totp, which in turn
	// renders the setup_authenticator_app.html page.
	html := testutil.Inst1UserClient.POST("/users/2fa_setup").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("PhoneNumber", testutil.Inst1User.PhoneNumber).
		WithFormField("AuthyStatus", constants.TwoFactorTOTP).
		Expect().Status(http.StatusOK).Body().Raw()

	expected := []string{
		"Set Up Authenticator App",
		"Show Setup Key",
		`action="/users/validate_totp"`,
	}
	testutil.AssertMatchesAll(t, html, expected)
}

func testSMSVerify(t *testing.T, targetURL string, successStrings, failureStrings []string) {
	testutil.InitHTTPTests(t)
	defer func() {
		testutil.Inst1User.AwaitingSecondFactor = false
		testutil.Inst1User.ClearOTPSecret()
	}()
	otp, err := testutil.Inst1User.CreateOTPToken()
	require.Nil(t, err)
	require.NotEmpty(t, otp)

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
