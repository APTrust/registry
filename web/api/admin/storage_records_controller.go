package admin_api

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
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
	pager, err := req.LoadResourceList(&storageRecords, "id", "desc")
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

// StorageRecordCreate creates a new StorageRecord. We only do this when
// ingesting a newer a version of a previously ingested file.
// On first ingest, we call GenericFileCreateBatch, and the
// initial checksum is saved there, as part of a batch transaction.
//
// POST /admin-api/v3/storage_records/create/:institution_id
func StorageRecordCreate(c *gin.Context) {
	gf, err := CreateStorageRecord(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gf)
}

func CreateStorageRecord(c *gin.Context) (*pgmodels.StorageRecord, error) {
	req := api.NewRequest(c)
	sr, err := StorageRecordFromJson(req)
	if err != nil {
		return nil, err
	}
	err = sr.Save()
	return sr, err
}

// StorageRecordFromJson returns the StorageRecord from the
// JSON in the request body and the existing file record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing file can't be found, or if
// changes made to the existing object are not allowed.
func StorageRecordFromJson(req *api.Request) (*pgmodels.StorageRecord, error) {
	sr := &pgmodels.StorageRecord{}
	err := req.GinContext.BindJSON(sr)
	if err != nil {
		return nil, err
	}
	// Updating storageRecords is not allowed.
	if sr.ID != 0 {
		return nil, common.ErrNotSupported
	}
	gf, err := pgmodels.GenericFileByID(sr.GenericFileID)
	if gf == nil {
		err = fmt.Errorf("Can't save storage record. GenericFile %d does not exist", sr.GenericFileID)
	}
	if err != nil {
		return nil, err
	}
	if gf.InstitutionID != req.Auth.ResourceInstID {
		err = fmt.Errorf("Can't save storage record. Institution ID mismatch")
		return nil, err
	}
	return sr, err
}
