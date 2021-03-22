package middleware

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := AuthorizeResource(c)
		if !auth.Checked {
			showNotCheckedError(c, auth)
			c.Abort()
			return
		}
		if !auth.Approved {
			showAuthFailedError(c, auth)
			c.Abort()
			return
		}
		c.Next()
	}
}

func showNotCheckedError(c *gin.Context, auth *ResourceAuthorization) {
	common.Context().Log.Error().Msgf(auth.GetError())
	c.HTML(http.StatusInternalServerError, "errors/show.html", gin.H{
		"suppressSideNav": true,
		"suppressTopNav":  false,
		"error":           fmt.Sprintf("Missing authorization check for %s", c.FullPath()),
	})
}

func showAuthFailedError(c *gin.Context, auth *ResourceAuthorization) {
	common.Context().Log.Error().Msgf(auth.GetNotAuthorizedMessage())
	c.HTML(http.StatusForbidden, "errors/show.html", gin.H{
		"suppressSideNav": true,
		"suppressTopNav":  false,
		"error":           fmt.Sprintf("Permission denied for %s", c.FullPath()),
	})
}
