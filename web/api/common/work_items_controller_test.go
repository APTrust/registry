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

func TestWorkItemShow(t *testing.T) {
	tu.InitHTTPTests(t)

	item, err := pgmodels.WorkItemByID(1)
	require.Nil(t, err)
	require.NotNil(t, item)

	// Sysadmin can read any item
	resp := tu.SysAdminClient.GET("/member-api/v3/items/show/{id}", item.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.WorkItem{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, item.ID, record.ID)
	assert.Equal(t, item.InstitutionID, record.InstitutionID)

	// Inst admin can read item from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/items/show/{id}", item.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, item.ID, record.ID)
	assert.Equal(t, item.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read item from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/items/show/{id}", item.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read item from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/items/show/{id}", item.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, item.ID, record.ID)
	assert.Equal(t, item.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read item from other institution
	tu.Inst2UserClient.GET("/member-api/v3/items/show/{id}", item.ID).
		Expect().
		Status(http.StatusForbidden)

}

func TestWorkItemIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all items and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/items").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.WorkItemViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 32, list.Count)
	assert.Equal(t, "/member-api/v3/items?page=3&per_page=5", list.Next)
	assert.Equal(t, "/member-api/v3/items?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters.
	resp = tu.SysAdminClient.GET("/member-api/v3/items").
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

	// Inst admin should see only his own institution's items.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/items").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 17, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 17, len(list.Results))
	for _, item := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, item.InstitutionID)
	}

	// Inst admin cannot see items belonging to other insitutions.
	tu.Inst2AdminClient.GET("/member-api/v3/items").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's items.
	resp = tu.Inst1UserClient.GET("/member-api/v3/items").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 17, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 17, len(list.Results))
	for _, item := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, item.InstitutionID)
	}

	// Inst user cannot see other institution's items.
	tu.Inst2UserClient.GET("/member-api/v3/items").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
