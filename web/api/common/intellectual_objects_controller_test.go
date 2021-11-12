package common_api_test

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

func TestIntellectualObjectShow(t *testing.T) {
	tu.InitHTTPTests(t)

	// In fixtures, obj 1 belongs to inst1
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)
	require.NotNil(t, obj)

	// Sysadmin can read any object
	resp := tu.SysAdminClient.GET("/member-api/v3/objects/show/{id}", obj.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.InstitutionID, record.InstitutionID)

	// Inst admin can read object from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/objects/show/{id}", obj.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read object from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/objects/show/{id}", obj.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read object from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/objects/show/{id}", obj.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read object from other institution
	tu.Inst2UserClient.GET("/member-api/v3/objects/show/{id}", obj.ID).
		Expect().
		Status(http.StatusForbidden)
}

func TestIntellectualObjectIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all objects and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/objects").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.IntellectualObjectList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 13, list.Count)
	assert.Equal(t, "/member-api/v3/objects?page=3&per_page=5", list.Next)
	assert.Equal(t, "/member-api/v3/objects?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters. This should return two results
	resp = tu.SysAdminClient.GET("/member-api/v3/objects").
		WithQuery("access", constants.AccessConsortia).
		WithQuery("state", "A").
		Expect().Status(http.StatusOK)

	list = api.IntellectualObjectList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, 2, len(list.Results))
	for _, object := range list.Results {
		assert.Equal(t, constants.AccessConsortia, object.Access)
	}

	// Inst admin should see only his own institution's objects.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/objects").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 6, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 6, len(list.Results))
	for _, obj := range list.Results {
		assert.Equal(t, tu.Inst1Admin.InstitutionID, obj.InstitutionID)
	}

	// Inst admin cannot see objects belonging to other insitutions.
	tu.Inst2AdminClient.GET("/member-api/v3/objects").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's objects.
	resp = tu.Inst1UserClient.GET("/member-api/v3/objects").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 6, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 6, len(list.Results))
	for _, obj := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, obj.InstitutionID)
	}

	// Inst user cannot see other institution's objects.
	tu.Inst2UserClient.GET("/member-api/v3/objects").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
