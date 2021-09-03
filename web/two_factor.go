package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
)

// UserTwoFactorChoose shows a list of radio button options so a user
// can choose their two-factor auth method (Authy, Backup Code, SMS).
// We show this page only to users who have enabled two-factor auth.
//
// GET /users/2fa_choose/
func UserTwoFactorChoose(c *gin.Context) {
	req := NewRequest(c)
	c.HTML(http.StatusOK, "users/choose_second_factor.html", req.TemplateData)
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
	fmt.Println("OTP token:", token)

	message := fmt.Sprintf("Your Registry one time password is: %s", token)
	err = common.Context().SNSClient.SendSMS(req.CurrentUser.PhoneNumber, message)
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
	// Send approval request and wait for response.
	// On approval, redirect to dashboard.
	// On rejection or timeout, log user out and redirect to sign-in
	req := NewRequest(c)
	if req.CurrentUser.AuthyID == "" {
		AbortIfError(c, common.ErrNoAuthyID)
		return
	}
	ok, err := common.Context().AuthyClient.AwaitOneTouch(
		req.CurrentUser.Email, req.CurrentUser.AuthyID)
	if AbortIfError(c, err) {
		return
	}
	if ok {
		// User approved login request
		req.CurrentUser.AwaitingSecondFactor = false
		req.CurrentUser.EncryptedOTPSecret = ""
		err := req.CurrentUser.Save()
		if AbortIfError(c, err) {
			return
		}
		c.Redirect(http.StatusFound, "/dashboard")
	} else {
		// Message expired or login request was rejected by user
		c.Redirect(http.StatusFound, "/users/sign_out")
	}
}

// UserTwoFactorResend resends the SMS two-factor auth code and then
// re-displays TwoFactorEnter. This is a post, because we don't want
// hackers spamming us with GETs. The post form includes a CSRF token.
//
// POST /users/2fa_resend/
func UserTwoFactorResend(c *gin.Context) {
	// This may re-send a push or SMS message.
	// Need to track which option user selected.
}

// UserTwoFactorVerify verifies the SMS or backup code that the user
// entered on TwoFactorEnter.
//
// POST /users/2fa_verify/
func UserTwoFactorVerify(c *gin.Context) {
	req := NewRequest(c)
	otp := c.PostForm("otp")
	method := c.PostForm("two_factor_method")

	tokenIsValid := false
	req.TemplateData["twoFactorMethod"] = method

	user := req.CurrentUser

	if method == constants.TwoFactorSMS {
		tokenIsValid = common.ComparePasswords(user.EncryptedOTPSecret, otp)
	} else {
		// Compare to backup code
	}
	if !tokenIsValid {
		// TODO: increment failed login attempt count
		// User model needs FailedLogins and LockoutUntil.
		// Then we need logic to enforce the lockout.
		req.TemplateData["flash"] = "One-time password is incorrect. Try again."
		c.HTML(http.StatusOK, "users/enter_auth_token.html", req.TemplateData)
	} else {
		user.AwaitingSecondFactor = false
		user.EncryptedOTPSecret = ""
		err := user.Save()
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
	// Show form with user phone number and radios for preferred
	// option: Text or Authy OneTouch.
	//
	// Use forms.TwoFactorSetupForm
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
	// Save user phone number, if changed.
	// Set EnabledTwoFactor to true.
	// If user chose Authy:
	//    set AuthyStatus to constants.TwoFactorAuthy
	//    If user has not registered for Authy:
	//        register user for Authy
	//    Else
	//        Send authy one touch and wait for confirmation
	// Else
	//    set AuthyStatus to constants.TwoFactorSMS
	//    send SMS code and redirect to UserConfirmPhone

	// -----------------------------------------------------------------
	//
	// TODO: Break this up and refactor sections that are reused!
	//
	// -----------------------------------------------------------------

	req := NewRequest(c)
	user := req.CurrentUser

	oldPhone := user.PhoneNumber
	oldMethod := user.AuthyStatus
	// Bind submitted form values in case we have to
	// re-display the form with an error message.
	c.ShouldBind(user)

	// Make sure the phone number is valid.
	// If not, send user back to the form.
	valError := user.Validate()
	if valError != nil && len(valError.Errors) > 0 {
		common.Context().Log.Error().Msgf("User validation error while completing two-factor setup: %s", valError.Error())
		req.TemplateData["form"] = forms.NewTwoFactorSetupForm(user)
		c.HTML(http.StatusBadRequest, "users/init_2fa_setup.html", req.TemplateData)
		return
	}

	phoneChanged := !(user.PhoneNumber == oldPhone)
	methodChanged := !(user.AuthyStatus == oldMethod)
	needsConfirmation := phoneChanged || methodChanged

	if !phoneChanged && !methodChanged {
		// Do nothing. Just redirect to My Account with flash
		// message saying nothing changed.
		helpers.SetFlashCookie(c, "Your two-factor preferences remain unchanged.")
		c.HTML(http.StatusFound, "/users/my_account", req.TemplateData)
		return
	}

	if phoneChanged {
		user.ConfirmedTwoFactor = false
	}

	if user.AuthyStatus == constants.TwoFactorNone {
		// Turn off two-factor, but only clear out ConfirmedTwoFactor
		// if the user changed their phone number.
		user.EnabledTwoFactor = false
		user.ConfirmedTwoFactor = !phoneChanged
		helpers.SetFlashCookie(c, "Two-factor authentication has been turned off for your account.")
		c.HTML(http.StatusFound, "/users/my_account", req.TemplateData)
		return
	}

	err := user.Save()
	if AbortIfError(c, err) {
		return
	}

	if needsConfirmation && user.AuthyStatus == constants.TwoFactorAuthy {
		// If user has no Authy ID, run UserAuthyRegister
		//    then set user.ConfirmedTwoFactor
		// If user has an Authy ID but it's not confirmed await OneTouch
		//    then set user.ConfirmedTwoFactor
		// If user has confirmed Authy ID, redirect to My Account with
		// flash message confirming the change
	} else if needsConfirmation && user.AuthyStatus == constants.TwoFactorSMS {
		// Send SMS code and redirect to UserConfirmPhone
		token, err := req.CurrentUser.CreateOTPToken()
		if AbortIfError(c, err) {
			return
		}
		// For dev work. You'll need this token to log in.
		fmt.Println("OTP token:", token)
		message := fmt.Sprintf("Your Registry one time password is: %s", token)
		err = common.Context().SNSClient.SendSMS(req.CurrentUser.PhoneNumber, message)
		if AbortIfError(c, err) {
			return
		}
		c.HTML(http.StatusOK, "users/confirm_phone.html", req.TemplateData)
		return
	}
}

// UserConfirmPhone accepts the form from UserComplete2FASetup.
//
// POST /users/confirm_phone
func UserConfirmPhone(c *gin.Context) {
	// If OTP is correct:
	// Set 2fa enabled and confirmed when done
	// If incorrect, show form. After three failures, lock for 5 min.

	// -----------------------------------------------------------------
	//
	// TODO: Refactor duplicated token verification code
	//
	// -----------------------------------------------------------------

	req := NewRequest(c)
	otp := c.PostForm("otp")
	user := req.CurrentUser
	if common.ComparePasswords(user.EncryptedOTPSecret, otp) {
		user.EnabledTwoFactor = true
		user.ConfirmedTwoFactor = true
		err := user.Save()
		if AbortIfError(c, err) {
			return
		}
		helpers.SetFlashCookie(c, "Thank you for confirming your phone number. Next time you log in, we'll send a one-time password to your phone to complete the login process.")
		c.Redirect(http.StatusFound, "/users/my_account")
		return
	} else {
		// Try again
		helpers.SetFlashCookie(c, "Oops! That wasn't the right code. Try again.")
		c.HTML(http.StatusOK, "users/confirm_phone.html", req.TemplateData)
	}
}

// UserRegisterWithAuthy registers a user with Authy.
//
// POST /users/authy_register
func UserAuthyRegister(c *gin.Context) {
	// Make sure they don't already have an AuthyID.
	// Set 2fa enabled when done
	req := NewRequest(c)
	user := req.CurrentUser
	if user.AuthyID != "" {
		AbortIfError(c, common.ErrAlreadyHasAuthyID)
		return
	}
	countryCode, phone, err := user.CountryCodeAndPhone()
	if AbortIfError(c, err) {
		return
	}
	authyID, err := common.Context().AuthyClient.RegisterUser(
		user.Email, int(countryCode), phone)
	if AbortIfError(c, err) {
		return
	}
	user.AuthyID = authyID
	err = user.Save()
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "users/authy_registration_complete.html", req.TemplateData)
}

// UserGenerateBackupCodes generates a set of five new, random backup
// codes and displays them to the user.
//
// POST /users/backup_codes
func UserGenerateBackupCodes(c *gin.Context) {
	// Generate six backup codes
	// Encrypt and save to DB
	// Display to user
}
