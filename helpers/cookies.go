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
		fmt.Println("Setting session cookie")
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
	fmt.Println("Deleting session cookie")
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
