package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

	// TODO: Ensure underlying logic, similar to IntelObj delete
	err = gf.Delete()
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}

// GenericFileCreate creates a new GenericFile.
//
// POST /admin-api/v3/files/create/:institution_id
func GenericFileCreate(c *gin.Context) {
	gf, err := CreateOrUpdateFile(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gf)
}

// GenericFileUpdate updates an existing GenericFile.
//
// PUT /admin-api/v3/files/update/:id
func GenericFileUpdate(c *gin.Context) {
	gf, err := CreateOrUpdateFile(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}

func CreateOrUpdateFile(c *gin.Context) (*pgmodels.GenericFile, error) {
	req := api.NewRequest(c)
	gf, err := GenericFileFromJson(req)
	if err != nil {
		return nil, err
	}
	err = gf.Save()
	return gf, err
}

// GenericFileFromJson returns the GenericFile from the
// JSON in the request body and the existing file record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing file can't be found, or if
// changes made to the existing object are not allowed.
func GenericFileFromJson(req *api.Request) (*pgmodels.GenericFile, error) {
	submittedFile := &pgmodels.GenericFile{}
	err := req.GinContext.BindJSON(submittedFile)
	if err != nil {
		return submittedFile, err
	}
	err = req.AssertValidIDs(submittedFile.ID, submittedFile.InstitutionID)
	if err != nil {
		return submittedFile, err
	}
	if req.Auth.ResourceID > 0 {
		existingFile, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
		if err != nil {
			return submittedFile, err
		}
		CoerceFileStorageOption(existingFile, submittedFile)
		err = existingFile.ValidateChanges(submittedFile)
	}
	return submittedFile, err
}

// CoerceFileStorageOption forces submittedFile.StorageOption to match
// existingFile.StorageOption if existingFile.State is Active. The reason
// for this is documented in the special note under allowed storage option
// values at
// https://aptrust.github.io/userguide/bagging/#allowed-storage-option-values
func CoerceFileStorageOption(existingFile, submittedFile *pgmodels.GenericFile) {
	if existingFile != nil && existingFile.State == constants.StateActive && existingFile.StorageOption != submittedFile.StorageOption {
		common.Context().Log.Warn().Msgf("Forcing storage option back to '%s' on GenericFile %d (%s)", existingFile.StorageOption, submittedFile.ID, submittedFile.Identifier)
		submittedFile.StorageOption = existingFile.StorageOption
	}
}
