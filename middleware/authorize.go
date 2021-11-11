package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/gin-gonic/gin"
)

// Authorize ensures that a user is authorized to commit a specific
// action on a specific resource. This function uses a
// ResourceAuthorization struct to figure out what's being requested,
// what action the user wants to take on the resource, and whether
// the user has sufficient permissions.
//
// With the exception of the login page and static resources such as
// images, scripts, and stylesheets, all endpoints require an authorization
// check. Failure to perform the check is itself an error.
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := AuthorizeResource(c)
		c.Set("ResourceAuthorization", auth)
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

func isAPIRequest(path string) bool {
	for _, prefix := range constants.APIPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func showNotCheckedError(c *gin.Context, auth *ResourceAuthorization) {
	common.Context().Log.Error().Msgf(auth.GetError())
	errMsg := fmt.Sprintf("Missing authorization check for %s", c.FullPath())
	if isAPIRequest(c.Request.URL.Path) {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": errMsg,
		})
	} else {
		c.HTML(http.StatusInternalServerError, "errors/show.html", gin.H{
			"suppressSideNav": true,
			"suppressTopNav":  false,
			"error":           errMsg,
		})
	}
}

func showAuthFailedError(c *gin.Context, auth *ResourceAuthorization) {
	common.Context().Log.Error().Msgf(auth.GetNotAuthorizedMessage())
	errMsg := fmt.Sprintf("Permission denied for %s (institution %d).", c.FullPath(), auth.ResourceInstID)
	if auth.Error != nil {
		errMsg = fmt.Sprintf("%s %s", errMsg, auth.Error.Error())
	}
	if isAPIRequest(c.Request.URL.Path) {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": errMsg,
		})
	} else {
		c.HTML(http.StatusForbidden, "errors/show.html", gin.H{
			"suppressSideNav": true,
			"suppressTopNav":  false,
			"error":           errMsg,
		})
	}
}
