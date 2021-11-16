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
	// Ensure the inst id in the JSON matches what's in the URL
	// Create the object.
	// Return the full object record.
	c.JSON(http.StatusOK, nil)
}

// IntellectualObjectUpdate updates an existing intellectual
// object record.
//
// PUT /admin-api/v3/objects/update/:id
func IntellectualObjectUpdate(c *gin.Context) {
	// Ensure the inst id in the JSON matches what's in the URL
	// Update the object, ensuring:
	//  - institution id can't change
	//  - storage option can't change
	// Return the full object record.

	c.JSON(http.StatusOK, nil)
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
	c.JSON(http.StatusOK, nil)
}

// IntellectualObjectFromJson returns the IntellectualObject from the
// JSON in the request body and the existing object record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing object can't be found, or if
// changes made to the existing object are not allowed.
func IntellectualObjectFromJson(req api.Request) (existingObject *pgmodels.IntellectualObject, submittedObject *pgmodels.IntellectualObject, err error) {
	err = req.GinContext.BindJSON(submittedObject)
	if err != nil {
		return existingObject, submittedObject, err
	}
	if req.Auth.ResourceID > 0 {
		existingObject, err = pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
		if err != nil {
			return existingObject, submittedObject, err
		}
		CoerceObjectStorageOption(existingObject, submittedObject)
		err = existingObject.ValidateChanges(submittedObject)
	}
	if err != nil {
		return existingObject, submittedObject, err
	}
	return existingObject, submittedObject, req.AssertValidIDs(submittedObject.ID, submittedObject.InstitutionID)
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
