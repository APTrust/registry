package middleware

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
)

var safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, _ := GetCSRFCookieToken(c)
		if !IsCSRFSafeMethod(c.Request.Method) && !ExemptFromAuth(c) && !IsAPIRequest(c) {
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
			AddTokenToContext(c, cookieToken)
		}
		c.Next()
	}
}

func abortWithError(c *gin.Context, err error) {
	common.Context().Log.Error().Msgf("CSRF Error: %v", err)
	templateVars := gin.H{
		"error":           err.Error(),
		"suppressSideNav": true,
		"suppressTopNav":  true,
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

// XorStrings scrambles the CSRF token that appears in the header
// and in forms on each request. This is for BREACH attack prevention.
func XorStrings(input, key string) string {
	output := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		output[i] = input[i] ^ key[i%len(key)]
	}
	return string(output)
}

// AddTokenToContext adds an xor'ed version of the CSRF token to
// the context, so we can pass it into forms. This is to thwart
// BREACH attacks.
//
// See http://breachattack.com/
func AddTokenToContext(c *gin.Context, cookieToken string) {
	key := common.RandomToken()
	xored := XorStrings(cookieToken, key)
	requestToken := fmt.Sprintf("%s$%s", hex.EncodeToString([]byte(xored)), key)
	c.Set("csrf_token", string(requestToken))
}

func CompareCSRFTokens(requestToken, cookieToken string) error {
	logErr := common.Context().Log.Error().Msgf
	parts := strings.Split(requestToken, "$")
	if len(parts) < 2 {
		logErr("CSRF token is missing '$': %s", requestToken)
		return common.ErrInvalidCSRFToken
	}
	token := parts[0]
	key := parts[1]
	deHexedReqToken, err := hex.DecodeString(token)
	if err != nil {
		logErr("Cannot hex decode CSRF token '%s': %v", requestToken, err)
		return common.ErrInvalidCSRFToken
	}
	finalReqToken := XorStrings(string(deHexedReqToken), key)
	if finalReqToken != cookieToken {
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
	ctx := common.Context()
	scheme := ctx.Config.HTTPScheme()
	host := c.Request.Host // host or host:port
	if referer.Scheme != scheme || referer.Host != host {
		ctx.Log.Warn().Msgf("Rejecting cross-origin request for '%s'. This host is '%s://%s', but referrer is '%s://%s'", c.Request.URL.String(), scheme, host, referer.Scheme, referer.Host)
		return common.ErrCrossOriginReferer
	}
	return nil
}
