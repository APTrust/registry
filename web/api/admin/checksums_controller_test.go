package admin_api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
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
	resp := tu.SysAdminClient.GET("/admin-api/v3/checksums/show/{id}", cs.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.ChecksumView{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, cs.ID, record.ID)
	assert.Equal(t, cs.Digest, record.Digest)
	assert.Equal(t, cs.InstitutionID, record.InstitutionID)

	// Non sys-admin can't access this endpoint
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/checksums/show/7").
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestChecksumIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all checksums and filters
	resp := tu.SysAdminClient.GET("/admin-api/v3/checksums").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		WithQuery("sort", "id__asc").
		Expect().Status(http.StatusOK)

	list := api.ChecksumViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 32, list.Count)
	assert.Equal(t, "/admin-api/v3/checksums?page=3&per_page=5&sort=id__asc", list.Next)
	assert.Equal(t, "/admin-api/v3/checksums?page=1&per_page=5&sort=id__asc", list.Previous)
	// We're on page 2 here, with 5 items per page, so ids will start at 6
	assert.Equal(t, int64(6), list.Results[0].ID)
	assert.Equal(t, int64(7), list.Results[1].ID)

	// Test some filters.
	resp = tu.SysAdminClient.GET("/admin-api/v3/checksums").
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

	// Non sys-admin can't access this endpoint
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/checksums").
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestChecksumCreate(t *testing.T) {
	// Reset DB after this test so we don't screw up others.
	defer db.ForceFixtureReload()
	tu.InitHTTPTests(t)
	gf, err := pgmodels.GenericFileByID(12)
	require.Nil(t, err)
	require.NotNil(t, gf)

	timestamp := time.Now().UTC()
	digest := "ManThose512sAreSomeLongAssDigests"
	checksum := &pgmodels.Checksum{
		GenericFileID: gf.ID,
		Algorithm:     constants.AlgSha512,
		Digest:        digest,
		DateTime:      timestamp,
	}

	resp := tu.SysAdminClient.POST("/admin-api/v3/checksums/create/{id}", gf.InstitutionID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(checksum).Expect()
	fmt.Println(resp.Body().Raw())
	resp.Status(http.StatusCreated)

	savedChecksum := &pgmodels.Checksum{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), savedChecksum)
	require.Nil(t, err)
	assert.True(t, savedChecksum.ID > int64(0))
	assert.Equal(t, gf.ID, savedChecksum.GenericFileID)
	assert.Equal(t, constants.AlgSha512, savedChecksum.Algorithm)
	assert.Equal(t, digest, savedChecksum.Digest)
	assert.Equal(t, timestamp, savedChecksum.DateTime)

	cs, err := pgmodels.ChecksumByID(savedChecksum.ID)
	require.Nil(t, err)
	require.NotNil(t, cs)

	// Also note that the new checksum should become the "latest"
	// sha512 for this file.
	gfView, err := pgmodels.GenericFileViewByID(12)
	require.Nil(t, err)
	require.NotNil(t, gfView)
	assert.EqualValues(t, 12, gfView.ID)
	assert.Equal(t, digest, gfView.Sha512)
}
