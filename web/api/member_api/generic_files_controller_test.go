package memberapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericFileShow(t *testing.T) {
	tu.InitHTTPTests(t)

	gf, err := pgmodels.GenericFileByID(1)
	require.Nil(t, err)
	require.NotNil(t, gf)

	// Sysadmin can read any file
	resp := tu.SysAdminClient.GET("/member-api/v3/files/show/{id}", gf.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.GenericFile{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, gf.ID, record.ID)
	assert.Equal(t, gf.InstitutionID, record.InstitutionID)

	// Inst admin can read file from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/files/show/{id}", gf.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, gf.ID, record.ID)
	assert.Equal(t, gf.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read file from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/files/show/{id}", gf.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read file from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/files/show/{id}", gf.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, gf.ID, record.ID)
	assert.Equal(t, gf.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read file from other institution
	tu.Inst2UserClient.GET("/member-api/v3/files/show/{id}", gf.ID).
		Expect().
		Status(http.StatusForbidden)

}

func TestGenericFileIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all files and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/files").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.GenericFileViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 61, list.Count)
	assert.Equal(t, "/member-api/v3/files/?page=3&per_page=5", list.Next)
	assert.Equal(t, "/member-api/v3/files/?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst1User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters. This object has 1 deleted, 4 active files.
	resp = tu.SysAdminClient.GET("/member-api/v3/files").
		WithQuery("intellectual_object_id", 3).
		WithQuery("state", "A").
		Expect().Status(http.StatusOK)

	list = api.GenericFileViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 4, list.Count)
	assert.Equal(t, 4, len(list.Results))
	for _, file := range list.Results {
		assert.Equal(t, int64(3), file.IntellectualObjectID)
		assert.Equal(t, "institution1.edu/glass", file.ObjectIdentifier)
		assert.Equal(t, "institution1.edu", file.InstitutionIdentifier)
		assert.Equal(t, constants.AccessConsortia, file.Access)
		assert.Equal(t, "A", file.State)
		assert.True(t, file.Size > 0)
	}

	// Inst admin should see only his own institution's files.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/files").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 18, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 18, len(list.Results))
	for _, gf := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, gf.InstitutionID)
	}

	// Inst admin cannot see files belonging to other insitutions.
	resp = tu.Inst2AdminClient.GET("/member-api/v3/files").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's files.
	resp = tu.Inst1UserClient.GET("/member-api/v3/files").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 18, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 18, len(list.Results))
	for _, gf := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, gf.InstitutionID)
	}

	// Inst user cannot see other institution's files.
	resp = tu.Inst2UserClient.GET("/member-api/v3/files").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
