package web_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL = "https://example.com"

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

	del, err := web.NewDeletionForFile(gf.ID, instAdmins[0], baseURL)
	require.Nil(t, err)
	require.NotNil(t, del)

	require.NotNil(t, del.DeletionRequest)
	assert.Equal(t, instAdmins[0].ID, del.DeletionRequest.RequestedByID)
	assert.NotEmpty(t, del.DeletionRequest.GenericFiles)

	assert.ElementsMatch(t, instAdmins, del.InstAdmins)
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
	del, err := web.NewDeletionForFile(gf.ID, &pgmodels.User{}, baseURL)
	assert.Nil(t, del)
	assert.Equal(t, common.ErrPendingWorkItems, err)
}

func TestNewDeletionForReview(t *testing.T) {
	db.LoadFixtures()
	admin, err := pgmodels.UserByEmail("admin@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, admin)

	del, err := web.NewDeletionForReview(1, admin, baseURL, confToken)
	require.Nil(t, err)
	require.NotNil(t, del)

	require.NotNil(t, del.DeletionRequest)
	assert.Equal(t, int64(1), del.DeletionRequest.ID)

	assert.Equal(t, 1, len(del.InstAdmins))
	assert.Equal(t, admin.ID, del.InstAdmins[0].ID)
}

func TestNewDeletionBadToken(t *testing.T) {
	db.LoadFixtures()
	admin, err := pgmodels.UserByEmail("admin@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, admin)

	del, err := web.NewDeletionForReview(1, admin, baseURL, "InvalidToken")
	require.Nil(t, del)
	assert.Equal(t, common.ErrInvalidToken, err)
}

func TestCreateAndQueueWorkItem(t *testing.T) {

}

func TestCreateRequestAlert(t *testing.T) {

}

func TestCreateApprovalAlert(t *testing.T) {

}

func TestCreateCancellationAlert(t *testing.T) {

}

func TestReviewURL(t *testing.T) {

}

func TestWorkItemURL(t *testing.T) {

}

func TestReadOnlyURL(t *testing.T) {

}
