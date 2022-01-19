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
	v "github.com/asaskevich/govalidator"
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
	resp := tu.SysAdminClient.GET("/admin-api/v3/files/show/{id}", gf.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().Status(http.StatusOK)
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
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		WithQuery("sort", "id__asc").
		Expect().Status(http.StatusOK)

	list := api.GenericFileViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 62, list.Count)
	assert.Equal(t, "/admin-api/v3/files?page=3&per_page=5&sort=id__asc", list.Next)
	assert.Equal(t, "/admin-api/v3/files?page=1&per_page=5&sort=id__asc", list.Previous)
	assert.Equal(t, tu.Inst1User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters. This object has 1 deleted, 4 active files.
	resp = tu.SysAdminClient.GET("/admin-api/v3/files").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
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

func TestFileCreateUpdateDelete(t *testing.T) {
	// Reset DB after this test so we don't screw up others.
	defer db.ForceFixtureReload()
	tu.InitHTTPTests(t)
	gf := testFileCreate(t)
	updatedFile := testFileUpdate(t, gf)

	createFileDeletionPreConditions(t, updatedFile)
	testFileDelete(t, updatedFile)
}

func testFileCreate(t *testing.T) *pgmodels.GenericFile {
	obj, err := pgmodels.IntellectualObjectGet(
		pgmodels.NewQuery().
			Where("institution_id", "=", 4).
			Limit(1))
	require.Nil(t, err)
	require.NotNil(t, obj)
	gf := pgmodels.RandomGenericFile(obj.ID, obj.Identifier)
	resp := tu.SysAdminClient.POST("/admin-api/v3/files/create/{id}", gf.InstitutionID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(gf).Expect()
	resp.Status(http.StatusCreated)

	savedFile := &pgmodels.GenericFile{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), savedFile)
	require.Nil(t, err)
	assert.True(t, savedFile.ID > int64(0))
	assert.Equal(t, gf.Identifier, savedFile.Identifier)
	assert.Equal(t, gf.InstitutionID, savedFile.InstitutionID)
	assert.Equal(t, gf.Size, savedFile.Size)
	assert.Equal(t, gf.FileFormat, savedFile.FileFormat)
	assert.Equal(t, gf.StorageOption, savedFile.StorageOption)
	assert.NotEmpty(t, savedFile.CreatedAt)
	assert.NotEmpty(t, savedFile.UpdatedAt)
	return savedFile
}

func testFileUpdate(t *testing.T, gf *pgmodels.GenericFile) *pgmodels.GenericFile {
	origUpdatedAt := gf.UpdatedAt
	copyOfGf := gf
	copyOfGf.Size = gf.Size + 200
	copyOfGf.FileFormat = "txt/screed"

	resp := tu.SysAdminClient.PUT("/admin-api/v3/files/update/{id}", gf.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(copyOfGf).Expect()
	resp.Status(http.StatusOK)

	updatedGf := &pgmodels.GenericFile{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), updatedGf)
	require.Nil(t, err)

	assert.Equal(t, copyOfGf.Size, updatedGf.Size)
	assert.Equal(t, copyOfGf.FileFormat, updatedGf.FileFormat)
	assert.Equal(t, gf.CreatedAt, updatedGf.CreatedAt)
	assert.True(t, updatedGf.UpdatedAt.After(origUpdatedAt))

	return updatedGf
}

// Registry business rules won't allow deletions without the following:
//
// - Ingest event at ingest.
// - Deletion request when a user clicks the delete file button
//   in the web UI.
// - WorkItem when an inst admin has approved the deletion request.
//
// Here, we create them just so we can complete our test.
func createFileDeletionPreConditions(t *testing.T, gf *pgmodels.GenericFile) {
	// Deletion checks for last ingest event on this object.
	event := pgmodels.RandomPremisEvent(constants.EventIngestion)
	event.IntellectualObjectID = gf.IntellectualObjectID
	event.GenericFileID = gf.ID
	event.GenericFileID = 0
	event.InstitutionID = gf.InstitutionID
	require.Nil(t, event.Save())

	// Also requires an approved Deletion work item
	item := pgmodels.RandomWorkItem(
		"TestBagName.tar",
		constants.ActionDelete,
		gf.IntellectualObjectID,
		gf.ID)
	item.User = "admin@test.edu"
	item.InstApprover = "admin@test.edu"
	item.Status = constants.StatusStarted
	require.Nil(t, item.Save())
	require.True(t, item.ID > 0)

	// Requires approved deletion request
	now := time.Now().UTC()
	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	req.GenericFiles = append(req.GenericFiles, gf)
	req.InstitutionID = gf.InstitutionID
	req.RequestedByID = 8 // admin@test.edu
	req.RequestedAt = now
	req.ConfirmedByID = 8
	req.ConfirmedAt = now
	req.WorkItemID = item.ID
	require.Nil(t, req.Save())
}

func testFileDelete(t *testing.T, gf *pgmodels.GenericFile) {
	resp := tu.SysAdminClient.DELETE("/admin-api/v3/files/delete/{id}", gf.ID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusOK)

	// Make sure we got the expected JSON response from the server.
	deletedFile := &pgmodels.GenericFile{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), deletedFile)
	require.Nil(t, err)

	assert.Equal(t, gf.ID, deletedFile.ID)
	assert.Equal(t, constants.StateDeleted, deletedFile.State)

	// Make sure the state was actually saved
	savedFile, err := pgmodels.GenericFileByID(gf.ID)
	require.Nil(t, err)
	assert.Equal(t, constants.StateDeleted, savedFile.State)

	// Test for deletion event
	event, err := gf.LastDeletionEvent()
	require.Nil(t, err)
	require.NotNil(t, event)
	require.True(t, event.DateTime.After(time.Now().UTC().Add(-5*time.Second)))
}

// POST /admin-api/v3/files/create_batch/:institution_id
func TestGenericFileCreateBatch(t *testing.T) {
	defer db.ForceFixtureReload()
	obj, files, err := pgmodels.RandomFileBatch()
	require.Nil(t, err)
	require.NotNil(t, obj)
	require.NotNil(t, files)
	require.NotEmpty(t, files)

	resp := tu.SysAdminClient.POST("/admin-api/v3/files/create_batch/{id}", obj.InstitutionID).
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithJSON(files).
		Expect()

	// Unless there's an error, we should get an empty JSON payload.
	// If there is an error, let's see it.
	fmt.Println(resp.Body())
	resp.Status(http.StatusCreated)

	// Make sure everything was saved correctly.
	query := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", obj.ID)
	savedFiles, err := pgmodels.GenericFileSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 20, len(savedFiles))
	for _, gf := range savedFiles {
		// Make sure file properties were set...
		assert.True(t, gf.ID > 0)
		assert.Equal(t, obj.ID, gf.IntellectualObjectID)
		assert.Equal(t, obj.InstitutionID, gf.InstitutionID)

		// Make sure all events were saved with correct properties.
		query := pgmodels.NewQuery().Where("generic_file_id", "=", gf.ID)
		events, err := pgmodels.PremisEventSelect(query)
		require.Nil(t, err)
		assert.Equal(t, 4, len(events))
		for _, event := range events {
			assert.Equal(t, gf.ID, event.GenericFileID)
			assert.Equal(t, gf.IntellectualObjectID, event.IntellectualObjectID)
			assert.Equal(t, gf.InstitutionID, event.InstitutionID)
			assert.Equal(t, constants.EventIngestion, event.EventType)
			assert.False(t, event.DateTime.IsZero())
			assert.False(t, event.CreatedAt.IsZero())
		}

		// Check checksums
		checksums, err := pgmodels.ChecksumSelect(query)
		require.Nil(t, err)
		assert.Equal(t, 4, len(checksums))
		for _, cs := range checksums {
			assert.Equal(t, gf.ID, cs.GenericFileID)
			assert.Equal(t, constants.AlgSha1, cs.Algorithm)
			assert.False(t, cs.DateTime.IsZero())
			assert.False(t, cs.CreatedAt.IsZero())
			assert.False(t, cs.UpdatedAt.IsZero())
		}

		// And storage records
		storageRecs, err := pgmodels.StorageRecordSelect(query)
		require.Nil(t, err)
		assert.Equal(t, 4, len(storageRecs))
		for _, sr := range storageRecs {
			assert.Equal(t, gf.ID, sr.GenericFileID)
			assert.True(t, v.IsURL(sr.URL))
		}
	}
}
