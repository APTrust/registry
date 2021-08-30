package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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
		helpers.DeleteFlashCookie(c)
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
		// increment failed login attempt count
		req.TemplateData["flash"] = "One-time password is incorrect. Try again."
		c.HTML(http.StatusOK, "users/enter_auth_token.html", req.TemplateData)
	} else {
		helpers.DeleteFlashCookie(c)
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

}

// UserComplete2FASetup receives a form from UserInit2FASetup.
// If user chooses SMS, we need to send them a code via SMS and have
// them enter it here to confirm. If they choose Authy, we need to
// register them if they're not already registered.
//
// POST /users/2fa_setup
func UserComplete2FASetup(c *gin.Context) {

}

// UserConfirmPhone accepts the form from UserComplete2FASetup.
//
// POST /users/confirm_phone
func UserConfirmPhone(c *gin.Context) {
	// If OTP is correct:
	// Set 2fa enabled and confirmed when done
	// If incorrect, show form. After three failures, lock for 5 min.
}

// UserRegisterWithAuthy registers a user with Authy.
func UserRegisterWithAuthy(c *gin.Context) {
	// Make sure they don't already have an AuthyID.
	// Set 2fa enabled when done
}

// UserGenerateBackupCodes generates a set of five new, random backup
// codes and displays them to the user.
func UserGenerateBackupCodes(c *gin.Context) {

}
