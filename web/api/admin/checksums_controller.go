package admin_api

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// ChecksumCreate creates a new Checksum. We only do this when
// ingesting a newer a version of a previously ingested file.
// On first ingest, we call GenericFileCreateBatch, and the
// initial checksum is saved there, as part of a batch transaction.
//
// POST /admin-api/v3/checsums/create/:institution_id
func ChecksumCreate(c *gin.Context) {
	gf, err := CreateChecksum(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gf)
}

func CreateChecksum(c *gin.Context) (*pgmodels.Checksum, error) {
	req := api.NewRequest(c)
	cs, err := ChecksumFromJson(req)
	if err != nil {
		return nil, err
	}
	err = cs.Save()
	return cs, err
}

// ChecksumFromJson returns the Checksum from the
// JSON in the request body and the existing file record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing file can't be found, or if
// changes made to the existing object are not allowed.
func ChecksumFromJson(req *api.Request) (*pgmodels.Checksum, error) {
	cs := &pgmodels.Checksum{}
	err := req.GinContext.BindJSON(cs)
	if err != nil {
		return nil, err
	}
	// Updating checksums is not allowed.
	if cs.ID != 0 {
		return nil, common.ErrNotSupported
	}
	gf, err := pgmodels.GenericFileByID(cs.GenericFileID)
	if gf == nil {
		err = fmt.Errorf("Can't save checksum. GenericFile %d does not exist", cs.GenericFileID)
	}
	if err != nil {
		return nil, err
	}
	if gf.InstitutionID != req.Auth.ResourceInstID {
		err = fmt.Errorf("Can't save checksum. Institution ID mismatch")
		return nil, err
	}
	return cs, err
}
