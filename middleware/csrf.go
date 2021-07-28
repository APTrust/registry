package middleware

import (
	"net/http"
	"net/url"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
	"golang.org/x/crypto/bcrypt"
)

var safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, _ := GetCSRFCookieToken(c)
		if !IsCSRFSafeMethod(c.Request.Method) && !ExemptFromAuth(c) {
			err := AssertSameOrigin(c)
			if err != nil {
				abortWithError(c, err)
			}
			requestToken := GetCSRFRequestToken(c)
			err = CompareCSRFTokens(requestToken, cookieToken)
			if err != nil {
				abortWithError(c, err)
			}
		}
		// Put an xor'ed csrf token into the context,
		// so controllers can add it to forms. We only do
		// this if this request has a logged-in user.
		// Otherwise, we get errors on path "/" and others
		// that don't require login.
		_, userLoggedIn := c.Get("CurrentUser")
		if userLoggedIn && cookieToken != "" {
			err := AddTokenToContext(c, cookieToken)
			if err != nil {
				abortWithError(c, err)
			}
		}
		c.Next()
	}
}

func abortWithError(c *gin.Context, err error) {
	common.Context().Log.Error().Msgf("CSRF Error: %v", err)
	templateVars := gin.H{
		"error":           err.Error(),
		"suppressSideNav": true,
		"suppressTopNav":  false,
	}
	c.HTML(http.StatusUnauthorized, "errors/show.html", templateVars)
	c.Abort()
}

func IsCSRFSafeMethod(method string) bool {
	return slice.Contains(safeMethods, method)
}

// GetCSRFRequestToken returns the token set in the request form or header.
func GetCSRFRequestToken(c *gin.Context) string {
	requestToken := c.Request.Header.Get(constants.CSRFHeaderName)
	if len(requestToken) == 0 {
		requestToken = c.Request.PostFormValue(constants.CSRFTokenName)
	}
	return requestToken
}

// GetCSRFCookieToken returns the csrf token set in the cookie.
func GetCSRFCookieToken(c *gin.Context) (string, error) {
	ctx := common.Context()
	value := ""
	cookie, err := c.Cookie(constants.CSRFCookieName)
	if err != nil {
		return value, err
	}
	if err = ctx.Config.Cookies.Secure.Decode(constants.CSRFCookieName, cookie, &value); err != nil {
		return "", common.ErrDecodeCookie
	}
	return value, nil
}

// AddTokenToContext adds a bcrypted version of the CSRF token to
// the context, so we can pass it into forms. We use a weak bcrypt
// here, with a cost of 2, because all we really need is to create
// enough entropy to thwart BREACH attacks.
//
// See http://breachattack.com/
func AddTokenToContext(c *gin.Context, cookieToken string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(cookieToken), 2)
	if err != nil {
		return err
	}
	c.Set("csrf_token", string(digest))
	return nil
}

func CompareCSRFTokens(requestToken, cookieToken string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(requestToken), []byte(cookieToken))
	if err != nil {
		// We want to log this, but we don't want to display this
		// to the user.
		common.Context().Log.Error().Msgf("bcrypt error csrf: %v", err)
		return common.ErrInvalidCSRFToken
	}
	return nil
}

func AssertSameOrigin(c *gin.Context) error {
	referer, err := url.Parse(c.Request.Referer())
	if err != nil || referer.String() == "" {
		return common.ErrMissingReferer
	}

	// Wrong scheme means possible man-in-the-middle when
	// switching from http to https. Wrong host means this
	// is a cross-origin request.
	scheme := common.Context().Config.HTTPScheme()
	host := c.Request.Host // host or host:port
	if referer.Scheme != scheme || referer.Host != host {
		return common.ErrCrossOriginReferer
	}
	return nil
}
