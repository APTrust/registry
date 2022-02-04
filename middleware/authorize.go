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

func IsAPIRequest(c *gin.Context) bool {
	// If we're going to bypass CSRF, make sure user authenticated
	// via API token, and not by cookie because cookies can be
	// hijacked by XSS attacks. API token won't exist in the browser
	// context, and can't be hijacked by a malicious script.
	isApiAuthenticated, exists := c.Get("UserIsApiAuthenticated")
	common.Context().Log.Debug().Msgf("IsAPIRequest - IsAPIAuthenticated: %t", isApiAuthenticated)
	if !exists || !isApiAuthenticated.(bool) {
		return false
	}
	return IsAPIRoute(c)
}

// IsAPIRoute returns true if the requested route matches one of our
// API prefixes. This uses c.Request.URL.Path because c.FullPath() can
// return an empty string if the path does not match any known routes.
func IsAPIRoute(c *gin.Context) bool {
	log := common.Context().Log
	path := c.Request.URL.Path // c.FullPath()
	for _, prefix := range constants.APIPrefixes {
		if strings.HasPrefix(path, prefix) {
			log.Debug().Msgf("IsAPIRoute - YES - %s", path)
			return true
		}
	}
	log.Debug().Msgf("IsAPIRoute - NO - %s", path)
	return false
}

func showNotCheckedError(c *gin.Context, auth *ResourceAuthorization) {
	common.Context().Log.Error().Msgf(auth.GetError())
	errMsg := fmt.Sprintf("Missing authorization check for %s", c.FullPath())
	if IsAPIRequest(c) {
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
	if IsAPIRoute(c) {
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
