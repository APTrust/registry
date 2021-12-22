package pgmodels_test

import (
	//"fmt"
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

var HammerTime, _ = time.Parse(time.RFC3339, "2022-01-02T11:42:00Z")

func TestFileIsGlacierOnly(t *testing.T) {
	gf := &pgmodels.GenericFile{}
	for _, option := range constants.GlacierOnlyOptions {
		gf.StorageOption = option
		assert.True(t, gf.IsGlacierOnly())
	}
	gf.StorageOption = constants.StorageOptionStandard
	assert.False(t, gf.IsGlacierOnly())
}

func TestIdForFileIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForFileIdentifier("institution2.edu/toads/toad19")
	require.Nil(t, err)
	assert.Equal(t, int64(42), id)

	id, err = pgmodels.IdForFileIdentifier("institution1.edu/gl-or/file1.epub")
	require.Nil(t, err)
	assert.Equal(t, int64(53), id)

	id, err = pgmodels.IdForFileIdentifier("bad identifier")
	require.NotNil(t, err)
}

func TestGenericFileByID(t *testing.T) {
	gf, err := pgmodels.GenericFileByID(1)
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)

	gf, err = pgmodels.GenericFileByID(-999)
	assert.NotNil(t, err)
	assert.Nil(t, gf)
}

func TestGenericFileByIdentifier(t *testing.T) {
	gf, err := pgmodels.GenericFileByIdentifier("institution2.edu/toads/toad16")
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, "institution2.edu/toads/toad16", gf.Identifier)

	gf, err = pgmodels.GenericFileByIdentifier("-- does not exist --")
	assert.NotNil(t, err)
	assert.Nil(t, gf)
}

func TestGenericFileSelect(t *testing.T) {
	db.ForceFixtureReload()
	query := pgmodels.NewQuery().
		Where("institution_id", "=", InstTest).
		OrderBy("identifier", "asc")
	files, err := pgmodels.GenericFileSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, files)
	assert.Equal(t, 7, len(files))
	assert.Equal(t, "test.edu/apt-test-restore/aptrust-info.txt", files[0].Identifier)
	assert.Equal(t, "test.edu/btr-512-test-restore/data/sample.xml", files[6].Identifier)
}

func TestFileSaveGetUpdate(t *testing.T) {
	gf, err := pgmodels.GenericFileByID(1)
	require.Nil(t, err)
	require.NotNil(t, gf)

	newFile := pgmodels.RandomGenericFile(1, "inst1.edu/photos/unit-test.txt")
	err = newFile.Save()
	require.Nil(t, err)
	assert.True(t, newFile.ID > 0)

	newFile, err = pgmodels.GenericFileByIdentifier(newFile.Identifier)
	require.Nil(t, err)
	require.NotNil(t, newFile)

	newFile.FileFormat = "this-has/been-changed"
	newFile.LastFixityCheck = HammerTime
	err = newFile.Save()
	require.Nil(t, err)

	reloadedFile, err := pgmodels.GenericFileByIdentifier(newFile.Identifier)
	require.Nil(t, err)
	require.NotNil(t, reloadedFile)
	assert.Equal(t, "this-has/been-changed", reloadedFile.FileFormat)
	assert.Equal(t, HammerTime, reloadedFile.LastFixityCheck)
}

func TestFileValidate(t *testing.T) {
	gf := &pgmodels.GenericFile{
		Size: -1,
	}
	err := gf.Validate()
	require.NotNil(t, err)

	assert.Equal(t, "FileFormat is required", err.Errors["FileFormat"])
	assert.Equal(t, "Identifier is required", err.Errors["Identifier"])
	assert.Equal(t, pgmodels.ErrInstState, err.Errors["State"])
	assert.Equal(t, "Size cannot be negative", err.Errors["Size"])
	assert.Equal(t, "Invalid institution id", err.Errors["InstitutionID"])
	assert.Equal(t, "Intellectual object ID is required", err.Errors["IntellectualObjectID"])
	assert.Equal(t, "Invalid storage option", err.Errors["StorageOption"])
	assert.Equal(t, "Valid UUID required", err.Errors["UUID"])

	gf.Size = 20
	gf.FileFormat = "text/html"
	gf.Identifier = "test.edu/some-html-file"
	gf.State = constants.StateActive
	gf.InstitutionID = InstTest
	gf.IntellectualObjectID = 20
	gf.StorageOption = constants.StorageOptionGlacierVA
	gf.UUID = "c464d6dd-9fa6-41d9-8cb5-cdc7b986d07d"

	err = gf.Validate()
	assert.Nil(t, err)
}

func TestObjectFileCount(t *testing.T) {
	count, err := pgmodels.ObjectFileCount(1, "", constants.StateActive)
	require.Nil(t, err)
	assert.Equal(t, 4, count)

	count, err = pgmodels.ObjectFileCount(1, "", constants.StateDeleted)
	require.Nil(t, err)
	assert.Equal(t, 0, count)

	count, err = pgmodels.ObjectFileCount(1, "picture", constants.StateActive)
	require.Nil(t, err)
	assert.Equal(t, 3, count)

	count, err = pgmodels.ObjectFileCount(1, "doc", constants.StateActive)
	require.Nil(t, err)
	assert.Equal(t, 0, count)

	count, err = pgmodels.ObjectFileCount(1, "9876543210", constants.StateActive)
	require.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestObjectFiles(t *testing.T) {
	files, err := pgmodels.ObjectFiles(1, "", constants.StateActive, 0, 20)
	require.Nil(t, err)
	assert.Equal(t, 4, len(files))

	files, err = pgmodels.ObjectFiles(1, "", constants.StateDeleted, 0, 20)
	require.Nil(t, err)
	assert.Equal(t, 0, len(files))

	files, err = pgmodels.ObjectFiles(1, "picture", constants.StateActive, 0, 20)
	require.Nil(t, err)
	assert.Equal(t, 3, len(files))

	files, err = pgmodels.ObjectFiles(1, "doc", constants.StateActive, 0, 20)
	require.Nil(t, err)
	assert.Equal(t, 0, len(files))

	files, err = pgmodels.ObjectFiles(1, "9876543210", constants.StateActive, 0, 20)
	require.Nil(t, err)
	assert.Equal(t, 1, len(files))

}

func TestFileDeletionPreConditions(t *testing.T) {
	defer db.ForceFixtureReload()

	gf, err := pgmodels.GenericFileByID(1)
	require.Nil(t, err)
	require.NotNil(t, gf)

	// Already marked as deleted...
	gf.State = constants.StateDeleted
	err = gf.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.Equal(t, "File is already in deleted state", err.Error())
	testGenericFileDeleteError(t, gf)
	gf.State = constants.StateActive

	// Has no deletion work item...
	err = gf.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.Equal(t, "Missing deletion request work item", err.Error())
	testGenericFileDeleteError(t, gf)

	item, err := gf.ActiveDeletionWorkItem()
	require.Nil(t, err)
	require.Nil(t, item)

	workItem := pgmodels.RandomWorkItem("BaggerVance.tar",
		constants.ActionDelete, gf.IntellectualObjectID, gf.ID)
	workItem.InstitutionID = gf.InstitutionID
	workItem.Status = constants.StatusStarted
	err = workItem.Save()
	require.Nil(t, err)

	item, err = gf.ActiveDeletionWorkItem()
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, workItem.ID, item.ID)

	// Have work item but it's not approved
	err = gf.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.Equal(t, "Deletion work item is missing institutional approver", err.Error())
	testGenericFileDeleteError(t, gf)

	workItem.InstApprover = "some-guy@example.com"
	err = workItem.Save()
	require.Nil(t, err)

	// Approved work item but no deletion request
	err = gf.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "No deletion request for work item"))
	testGenericFileDeleteError(t, gf)

	testFileDeletionRequest(t, gf, workItem.ID)

	// Request not yet approved
	err = gf.AssertDeletionPreconditions()
	require.NotNil(t, err)
	assert.True(t, strings.HasSuffix(err.Error(), "has no approver"))
	testGenericFileDeleteError(t, gf)

	// Add approver & test
	query := pgmodels.NewQuery().
		Where("work_item_id", "=", workItem.ID)
	req, err := pgmodels.DeletionRequestGet(query)
	require.Nil(t, err)
	require.NotNil(t, req)
	require.NotEqual(t, int64(0), req.ID)

	req.ConfirmedByID = req.RequestedByID
	err = req.Save()
	require.Nil(t, err)

	err = gf.AssertDeletionPreconditions()
	require.Nil(t, err)

	testGenericFileDeleteSuccess(t, gf)
}

func testFileDeletionRequest(t *testing.T, gf *pgmodels.GenericFile, workItemID int64) {
	// Request doesn't exist yet.
	reqView, err := gf.DeletionRequest(workItemID)
	require.NotNil(t, err)
	assert.True(t, pgmodels.IsNoRowError(err))
	require.Nil(t, reqView)

	// Now make the request and we should find it.
	files := []*pgmodels.GenericFile{gf}
	req, err := pgmodels.CreateDeletionRequest(nil, files)
	require.Nil(t, err)
	require.NotNil(t, req)
	req.WorkItemID = workItemID
	req.ConfirmedByID = 0
	require.Nil(t, req.Save())

	reqView, err = gf.DeletionRequest(workItemID)
	require.Nil(t, err)
	require.NotNil(t, reqView)
	assert.Equal(t, req.ID, reqView.ID)
}

// This test is called a number of times above. In each case, we're
// missing some deletion pre-condition, and the Delete() call should
// fail. We check specific errors above. This test just ensures
// that nothing slips through.
func testGenericFileDeleteError(t *testing.T, gf *pgmodels.GenericFile) {
	err := gf.Delete()
	require.NotNil(t, err)
}

// We call this above once we've set up all necessary pre-conditions for
// file deletion, including an approved deletion request and an approved
// work item. This call should succeed.
//
// Note that this also implicitly tests NewFileDeletionEvent. We want to
// make sure 1) the deletion event was created and saved, and 2) it contains
// the right data.
func testGenericFileDeleteSuccess(t *testing.T, gf *pgmodels.GenericFile) {
	err := gf.Delete()
	require.Nil(t, err)

	reloadedFile, err := pgmodels.GenericFileByID(gf.ID)
	require.Nil(t, err)
	require.NotNil(t, reloadedFile)
	assert.Equal(t, constants.StateDeleted, reloadedFile.State)
	assert.True(t, reloadedFile.UpdatedAt.After(time.Now().UTC().Add(-5*time.Second)))

	// Make sure the required event was created.
	deletionEvent, err := reloadedFile.LastDeletionEvent()
	require.Nil(t, err)
	require.NotNil(t, deletionEvent)
	testFileDeletionEventProperties(t, gf, deletionEvent)
}

func testFileDeletionEventProperties(t *testing.T, gf *pgmodels.GenericFile, event *pgmodels.PremisEvent) {
	assert.Equal(t, "APTrust preservation services", event.Agent)
	assert.True(t, event.DateTime.After(time.Now().UTC().Add(-5*time.Second)))
	assert.Equal(t, "File deleted from preservation storage", event.Detail)
	assert.Equal(t, constants.EventDeletion, event.EventType)
	assert.True(t, common.LooksLikeUUID(event.Identifier))
	assert.Equal(t, gf.InstitutionID, event.InstitutionID)
	assert.Equal(t, gf.IntellectualObjectID, event.IntellectualObjectID)
	assert.Equal(t, gf.ID, event.GenericFileID)
	assert.Equal(t, "Minio S3 library", event.Object)
	assert.Equal(t, constants.OutcomeSuccess, event.Outcome)
	assert.Equal(t, "user@test.edu", event.OutcomeDetail)
	assert.Equal(t, "File deleted at the request of user@test.edu. Institutional approver: user@test.edu.", event.OutcomeInformation)
}
