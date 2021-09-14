package memberapi

import (
	//"fmt"
	"net/http"
	//"time"

	"github.com/APTrust/registry/api/core"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
	//"github.com/APTrust/registry/forms"
	//"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// GenericFileIndex shows list of objects.
// GET /member-api/v3/files
func GenericFileIndex(c *gin.Context) {
	req := core.NewRequest(c)
	var files []*pgmodels.GenericFile
	pager, err := req.LoadResourceList(&files, "updated_at desc")
	if core.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, core.NewJsonList(files, pager))
}

// GenericFileShow returns the object with the specified id.
// GET /member-api/v3/files/show/:id
func GenericFileShow(c *gin.Context) {
	req := core.NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if core.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}
