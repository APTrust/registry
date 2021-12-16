package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var timeLayout = "2006-01-02 15:04:05 -0700 MST"

func TestObjIsGlacierOnly(t *testing.T) {
	obj := &pgmodels.IntellectualObject{}
	for _, option := range constants.GlacierOnlyOptions {
		obj.StorageOption = option
		assert.True(t, obj.IsGlacierOnly())
	}
	obj.StorageOption = constants.StorageOptionStandard
	assert.False(t, obj.IsGlacierOnly())
}

func TestIdForObjIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForObjIdentifier("institution1.edu/photos")
	require.Nil(t, err)
	assert.Equal(t, int64(1), id)

	id, err = pgmodels.IdForObjIdentifier("institution2.edu/coal")
	require.Nil(t, err)
	assert.Equal(t, int64(5), id)

	id, err = pgmodels.IdForObjIdentifier("bad identifier")
	require.NotNil(t, err)
}

func TestObjHasActiveFiles(t *testing.T) {
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)
	require.NotNil(t, obj)
	hasFiles, err := obj.HasActiveFiles()
	require.Nil(t, err)
	assert.True(t, hasFiles)

	// This fixture object has only one file, with state = "D"
	obj, err = pgmodels.IntellectualObjectByID(14)
	require.Nil(t, err)
	require.NotNil(t, obj)
	hasFiles, err = obj.HasActiveFiles()
	require.Nil(t, err)
	assert.False(t, hasFiles)

}

func TestObjLastIngestEvent(t *testing.T) {
	obj, err := pgmodels.IntellectualObjectByID(6)
	require.Nil(t, err)
	require.NotNil(t, obj)

	event, err := obj.LastIngestEvent()
	require.Nil(t, err)
	require.NotNil(t, event)

	assert.Equal(t, obj.ID, event.IntellectualObjectID)
	assert.Equal(t, int64(0), event.GenericFileID)
	assert.Equal(t, constants.EventIngestion, event.EventType)
	assert.Equal(t, "bbe16041-8887-4739-af04-9d35e5cab4dc", event.Identifier)
	assert.Equal(t, int64(53), event.ID)
}

func TestObjLastDeletionEvent(t *testing.T) {
	obj, err := pgmodels.IntellectualObjectByID(6)
	require.Nil(t, err)
	require.NotNil(t, obj)

	event, err := obj.LastDeletionEvent()
	require.Nil(t, err)
	require.NotNil(t, event)

	assert.Equal(t, obj.ID, event.IntellectualObjectID)
	assert.Equal(t, int64(0), event.GenericFileID)
	assert.Equal(t, constants.EventDeletion, event.EventType)
	assert.Equal(t, "775af09e-87d1-42be-9fcd-4315c5836099", event.Identifier)
	assert.Equal(t, int64(54), event.ID)
}

func TestObjValidate(t *testing.T) {
	obj1 := pgmodels.GetTestObject()
	assert.Nil(t, obj1.Validate())

	obj1.Title = " "
	obj1.Identifier = " \t\t\n "
	obj1.State = "X"
	obj1.Access = "not a real access value"
	obj1.InstitutionID = 0
	obj1.StorageOption = "not a real option"
	obj1.BagItProfileIdentifier = "  "
	obj1.SourceOrganization = ""

	err := obj1.Validate()
	assert.NotNil(t, err)

	keys := []string{
		"Title",
		"Identifier",
		"State",
		"Access",
		"InstitutionID",
		"StorageOption",
		"BagItProfileIdentifier",
		"SourceOrganization",
	}

	for _, key := range keys {
		assert.NotEmpty(t, err.Errors[key], key)
	}
}

func TestObjValidateChanges(t *testing.T) {
	obj1 := pgmodels.GetTestObject()
	obj2 := pgmodels.GetTestObject()

	assert.Nil(t, obj1.ValidateChanges(obj2))

	obj1.ID = 999
	obj2.ID = 1000
	err := obj1.ValidateChanges(obj2)
	assert.Equal(t, common.ErrIDMismatch, err)
	obj2.ID = obj1.ID

	obj2.InstitutionID = 4500
	err = obj1.ValidateChanges(obj2)
	assert.Equal(t, common.ErrInstIDChange, err)
	obj2.InstitutionID = obj1.InstitutionID

	obj2.Identifier = "test.edu/changed"
	err = obj1.ValidateChanges(obj2)
	assert.Equal(t, common.ErrIdentifierChange, err)
	obj2.Identifier = obj1.Identifier

	obj2.StorageOption = constants.StorageOptionGlacierOH
	err = obj1.ValidateChanges(obj2)
	assert.Equal(t, common.ErrStorageOptionChange, err)
}

func TestObjInsertAndUpdate(t *testing.T) {
	// Insert
	obj := pgmodels.GetTestObject()
	err := obj.Save()
	assert.Nil(t, err)
	assert.True(t, obj.ID > 0)
	assert.NotEmpty(t, obj.CreatedAt)
	assert.NotEmpty(t, obj.UpdatedAt)

	// Update
	obj, err = pgmodels.IntellectualObjectByIdentifier(obj.Identifier)
	require.Nil(t, err)
	require.NotNil(t, obj)

	origUpdatedAt := obj.UpdatedAt
	obj.Description = "Updated description of obj1"

	err = obj.Save()
	require.Nil(t, err)
	assert.True(t, obj.UpdatedAt.After(origUpdatedAt))
}

// TODO: Test the following
//
// AssertDeletionPreconditions
// LatestDeletionWorkItem
// DeletionRequest
// NewDeletionEvent

func TestAssertObjDeletionPreconditions(t *testing.T) {
	db.LoadFixtures()
	defer db.ForceFixtureReload()
	obj, err := pgmodels.CreateObjectWithRelations()
	require.Nil(t, err)
	require.NotNil(t, obj)

	testLastObjDeletionWorkItem(t, obj)
	testObjectDeletionRequest(t, obj)
}

func testLastObjDeletionWorkItem(t *testing.T, obj *pgmodels.IntellectualObject) {
	item, err := obj.ActiveDeletionWorkItem()
	require.Nil(t, err)
	assert.Nil(t, item)

	// This one should NOT be returned because its a file deletion
	gfWorkItem := pgmodels.RandomWorkItem(obj.BagName, constants.ActionDelete, obj.ID, obj.GenericFiles[0].ID)
	require.Nil(t, gfWorkItem.Save())
	item, err = obj.ActiveDeletionWorkItem()
	require.Nil(t, err)
	assert.Nil(t, item)

	// This one should NOT be returned because item status is not Started
	objWorkItem := pgmodels.RandomWorkItem(obj.BagName, constants.ActionDelete, obj.ID, 0)
	require.Nil(t, objWorkItem.Save())
	item, err = obj.ActiveDeletionWorkItem()
	require.Nil(t, err)
	require.Nil(t, item)

	// This one SHOULD be returned because it's an object deletion
	// and it has been started.
	objWorkItem.Status = constants.StatusStarted
	require.Nil(t, objWorkItem.Save())
	item, err = obj.ActiveDeletionWorkItem()
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, objWorkItem.ID, item.ID)

}

func testObjectDeletionRequest(t *testing.T, obj *pgmodels.IntellectualObject) {
	// Initially, there's no deletion request for this object
	reqView, err := obj.DeletionRequest(99999999999)
	require.NotNil(t, err)
	assert.Empty(t, reqView.ID) // Should we even get this back??

	// Figure out the work item id. That will lead us back to
	// the original deletion request.
	item, err := obj.ActiveDeletionWorkItem()
	require.Nil(t, err)
	require.NotNil(t, item)

	objects := []*pgmodels.IntellectualObject{obj}
	req, err := pgmodels.CreateDeletionRequest(objects, nil)
	require.Nil(t, err)
	require.NotNil(t, req)
	req.WorkItemID = item.ID
	require.Nil(t, req.Save())

	deletionReqView, err := obj.DeletionRequest(item.ID)
	require.Nil(t, err)
	require.NotNil(t, deletionReqView)

	assert.Equal(t, req.ID, deletionReqView.ID)
}

func TestNewObjDeletionEvent(t *testing.T) {

}
