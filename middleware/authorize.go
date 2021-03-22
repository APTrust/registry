package middleware

import (
	// "fmt"
	// "net/http"
	// "strconv"
	// "strings"

	// "github.com/APTrust/registry/common"
	// "github.com/APTrust/registry/constants"
	// "github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := AuthorizeResource(c)
		if !auth.Checked {
			// return 500 / Internal Server Error
			// auth MUST be checked
			// c.Abort()
		}
		if !auth.Approved {
			// return 403 / Forbidden
			// c.Abort()
		}
		c.Next()
	}
}
