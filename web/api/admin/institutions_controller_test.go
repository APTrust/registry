package admin_api_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionShow(t *testing.T) {
	tu.InitHTTPTests(t)

	instIds := []int64{1, 2, 3, 4, 5}

	// Sysadmin can read any institution
	for _, id := range instIds {
		resp := tu.SysAdminClient.GET("/admin-api/v3/institutions/show/{id}", id).Expect().Status(http.StatusOK)
		inst := &pgmodels.Institution{}
		err := json.Unmarshal([]byte(resp.Body().Raw()), inst)
		require.Nil(t, err)
		assert.Equal(t, id, inst.ID)
	}

	// Non sys-admin cannot access this endpoint at all,
	// even for looking up their own institution.
	for _, id := range instIds {
		tu.Inst1AdminClient.GET("/admin-api/v3/institutions/show/{id}", id).
			Expect().
			Status(http.StatusForbidden)
		tu.Inst1UserClient.GET("/admin-api/v3/institutions/show/{id}", id).
			Expect().
			Status(http.StatusForbidden)
		tu.Inst2AdminClient.GET("/admin-api/v3/institutions/show/{id}", id).
			Expect().
			Status(http.StatusForbidden)
		tu.Inst2UserClient.GET("/admin-api/v3/institutions/show/{id}", id).
			Expect().
			Status(http.StatusForbidden)
	}
}

func TestInstitutionIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all institutions and filters
	resp := tu.SysAdminClient.GET("/admin-api/v3/institutions").
		WithQuery("page", 2).
		WithQuery("per_page", 2).
		Expect().Status(http.StatusOK)

	list := api.InstitutionViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 5, list.Count)
	assert.Equal(t, "/admin-api/v3/institutions?page=3&per_page=2", list.Next)
	assert.Equal(t, "/admin-api/v3/institutions?page=1&per_page=2", list.Previous)
	assert.Equal(t, int64(2), list.Results[0].ID) // order by name, top of page 2

	// Test some filters. This should return "Institution One",
	// "Institution Two", "Test Institution", and "Example Institution".
	resp = tu.SysAdminClient.GET("/admin-api/v3/institutions").
		WithQuery("name__contains", "Institution").
		Expect().Status(http.StatusOK)

	list = api.InstitutionViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 4, list.Count)
	assert.Equal(t, 4, len(list.Results))
	for _, institution := range list.Results {
		assert.True(t, strings.Contains(institution.Name, "Institution"))
	}

	// Non sys admins cannot access this endpoint
	tu.Inst1AdminClient.GET("/admin-api/v3/institutions").
		Expect().Status(http.StatusForbidden)
	tu.Inst1UserClient.GET("/admin-api/v3/institutions").
		Expect().Status(http.StatusForbidden)
	tu.Inst2AdminClient.GET("/admin-api/v3/institutions").
		Expect().Status(http.StatusForbidden)
	tu.Inst2UserClient.GET("/admin-api/v3/institutions").
		Expect().Status(http.StatusForbidden)
}
