package admin_api_test

import (
	"encoding/json"
	"net/http"
	"testing"
	//"time"

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

	// Sysadmin should be able to get this.
	// This is a pass-through to the common api endpoint,
	// but we want to make sure it's available at this URL.
	resp := tu.SysAdminClient.GET("/admin-api/v3/files/show/{id}", gf.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.GenericFile{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, gf.ID, record.ID)
	assert.Equal(t, gf.InstitutionID, record.InstitutionID)

	// Non-admins should get an error. They have to go through
	// the member API.
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/files/show/{id}", gf.ID).
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestGenericFileIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// This endpoint should work for sys admin
	// but not for others.
	resp := tu.SysAdminClient.GET("/admin-api/v3/files").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.GenericFileViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 62, list.Count)
	assert.Equal(t, "/admin-api/v3/files?page=3&per_page=5", list.Next)
	assert.Equal(t, "/admin-api/v3/files?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst1User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters. This object has 1 deleted, 4 active files.
	resp = tu.SysAdminClient.GET("/admin-api/v3/files").
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

	// Non-admins are forbidden. They have to go through
	// the member API.
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/files").
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestFileCreate(t *testing.T) {

}

func TestFileUpdate(t *testing.T) {

}
