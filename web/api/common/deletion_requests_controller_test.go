package common_api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
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

func TestDeletionRequestIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all deletions and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/deletions").
		WithQuery("page", 2).
		WithQuery("per_page", 1).
		Expect().Status(http.StatusOK)

	list := api.DeletionRequestViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, "/member-api/v3/deletions?page=3&per_page=1", list.Next)
	assert.Equal(t, "/member-api/v3/deletions?page=1&per_page=1", list.Previous)
	assert.Equal(t, tu.Inst1User.ID, list.Results[0].RequestedByID)

	// Inst admin should see only his own institution's deletions.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/deletions").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 3, len(list.Results))
	for _, deletion := range list.Results {
		assert.Equal(t, tu.Inst1User.ID, deletion.RequestedByID)
	}

	// Inst admin cannot see results for other insitutions.
	tu.Inst2AdminClient.GET("/member-api/v3/deletions").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's deletions.
	resp = tu.Inst1UserClient.GET("/member-api/v3/deletions").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 3, len(list.Results))
	for _, deletion := range list.Results {
		assert.Equal(t, tu.Inst1User.ID, deletion.RequestedByID)
	}

	// Inst user cannot see other institution's deletions.
	tu.Inst2UserClient.GET("/member-api/v3/deletions").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
