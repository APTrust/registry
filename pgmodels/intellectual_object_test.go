package pgmodels_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var timeLayout = "2006-01-02 15:04:05 -0700 MST"

func TestIntellectualObjectByID(t *testing.T) {
	obj, err := pgmodels.IntellectualObjectByID(1)
	require.Nil(t, err)
	require.NotNil(t, obj)
	assert.Equal(t, int64(1), obj.ID)
	assert.Equal(t, "institution1.edu/photos", obj.Identifier)
}

func TestIntellectualObjectByIdentifier(t *testing.T) {
	obj, err := pgmodels.IntellectualObjectByIdentifier("institution1.edu/photos")
	require.Nil(t, err)
	require.NotNil(t, obj)
	assert.Equal(t, int64(1), obj.ID)
	assert.Equal(t, "institution1.edu/photos", obj.Identifier)
}

func TestIntellectualObjectSelect(t *testing.T) {
	query := pgmodels.NewQuery().
		Where("institution_id", "=", InstOne)
	objects, err := pgmodels.IntellectualObjectSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 6, len(objects))
	for _, obj := range objects {
		assert.Equal(t, InstOne, obj.InstitutionID)
	}

	query.Where("access", "=", constants.AccessConsortia)
	objects, err = pgmodels.IntellectualObjectSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 1, len(objects))
	for _, obj := range objects {
		assert.Equal(t, InstOne, obj.InstitutionID)
		assert.Equal(t, constants.AccessConsortia, obj.Access)
	}
}

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

func TestAssertObjDeletionPreconditions(t *testing.T) {
	defer db.ForceFixtureReload()
	obj, err := pgmodels.CreateObjectWithRelations()
	require.Nil(t, err)
	require.NotNil(t, obj)

	testLastObjDeletionWorkItem(t, obj)
	testObjectDeletionRequest(t, obj)

	// Hit the following underlying methods:
	// assertNoActiveFiles
	// assertNotAlreadyDeleted
	// assertDeletionApproved

	// First pre-condition: No active files.
	err = obj.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.Equal(t, common.ErrActiveFiles, err)

	// Mark files deleted...
	for _, gf := range obj.GenericFiles {
		gf.State = constants.StateDeleted
		err = gf.Save()
		require.Nil(t, err)
	}

	// Next pre-condition: Object isn't already deleted.
	obj.State = constants.StateDeleted
	err = obj.Save()
	require.Nil(t, err)

	err = obj.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.Equal(t, "Object is already in deleted state", err.Error())

	obj.State = constants.StateActive
	err = obj.Save()
	require.Nil(t, err)

	// Create a deletion work item for this object
	workItem := pgmodels.RandomWorkItem(obj.BagName, constants.ActionDelete, obj.ID, 0)
	workItem.Status = constants.StatusStarted
	require.Nil(t, workItem.Save())

	err = obj.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Deletion work item is missing institutional approver"), err.Error())

	// Approve it
	workItem.InstApprover = "someone@example.com"
	require.Nil(t, workItem.Save())

	err = obj.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "No deletion request for work item"), err.Error())

	// Last pre-condition: Object requires an approved deletion request
	objects := []*pgmodels.IntellectualObject{obj}
	req, err := pgmodels.CreateDeletionRequest(objects, nil)
	require.Nil(t, err)
	require.NotNil(t, req)
	req.WorkItemID = workItem.ID
	require.Nil(t, req.Save())

	err = obj.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Deletion request"), err.Error())
	assert.True(t, strings.Contains(err.Error(), "has no approver"), err.Error())

	// Add an approver
	req.ConfirmedByID = req.RequestedByID
	req.ConfirmedAt = time.Now().UTC()
	require.Nil(t, req.Save())

	// Now we should be OK
	err = obj.AssertDeletionPreconditions()
	assert.Nil(t, err)
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
	assert.Nil(t, reqView)

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
	defer db.ForceFixtureReload()
	obj, err := pgmodels.CreateObjectWithRelations()
	require.Nil(t, err)
	require.NotNil(t, obj)

	event, err := obj.NewDeletionEvent()
	assert.Nil(t, event)
	require.NotNil(t, err)
	assert.Equal(t, "Missing deletion request work item", err.Error())

	// Create a deletion work item for this object
	workItem := pgmodels.RandomWorkItem(obj.BagName, constants.ActionDelete, obj.ID, 0)
	workItem.Status = constants.StatusStarted
	require.Nil(t, workItem.Save())

	// But it's not approved yet...
	event, err = obj.NewDeletionEvent()
	assert.Nil(t, event)
	require.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Deletion work item is missing institutional approver"), err.Error())

	// OK, so approve it
	workItem.InstApprover = "someone@example.com"
	require.Nil(t, workItem.Save())

	// Now we have an approved work item, but no deletion request.
	// That should produce the following error...
	event, err = obj.NewDeletionEvent()
	assert.Nil(t, event)
	require.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "No deletion request for work item"), err.Error())

	// Now we have a deletion request with no approver
	objects := []*pgmodels.IntellectualObject{obj}
	req, err := pgmodels.CreateDeletionRequest(objects, nil)
	require.Nil(t, err)
	require.NotNil(t, req)
	req.WorkItemID = workItem.ID
	require.Nil(t, req.Save())

	event, err = obj.NewDeletionEvent()
	assert.Nil(t, event)
	require.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Deletion request"), err.Error())
	assert.True(t, strings.Contains(err.Error(), "has no approver"), err.Error())

	// Add an approver
	req.ConfirmedByID = req.RequestedByID
	req.ConfirmedAt = time.Now().UTC()
	require.Nil(t, req.Save())

	// Now we're getting somewhere
	event, err = obj.NewDeletionEvent()
	require.Nil(t, err)
	assert.NotNil(t, event)

	// Get the deletion request view, which has some info we'll need
	// to verify premis event details below.
	reqView, err := pgmodels.DeletionRequestViewByID(req.ID)
	require.Nil(t, err)
	require.NotNil(t, reqView)

	assert.Equal(t, "APTrust preservation services", event.Agent)
	assert.True(t, event.DateTime.After(time.Now().UTC().Add(-5*time.Second)))
	assert.Equal(t, "Object deleted from preservation storage", event.Detail)
	assert.Equal(t, constants.EventDeletion, event.EventType)
	assert.True(t, common.LooksLikeUUID(event.Identifier))
	assert.Equal(t, obj.InstitutionID, event.InstitutionID)
	assert.Equal(t, obj.ID, event.IntellectualObjectID)
	assert.Equal(t, "Minio S3 library", event.Object)
	assert.Equal(t, constants.OutcomeSuccess, event.Outcome)
	assert.Equal(t, reqView.RequestedByEmail, event.OutcomeDetail)
	assert.Equal(t,
		fmt.Sprintf("Object deleted at the request of %s. Institutional approver: %s.",
			reqView.RequestedByEmail, reqView.ConfirmedByEmail),
		event.OutcomeInformation)
}

func TestObjectDelete(t *testing.T) {

}
