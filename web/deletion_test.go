package web_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var exampleURL = "https://example.com"

// As documented in db/fixtures/README, this is the confirmation
// token for all DeletionRequests in the fixture data.
var confToken = "ConfirmationToken"

func getFileAndInstAdmins() (*pgmodels.GenericFile, []*pgmodels.User, error) {
	db.LoadFixtures()
	gf, err := pgmodels.GenericFileByIdentifier("institution1.edu/glass/shard3")
	if err != nil {
		return nil, nil, err
	}
	instAdminQuery := pgmodels.NewQuery().Where("institution_id", "=", 2).Where("role", "=", constants.RoleInstAdmin)
	instAdmins, err := pgmodels.UserSelect(instAdminQuery)
	if err != nil {
		return nil, nil, err
	}
	return gf, instAdmins, nil
}

func TestNewDeletionForFile(t *testing.T) {
	db.LoadFixtures()
	gf, instAdmins, err := getFileAndInstAdmins()
	require.Nil(t, err)
	require.True(t, len(instAdmins) > 0)

	del, err := web.NewDeletionForFile(gf.ID, instAdmins[0], exampleURL)
	require.Nil(t, err)
	require.NotNil(t, del)

	require.NotNil(t, del.DeletionRequest)
	assert.Equal(t, instAdmins[0].ID, del.DeletionRequest.RequestedByID)
	assert.NotEmpty(t, del.DeletionRequest.GenericFiles)

	assert.ElementsMatch(t, instAdmins, del.InstAdmins)

	expectedReviewURL := fmt.Sprintf("https://example.com/deletions/review/4?token=%s", del.DeletionRequest.ConfirmationToken)
	actualReviewURL, err := del.ReviewURL()
	require.Nil(t, err)
	assert.Equal(t, expectedReviewURL, actualReviewURL)

	// CreateRequestAlert only works on new deletion requests.
	testCreateRequestAlert(t, del)
}

func TestNewDeletionForFileWithPendingItems(t *testing.T) {
	db.LoadFixtures()

	// Generic file test fixture 49 has a pending restoration
	// WorkItem. This should prevent us from initializing a
	// deletion request, since the deletion would conflict with
	// the restoration.
	gf, err := pgmodels.GenericFileByID(49)
	require.Nil(t, err)
	require.NotNil(t, gf)

	// The user param doesn't matter here, because we should get
	// ErrPendingWorkItems before the function even looks at the user.
	del, err := web.NewDeletionForFile(gf.ID, &pgmodels.User{}, exampleURL)
	assert.Nil(t, del)
	assert.Equal(t, common.ErrPendingWorkItems, err)
}

func TestNewDeletionBadToken(t *testing.T) {
	db.LoadFixtures()
	admin, err := pgmodels.UserByEmail("admin@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, admin)

	del, err := web.NewDeletionForReview(1, admin, exampleURL, "InvalidToken")
	require.Nil(t, del)
	assert.Equal(t, common.ErrInvalidToken, err)
}

func TestNewDeletionForReview(t *testing.T) {
	db.LoadFixtures()
	admin, err := pgmodels.UserByEmail("admin@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, admin)

	del, err := web.NewDeletionForReview(1, admin, exampleURL, confToken)
	require.Nil(t, err)
	require.NotNil(t, del)

	require.NotNil(t, del.DeletionRequest)
	assert.Equal(t, int64(1), del.DeletionRequest.ID)

	assert.Equal(t, 1, len(del.InstAdmins))
	assert.Equal(t, admin.ID, del.InstAdmins[0].ID)

	del.DeletionRequest.Cancel(admin)
	err = del.DeletionRequest.Save()
	require.Nil(t, err)
	testCreateCancellationAlert(t, del)

	del.DeletionRequest.CancelledByID = 0
	del.DeletionRequest.CancelledAt = time.Time{}
	del.DeletionRequest.Confirm(admin)
	err = del.DeletionRequest.Save()
	require.Nil(t, err)

	testCreateAndQueueWorkItem(t, del)
	testCreateApprovalAlert(t, del)

	readOnlyURL := fmt.Sprintf("https://example.com/deletions/show/%d", del.DeletionRequest.ID)
	assert.Equal(t, readOnlyURL, del.ReadOnlyURL())

	expectedWorkItemURL := fmt.Sprintf("https://example.com/work_items/show/%d", del.DeletionRequest.WorkItemID)
	actualWorkItemURL, err := del.WorkItemURL()
	require.Nil(t, err)
	assert.Equal(t, expectedWorkItemURL, actualWorkItemURL)
}

func testCreateAndQueueWorkItem(t *testing.T, del *web.Deletion) {
	item, err := del.CreateAndQueueWorkItem()
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.True(t, item.ID > 0)
	assert.Equal(t, del.DeletionRequest.GenericFiles[0].ID, item.GenericFileID)
	assert.Equal(t, constants.ActionDelete, item.Action)
}

func testCreateRequestAlert(t *testing.T, del *web.Deletion) {
	alert, err := del.CreateRequestAlert()
	require.Nil(t, err)
	require.NotNil(t, alert)
	assert.Equal(t, constants.AlertDeletionRequested, alert.Type)
	assert.Equal(t, del.DeletionRequest.ID, alert.DeletionRequestID)
	assert.Equal(t, del.DeletionRequest.InstitutionID, alert.InstitutionID)
	assert.True(t, len(alert.Content) > 100)
	assert.True(t, len(alert.Users) > 0)
	for _, recipient := range alert.Users {
		assert.Equal(t, constants.RoleInstAdmin, recipient.Role)
		assert.Equal(t, del.DeletionRequest.InstitutionID, recipient.InstitutionID)
	}

	reviewURL, err := del.ReviewURL()
	require.Nil(t, err)
	assert.Contains(t, alert.Content, reviewURL)
}

func testCreateApprovalAlert(t *testing.T, del *web.Deletion) {
	alert, err := del.CreateApprovalAlert()
	require.Nil(t, err)
	require.NotNil(t, alert)
	assert.Equal(t, constants.AlertDeletionConfirmed, alert.Type)
	assert.Equal(t, del.DeletionRequest.ID, alert.DeletionRequestID)
	assert.Equal(t, del.DeletionRequest.InstitutionID, alert.InstitutionID)
	assert.True(t, len(alert.Content) > 100)
	assert.True(t, len(alert.Users) > 0)
	for _, recipient := range alert.Users {
		assert.Equal(t, constants.RoleInstAdmin, recipient.Role)
		assert.Equal(t, del.DeletionRequest.InstitutionID, recipient.InstitutionID)
	}

	workItemURL, err := del.WorkItemURL()
	require.Nil(t, err)
	assert.Contains(t, alert.Content, workItemURL)
}

func testCreateCancellationAlert(t *testing.T, del *web.Deletion) {
	alert, err := del.CreateCancellationAlert()
	require.Nil(t, err)
	require.NotNil(t, alert)
	assert.Equal(t, constants.AlertDeletionCancelled, alert.Type)
	assert.Equal(t, del.DeletionRequest.ID, alert.DeletionRequestID)
	assert.Equal(t, del.DeletionRequest.InstitutionID, alert.InstitutionID)
	assert.True(t, len(alert.Content) > 100)
	assert.True(t, len(alert.Users) > 0)
	for _, recipient := range alert.Users {
		assert.Equal(t, constants.RoleInstAdmin, recipient.Role)
		assert.Equal(t, del.DeletionRequest.InstitutionID, recipient.InstitutionID)
	}
	assert.Contains(t, alert.Content, del.ReadOnlyURL())
}
