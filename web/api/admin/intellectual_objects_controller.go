package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// IntellectualObjectCreate creates a new object record.
//
// POST /admin-api/v3/objects/create/:institution_id
func IntellectualObjectCreate(c *gin.Context) {
	obj, err := CreateOrUpdateObject(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, obj)
}

// IntellectualObjectUpdate updates an existing intellectual
// object record.
//
// PUT /admin-api/v3/objects/update/:id
func IntellectualObjectUpdate(c *gin.Context) {
	obj, err := CreateOrUpdateObject(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, obj)
}

// IntellectualObjectDelete marks an object record as deleted.
//
// DELETE /admin-api/v3/objects/delete/:id
func IntellectualObjectDelete(c *gin.Context) {
	// We should probably not allow the object to be deleted
	// unless all of its files have been deleted. Double check
	// the business logic in Pharos.
	//
	// Object deletion changes the state from "A" to "D".
	//
	// We should also ensure a Premis Event exists or is created
	// the documents who deleted this and when.
	//
	// Check the Pharos logic on that too. It may be the Go
	// worker's responsibility to ensure this, or it may be
	// registry's responsibility.
	//
	// Return the full object record.
	//
	// See https://github.com/APTrust/pharos/blob/03dda1f57a499c214691f6a739c22884ea2d2f4b/app/controllers/intellectual_objects_controller.rb#L171-L202
	//
	// And https://github.com/APTrust/pharos/blob/03dda1f57a499c214691f6a739c22884ea2d2f4b/app/models/intellectual_object.rb#L110-L127
	//
	req := api.NewRequest(c)
	obj, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	err = obj.Delete()
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateOrUpdateObject(c *gin.Context) (*pgmodels.IntellectualObject, error) {
	req := api.NewRequest(c)
	obj, err := IntellectualObjectFromJson(req)
	if err != nil {
		return nil, err
	}
	err = obj.Save()
	return obj, err
}

// IntellectualObjectFromJson returns the IntellectualObject from the
// JSON in the request body and the existing object record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing object can't be found, or if
// changes made to the existing object are not allowed.
func IntellectualObjectFromJson(req *api.Request) (*pgmodels.IntellectualObject, error) {
	submittedObject := &pgmodels.IntellectualObject{}
	err := req.GinContext.BindJSON(submittedObject)
	if err != nil {
		return submittedObject, err
	}
	err = req.AssertValidIDs(submittedObject.ID, submittedObject.InstitutionID)
	if err != nil {
		return submittedObject, err
	}
	if req.Auth.ResourceID > 0 {
		existingObject, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
		if err != nil {
			return submittedObject, err
		}
		CoerceObjectStorageOption(existingObject, submittedObject)
		err = existingObject.ValidateChanges(submittedObject)
	}
	return submittedObject, err
}

// CoerceObjectStorageOption forces submittedObject.StorageOption to match
// existingObject.StorageOption if existingObject.State is Active. The reason
// for this is documented in the special note under allowed storage option
// values at
// https://aptrust.github.io/userguide/bagging/#allowed-storage-option-values
func CoerceObjectStorageOption(existingObject, submittedObject *pgmodels.IntellectualObject) {
	if existingObject != nil && existingObject.State == constants.StateActive && existingObject.StorageOption != submittedObject.StorageOption {
		common.Context().Log.Warn().Msgf("Forcing storage option back to '%s' on IntellectualObject %d (%s)", existingObject.StorageOption, submittedObject.ID, submittedObject.Identifier)
		submittedObject.StorageOption = existingObject.StorageOption
	}
}
