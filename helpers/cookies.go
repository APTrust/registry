package helpers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

func SetSessionCookie(c *gin.Context, user *models.User) error {
	c.SetSameSite(http.SameSiteStrictMode)
	ctx := common.Context()
	id := fmt.Sprintf("%d", user.ID)
	if encoded, err := ctx.Config.Cookies.Secure.Encode(ctx.Config.Cookies.SessionCookie, id); err == nil {
		ctx.Log.Info().Msgf("Setting session cookie for %s", user.Email)
		c.SetCookie(
			ctx.Config.Cookies.SessionCookie,
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

func DeleteSessionCookie(c *gin.Context) {
	ctx := common.Context()
	user := CurrentUser(c)
	email := "unknown@not-logged-in.edu"
	if user != nil {
		email = user.Email
	}
	ctx.Log.Info().Msgf("Deleting session cookie for %s", email)
	c.SetCookie(
		ctx.Config.Cookies.SessionCookie,
		"",
		-1, // -1 means expire immediately, which equals delete
		"/",
		ctx.Config.Cookies.Domain,
		ctx.Config.Cookies.HTTPSOnly, // set only via HTTPS?
		true,                         // http-only: javascript can't access
	)
}

func CurrentUser(c *gin.Context) *models.User {
	if currentUser, ok := c.Get("CurrentUser"); ok && currentUser != nil {
		return currentUser.(*models.User)
	}
	return nil
}
