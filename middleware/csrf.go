package middleware

import (
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
		cookieToken, err := GetCSRFCookieToken(c)
		if err != nil {
			abortWithError(c, err)
		}
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
		// so controllers can add it to forms.
		AddTokenToContext(c, cookieToken)
		c.Next()
	}
}

func abortWithError(c *gin.Context, err error) {
	common.Context().Log.Error().Msgf("CSRF Error: %v", err)
	templateVars := gin.H{"error": err.Error()}
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

func AddTokenToContext(c *gin.Context, cookieToken string) {
	xorSalt := common.RandomToken()
	encToken := XorStrings(cookieToken, xorSalt)
	xorToken := fmt.Sprintf("%s$%s", encToken, xorSalt)
	c.Set("csrf_token", xorToken)
}

func CompareCSRFTokens(requestToken, cookieToken string) error {
	// Split request token at $
	// Xor with salt to decrypt
	// Compare cookie token to decrypted request token
	decryptedRequestToken, err := DecryptRequestToken(requestToken)
	if err != nil {
		return err
	}
	if decryptedRequestToken != cookieToken {
		return common.ErrInvalidToken
	}
	return nil
}

// DecryptRequestToken converts the xor'ed request token back
// to the plaintext version, which should match the csrf token
// in the cookie.
func DecryptRequestToken(requestToken string) (string, error) {
	parts := strings.Split(requestToken, "$")
	if len(parts) != 2 {
		return "", common.ErrInvalidToken
	}
	// First part is xor'ed token, second is salt/key
	return XorStrings(parts[0], parts[1]), nil
}

// XorString lets us alter the csrf token on each request to
// protect against BREACH attacks. See http://breachattack.com/
func XorStrings(input, key string) string {
	output := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		output[i] = input[i] ^ key[i%len(key)]
	}
	return fmt.Sprintf("%x", output)
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
