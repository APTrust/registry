package admin_api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/APTrust/registry/web/testutil"
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
	resp := tu.SysAdminClient.GET("/admin-api/v3/items/show/{id}", item.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusOK)
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
		tu.Inst2AdminClient.GET("/admin-api/v3/items/show/{id}", item.ID).
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestWorkItemIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all items and filters
	resp := tu.SysAdminClient.GET("/admin-api/v3/items").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
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
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
		Expect().Status(http.StatusOK)

	list = api.WorkItemViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 16, list.Count)
	assert.Equal(t, 16, len(list.Results))
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
		tu.Inst2AdminClient.GET("/admin-api/v3/items").
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestWorkItemRequeue(t *testing.T) {
	tu.InitHTTPTests(t)

	workItem := testutil.CreateWorkItem(t, "unit_test_bag2.tar")

	// SysAdmin can requeue
	tu.SysAdminClient.PUT("/admin-api/v3/items/requeue/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithFormField("stage", constants.StageReingestCheck).
		Expect().Status(http.StatusOK)

	// Make sure that worked. It should set not only the stage,
	// but the other values noted below.
	item, err := pgmodels.WorkItemByID(workItem.ID)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, constants.StageReingestCheck, item.Stage)
	assert.True(t, item.Retry)
	assert.False(t, item.NeedsAdminReview)
	assert.Empty(t, item.PID)
	assert.Empty(t, item.Node)
	assert.Equal(t, constants.StatusPending, item.Status)

	// Make sure other roles cannot requeue
	tu.Inst1AdminClient.PUT("/admin-api/v3/items/requeue/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithFormField("stage", constants.StageReingestCheck).
		Expect().Status(http.StatusForbidden)
	tu.Inst1UserClient.PUT("/admin-api/v3/items/requeue/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithFormField("stage", constants.StageReingestCheck).
		Expect().Status(http.StatusForbidden)
}

func TestWorkItemRedisDelete(t *testing.T) {
	tu.InitHTTPTests(t)
	workItem := testutil.CreateWorkItem(t, "unit_test_bag2.tar")

	// Make sure non-admins cannot delete the Redis record
	tu.Inst1AdminClient.DELETE("/admin-api/v3/items/redis_delete/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)
	tu.Inst1UserClient.DELETE("/admin-api/v3/items/redis_delete/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusForbidden)

	// SysAdmin can delete the Redis data
	tu.SysAdminClient.DELETE("/admin-api/v3/items/redis_delete/{id}", workItem.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusOK)
}

func TestItemCreateAndUpdate(t *testing.T) {
	item := testItemCreate(t)
	testItemUpdate(t, item)
}

func testItemCreate(t *testing.T) *pgmodels.WorkItem {
	obj, err := pgmodels.IntellectualObjectGet(
		pgmodels.NewQuery().
			Where("institution_id", "=", 4).
			Limit(1))
	require.Nil(t, err)
	require.NotNil(t, obj)
	item := pgmodels.RandomWorkItem(obj.BagName, constants.ActionIngest, 1, 1)
	resp := tu.SysAdminClient.POST("/admin-api/v3/items/create/{id}", item.InstitutionID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(item).Expect()
	resp.Status(http.StatusCreated)

	savedItem := &pgmodels.WorkItem{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), savedItem)
	require.Nil(t, err)
	assert.True(t, savedItem.ID > int64(0))
	assert.Equal(t, item.Name, savedItem.Name)
	assert.Equal(t, item.InstitutionID, savedItem.InstitutionID)
	assert.Equal(t, item.Size, savedItem.Size)
	assert.Equal(t, item.Action, savedItem.Action)
	assert.Equal(t, item.Outcome, savedItem.Outcome)
	assert.NotEmpty(t, savedItem.CreatedAt)
	assert.NotEmpty(t, savedItem.UpdatedAt)
	return savedItem
}

func testItemUpdate(t *testing.T, item *pgmodels.WorkItem) *pgmodels.WorkItem {
	origUpdatedAt := item.UpdatedAt
	copyOfItem := item
	copyOfItem.Size = item.Size + 200
	copyOfItem.Outcome = "This outcome has been edited"

	resp := tu.SysAdminClient.PUT("/admin-api/v3/items/update/{id}", item.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(copyOfItem).Expect()
	resp.Status(http.StatusOK)

	updatedItem := &pgmodels.WorkItem{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), updatedItem)
	require.Nil(t, err)

	assert.Equal(t, copyOfItem.Size, updatedItem.Size)
	assert.Equal(t, copyOfItem.Outcome, updatedItem.Outcome)
	assert.Equal(t, item.CreatedAt, updatedItem.CreatedAt)
	assert.True(t, updatedItem.UpdatedAt.After(origUpdatedAt))

	return updatedItem
}
