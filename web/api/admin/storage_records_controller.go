package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// Delete is omitted on purpose. To delete a storage record, you
// have to delete the generic file, and there's a whole of requirements
// wrapped aroud that.

// StorageRecordIndex shows list of objects.
//
// GET /admin-api/v3/storage_records
func StorageRecordIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var storageRecords []*pgmodels.StorageRecord
	pager, err := req.LoadResourceList(&storageRecords, "datetime", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(storageRecords, pager))
}

// StorageRecordShow returns the object with the specified id.
//
// GET /admin-api/v3/storage_records/show/:id
func StorageRecordShow(c *gin.Context) {
	req := api.NewRequest(c)
	sr, err := pgmodels.StorageRecordByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, sr)
}
