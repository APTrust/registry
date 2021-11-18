package admin_api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
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

func TestObjectCreateUpdateDelete(t *testing.T) {
	tu.InitHTTPTests(t)

	// TODO: Split into three sub-tests

	// Random objects use inst id 4 -> test.edu
	obj := pgmodels.RandomObject()
	resp := tu.SysAdminClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).WithJSON(obj).Expect()
	resp.Status(http.StatusCreated)

	savedObj := &pgmodels.IntellectualObject{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), savedObj)
	require.Nil(t, err)
	assert.True(t, savedObj.ID > int64(0))
	assert.Equal(t, obj.Identifier, savedObj.Identifier)
	assert.Equal(t, obj.InstitutionID, savedObj.InstitutionID)
	assert.Equal(t, obj.BagName, savedObj.BagName)
	assert.Equal(t, obj.ETag, savedObj.ETag)
	assert.Equal(t, obj.StorageOption, savedObj.StorageOption)
	assert.NotEmpty(t, savedObj.CreatedAt)
	assert.NotEmpty(t, savedObj.UpdatedAt)

	origUpdatedAt := savedObj.UpdatedAt

	copyOfSaved := savedObj
	copyOfSaved.Access = constants.AccessConsortia
	copyOfSaved.Title = "Updated Title"
	copyOfSaved.ETag = "UpdatedETag"
	resp = tu.SysAdminClient.PUT("/admin-api/v3/objects/update/{id}", savedObj.ID).WithJSON(copyOfSaved).Expect()
	resp.Status(http.StatusOK)

	updatedObj := &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), updatedObj)
	require.Nil(t, err)

	assert.Equal(t, copyOfSaved.Access, updatedObj.Access)
	assert.Equal(t, copyOfSaved.Title, updatedObj.Title)
	assert.Equal(t, copyOfSaved.ETag, updatedObj.ETag)
	assert.Equal(t, savedObj.CreatedAt, updatedObj.CreatedAt)
	assert.True(t, updatedObj.UpdatedAt.After(origUpdatedAt))

	resp = tu.SysAdminClient.DELETE("/admin-api/v3/objects/delete/{id}", savedObj.ID).Expect()
	fmt.Println(resp.Body())
	resp.Status(http.StatusOK)

	deletedObj := &pgmodels.IntellectualObject{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), deletedObj)
	require.Nil(t, err)

	assert.Equal(t, savedObj.ID, deletedObj.ID)
	assert.Equal(t, constants.StateDeleted, deletedObj.State)
}

func TestObjectCreateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot create objects, even for their
	// own institutions.
	obj := pgmodels.RandomObject()
	obj.InstitutionID = tu.Inst1Admin.InstitutionID

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/create/{id}", obj.InstitutionID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}

func TestObjectUpdateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot update objects, even for their
	// own institutions.
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/update/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/update/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}

func TestObjectDeleteUpdateUnauthorized(t *testing.T) {
	tu.InitHTTPTests(t)

	// Non sysadmins cannot delete objects, even for their
	// own institutions. (Not through the API anyway.)
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)

	resp := tu.Inst1AdminClient.POST("/admin-api/v3/objects/delete/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)

	resp = tu.Inst1UserClient.POST("/admin-api/v3/objects/delete/{id}", obj.ID).WithJSON(obj).Expect()
	resp.Status(http.StatusForbidden)
}
