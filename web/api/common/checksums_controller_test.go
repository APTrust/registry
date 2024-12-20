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

func TestChecksumShow(t *testing.T) {
	tu.InitHTTPTests(t)

	cs, err := pgmodels.ChecksumViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, cs)
	require.True(t, cs.ID > 0)

	// Sysadmin can see any checksum
	resp := tu.SysAdminClient.GET("/member-api/v3/checksums/show/{id}", cs.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.ChecksumView{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, cs.ID, record.ID)
	assert.Equal(t, cs.Digest, record.Digest)
	assert.Equal(t, cs.InstitutionID, record.InstitutionID)

	// Inst admin can read checksum from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/checksums/show/{id}", cs.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, cs.ID, record.ID)
	assert.Equal(t, cs.Digest, record.Digest)
	assert.Equal(t, cs.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read checksum from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/checksums/show/{id}", cs.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read checksum from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/checksums/show/{id}", cs.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, cs.ID, record.ID)
	assert.Equal(t, cs.Digest, record.Digest)
	assert.Equal(t, cs.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read checksum from other institution
	tu.Inst2UserClient.GET("/member-api/v3/checksums/show/{id}", cs.ID).
		Expect().
		Status(http.StatusForbidden)

}

func TestChecksumIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all checksums and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/checksums").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		WithQuery("sort", "id__asc").
		Expect().Status(http.StatusOK)

	list := api.ChecksumViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 32, list.Count)
	assert.Equal(t, "/member-api/v3/checksums?page=3&per_page=5&sort=id__asc", list.Next)
	assert.Equal(t, "/member-api/v3/checksums?page=1&per_page=5&sort=id__asc", list.Previous)
	// We're on page 2 here, with 5 items per page, so ids will start at 6
	assert.Equal(t, int64(6), list.Results[0].ID)
	assert.Equal(t, int64(7), list.Results[1].ID)

	// Test some filters.
	resp = tu.SysAdminClient.GET("/member-api/v3/checksums").
		WithQuery("intellectual_object_id", 7).
		WithQuery("state", "A").
		Expect()
	resp.Status(http.StatusOK)

	list = api.ChecksumViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 12, list.Count)
	assert.Equal(t, 12, len(list.Results))
	for _, cs := range list.Results {
		assert.Equal(t, int64(7), cs.IntellectualObjectID)
		assert.Equal(t, "A", cs.State)
	}

	// Inst admin should see only his own institution's checksums.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/checksums").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 4, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 4, len(list.Results))
	for _, cs := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, cs.InstitutionID)
	}

	// Inst admin cannot see checksums belonging to other insitutions.
	tu.Inst2AdminClient.GET("/member-api/v3/checksums").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's checksums.
	resp = tu.Inst1UserClient.GET("/member-api/v3/checksums").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 4, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 4, len(list.Results))
	for _, cs := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, cs.InstitutionID)
	}

	// Inst user cannot see other institution's checksums.
	tu.Inst2UserClient.GET("/member-api/v3/checksums").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
