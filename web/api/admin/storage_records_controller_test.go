package admin_api_test

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

func TestStorageRecordShow(t *testing.T) {
	tu.InitHTTPTests(t)

	sr, err := pgmodels.StorageRecordByID(1)
	require.Nil(t, err)
	require.NotNil(t, sr)
	require.True(t, sr.ID > 0)

	// Sysadmin can see any checksum
	resp := tu.SysAdminClient.GET("/admin-api/v3/storage_records/show/{id}", sr.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.StorageRecord{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, sr.ID, record.ID)
	assert.Equal(t, sr.GenericFileID, record.GenericFileID)
	assert.Equal(t, sr.URL, record.URL)

	// Non sys-admin can't access this endpoint
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/storage_records/show/1").
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestStorageRecordIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all checksums and filters
	resp := tu.SysAdminClient.GET("/admin-api/v3/storage_records").
		WithQuery("generic_file_id", 1).
		WithQuery("page", 1).
		WithQuery("per_page", 1).
		WithQuery("sort", "id__asc").
		Expect().Status(http.StatusOK)

	list := api.StorageRecordList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, "/admin-api/v3/storage_records?generic_file_id=1&page=2&per_page=1&sort=id__asc", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, int64(1), list.Results[0].ID)

	// Non sys-admin can't access this endpoint
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/storage_records").
			Expect().
			Status(http.StatusForbidden)
	}
}
