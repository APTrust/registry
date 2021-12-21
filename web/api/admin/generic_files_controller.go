package admin_api

import (
	"net/http"

	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// GenericFileDelete marks a generic file record as deleted.
// It also creates a deletion premis event. Before it does any of
// that, it checks a number of pre-conditions. See the
// GenericFile model for more info.
//
// DELETE /admin-api/v3/files/delete/:id
func GenericFileDelete(c *gin.Context) {
	req := api.NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	err = gf.Delete()
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}
