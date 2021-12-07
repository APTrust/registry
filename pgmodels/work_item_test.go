package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkItemValidation(t *testing.T) {
	item := &pgmodels.WorkItem{}
	err := item.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrItemName, err.Errors["Name"])
	assert.Equal(t, pgmodels.ErrItemETag, err.Errors["ETag"])
	assert.Equal(t, pgmodels.ErrItemBagDate, err.Errors["BagDate"])
	assert.Equal(t, pgmodels.ErrItemBucket, err.Errors["Bucket"])
	assert.Equal(t, pgmodels.ErrItemUser, err.Errors["User"])
	assert.Equal(t, pgmodels.ErrItemInstID, err.Errors["InstitutionID"])
	assert.Equal(t, pgmodels.ErrItemDateProcessed, err.Errors["DateProcessed"])
	assert.Equal(t, pgmodels.ErrItemNote, err.Errors["Note"])
	assert.Equal(t, pgmodels.ErrItemAction, err.Errors["Action"])
	assert.Equal(t, pgmodels.ErrItemStage, err.Errors["Stage"])
	assert.Equal(t, pgmodels.ErrItemStatus, err.Errors["Status"])
	assert.Equal(t, pgmodels.ErrItemOutcome, err.Errors["Outcome"])
}

func TestWorkItemGetID(t *testing.T) {
	item := &pgmodels.WorkItem{
		ID: 199,
	}
	assert.Equal(t, int64(199), item.GetID())
}

func TestWorkItemByID(t *testing.T) {
	db.LoadFixtures()
	item, err := pgmodels.WorkItemByID(int64(23))
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, int64(23), item.ID)
}

func TestWorkItemGet(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where("name", "=", "pdfs.tar")
	item, err := pgmodels.WorkItemGet(query)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, "pdfs.tar", item.Name)
}

func TestWorkItemSelect(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery()
	query.Where("name", "!=", "pdfs.tar")
	query.Where("name", "!=", "coal.tar")
	query.OrderBy("name", "asc")
	items, err := pgmodels.WorkItemSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, items)
	assert.True(t, (len(items) > 25 && len(items) < 35))
	for _, item := range items {
		assert.NotEqual(t, "pdfs.tar", item)
		assert.NotEqual(t, "coal.tar", item)
	}
}

func TestWorkItemSave(t *testing.T) {
	db.LoadFixtures()
	item := &pgmodels.WorkItem{
		Name:          "unit_00001.tar",
		ETag:          "12345678901234567890123456789099",
		InstitutionID: 4,
		User:          "system@aptrust.org",
		Bucket:        "aptrust.receiving.test.test.edu",
		Action:        constants.ActionIngest,
		Stage:         constants.StageRequested,
		Status:        constants.StatusPending,
		Note:          "Item is awaiting ingest.",
		Outcome:       "I said item is awaiting ingest.",
		BagDate:       TestDate,
		DateProcessed: TestDate,
		Retry:         true,
		Size:          8000,
	}
	err := item.Save()
	require.Nil(t, err)

	// pg library should set ID, BeforeInsert hook should set other values
	assert.True(t, item.ID > int64(0))
	assert.Equal(t, "unit_00001.tar", item.Name)
	assert.Equal(t, int64(4), item.InstitutionID)
	assert.NotEmpty(t, item.CreatedAt)
	assert.NotEmpty(t, item.UpdatedAt)
}

func TestWorkItemHasCompleted(t *testing.T) {
	item := &pgmodels.WorkItem{}
	for _, status := range constants.IncompleteStatusValues {
		item.Status = status
		assert.False(t, item.HasCompleted())
	}
	for _, status := range constants.CompletedStatusValues {
		item.Status = status
		assert.True(t, item.HasCompleted())
	}
}

func TestWorkItemSetForRequeue(t *testing.T) {
	db.LoadFixtures()
	item := &pgmodels.WorkItem{
		Name:          "unit_00002.tar",
		ETag:          "12345678901234567890123456789022",
		InstitutionID: 4,
		User:          "system@aptrust.org",
		Bucket:        "aptrust.receiving.test.test.edu",
		Action:        constants.ActionIngest,
		Stage:         constants.StageStore,
		Status:        constants.StatusStarted,
		Note:          "Item is being stored.",
		Outcome:       "I said item is being stored.",
		BagDate:       TestDate,
		DateProcessed: TestDate,
		Retry:         true,
		Size:          8000,
	}
	err := item.Save()
	require.Nil(t, err)

	err = item.SetForRequeue(constants.StageFormatIdentification)
	require.Nil(t, err)

	assert.Equal(t, constants.StageFormatIdentification, item.Stage)
	assert.Equal(t, constants.StatusPending, item.Status)
	assert.True(t, item.Retry)
	assert.False(t, item.NeedsAdminReview)
	assert.Empty(t, item.Node)
	assert.Empty(t, item.PID)
	assert.Equal(t, "Requeued for Format Identification", item.Note)

	// This should fail, because Package is not a valid stage for Ingest.
	err = item.SetForRequeue(constants.StagePackage)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, common.ErrInvalidRequeue)
}

func TestWorkItemsPendingForObject(t *testing.T) {
	db.LoadFixtures()

	item := &pgmodels.WorkItem{
		Name:          "pending.tar",
		ETag:          "12345678901234567890123456789022",
		InstitutionID: 4,
		User:          "system@aptrust.org",
		Bucket:        "aptrust.receiving.test.test.edu",
		Action:        constants.ActionIngest,
		Stage:         constants.StageStore,
		Status:        constants.StatusStarted,
		Note:          "Item is being stored.",
		Outcome:       "I said item is being stored.",
		BagDate:       TestDate,
		DateProcessed: TestDate,
		Retry:         true,
		Size:          8000,
	}
	err := item.Save()
	require.Nil(t, err)

	// Should return nothing, because inst ID doesn't match.
	itemsInProgress, err := pgmodels.WorkItemsPendingForObject(3, "pending.tar")
	require.Nil(t, err)
	assert.Equal(t, 0, len(itemsInProgress))

	// This should get the item above
	itemsInProgress, err = pgmodels.WorkItemsPendingForObject(4, "pending.tar")
	require.Nil(t, err)
	assert.Equal(t, 1, len(itemsInProgress))

	item = itemsInProgress[0]
	item.Status = constants.StatusCancelled
	err = item.Save()
	require.Nil(t, err)

	// It should not come back this time because it has a completed status.
	itemsInProgress, err = pgmodels.WorkItemsPendingForObject(4, "pending.tar")
	require.Nil(t, err)
	assert.Equal(t, 0, len(itemsInProgress))

}

func TestWorkItemsPendingForFile(t *testing.T) {
	db.LoadFixtures()

	// File 1 from fixtures is institution1.edu/photos/picture1
	item := &pgmodels.WorkItem{
		Name:                 "photos.tar",
		ETag:                 "99995678901234567890123456789999",
		InstitutionID:        1,
		IntellectualObjectID: 1,
		GenericFileID:        1,
		User:                 "system@aptrust.org",
		Bucket:               "aptrust.receiving.test.test.edu",
		Action:               constants.ActionRestoreFile,
		Stage:                constants.StageRequested,
		Status:               constants.StatusPending,
		Note:                 "Test restoration item (file)",
		Outcome:              "Item is pending",
		BagDate:              TestDate,
		DateProcessed:        TestDate,
		Retry:                true,
		Size:                 8000,
	}
	err := item.Save()
	require.Nil(t, err)

	// This should get the item above
	itemsInProgress, err := pgmodels.WorkItemsPendingForFile(1)
	require.Nil(t, err)
	require.Equal(t, 1, len(itemsInProgress))

	item = itemsInProgress[0]
	item.Stage = constants.StageAvailableInS3
	item.Status = constants.StatusSuccess
	item.Note = "This thing is done. Look in S3."
	err = item.Save()
	require.Nil(t, err)

	// It should not come back this time because it has a completed status.
	itemsInProgress, err = pgmodels.WorkItemsPendingForFile(1)
	require.Nil(t, err)
	assert.Equal(t, 0, len(itemsInProgress))
}

func TestLastSuccessfulIngest(t *testing.T) {
	db.LoadFixtures()

	// This item in our fixtures is the last successful ingest
	// of object id #4 from fixtures, institution2.edu/chocolate
	item, err := pgmodels.WorkItemByID(26)
	require.Nil(t, err)
	require.NotNil(t, item)

	// Make sure LastSuccessfulIngest returns the
	// expected WorkItem.
	lastIngest, err := pgmodels.LastSuccessfulIngest(item.IntellectualObjectID)
	require.Nil(t, err)
	assert.Equal(t, item.ID, lastIngest.ID)

	// Save work item with a later ingest of same object
	var copyOfItem pgmodels.WorkItem
	err = copier.Copy(&copyOfItem, item)
	require.Nil(t, err)

	copyOfItem.ETag = "aaaabbbb4939474aa6f5f77bf56faaaa"
	copyOfItem.DateProcessed = time.Now().UTC()
	err = copyOfItem.Save()
	require.Nil(t, err)

	// Now we should get that later work item.
	lastIngest, err = pgmodels.LastSuccessfulIngest(item.IntellectualObjectID)
	require.Nil(t, err)
	assert.Equal(t, copyOfItem.ID, lastIngest.ID)
	assert.Equal(t, copyOfItem.ETag, lastIngest.ETag)
}

func TestNewRestorationItem(t *testing.T) {
	// Object id #4 from fixtures, institution2.edu/chocolate
	// has at least one successful ingest WorkItem.
	// File id #11, institution2.edu/chocolate/picture1,
	// belongs to that object.

	obj, err := pgmodels.IntellectualObjectByID(4)
	require.Nil(t, err)
	require.NotNil(t, obj)

	file, err := pgmodels.GenericFileByID(11)
	require.Nil(t, err)
	require.NotNil(t, file)

	user := &pgmodels.User{
		Email: "unittest@example.com",
	}

	// Object restoration
	item, err := pgmodels.NewRestorationItem(obj, nil, user)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.True(t, item.ID > 0)
	assert.Equal(t, obj.ID, item.IntellectualObjectID)
	assert.EqualValues(t, 0, item.GenericFileID)
	assert.Equal(t, constants.ActionRestoreObject, item.Action)
	assert.Equal(t, user.Email, item.User)
	assert.Empty(t, item.Node)
	assert.Empty(t, item.PID)
	assert.True(t, item.Retry)
	assert.False(t, item.NeedsAdminReview)
	assert.Empty(t, item.QueuedAt)
	assert.Equal(t, constants.StageRequested, item.Stage)
	assert.Equal(t, constants.StatusPending, item.Status)

	// File restoration
	item, err = pgmodels.NewRestorationItem(obj, file, user)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, obj.ID, item.IntellectualObjectID)
	assert.Equal(t, file.ID, item.GenericFileID)
	assert.Equal(t, constants.ActionRestoreFile, item.Action)

	// Now check Glacier restoration. This is on an object.
	obj.StorageOption = constants.StorageOptionGlacierDeepOH
	item, err = pgmodels.NewRestorationItem(obj, nil, user)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, obj.ID, item.IntellectualObjectID)
	assert.EqualValues(t, 0, item.GenericFileID)
	assert.Equal(t, constants.ActionGlacierRestore, item.Action)

	// Restoring a file from Glacier should also result
	// in action being GlacierRestoration
	obj.StorageOption = constants.StorageOptionGlacierDeepOH
	item, err = pgmodels.NewRestorationItem(obj, file, user)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, obj.ID, item.IntellectualObjectID)
	assert.EqualValues(t, file.ID, item.GenericFileID)
	assert.Equal(t, constants.ActionGlacierRestore, item.Action)
}
