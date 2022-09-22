package helpers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	//"github.com/gin-gonic/gin"
)

// CookieSetter defines the methods used by gin.Context
// to set cookies. We define it as an interface so we don't
// have to mock the entire sprawling gin.Context when testing.
//
// This also defines gin.Context's Get() method, which retrieves
// context values and SetSameSite, which is a security directive
// for cookie access.
type CookieSetter interface {
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetSameSite(samesite http.SameSite)
	Get(key string) (value interface{}, exists bool)
}

func SetCookie(c CookieSetter, name, value string) error {
	c.SetSameSite(http.SameSiteStrictMode)
	ctx := common.Context()
	if encoded, err := ctx.Config.Cookies.Secure.Encode(name, value); err == nil {
		c.SetCookie(
			name,
			encoded,
			ctx.Config.Cookies.MaxAge,
			"/",
			ctx.Config.Cookies.Domain,
			ctx.Config.Cookies.HTTPSOnly, // set only via HTTPS?
			true,                         // http-only: javascript can't access
		)
	}
	return nil
}

func DeleteCookie(c CookieSetter, name string) {
	ctx := common.Context()
	c.SetCookie(
		name,
		"",
		-1, // -1 means expire immediately, which equals delete
		"/",
		ctx.Config.Cookies.Domain,
		ctx.Config.Cookies.HTTPSOnly, // set only via HTTPS?
		true,                         // http-only: javascript can't access
	)
}

func SetSessionCookie(c CookieSetter, user *pgmodels.User) error {
	ctx := common.Context()
	id := fmt.Sprintf("%d", user.ID)
	return SetCookie(c, ctx.Config.Cookies.SessionCookie, id)
}

func DeleteSessionCookie(c CookieSetter) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.SessionCookie)
}

func SetFlashCookie(c CookieSetter, value string) error {
	ctx := common.Context()
	return SetCookie(c, ctx.Config.Cookies.FlashCookie, value)
}

func DeleteFlashCookie(c CookieSetter) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.FlashCookie)
}

func SetPrefsCookie(c CookieSetter, value string) error {
	ctx := common.Context()
	return SetCookie(c, ctx.Config.Cookies.PrefsCookie, value)
}

func DeletePrefsCookie(c CookieSetter) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.PrefsCookie)
}

func SetCSRFCookie(c CookieSetter) (string, error) {
	token := common.RandomToken()
	return token, SetCookie(c, constants.CSRFCookieName, token)
}

func DeleteCSRFCookie(c CookieSetter) {
	DeleteCookie(c, constants.CSRFCookieName)
}

func CurrentUser(c CookieSetter) *pgmodels.User {
	if currentUser, ok := c.Get("CurrentUser"); ok && currentUser != nil {
		return currentUser.(*pgmodels.User)
	}
	return nil
}
