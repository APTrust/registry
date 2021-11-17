package admin_api_test

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"testing"

	// "github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectShow(t *testing.T) {
	tu.InitHTTPTests(t)

	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)
	require.NotNil(t, obj)

	// Sysadmin can get this object through the admin API.
	resp := tu.SysAdminClient.GET("/admin-api/v3/objects/show/{id}", obj.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.InstitutionID, record.InstitutionID)

	// Sysadmin should also be able to find by identifier
	resp = tu.SysAdminClient.GET("/admin-api/v3/objects/show/{id}", obj.Identifier).Expect().Status(http.StatusOK)
	record = &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, obj.ID, record.ID)
	assert.Equal(t, obj.Identifier, record.Identifier)

	// Non-admins should get an error message telling them to
	// use the Member API
	resp = tu.Inst1AdminClient.GET("/admin-api/v3/objects/show/{id}", obj.ID).Expect()
	resp.Status(http.StatusForbidden)
	assert.Equal(t, `{"error":"Permission denied for /admin-api/v3/objects/show/*id (institution 0). non-admins must use the member api"}`, resp.Body().Raw())
}

func TestObjectIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Admin can see this page.
	resp := tu.SysAdminClient.GET("/admin-api/v3/objects").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.IntellectualObjectList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 14, list.Count)
	assert.Equal(t, "/admin-api/v3/objects?page=3&per_page=5", list.Next)
	assert.Equal(t, "/admin-api/v3/objects?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Non-admins can't see this page
	resp = tu.Inst1UserClient.GET("/admin-api/v3/objects").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect()

	resp.Status(http.StatusForbidden)
	assert.Equal(t, `{"error":"Permission denied for /admin-api/v3/objects (institution 0). non-admins must use the member api"}`, resp.Body().Raw())
}

func TestObjectCreate(t *testing.T) {
	tu.InitHTTPTests(t)

}

func TestObjectUpdate(t *testing.T) {
	tu.InitHTTPTests(t)

}

func TestObjectDelete(t *testing.T) {
	tu.InitHTTPTests(t)

}
