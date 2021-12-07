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
	assert.NotEmpty(t, obj.CreatedAt)
	assert.NotEmpty(t, obj.UpdatedAt)

	// Update
	obj, err = pgmodels.IntellectualObjectByIdentifier("test.edu/obj1")
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
