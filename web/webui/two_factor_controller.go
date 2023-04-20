package webui

import (
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
)

// UserTwoFactorChoose shows a list of radio button options so a user
// can choose their two-factor auth method (Authy, Backup Code, SMS).
// We show this page after a user has entered their email and password,
// if they have two-factor enabled. This is part of the login process,
// not part of the setup process.
//
// GET /users/2fa_choose/
func UserTwoFactorChoose(c *gin.Context) {
	req := NewRequest(c)
	c.HTML(http.StatusOK, "users/choose_second_factor.html", req.TemplateData)
}

// UserTwoFactorBackup shows the form on which the user can enter a
// backup code to complete two-factor authentication.
//
// GET /users/2fa_backup/
func UserTwoFactorBackup(c *gin.Context) {
	req := NewRequest(c)
	req.TemplateData["twoFactorMethod"] = "Backup Code"
	c.HTML(http.StatusOK, "users/enter_auth_token.html", req.TemplateData)
}

// UserTwoFactorGenerateSMS generates an OTP and sends it via SMS
// the user.
//
// POST /users/2fa_sms
func UserTwoFactorGenerateSMS(c *gin.Context) {
	req := NewRequest(c)
	token, err := req.CurrentUser.CreateOTPToken()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["twoFactorMethod"] = constants.TwoFactorSMS

	// For dev work. You'll need this token to log in.
	common.ConsoleDebug("OTP token: %s", token)

	user := req.CurrentUser

	message := fmt.Sprintf("Your Registry one time password is %s", token)
	err = common.Context().SNSClient.SendSMS(user.PhoneNumber, message)
	if AbortIfError(c, err) {
		return
	}
	user.EncryptedOTPSentAt = time.Now().UTC()
	err = user.Save()
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "users/enter_auth_token.html", req.TemplateData)
}

// UserTwoFactorPush initiates a push request to the user's authentication
// app asking them to approve the login. This method waits for a response
// from the authentication service. It's a POST to avoid GET spam.
// POST form includes CSRF token.
//
// POST /users/2fa_push/
func UserTwoFactorPush(c *gin.Context) {
	req := NewRequest(c)
	approved, err := userSendAuthyOneTouch(req)
	if AbortIfError(c, err) {
		return
	}
	if approved {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.Redirect(http.StatusFound, "/users/sign_out")
}

// UserTwoFactorVerify verifies the SMS or backup code that the user
// entered on TwoFactorEnter.
//
// POST /users/2fa_verify/
func UserTwoFactorVerify(c *gin.Context) {
	req := NewRequest(c)
	otp := c.PostForm("otp")
	method := c.PostForm("two_factor_method")

	var err error
	tokenIsValid := false
	req.TemplateData["twoFactorMethod"] = method

	user := req.CurrentUser

	if method == constants.TwoFactorSMS {
		if OTPTokenIsExpired(user.EncryptedOTPSentAt) {
			helpers.SetFlashCookie(c, "Your one-time password expired. Please sign in again.")
			c.Redirect(http.StatusFound, "/users/sign_out")
			return
		}
		tokenIsValid = common.ComparePasswords(user.EncryptedOTPSecret, otp)
	} else {
		tokenIsValid, err = userVerifyBackupCode(req, otp)
	}
	if AbortIfError(c, err) {
		return
	}
	if !tokenIsValid {
		// TODO: increment failed login attempt count
		// User model needs FailedLogins and LockoutUntil.
		// Then we need logic to enforce the lockout.
		msg := "Backup code is incorrect. Try again."
		if method == constants.TwoFactorSMS {
			msg = "One-time password is incorrect. Try again."
		}
		req.TemplateData["flash"] = msg
		c.HTML(http.StatusBadRequest, "users/enter_auth_token.html", req.TemplateData)
	} else {
		// Note that call to ClearOTPSecret saves user record to db.
		user.AwaitingSecondFactor = false
		err := user.ClearOTPSecret()
		if AbortIfError(c, err) {
			return
		}
		c.Redirect(http.StatusFound, "/dashboard")
	}
}

// UserInit2FASetup shows a page on which a user chooses their preferred
// two-factor auth method (or they can choose None and turn off two-factor).
//
// GET /users/2fa_setup
func UserInit2FASetup(c *gin.Context) {
	req := NewRequest(c)
	req.TemplateData["form"] = forms.NewTwoFactorSetupForm(req.CurrentUser)
	c.HTML(http.StatusOK, "users/init_2fa_setup.html", req.TemplateData)
}

// UserComplete2FASetup receives a form from UserInit2FASetup.
// If user chooses SMS, we need to send them a code via SMS and have
// them enter it here to confirm. If they choose Authy, we need to
// register them if they're not already registered.
//
// POST /users/2fa_setup
func UserComplete2FASetup(c *gin.Context) {
	req := NewRequest(c)
	user := req.CurrentUser
	prefs, err := NewTwoFactorPreferences(req)
	if err != nil {
		// errors.Is doesn't work here (???)
		if _, ok := err.(*common.ValidationError); ok {
			req.TemplateData["form"] = forms.NewTwoFactorSetupForm(user)
			c.HTML(http.StatusBadRequest, "users/init_2fa_setup.html", req.TemplateData)
		} else {
			AbortIfError(c, err)
		}
		return
	}

	if prefs.NothingChanged() {
		helpers.SetFlashCookie(c, "Your two-factor preferences remain unchanged.")
		c.Redirect(http.StatusFound, "/users/my_account")
		return
	}

	if prefs.PhoneChanged() {
		user.ConfirmedTwoFactor = false
	}

	if prefs.DoNotUseTwoFactor() {
		user.EnabledTwoFactor = false
		helpers.SetFlashCookie(c, "Two-factor authentication has been turned off for your account.")
		c.Redirect(http.StatusFound, "/users/my_account")
		return
	}

	err = user.Save()
	if AbortIfError(c, err) {
		return
	}

	if prefs.UseAuthy() {
		ok, err := userCompleteAuthySetup(req, prefs)
		if AbortIfError(c, err) {
			return
		}
		if ok {
			helpers.SetFlashCookie(c, "Your two-factor setup is complete. Next time you log in, you'll receive a push notification from Authy to complete the sign-in process.")
			c.Redirect(http.StatusFound, "/users/my_account")
			return
		} else {
			// User did not approve
			c.Redirect(http.StatusFound, "/users/sign_out")
			return
		}
	}

	if prefs.NeedsSMSConfirmation() {
		// Send SMS code and redirect to UserConfirmPhone
		err = UserCompleteSMSSetup(req)
		if AbortIfError(c, err) {
			return
		}
		c.HTML(http.StatusOK, "users/confirm_phone.html", req.TemplateData)
		return
	} else {
		if prefs.PhoneChanged() {
			helpers.SetFlashCookie(c, "Your phone number has been updated.")
		}
	}
	c.Redirect(http.StatusFound, "/users/my_account")
}

// UserConfirmPhone accepts the form from UserComplete2FASetup.
//
// POST /users/confirm_phone
func UserConfirmPhone(c *gin.Context) {
	req := NewRequest(c)
	otp := c.PostForm("otp")
	user := req.CurrentUser
	if common.ComparePasswords(user.EncryptedOTPSecret, otp) {
		user.EnabledTwoFactor = true
		user.ConfirmedTwoFactor = true
		err := user.ClearOTPSecret()
		if AbortIfError(c, err) {
			return
		}
		helpers.SetFlashCookie(c, "Thank you for confirming your phone number. Next time you log in, we'll send a one-time password to your phone to complete the login process.")
		c.Redirect(http.StatusFound, "/users/my_account")
		return
	} else {
		// Try again
		helpers.SetFlashCookie(c, "Oops! That wasn't the right code. Try again.")
		c.HTML(http.StatusBadRequest, "users/confirm_phone.html", req.TemplateData)
	}
}

// UserRegisterWithAuthy registers a user with Authy.
func UserAuthyRegister(req *Request) error {
	user := req.CurrentUser
	if user.AuthyID != "" {
		return common.ErrAlreadyHasAuthyID
	}
	countryCode, phone, err := user.CountryCodeAndPhone()
	if err != nil {
		return err
	}
	authyID, err := common.Context().AuthyClient.RegisterUser(
		user.Email, int(countryCode), phone)
	if err != nil {
		return err
	}
	user.AuthyID = authyID
	err = user.Save()
	if err != nil {
		return err
	}
	return nil
}

// UserGenerateBackupCodes generates a set of five new, random backup
// codes and displays them to the user.
//
// POST /users/backup_codes
func UserGenerateBackupCodes(c *gin.Context) {
	req := NewRequest(c)
	backupCodes := make([]string, 6)
	encCodes := make([]string, 6)
	for i := 0; i < 6; i++ {
		var err error
		code := common.RandomToken()
		backupCodes[i] = code[4:14]
		encCodes[i], err = common.EncryptPassword(backupCodes[i])
		if AbortIfError(c, err) {
			return
		}
	}
	req.CurrentUser.OTPBackupCodes = encCodes
	err := req.CurrentUser.Save()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["backupCodes"] = backupCodes
	req.TemplateData["showAsModal"] = true
	c.HTML(http.StatusOK, "users/backup_codes.html", req.TemplateData)
}

// Send an Authy OneTouch message to the user and await the
// response. This returns a boolean indicating whether the
// user approved the login.
func userSendAuthyOneTouch(req *Request) (bool, error) {
	if req.CurrentUser.AuthyID == "" {
		return false, common.ErrNoAuthyID
	}
	ok, err := common.Context().AuthyClient.AwaitOneTouch(
		req.CurrentUser.Email, req.CurrentUser.AuthyID)
	if err != nil {
		return false, err
	}
	if ok {
		// User approved login request
		req.CurrentUser.AwaitingSecondFactor = false
		req.CurrentUser.EncryptedOTPSecret = ""
		err := req.CurrentUser.Save()
		if err != nil {
			return false, err
		}
	}
	return ok, err
}

// Send an Authy push message to the user so they can confirm
// that Authy works.
func userConfirmAuthy(req *Request) (bool, error) {
	approved, err := userSendAuthyOneTouch(req)
	if approved {
		req.CurrentUser.ConfirmedTwoFactor = true
		err = req.CurrentUser.Save()
	}
	if err == nil && !approved {
		common.Context().Log.Warn().Msgf("User %s rejected Authy confirmation", req.CurrentUser.Email)
	}
	return approved, err
}

// This compares the user-supplied one-time password against all of the
// user's backup codes. If one matches, we delete that backup code and
// return true.
func userVerifyBackupCode(req *Request, otp string) (bool, error) {
	var err error
	tokenIsValid := false
	user := req.CurrentUser
	for _, encCode := range user.OTPBackupCodes {
		if common.ComparePasswords(encCode, otp) {
			user.OTPBackupCodes = slice.MinusStrings(user.OTPBackupCodes, []string{encCode})
			err = user.Save()
			tokenIsValid = true
			break
		}
	}
	return tokenIsValid, err
}

func userCompleteAuthySetup(req *Request, prefs *TwoFactorPreferences) (ok bool, err error) {
	if prefs.NeedsAuthyRegistration() {
		err = UserAuthyRegister(req)
		if err != nil {
			return false, err
		}
	}
	if prefs.NeedsAuthyConfirmation() {
		ok, err = userConfirmAuthy(req)
		if err != nil {
			return false, err
		}
	}
	return ok, err
}

func UserCompleteSMSSetup(req *Request) error {
	// Send SMS code and redirect to UserConfirmPhone
	user := req.CurrentUser
	token, err := user.CreateOTPToken()
	if err != nil {
		return err
	}
	// For dev work. You'll need this token to log in.
	common.ConsoleDebug("OTP token: %s", token)
	message := fmt.Sprintf("Your Registry one time password is %s", token)
	err = common.Context().SNSClient.SendSMS(user.PhoneNumber, message)
	if err != nil {
		return err
	}
	user.EncryptedOTPSentAt = time.Now().UTC()
	return user.Save()
}

func OTPTokenIsExpired(tokenSentAt time.Time) bool {
	expiration := tokenSentAt.Add(common.Context().Config.TwoFactor.OTPExpiration)
	return time.Now().After(expiration)
}
