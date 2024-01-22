package admin_api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	admin_api "github.com/APTrust/registry/web/api/admin"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectShow(t *testing.T) {
	tu.InitHTTPTests(t)

	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)
	require.NotNil(t, obj)

	// Sysadmin can get this object through the admin API.
	resp := tu.SysAdminClient.GET("/admin-api/v3/objects/show/{id}", obj.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusOK)
	record := &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.InstitutionID, record.InstitutionID)

	// Sysadmin should also be able to find by identifier
	resp = tu.SysAdminClient.GET("/admin-api/v3/objects/show/{id}", obj.Identifier).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusOK)
	record = &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.Identifier, record.Identifier)

	// Non-admins should get an error message telling them to
	// use the Member API
	resp = tu.Inst1AdminClient.GET("/admin-api/v3/objects/show/{id}", obj.ID).Expect()
	resp.Status(http.StatusForbidden)
	assert.Equal(t, `{"error":"Permission denied for /admin-api/v3/objects/show/*id (institution 0). non-admins must use the member api"}`, resp.Body().Raw())
}

func TestObjectIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Admin can see this page.
	resp := tu.SysAdminClient.GET("/admin-api/v3/objects").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.IntellectualObjectList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 14, list.Count)
	assert.Equal(t, "/admin-api/v3/objects?page=3&per_page=5", list.Next)
	assert.Equal(t, "/admin-api/v3/objects?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Non-admins can't see this page
	resp = tu.Inst1UserClient.GET("/admin-api/v3/objects").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect()

	resp.Status(http.StatusForbidden)
	assert.Equal(t, `{"error":"Permission denied for /admin-api/v3/objects (institution 0). non-admins must use the member api"}`, resp.Body().Raw())
}

func TestObjectCreateUpdateDelete(t *testing.T) {

	if common.Context().Config.EnvName != "test" {
		// For security reasons, the deletion setup endpoint works
		// only in the test env. In all others, it returns an error.
		// Devs may sometimes run unit tests in dev mode, just for
		// the side effect of populating the DB. Skipping this test
		// when APT_ENV=dev prevents an error.
		//
		// See web/admin/api/dummy_controller.go
		fmt.Println("Skipping TestObjectCreateUpdateDelete because env is not test")
		return
	}

	tu.InitHTTPTests(t)
	obj := testObjectCreate(t)
	updatedObj := testObjectUpdate(t, obj)

	createObjectDeletionPreConditions(t, obj)
	testObjectDelete(t, updatedObj)
}

// The registry won't allow deletions without the pre-conditions
// below. In reality, the supporting object are created during
// actual workflows.
//
//   - Ingest event at ingest.
//   - Deletion request when a user clicks the delete object button
//     in the web UI.
//   - WorkItem when an inst admin has approved the deletion request.
//
// Here, we create them just so we can complete our test.
func createObjectDeletionPreConditions(t *testing.T, obj *pgmodels.IntellectualObject) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/prepare_object_delete/{id}", obj.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusOK)
}

func testObjectCreate(t *testing.T) *pgmodels.IntellectualObject {
	// Random objects use inst id 4 -> test.edu
	obj := pgmodels.RandomObject()
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(obj).Expect()
	resp.Status(http.StatusCreated)

	savedObj := &pgmodels.IntellectualObject{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), savedObj)
	require.Nil(t, err)
	assert.True(t, savedObj.ID > int64(0))
	assert.Equal(t, obj.Identifier, savedObj.Identifier)
	assert.Equal(t, obj.InstitutionID, savedObj.InstitutionID)
	assert.Equal(t, obj.BagName, savedObj.BagName)
	assert.Equal(t, obj.ETag, savedObj.ETag)
	assert.Equal(t, obj.StorageOption, savedObj.StorageOption)
	assert.NotEmpty(t, savedObj.CreatedAt)
	assert.NotEmpty(t, savedObj.UpdatedAt)
	return savedObj
}

func testObjectUpdate(t *testing.T, obj *pgmodels.IntellectualObject) *pgmodels.IntellectualObject {
	origUpdatedAt := obj.UpdatedAt
	copyOfObj := obj
	copyOfObj.Access = constants.AccessConsortia
	copyOfObj.Title = "Updated Title"
	copyOfObj.ETag = "UpdatedETag"

	resp := tu.SysAdminClient.PUT("/admin-api/v3/objects/update/{id}", obj.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(copyOfObj).Expect()
	resp.Status(http.StatusOK)

	updatedObj := &pgmodels.IntellectualObject{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), updatedObj)
	require.Nil(t, err)

	assert.Equal(t, copyOfObj.Access, updatedObj.Access)
	assert.Equal(t, copyOfObj.Title, updatedObj.Title)
	assert.Equal(t, copyOfObj.ETag, updatedObj.ETag)
	assert.InDelta(t, obj.CreatedAt.Unix(), updatedObj.CreatedAt.Unix(), 1)
	assert.True(t, updatedObj.UpdatedAt.After(origUpdatedAt))

	return updatedObj
}

func testObjectDelete(t *testing.T, obj *pgmodels.IntellectualObject) {
	resp := tu.SysAdminClient.DELETE("/admin-api/v3/objects/delete/{id}", obj.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusOK)

	deletedObj := &pgmodels.IntellectualObject{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), deletedObj)
	require.Nil(t, err)

	assert.Equal(t, obj.ID, deletedObj.ID)
	assert.Equal(t, constants.StateDeleted, deletedObj.State)

	// Make sure the state was actually saved
	savedObj, err := pgmodels.IntellectualObjectByID(obj.ID)
	require.Nil(t, err)
	assert.Equal(t, constants.StateDeleted, savedObj.State)

	// Test for deletion event
	event, err := obj.LastDeletionEvent()
	require.Nil(t, err)
	require.NotNil(t, event)
	require.True(t, event.DateTime.After(time.Now().UTC().Add(-5*time.Second)))
}

func TestObjectCreateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot create objects, even for their
	// own institutions.
	obj := pgmodels.RandomObject()
	obj.InstitutionID = tu.Inst1Admin.InstitutionID

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}

func TestObjectUpdateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot update objects, even for their
	// own institutions.
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/update/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/update/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}

func TestObjectDeleteUpdateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot delete objects, even for their
	// own institutions. (Not through the API anyway.)
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/delete/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/delete/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}

func TestCoerceObjectStorageOption(t *testing.T) {
	existingObj := &pgmodels.IntellectualObject{
		State:         constants.StateActive,
		StorageOption: constants.StorageOptionGlacierDeepOR,
	}
	submittedObj := &pgmodels.IntellectualObject{
		State:         constants.StateActive,
		StorageOption: constants.StorageOptionWasabiVA,
	}
	assert.NotEqual(t, existingObj.StorageOption, submittedObj.StorageOption)

	// Per APTrust rules, new obj storage option must be set to match
	// existing object storage option if existing obj is still active.
	// https://aptrust.github.io/userguide/bagging/#allowed-storage-option-values
	admin_api.CoerceObjectStorageOption(existingObj, submittedObj)
	assert.Equal(t, existingObj.StorageOption, submittedObj.StorageOption)

	// But if existing obj is not active, this should not coerce.
	// Documented way of changing storage option is to delete the old
	// copy, then upload new copy with desired storage option.
	existingObj.State = constants.StateDeleted
	submittedObj.StorageOption = constants.StorageOptionWasabiVA
	assert.NotEqual(t, existingObj.StorageOption, submittedObj.StorageOption)
	admin_api.CoerceObjectStorageOption(existingObj, submittedObj)
	assert.NotEqual(t, existingObj.StorageOption, submittedObj.StorageOption)
}

func TestObjectInitRestore(t *testing.T) {
	// Force fixture reload to prevent "pending work item"
	// error when requesting restoration.
	err := db.ForceFixtureReload()
	require.Nil(t, err)
	tu.InitHTTPTests(t)

	tu.SysAdminClient.POST("/admin-api/v3/objects/init_restore/{id}", 6).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusCreated)

	tu.Inst1AdminClient.POST("/admin-api/v3/objects/init_restore/{id}", 6).
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/objects/init_restore/{id}", 6).
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)
}

func TestObjectBatchDelete(t *testing.T) {
	// START HERE

	// Test permissions. Only APTrust admin should be allowed to do this.

	// Ensure that we get failure if we include an object
	// with a pending WorkItem.

	// Ensure that we get failure if we include an object
	// that belongs to another institution.

	// Ensure that we get failure if requestorID belongs
	// to an inst user rather than inst admin.

	// Ensure we get success with valid params:
	// inst id, user id, object ids.

	// Check post conditions. There should be a deletion
	// request with all expected properties and with the
	// right list of object IDs.

	// Check the text of the alert. It should include
	// all of the object identifiers.

	// Confirm the alert, and then test that the correct
	// WorkItems were created and that no spurious work
	// items were created.

	// TODO: Create & test bulk delete ENV token from
	//       parameter store?
}
