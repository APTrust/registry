package helpers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string) error {
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

func DeleteCookie(c *gin.Context, name string) {
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

func SetSessionCookie(c *gin.Context, user *pgmodels.User) error {
	ctx := common.Context()
	id := fmt.Sprintf("%d", user.ID)
	return SetCookie(c, ctx.Config.Cookies.SessionCookie, id)
}

func DeleteSessionCookie(c *gin.Context) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.SessionCookie)
}

func SetFlashCookie(c *gin.Context, value string) error {
	ctx := common.Context()
	return SetCookie(c, ctx.Config.Cookies.FlashCookie, value)
}

func DeleteFlashCookie(c *gin.Context) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.FlashCookie)
}

func SetPrefsCookie(c *gin.Context, value string) error {
	ctx := common.Context()
	return SetCookie(c, ctx.Config.Cookies.PrefsCookie, value)
}

func DeletePrefsCookie(c *gin.Context) {
	ctx := common.Context()
	DeleteCookie(c, ctx.Config.Cookies.PrefsCookie)
}

func SetCSRFCookie(c *gin.Context) error {
	token := common.RandomToken()
	return SetCookie(c, constants.CSRFCookieName, token)
}

func DeleteCSRFCookie(c *gin.Context) {
	DeleteCookie(c, constants.CSRFCookieName)
}

func CurrentUser(c *gin.Context) *pgmodels.User {
	if currentUser, ok := c.Get("CurrentUser"); ok && currentUser != nil {
		return currentUser.(*pgmodels.User)
	}
	return nil
}
