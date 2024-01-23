package admin_api_test

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	admin_api "github.com/APTrust/registry/web/api/admin"
	tu "github.com/APTrust/registry/web/testutil"
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

	// validParams := url.Values{}
	// validParams.Set("institutionID", strconv.FormatInt(tu.Inst2Admin.InstitutionID, 10))
	// validParams.Set("requestorID", strconv.FormatInt(tu.Inst2Admin.ID, 10))
	// for _, id := range idsThatCanBeDeleted {
	// 	validParams.Add("objectID", strconv.FormatInt(id, 10))
	// }

	validParams := admin_api.ObjectBatchDeleteParams{
		InstitutionID: tu.Inst2Admin.InstitutionID,
		RequestorID:   tu.Inst2Admin.ID,
		ObjectIDs:     idsThatCanBeDeleted,
		SecretKey:     common.Context().Config.BatchDeletionKey,
	}

	// Inst admin can request batch deletion, but inst user cannot.
	paramsBadRequestorRole := url.Values{}
	paramsBadRequestorRole.Set("institutionID", strconv.FormatInt(tu.Inst2User.InstitutionID, 10))
	paramsBadRequestorRole.Set("requestorID", strconv.FormatInt(tu.Inst2User.ID, 10))
	for _, id := range idsThatCanBeDeleted {
		paramsBadRequestorRole.Add("objectID", strconv.FormatInt(id, 10))
	}

	// This inst admin is at the wrong institution
	paramsBadRequestorInst := url.Values{}
	paramsBadRequestorInst.Set("institutionID", strconv.FormatInt(tu.Inst1Admin.InstitutionID, 10))
	paramsBadRequestorInst.Set("requestorID", strconv.FormatInt(tu.Inst1Admin.ID, 10))
	for _, id := range idsThatCanBeDeleted {
		paramsBadRequestorInst.Add("objectID", strconv.FormatInt(id, 10))
	}

	// This batch contains one id referring to an object that
	// has already been deleted.
	paramsBadIdAlreadyDeleted := url.Values{}
	paramsBadIdAlreadyDeleted.Set("institutionID", strconv.FormatInt(tu.Inst2Admin.InstitutionID, 10))
	paramsBadIdAlreadyDeleted.Set("requestorID", strconv.FormatInt(tu.Inst2Admin.ID, 10))
	for _, id := range idsThatCanBeDeleted {
		paramsBadIdAlreadyDeleted.Add("objectID", strconv.FormatInt(id, 10))
	}
	paramsBadIdAlreadyDeleted.Add("objectID", strconv.FormatInt(idAlreadyDeleted, 10))

	// This batch contains one id that has pending work items.
	paramsBadPendingWorkItem := url.Values{}
	paramsBadPendingWorkItem.Set("institutionID", strconv.FormatInt(tu.Inst2Admin.InstitutionID, 10))
	paramsBadPendingWorkItem.Set("requestorID", strconv.FormatInt(tu.Inst2Admin.ID, 10))
	for _, id := range idsThatCanBeDeleted {
		paramsBadPendingWorkItem.Add("objectID", strconv.FormatInt(id, 10))
	}
	paramsBadPendingWorkItem.Add("objectID", strconv.FormatInt(idWithPendingWorkItems, 10))

	// This batch contains one object that belongs to another institution.
	paramsBadOtherInstItem := url.Values{}
	paramsBadOtherInstItem.Set("institutionID", strconv.FormatInt(tu.Inst2Admin.InstitutionID, 10))
	paramsBadOtherInstItem.Set("requestorID", strconv.FormatInt(tu.Inst2Admin.ID, 10))
	for _, id := range idsThatCanBeDeleted {
		paramsBadOtherInstItem.Add("objectID", strconv.FormatInt(id, 10))
	}
	paramsBadOtherInstItem.Add("objectID", strconv.FormatInt(idBelongingToOtherInst, 10))

	testObjectBatchDeletePermissions(t, validParams)

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

func testObjectBatchDeletePermissions(t *testing.T, params admin_api.ObjectBatchDeleteParams) {

	tu.SysAdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect().Status(http.StatusCreated)

	tu.Inst1AdminClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(params).
		Expect().Status(http.StatusForbidden)

	tu.Inst1UserClient.POST("/admin-api/v3/objects/init_batch_delete").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithJSON(params).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)

}
