package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// GenericFileIndex shows list of objects.
//
// GET /member-api/v3/files
// GET /admin-api/v3/files
func GenericFileIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var files []*pgmodels.GenericFileView
	pager, err := req.LoadResourceList(&files, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	// This sucks. A late hack to return storage records with
	// the files for a call in the bag restorer, which needs to know
	// where to find these files in preservation storage.
	// This option is not documented for the member API, and shouldn't
	// be, as we may move this functionality elsewhere when the dust
	// settles.
	if c.Query("include_storage_records") == "true" {
		for _, gf := range files {
			query := pgmodels.NewQuery().Where("generic_file_id", "=", gf.ID)
			recs, err := pgmodels.StorageRecordSelect(query)
			if api.AbortIfError(c, err) {
				return
			}
			gf.StorageRecords = recs
		}
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
