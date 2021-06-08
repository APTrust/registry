package web_test

import (
	"testing"

	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL = "https://example.com"

func getFileAndInstAdmins() (*pgmodels.GenericFile, []*pgmodels.User, error) {
	db.LoadFixtures()
	gfQuery := pgmodels.NewQuery().Where("institution_id", "=", 2).Limit(1).Offset(0)
	gf, err := pgmodels.GenericFileGet(gfQuery)
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

func TestNewDeletionForReview(t *testing.T) {
	db.LoadFixtures()
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
