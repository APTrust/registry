package admin_api_test

import (
	"encoding/json"
	"net/http"
	"testing"
	//"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkItemShow(t *testing.T) {
	tu.InitHTTPTests(t)

	item, err := pgmodels.WorkItemByID(1)
	require.Nil(t, err)
	require.NotNil(t, item)

	// Sysadmin can access any item through this endpoint
	resp := tu.SysAdminClient.GET("/admin-api/v3/items/show/{id}", item.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.WorkItem{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, item.ID, record.ID)
	assert.Equal(t, item.InstitutionID, record.InstitutionID)

	// Non-admins should get an error. They have to go through
	// the member API.
	for _, client := range tu.AllClients {
		if client == tu.SysAdminClient {
			continue
		}
		tu.Inst2AdminClient.GET("/admin-api/v3/items/show/{id}", gf.ID).
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestWorkItemIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all items and filters
	resp := tu.SysAdminClient.GET("/admin-api/v3/items").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.WorkItemViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 32, list.Count)
	assert.Equal(t, "/admin-api/v3/items?page=3&per_page=5", list.Next)
	assert.Equal(t, "/admin-api/v3/items?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters.
	resp = tu.SysAdminClient.GET("/admin-api/v3/items").
		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
		Expect().Status(http.StatusOK)

	list = api.WorkItemViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 15, list.Count)
	assert.Equal(t, 15, len(list.Results))
	for _, item := range list.Results {
		assert.Equal(t, tu.Inst2Admin.InstitutionID, item.InstitutionID)
		assert.Equal(t, "institution2.edu", item.InstitutionIdentifier)
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

// TODO: Test create and update
