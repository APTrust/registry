package middleware

import (
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

		c.Next()
	}
}

// TODO: Check permission based on user, resource, and institution id.
//
// - Check permission and set context vars saying permission was checked
//   and whether check passed.
// - If permission was not checked, or if error occurred, return status 500
//   with missing check error.
// - If not authorized, return 403.
// - If authorized, proceed to handler.
// - Index methods should still force user's institution ID into query.
// - Some resources, like checksums and storage records, will have to
//   get the InstitutionID from the parent object.
