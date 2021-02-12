package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntellectualObjectGetID(t *testing.T) {
	obj, err := models.IntellectualObjectFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, obj)
	assert.Equal(t, int64(1), obj.GetID())
}

func TestIntellectualObjectAuthorize(t *testing.T) {
	sysAdmin, err := ds.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := ds.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := ds.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := ds.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// This object belongs to same inst as instAdmin, instUser
	// and inactiveUser
	obj := &models.IntellectualObject{
		InstitutionID: InstOne,
		State:         "A",
	}

	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionRead))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionDelete))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionRequestDelete))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionApproveDelete))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionFinishBulkDelete))
	assert.Nil(t, obj.Authorize(sysAdmin, constants.ActionRestore))

	// Users other than SysAdmin cannot create or update
	// intellectual objects, and cannot finish bulk delete.
	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(user, constants.ActionFinishBulkDelete))
	}

	// Inst user and admin can read their own institution's objects.
	assert.Nil(t, obj.Authorize(instUser, constants.ActionRead))
	assert.Nil(t, obj.Authorize(instAdmin, constants.ActionRead))

	// Inst admin can restore and delete
	assert.Nil(t, obj.Authorize(instAdmin, constants.ActionRestore))
	assert.Nil(t, obj.Authorize(instAdmin, constants.ActionDelete))
	assert.Nil(t, obj.Authorize(instAdmin, constants.ActionRequestDelete))
	assert.Nil(t, obj.Authorize(instAdmin, constants.ActionApproveDelete))

	// Inst user cannot restore or delete
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(instUser, constants.ActionRestore))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(instUser, constants.ActionDelete))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(instUser, constants.ActionRequestDelete))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(instUser, constants.ActionApproveDelete))

	// Inactive User cannot read, delete, or restore
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(inactiveUser, constants.ActionRead))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(inactiveUser, constants.ActionRestore))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(inactiveUser, constants.ActionDelete))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(inactiveUser, constants.ActionRequestDelete))
	assert.Equal(t, common.ErrPermissionDenied, obj.Authorize(inactiveUser, constants.ActionApproveDelete))

	// No users other than SysAdmin can access objects from
	// other institutions.
	foreignObject := &models.IntellectualObject{
		InstitutionID: InstTwo,
		State:         "A",
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionFinishBulkDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionRestore))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionRequestDelete))
		assert.Equal(t, common.ErrPermissionDenied, foreignObject.Authorize(user, constants.ActionApproveDelete))
	}

}

func TestIntellectualObjectDeleteIsFobidden(t *testing.T) {
	obj := &models.IntellectualObject{}
	assert.False(t, obj.DeleteIsForbidden())
}

func TestIntellectualObjectUpdateIsFobidden(t *testing.T) {
	obj := &models.IntellectualObject{}
	assert.False(t, obj.UpdateIsForbidden())
}

func TestIntellectualObjectIsReadOnly(t *testing.T) {
	obj := &models.IntellectualObject{}
	assert.False(t, obj.IsReadOnly())
}

func TestIntellectualObjectSupportsSoftDelete(t *testing.T) {
	obj := &models.IntellectualObject{}
	assert.True(t, obj.SupportsSoftDelete())
}

func TestIntellectualObjectSoftDeleteAttributes(t *testing.T) {
	obj := &models.IntellectualObject{
		State: "A",
	}
	assert.Equal(t, "A", obj.State)

	obj.SetSoftDeleteAttributes(&models.User{})
	assert.Equal(t, "D", obj.State)

	obj.ClearSoftDeleteAttributes()
	assert.Equal(t, "A", obj.State)
}

func TestIntellectualObjectSetTimestamps(t *testing.T) {
	obj := &models.IntellectualObject{}
	assert.True(t, obj.CreatedAt.IsZero())
	assert.True(t, obj.UpdatedAt.IsZero())

	obj.SetTimestamps()
	assert.False(t, obj.CreatedAt.IsZero())
	assert.False(t, obj.UpdatedAt.IsZero())
}

func TestIntellectualObjectFind(t *testing.T) {
	obj, err := models.IntellectualObjectFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, obj)
	assert.Equal(t, int64(1), obj.ID)
	assert.Equal(t, "institution1.edu/photos", obj.Identifier)
	assert.Equal(t, "First Object for Institution One", obj.Title)
}

func TestIntellectualObjectFindByIdentifier(t *testing.T) {
	obj, err := models.IntellectualObjectFindByIdentifier("institution1.edu/photos")
	require.Nil(t, err)
	require.NotNil(t, obj)
	assert.Equal(t, int64(1), obj.ID)
	assert.Equal(t, "institution1.edu/photos", obj.Identifier)
	assert.Equal(t, "First Object for Institution One", obj.Title)
}
