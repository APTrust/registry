package pgmodels_test

import (
	"testing"

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

func TestObjValidateChanges(t *testing.T) {
	// START HERE
}

func TestObjInsert(t *testing.T) {

}

func TestObjUpdate(t *testing.T) {

}
