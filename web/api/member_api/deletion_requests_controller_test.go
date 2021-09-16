package memberapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	//"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	//"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletionRequestShow(t *testing.T) {
	tu.InitHTTPTests(t)

	deletion, err := pgmodels.DeletionRequestByID(2)
	require.Nil(t, err)
	require.NotNil(t, deletion)

	// Sysadmin can read any deletion
	resp := tu.SysAdminClient.GET("/member-api/v3/deletions/show/{id}", deletion.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.DeletionRequest{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, deletion.ID, record.ID)
	assert.Equal(t, deletion.InstitutionID, record.InstitutionID)

	// Inst admin can read own deletion from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/deletions/show/{id}", deletion.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, deletion.ID, record.ID)
	assert.Equal(t, deletion.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read deletion from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/deletions/show/{id}", deletion.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read own deletion from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/deletions/show/{id}", deletion.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, deletion.ID, record.ID)
	assert.Equal(t, deletion.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read deletion from other institution
	tu.Inst2UserClient.GET("/member-api/v3/deletions/show/{id}", deletion.ID).
		Expect().
		Status(http.StatusForbidden)

}

// func TestDeletionRequestIndex(t *testing.T) {
// 	tu.InitHTTPTests(t)

// 	// Sys Admin should see all deletions and filters
// 	resp := tu.SysAdminClient.GET("/member-api/v3/deletions").
// 		WithQuery("page", 2).
// 		WithQuery("per_page", 5).
// 		Expect().Status(http.StatusOK)

// 	list := api.DeletionRequestViewList{}
// 	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 15, list.Count)
// 	assert.Equal(t, "/member-api/v3/deletions?page=3&per_page=5", list.Next)
// 	assert.Equal(t, "/member-api/v3/deletions?page=1&per_page=5", list.Previous)
// 	assert.Equal(t, "Deletion Confirmed", list.Results[0].Subject)

// 	// Make sure filters work. Should be 3 deletion requested
// 	// deletions for the inst 1 admin.
// 	resp = tu.SysAdminClient.GET("/member-api/v3/deletions").
// 		WithQuery("user_id", tu.Inst1Admin.ID).
// 		WithQuery("type", constants.DeletionRequestDeletionRequested).
// 		Expect().Status(http.StatusOK)

// 	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 3, list.Count)
// 	assert.Equal(t, "", list.Next)
// 	assert.Equal(t, "", list.Previous)
// 	assert.Equal(t, 3, len(list.Results))
// 	assert.Equal(t, "Deletion Requested", list.Results[0].Subject)

// 	// Inst admin should see only his own deletions.
// 	resp = tu.Inst1AdminClient.GET("/member-api/v3/deletions").
// 		Expect().Status(http.StatusOK)
// 	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 6, list.Count)
// 	assert.Equal(t, "", list.Next)
// 	assert.Equal(t, "", list.Previous)
// 	assert.Equal(t, 6, len(list.Results))
// 	for _, deletion := range list.Results {
// 		assert.Equal(t, tu.Inst1Admin.ID, deletion.UserID)
// 	}

// 	// Inst admin cannot see results for other insitutions.
// 	resp = tu.Inst1AdminClient.GET("/member-api/v3/deletions").
// 		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
// 		Expect().Status(http.StatusForbidden)

// 	// Inst admin cannot see results for other users. Technically,
// 	// this should return 403. For now, it returns OK with zero results.
// 	resp = tu.Inst1AdminClient.GET("/member-api/v3/deletions").
// 		WithQuery("user_id", tu.Inst1User.ID).
// 		Expect().Status(http.StatusOK)
// 	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 0, list.Count)
// 	assert.Equal(t, 0, len(list.Results))

// 	// Inst user should see only his own deletions.
// 	resp = tu.Inst1UserClient.GET("/member-api/v3/deletions").
// 		WithQuery("institution_id", tu.Inst1User.InstitutionID).
// 		WithQuery("user_id", tu.Inst1User.ID).
// 		Expect().Status(http.StatusOK)
// 	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 2, list.Count)
// 	assert.Equal(t, "", list.Next)
// 	assert.Equal(t, "", list.Previous)
// 	assert.Equal(t, 2, len(list.Results))
// 	for _, deletion := range list.Results {
// 		assert.Equal(t, tu.Inst1User.ID, deletion.UserID)
// 	}

// 	// Inst user cannot see other institution's deletions.
// 	resp = tu.Inst1UserClient.GET("/member-api/v3/deletions").
// 		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
// 		Expect().Status(http.StatusForbidden)

// 	// Inst user cannot see results for other users. Technically,
// 	// this should return 403. For now, it returns OK with zero results.
// 	resp = tu.Inst1UserClient.GET("/member-api/v3/deletions").
// 		WithQuery("user_id", tu.Inst1Admin.ID).
// 		Expect().Status(http.StatusOK)
// 	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
// 	require.Nil(t, err)
// 	assert.Equal(t, 0, list.Count)
// 	assert.Equal(t, 0, len(list.Results))
// }
