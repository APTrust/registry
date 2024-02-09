package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// Authenticate ensures the current user is logged in for all requests other
// than those going to "/" or static resources.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		LoadCookies(c)
		SetDefaultHeaders(c)
		if common.Context().Config.MaintenanceMode && c.Request.URL.Path != "/maintenance" {
			redirectTo := "/maintenance"
			if IsAPIRequest(c) || IsAPIRoute(c) {
				redirectTo += "?format=json"
			}
			c.Redirect(http.StatusTemporaryRedirect, redirectTo)
			c.Abort()
		}
		var user *pgmodels.User
		var err error
		if !ExemptFromAuth(c) {
			user, err = GetUser(c)
			if err != nil {
				respondToAuthError(c, err)
				c.Abort()
			} else {
				c.Set("CurrentUser", user)
			}
		}
		if forceCompletionOfPasswordChange(c, user) {
			logPasswordChangeIncomplete(c, user)
			c.HTML(http.StatusUnauthorized, "errors/show.html", gin.H{
				"suppressSideNav": true,
				"suppressTopNav":  true,
				"error":           "Please finish changing your password",
			})
			c.Abort()
		}
		if forceCompletionOfTwoFactorAuth(c, user) {
			log2FAIncomplete(c, user)
			c.Redirect(http.StatusFound, "/users/2fa_choose/")
			c.Abort()
		}
		c.Next()
	}
}

func GetUser(c *gin.Context) (user *pgmodels.User, err error) {
	apiUserEmail := c.Request.Header.Get(constants.APIUserHeader)
	apiUserKey := c.Request.Header.Get(constants.APIKeyHeader)
	if apiUserEmail != "" && apiUserKey != "" {
		return GetUserFromAPIHeaders(c)
	}
	return GetUserFromSession(c)
}

// GetUserFromSession returns the User for the current session.
func GetUserFromSession(c *gin.Context) (user *pgmodels.User, err error) {
	ctx := common.Context()
	cookie, err := c.Cookie(ctx.Config.Cookies.SessionCookie)
	if err == nil {
		value := ""
		if err = ctx.Config.Cookies.Secure.Decode(ctx.Config.Cookies.SessionCookie, cookie, &value); err != nil {
			ctx.Log.Error().Msgf("GetUserFromSession: Error decoding session cookie: %v", err)
			return nil, common.ErrDecodeCookie
		}
		var userID int64
		userID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			ctx.Log.Error().Msgf("GetUserFromSession: Session cookie contains non-numeric user id: %s", value)
			return nil, common.ErrWrongDataType
		}
		user, err = pgmodels.UserByID(userID)
		if err != nil {
			ctx.Log.Error().Msgf("GetUserFromSession: Got user id from session cookie but user lookup returned error: %v", err)
		}
	} else {
		// This is for a specific and recurrent auth error.
		// https://trello.com/c/rHKhPkau
		err = fmt.Errorf("%v - missing cookie is %s", err, ctx.Config.Cookies.SessionCookie)
	}
	return user, err
}

// GetUserFromAPIHeaders returns the current user based on the API
// auth headers.
func GetUserFromAPIHeaders(c *gin.Context) (user *pgmodels.User, err error) {
	ctx := common.Context()
	apiUserEmail := c.Request.Header.Get(constants.APIUserHeader)
	apiUserKey := c.Request.Header.Get(constants.APIKeyHeader)
	user, err = pgmodels.UserByEmail(apiUserEmail)
	if err != nil {
		ctx.Log.Error().Msgf("GetUserFromAPIHeaders: Attempt to look up user %s failed with error %v", apiUserEmail, err)
		return nil, err
	}
	if common.ComparePasswords(user.EncryptedAPISecretKey, apiUserKey) {
		// Set this because API requests bypass CSRF protection and
		// we want to ensure user passed valid auth headers. This
		// prevents a CSRF hijack where bad actor sends XHR PUT/POST
		// with valid user cookie. For API, we accept the token header
		// only, which won't be present in the browser environment and
		// therefore can't be hijacked.
		c.Set("UserIsApiAuthenticated", true)
		return user, nil
	}
	ctx.Log.Warn().Msgf("Invalid API token from user %s at %s.", apiUserEmail, c.Request.RemoteAddr)
	helpers.DeleteSessionCookie(c) // just to be extra safe
	return nil, common.ErrInvalidAPICredentials
}

// LoadCookies loads the user's flash and preference cookes into
// the request context.
func LoadCookies(c *gin.Context) error {
	ctx := common.Context()
	err := LoadCookie(c, ctx.Config.Cookies.FlashCookie)
	if err != nil && err != http.ErrNoCookie {
		return err
	}
	err = LoadCookie(c, ctx.Config.Cookies.PrefsCookie)
	if err != nil && err != http.ErrNoCookie {
		return err
	}
	return nil
}

// LoadCookie loads a cookie's value into the request context.
func LoadCookie(c *gin.Context, name string) error {
	ctx := common.Context()
	cookie, err := c.Cookie(name)
	if err != nil {
		return err
	}
	value := ""
	if err = ctx.Config.Cookies.Secure.Decode(name, cookie, &value); err != nil {
		return common.ErrDecodeCookie
	}
	c.Set(name, value)
	return nil
}

// SetDefaultHeaders sets headers that we want to include with every response.
// Note that it's OK for client to cache and store static resources such as
// images, scripts and stylesheets. Those are public resources containing no
// sensitive info. All other resources must use no-cache/no-store.
func SetDefaultHeaders(c *gin.Context) {
	p := c.FullPath()
	if !strings.HasPrefix(p, "/static") && !strings.HasPrefix(p, "/favicon") {
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Pragma", "no-store")
	}

	// This header tells the client to use HTTPS only. We can only set
	// this on systems that have certificates (staging, demo, production);
	// it won't work on dev/localhost setups, so check to see if config
	// says to run in HTTPSOnly mode.
	if common.Context().Config.Cookies.HTTPSOnly {
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000")
	}

	// Set these security headers on all responses.
	c.Writer.Header().Set("X-XSS-Protection", "1")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; font-src 'self' fonts.gstatic.com; style-src 'self' 'unsafe-inline' fonts.googleapis.com; script-src 'self' 'unsafe-inline'")
}

func forceCompletionOfPasswordChange(c *gin.Context, currentUser *pgmodels.User) bool {
	if currentUser == nil {
		return false // user isn't even signed in
	}
	p := c.FullPath()
	return currentUser.ResetPasswordToken != "" &&
		!strings.HasPrefix(p, "/users/change_password/") &&
		!strings.HasPrefix(p, "/errors/show/")
}

func logPasswordChangeIncomplete(c *gin.Context, currentUser *pgmodels.User) {
	common.Context().Log.Warn().Msgf("Password change incomplete. User %s tried to access URL [%s]. Forcing user to complete password change.", currentUser.Email, c.Request.RequestURI)
}

func forceCompletionOfTwoFactorAuth(c *gin.Context, currentUser *pgmodels.User) bool {
	if currentUser == nil {
		return false // user isn't even signed in
	}
	p := c.FullPath()
	return currentUser.AwaitingSecondFactor &&
		!strings.HasPrefix(p, "/users/2fa_backup") &&
		!strings.HasPrefix(p, "/users/2fa_choose") &&
		!strings.HasPrefix(p, "/users/2fa_push") &&
		!strings.HasPrefix(p, "/users/2fa_sms") &&
		!strings.HasPrefix(p, "/users/2fa_verify")
}

func log2FAIncomplete(c *gin.Context, currentUser *pgmodels.User) {
	common.Context().Log.Warn().Msgf("Two-factor auth incomplete. User %s tried to access URL [%s]. Forcing user to complete two-factor authentication.", currentUser.Email, c.Request.RequestURI)
}

func respondToAuthError(c *gin.Context, err error) {
	if IsAPIRequest(c) || IsAPIRoute(c) {
		msg := "API credentials are missing or invalid."
		common.Context().Log.Warn().Msgf("AuthError: %s", msg)
		obj := map[string]interface{}{
			"StatusCode": http.StatusUnauthorized,
			"Error":      msg,
		}
		c.JSON(http.StatusUnauthorized, obj)
	} else {
		common.Context().Log.Warn().Msgf("AuthError: %s. IP: %s, URL: %s, Agent: %s, Referer: %s", err.Error(), c.Request.RemoteAddr, c.Request.RequestURI, c.Request.UserAgent(), c.Request.Referer())
		c.HTML(http.StatusUnauthorized, "errors/show.html", gin.H{
			"suppressSideNav": true,
			"suppressTopNav":  true,
			"error":           "Please log in",
			"redirectURL":     fmt.Sprintf("/?requrl=%s", c.Request.URL),
		})
	}
}

func ExemptFromAuth(c *gin.Context) bool {
	p := c.FullPath()
	return p == "/" ||
		p == "/users/sign_in" ||
		p == "/users/sign_out" ||
		p == "/users/forgot_password" ||
		p == "/ui_components" ||
		p == "/maintenance" ||
		strings.HasPrefix(p, "/static") ||
		strings.HasPrefix(p, "/favicon") ||
		strings.HasPrefix(p, "/error") ||
		strings.HasPrefix(p, "/users/complete_password_reset/")
}
