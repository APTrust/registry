package admin_api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/APTrust/registry/web/webui"
	"github.com/gin-gonic/gin"
)

// ObjectBatchDeleteParams contains info about which objects to
// delete in an object batch delete operation.
//
// We use this struct for two reasons:
//
// 1. JSON is easier to craft than a form with a thousand values.
//
//  2. Because the httptest library is lame and cannot properly
//     create multiple form values with the same name due to a
//     problem with the underlying github.com/ajg/form library.
//     It turns all params into a flat map. This is documented.
//     So instead of getting objectIds = [1,2,3,4] as url.Values
//     would give it to us, we get objectIds.0 = 1, objectIds.1 = 2,
//     objectIds.3 = 2, etc. That's worthless in a testing library
//     that needs to be able to pass values in the standard format
//     that the back end expects.
type ObjectBatchDeleteParams struct {
	InstitutionID int64   `json:"institutionId"`
	RequestorID   int64   `json:"requestorId"`
	ObjectIDs     []int64 `json:"objectIds"`
	SecretKey     string  `json:"secretKey"`
}

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
// Note that becaue this is part of the admin API, access to this call
// is restricted to APTrust admins.
//
// POST /objects/init_batch_delete
func IntellectualObjectInitBatchDelete(c *gin.Context) {
	req := api.NewRequest(c)

	// We want to log this because this is a dangerous operation.
	// We should never hit this line unless the request was submitted
	// by an APTrust admin.
	common.Context().Log.Warn().Msgf("Got batch deletion request from user %s, IP address %s", req.CurrentUser.Email, c.RemoteIP())

	// If the batch deletion key is not set in the config, bail
	// because this is unsafe.
	if !common.LooksLikeUUID(common.Context().Config.BatchDeletionKey) {
		message := "Configuration setting for BatchDeletionKey is missing or invalid"
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Rejecting object batch deletion request: %s", message)
		api.AbortIfError(c, errors.New(message))
		return
	}

	// The request body will be JSON, not a normal post form.
	// See note on ObjectBatchDeleteParams above.
	// First, read the request body into a byte slice.
	requestJson, err := io.ReadAll(c.Request.Body)
	if api.AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Could not read JSON from request body: %v", err)
		return
	}

	// Now parse the json bytes.
	params := ObjectBatchDeleteParams{}
	err = json.Unmarshal(requestJson, &params)
	if api.AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Error parsing JSON from request body: %v", err)
		return
	}

	// OK, if the request doesn't include the secret key, reject it.
	// We don't want anyone maliciously deleting files.
	if params.SecretKey != common.Context().Config.BatchDeletionKey {
		message := "Request secret key does not match configuration's BatchDeletionKey"
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Rejecting object batch deletion request: %s", message)
		api.AbortIfError(c, common.ErrInvalidToken)
		return
	}

	// If we made it this far, we're going to proceed with the request.
	// Let the logs know.
	common.Context().Log.Warn().Msgf("IntellectualObjectInitBatchDelete: Creating batch deletion request on behalf of user %d for %d objects belonging to institution %d. Current user is %s.",
		params.RequestorID, len(params.ObjectIDs), params.InstitutionID, req.CurrentUser.Email)

	// Create the batch deletion request. Note that this will fail if
	// certain internal checks fail. E.g. RequestorID does not belong
	// to an inst admin, one or more files in the batch belongs to another
	// institution, or has already been deleted.
	del, err := webui.NewDeletionForObjectBatch(params.RequestorID, params.InstitutionID, params.ObjectIDs, req.BaseURL())
	if api.AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Creating batch deletion failed: %v", err)
		return
	}

	// Now create the alert email to the institutional admin so they
	// can review and approve or cancel the request. This last line
	// of human intervention ensures batch deletions won't happen
	// silently or without explicit approval.
	_, err = del.CreateRequestAlert()
	if api.AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("IntellectualObjectInitBatchDelete: Batch deletion created, but creation of confirmation email failed: %v", err)
		return
	}

	// Now send the JSON response. This will be a pretty hefty chunk of JSON,
	// since it will include an object record for each object in the batch.
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
