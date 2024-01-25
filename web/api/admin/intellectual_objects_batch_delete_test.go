package admin_api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	admin_api "github.com/APTrust/registry/web/api/admin"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectBatchDelete(t *testing.T) {

	os.Setenv("APT_ENV", "test")

	err := db.ForceFixtureReload()
	require.Nil(t, err)
	tu.InitHTTPTests(t)

	idsThatCanBeDeleted := []int64{5, 6, 12, 13}
	idAlreadyDeleted := int64(14)
	idWithPendingWorkItems := int64(4)
	idBelongingToOtherInst := int64(1)

	// These params are valid and will create a bulk deletion request
	// if submitted by APTrust admin user.
	validParams := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst2Admin.ID,
		ObjectIDs:     idsThatCanBeDeleted,
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// Inst admin can request batch deletion, but inst user cannot.
	paramsBadRequestorRole := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2User.InstitutionID,
		RequestorID:   tu.Inst2User.ID,
		ObjectIDs:     idsThatCanBeDeleted,
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// This inst admin is at the wrong institution.
	// He's requesting deletion of items belonging to Inst2,
	// but he belongs to Inst1
	paramsBadRequestorInst := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst1Admin.ID,
		ObjectIDs:     idsThatCanBeDeleted,
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// This batch contains one id referring to an object that
	// has already been deleted.
	paramsBadIdAlreadyDeleted := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst2Admin.ID,
		ObjectIDs:     append(idsThatCanBeDeleted, idAlreadyDeleted),
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// This batch contains one id that has pending work items.
	paramsBadPendingWorkItem := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst2Admin.ID,
		ObjectIDs:     append(idsThatCanBeDeleted, idWithPendingWorkItems),
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// // This batch contains one object that belongs to another institution.
	paramsBadOtherInstItem := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst2Admin.ID,
		ObjectIDs:     append(idsThatCanBeDeleted, idBelongingToOtherInst),
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// Test permissions. Only APTrust admin should be allowed to do this.
	testObjectBatchDeletePermissions(t, validParams)

	// Ensure we get success with valid params and
	// make sure a successful request creates the expected
	// records in the DB.
	testObjectBatchDeleteCreatesExpectedRecords(t, validParams)

	// Ensure that we get failure if we include an object
	// with a pending WorkItem.
	testObjectBatchDeleteWithPendingWorkItem(t, paramsBadPendingWorkItem)

	// Ensure that we get failure if we include an object
	// that belongs to another institution.
	testObjectBatchDeleteWithOtherInstItem(t, paramsBadOtherInstItem)

	// Ensure that we get failure if requestorID belongs
	// to an inst user rather than inst admin.
	testObjectBatchDeleteWithWrongRole(t, paramsBadRequestorRole)

	// Ensure failure if requestor id belongs to someone at
	// a different institution. An institution that does not own
	// the objects.
	testObjectBatchDeleteWrongRequestorInst(t, paramsBadRequestorInst)

	// This should fail because one of the objects in the list
	// has already been deleted.
	testObjectBatchDeleteAlreadyDeleted(t, paramsBadIdAlreadyDeleted)
}

func testObjectBatchDeletePermissions(t *testing.T, params admin_api.ObjectBatchDeleteParams) {

	// No institutional admin can create a bulk deletion request. Period.
	tu.Inst1AdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect().Status(http.StatusForbidden)

	// No institutional user can create a bulk deletion request. Period.
	tu.Inst1UserClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithJSON(params).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)

	// APTrust SysAdmin can create bulk deletion requests on behalf
	// of a local institutional admin.
	tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect().Status(http.StatusCreated)

}

// Run this test with validaParams. All other param sets are invalid
// and will fail.
func testObjectBatchDeleteCreatesExpectedRecords(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()

	// Make sure we got the expected status.
	resp.Status(http.StatusCreated)
	respData, err := io.ReadAll(resp.Raw().Body)
	resp.Raw().Body.Close()
	require.NoError(t, err)

	deletionRequest := &pgmodels.DeletionRequest{}
	err = json.Unmarshal(respData, deletionRequest)
	require.NoError(t, err)

	// Make sure the deletion request has the right info.
	// We're not testing RequestedByID here because we don't
	// serialize that value to JSON. I can't remember why not,
	// but I think it's to keep from exposing user IDs. We
	// will test RequestedByID below, when we retrieve the
	// database record.
	assert.Equal(t, params.InstitutionID, deletionRequest.InstitutionID)
	assert.Equal(t, len(params.ObjectIDs), len(deletionRequest.IntellectualObjects))
	for _, objId := range params.ObjectIDs {
		found := false
		for _, obj := range deletionRequest.IntellectualObjects {
			if obj.ID == objId {
				found = true
				break
			}
		}
		assert.True(t, found, objId)
	}

	// Make sure this deletion request was saved to the DB as well.
	// Also ensure that the request is linked to all four requested
	// objects, and no other objects.
	savedDelRequest, err := pgmodels.DeletionRequestByID(deletionRequest.ID)
	require.NoError(t, err)
	require.NotNil(t, savedDelRequest)
	assert.Equal(t, params.InstitutionID, savedDelRequest.InstitutionID)
	assert.Equal(t, params.RequestorID, savedDelRequest.RequestedByID)
	assert.Equal(t, len(params.ObjectIDs), len(savedDelRequest.IntellectualObjects))
	for _, objId := range params.ObjectIDs {
		found := false
		for _, obj := range savedDelRequest.IntellectualObjects {
			if obj.ID == objId {
				found = true
				break
			}
		}
		assert.True(t, found, objId)
	}

	// Because this deletion request has not yet been approved,
	// there should be no associated WorkItem. We create the WorkItem
	// only after the deletion request has been approved by the
	// local institutional admin.
	assert.Empty(t, deletionRequest.WorkItemID)

	// Make sure the deletion request alert was created in the DB.
	query := pgmodels.NewQuery().Where("deletion_request_id", "=", deletionRequest.ID)
	alerts, err := pgmodels.AlertSelect(query)
	require.NoError(t, err)
	assert.Equal(t, 1, len(alerts))
	alert := alerts[0]
	assert.Equal(t, params.InstitutionID, alert.InstitutionID)
	assert.Equal(t, "Deletion Requested", alert.Subject)
	assert.Equal(t, "Deletion Requested", alert.Type)

	// The link to review the deletion request should also contain a token,
	// but we don't know offhand what it is. Other tests dig into this.
	// Here, we just want to be sure we're sending the right email message.
	expectedReviewLink := fmt.Sprintf("http://localhost/deletions/review/%d?token=", deletionRequest.ID)
	assert.Contains(t, alert.Content, expectedReviewLink)

	// Should be no work items because request has not yet been approved.
	assert.Empty(t, alert.WorkItems)
}

func testObjectBatchDeleteWithPendingWorkItem(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()
	resp.Status(http.StatusConflict)
	assert.Equal(t, `{"StatusCode":409,"Error":"task cannot be completed because this object has pending work items"}`, resp.Body().Raw())
}

func testObjectBatchDeleteWithOtherInstItem(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()
	resp.Status(http.StatusBadRequest)
	assert.Equal(t, `{"StatusCode":400,"Error":"one or more object ids is invalid"}`, resp.Body().Raw())
}

func testObjectBatchDeleteWithWrongRole(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()
	resp.Status(http.StatusBadRequest)
	assert.Equal(t, `{"StatusCode":400,"Error":"invalid requestor id"}`, resp.Body().Raw())
}

func testObjectBatchDeleteWrongRequestorInst(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()
	resp.Status(http.StatusBadRequest)
	assert.Equal(t, `{"StatusCode":400,"Error":"invalid requestor id"}`, resp.Body().Raw())
}

func testObjectBatchDeleteAlreadyDeleted(t *testing.T, params admin_api.ObjectBatchDeleteParams) {
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect()
	resp.Status(http.StatusBadRequest)
	assert.Equal(t, `{"StatusCode":400,"Error":"one or more object ids is invalid"}`, resp.Body().Raw())
}
