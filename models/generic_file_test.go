package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericFileGetID(t *testing.T) {
	gf, err := models.GenericFileFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.GetID())
}

func TestGenericFileAuthorize(t *testing.T) {
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := models.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := models.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := models.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// This file belongs to same inst as instAdmin, instUser
	// and inactiveUser
	gf := &models.GenericFile{
		InstitutionID: InstOne,
		State:         "A",
	}

	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionRead))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionDelete))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionRequestDelete))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionApproveDelete))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionFinishBulkDelete))
	assert.Nil(t, gf.Authorize(sysAdmin, constants.ActionRestore))

	// Users other than SysAdmin cannot create or update
	// files, and cannot finish bulk delete.
	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(user, constants.ActionFinishBulkDelete))
	}

	// Inst user and admin can read their own institution's files.
	assert.Nil(t, gf.Authorize(instUser, constants.ActionRead))
	assert.Nil(t, gf.Authorize(instAdmin, constants.ActionRead))

	// Inst admin can restore and delete
	assert.Nil(t, gf.Authorize(instAdmin, constants.ActionRestore))
	assert.Nil(t, gf.Authorize(instAdmin, constants.ActionDelete))
	assert.Nil(t, gf.Authorize(instAdmin, constants.ActionRequestDelete))
	assert.Nil(t, gf.Authorize(instAdmin, constants.ActionApproveDelete))

	// Inst user cannot restore or delete
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(instUser, constants.ActionRestore))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(instUser, constants.ActionDelete))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(instUser, constants.ActionRequestDelete))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(instUser, constants.ActionApproveDelete))

	// Inactive User cannot read, delete, or restore
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(inactiveUser, constants.ActionRead))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(inactiveUser, constants.ActionRestore))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(inactiveUser, constants.ActionDelete))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(inactiveUser, constants.ActionRequestDelete))
	assert.Equal(t, common.ErrPermissionDenied, gf.Authorize(inactiveUser, constants.ActionApproveDelete))

	// No users other than SysAdmin can access files from
	// other institutions.
	foreignFile := &models.GenericFile{
		InstitutionID: InstTwo,
		State:         "A",
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionFinishBulkDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionRestore))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionRequestDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignFile.Authorize(user, constants.ActionApproveDelete))
	}

}

func TestGenericFileDeleteIsFobidden(t *testing.T) {
	gf := &models.GenericFile{}
	assert.False(t, gf.DeleteIsForbidden())
}

func TestGenericFileUpdateIsFobidden(t *testing.T) {
	gf := &models.GenericFile{}
	assert.False(t, gf.UpdateIsForbidden())
}

func TestGenericFileIsReadOnly(t *testing.T) {
	gf := &models.GenericFile{}
	assert.False(t, gf.IsReadOnly())
}

func TestGenericFileSupportsSoftDelete(t *testing.T) {
	gf := &models.GenericFile{}
	assert.True(t, gf.SupportsSoftDelete())
}

func TestGenericFileSoftDeleteAttributes(t *testing.T) {
	gf := &models.GenericFile{
		State: "A",
	}
	assert.Equal(t, "A", gf.State)

	gf.SetSoftDeleteAttributes(&models.User{})
	assert.Equal(t, "D", gf.State)

	gf.ClearSoftDeleteAttributes()
	assert.Equal(t, "A", gf.State)
}

func TestGenericFileSetTimestamps(t *testing.T) {
	gf := &models.GenericFile{}
	assert.True(t, gf.CreatedAt.IsZero())
	assert.True(t, gf.UpdatedAt.IsZero())

	gf.SetTimestamps()
	assert.False(t, gf.CreatedAt.IsZero())
	assert.False(t, gf.UpdatedAt.IsZero())
}

func TestGenericFileFind(t *testing.T) {
	gf, err := models.GenericFileFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)
	assert.Equal(t, "institution1.edu/photos/picture1", gf.Identifier)
	assert.Equal(t, int64(48771), gf.Size)
}

func TestGenericFileFindByIdentifier(t *testing.T) {
	gf, err := models.GenericFileFindByIdentifier("institution1.edu/photos/picture1")
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)
	assert.Equal(t, "institution1.edu/photos/picture1", gf.Identifier)
	assert.Equal(t, int64(48771), gf.Size)
}
