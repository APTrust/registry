package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// GenericFileIndex shows list of files. Unlike the admin API,
// the member API version returns a list of GenericFileView objects.
//
// GET /member-api/v3/files
func GenericFileIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var files []*pgmodels.GenericFileView
	pager, err := req.LoadResourceList(&files, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(files, pager))
}

// GenericFileShow returns the object with the specified id.
//
// GET /member-api/v3/files/show/:id
// GET /admin-api/v3/files/show/:id
func GenericFileShow(c *gin.Context) {
	req := api.NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}
