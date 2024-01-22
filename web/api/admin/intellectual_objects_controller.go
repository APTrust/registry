package admin_api

import (
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/APTrust/registry/web/webui"
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

// IntellectualObjectInitRestore creates an object restoration request,
// which is really just a WorkItem that gets queued. Restoration can take
// seconds or hours, depending on where the object is stored and how big it is.
// POST /admin-api/v3/objects/init_restore/:id
func IntellectualObjectInitRestore(c *gin.Context) {
	req := api.NewRequest(c)
	obj, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}

	workItem, err := webui.InitObjectRestoration(obj, req.CurrentUser)
	if api.AbortIfError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, workItem)
}

// IntellectualObjectInitBatchDelete creates an deletion request for
// multiple objects. This request must be approved by an administrator
// at the depositing institution before the deletion will actually be queued.
//
// POST /objects/init_batch_delete/:id
func IntellectualObjectInitBatchDelete(c *gin.Context) {
	req := api.NewRequest(c)
	objectIDs, err := StringSliceToInt64Slice(c.Request.PostForm["objectID"])
	if api.AbortIfError(c, err) {
		return
	}
	institutionID, err := strconv.ParseInt(c.Request.PostFormValue("institutionID"), 10, 64)
	if api.AbortIfError(c, err) {
		return
	}
	requestorID, err := strconv.ParseInt(c.Request.PostFormValue("requestorID"), 10, 64)
	if api.AbortIfError(c, err) {
		return
	}

	common.Context().Log.Warn().Msgf("Creating batch deletion request on behalf of user %d for %d objects belonging to institution %d. Current user is %s.",
		requestorID, len(objectIDs), institutionID, req.CurrentUser.Email)

	del, err := webui.NewDeletionForObjectBatch(requestorID, institutionID, objectIDs, req.BaseURL())
	if api.AbortIfError(c, err) {
		return
	}
	_, err = del.CreateRequestAlert()
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, del.DeletionRequest)
}

// IntellectualObjectDelete marks an object record as deleted.
// It also creates a deletion premis event. Before it does any of
// that, it checks a number of pre-conditions. See the
// IntellectualObject model for more info.
//
// DELETE /admin-api/v3/objects/delete/:id
func IntellectualObjectDelete(c *gin.Context) {
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
	if err != nil {
		return nil, err
	}
	return obj, nil
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

		// Preservation services won't send the CreatedAt
		// timestamp, so we have to set this.
		submittedObject.CreatedAt = existingObject.CreatedAt

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

func StringSliceToInt64Slice(strSlice []string) ([]int64, error) {
	var err error
	ints := make([]int64, len(strSlice))
	for i, strValue := range strSlice {
		ints[i], err = strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			break
		}
	}
	return ints, err
}
