package models_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGetID(t *testing.T) {
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	assert.Equal(t, int64(1), sysAdmin.GetID())
}

func TestUserDeleteIsForbidden(t *testing.T) {
	user := &models.User{}
	assert.False(t, user.DeleteIsForbidden())
}

func TestUserIsReadOnly(t *testing.T) {
	user := &models.User{}
	assert.False(t, user.IsReadOnly())
}

func TestUserSupportsSoftDelete(t *testing.T) {
	user := &models.User{}
	assert.True(t, user.SupportsSoftDelete())
}

func TestUserSoftDeleteAttributes(t *testing.T) {
	user := &models.User{}
	assert.True(t, user.DeactivatedAt.IsZero())

	user.SetSoftDeleteAttributes(&models.User{})
	assert.False(t, user.DeactivatedAt.IsZero())

	user.ClearSoftDeleteAttributes()
	assert.True(t, user.DeactivatedAt.IsZero())
}

func TestUserSetTimestamps(t *testing.T) {
	user := &models.User{}
	assert.True(t, user.CreatedAt.IsZero())
	assert.True(t, user.UpdatedAt.IsZero())

	user.SetTimestamps()
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserFind(t *testing.T) {
	user, err := models.UserFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
}

func TestSignInUser_Valid(t *testing.T) {
	db.LoadFixtures()
	// Constants below are defined in models/common_test.go
	users := []string{
		SysAdmin,
		InstAdmin,
		InstUser,
	}
	for _, email := range users {
		user, err := models.SignInUser(email, Password, "1.1.1.1")
		require.Nil(t, err)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "1.1.1.1", user.CurrentSignInIP)
		assert.True(t, user.SignInCount > 0)
		assert.InDelta(t, time.Now().Unix(), user.CurrentSignInAt.Unix(), 10)
		assert.NotEmpty(t, user.Role)
		oldSignInTime := user.CurrentSignInAt
		oldSignInCount := user.SignInCount

		user, err = models.SignInUser(email, Password, "2.2.2.2")
		require.Nil(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "2.2.2.2", user.CurrentSignInIP)
		assert.True(t, user.SignInCount > oldSignInCount)
		assert.True(t, user.CurrentSignInAt.After(oldSignInTime))
	}
}

func TestSignInUser_Invalid(t *testing.T) {
	db.LoadFixtures()
	// User does not exist
	user, err := models.SignInUser("noone@example.com", "xyz", "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrInvalidLogin, err)

	// User exists, wrong password
	user, err = models.SignInUser(SysAdmin, "xyz", "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrInvalidLogin, err)
}

func TestSignInUser_Deactivated(t *testing.T) {
	db.LoadFixtures()
	user, err := models.SignInUser(InactiveUser, Password, "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrAccountDeactivated, err)
}

func TestUserDeleteUndelete(t *testing.T) {
	db.LoadFixtures()

	admin, err := getUser()
	require.Nil(t, err)
	require.NotNil(t, admin)
	admin.Role = constants.RoleSysAdmin

	regUser, err := getUser()
	require.Nil(t, err)
	require.NotNil(t, regUser)
	regUser.Role = constants.RoleInstUser

	user, err := getUser()
	require.Nil(t, err)
	err = models.Save(user, admin)
	require.Nil(t, err)
	assert.True(t, user.ID > int64(0))
	assert.True(t, user.DeactivatedAt.IsZero())

	// This should raise an error, since regular user cannot
	// delete users.
	err = models.Delete(user, regUser)
	assert.Equal(t, common.ErrPermissionDenied, err)

	// We don't hard-delete users. We set a timestamp on
	// User.DeactivatedAt to indicate they're no longer active.
	err = models.Delete(user, admin)
	require.Nil(t, err)

	// Reload deleted user. They should exist with a
	// DeactivatedAt timestamp.
	err = models.Find(user, user.ID, admin)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.False(t, user.DeactivatedAt.IsZero())

	// Undelete the bastard.
	err = models.Undelete(user, admin)
	require.Nil(t, err)

	// His deactivation timestamp should be cleared.
	err = models.Find(user, user.ID, admin)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.True(t, user.DeactivatedAt.IsZero())
}

func TestUserHasPermission(t *testing.T) {
	ownInst := int64(2)
	otherInst := int64(3)
	sysAdmin := &models.User{
		Role:          constants.RoleSysAdmin,
		InstitutionID: int64(1),
	}
	instAdmin := &models.User{
		Role:          constants.RoleInstAdmin,
		InstitutionID: ownInst,
	}
	instUser := &models.User{
		Role:          constants.RoleInstUser,
		InstitutionID: ownInst,
	}
	nobody := &models.User{
		Role:          constants.RoleNone,
		InstitutionID: ownInst,
	}

	// Sys Admin has privileges across all institutions
	assert.True(t, sysAdmin.HasPermission(constants.ObjectUpdate, 0))
	assert.True(t, sysAdmin.HasPermission(constants.ObjectUpdate, 1))
	assert.True(t, sysAdmin.HasPermission(constants.ObjectUpdate, 2))
	assert.True(t, sysAdmin.HasPermission(constants.UserUpdate, 1))
	assert.True(t, sysAdmin.HasPermission(constants.InstitutionUpdate, 1))

	// Inst admin can read and update some things at their own institution.
	assert.True(t, instAdmin.HasPermission(constants.ObjectRead, ownInst))
	assert.True(t, instAdmin.HasPermission(constants.UserUpdate, ownInst))
	assert.True(t, instAdmin.HasPermission(constants.ObjectRequestDelete, ownInst))
	assert.True(t, instAdmin.HasPermission(constants.ObjectApproveDelete, ownInst))
	assert.True(t, instAdmin.HasPermission(constants.ObjectRestore, ownInst))
	assert.True(t, instAdmin.HasPermission(constants.FileRestore, ownInst))

	// Inst admin cannot read and update items at their other institutions.
	assert.False(t, instAdmin.HasPermission(constants.ObjectRead, otherInst))
	assert.False(t, instAdmin.HasPermission(constants.UserUpdate, otherInst))
	assert.False(t, instAdmin.HasPermission(constants.ObjectRequestDelete, otherInst))
	assert.False(t, instAdmin.HasPermission(constants.ObjectApproveDelete, otherInst))
	assert.False(t, instAdmin.HasPermission(constants.ObjectRestore, otherInst))
	assert.False(t, instAdmin.HasPermission(constants.FileRestore, otherInst))

	// Inst user can read records at their own institution...
	assert.True(t, instUser.HasPermission(constants.ObjectRead, ownInst))
	assert.True(t, instUser.HasPermission(constants.FileRead, ownInst))
	assert.True(t, instUser.HasPermission(constants.EventRead, ownInst))
	assert.True(t, instUser.HasPermission(constants.ChecksumRead, ownInst))
	assert.True(t, instUser.HasPermission(constants.StorageRecordRead, ownInst))

	// Inst user cannot read records at other institutions.
	assert.False(t, instUser.HasPermission(constants.ObjectRead, otherInst))
	assert.False(t, instUser.HasPermission(constants.FileRead, otherInst))
	assert.False(t, instUser.HasPermission(constants.EventRead, otherInst))
	assert.False(t, instUser.HasPermission(constants.ChecksumRead, otherInst))
	assert.False(t, instUser.HasPermission(constants.StorageRecordRead, otherInst))

	// Inst user can't edit much at any institution
	assert.False(t, instUser.HasPermission(constants.UserUpdate, ownInst))
	assert.False(t, instUser.HasPermission(constants.ObjectRequestDelete, ownInst))
	assert.False(t, instUser.HasPermission(constants.ObjectApproveDelete, ownInst))
	assert.False(t, instUser.HasPermission(constants.ObjectRestore, ownInst))
	assert.False(t, instUser.HasPermission(constants.FileRestore, ownInst))

	assert.False(t, instUser.HasPermission(constants.UserUpdate, otherInst))
	assert.False(t, instUser.HasPermission(constants.ObjectRequestDelete, otherInst))
	assert.False(t, instUser.HasPermission(constants.ObjectApproveDelete, otherInst))
	assert.False(t, instUser.HasPermission(constants.ObjectRestore, otherInst))
	assert.False(t, instUser.HasPermission(constants.FileRestore, otherInst))

	// User with no role has no permissions.
	// RoleNone is set by default until we can determine from the database
	// what the user's actual role is. See User.loadRole().
	assert.False(t, nobody.HasPermission(constants.ObjectRead, otherInst))
	assert.False(t, nobody.HasPermission(constants.FileRead, otherInst))
	assert.False(t, nobody.HasPermission(constants.EventRead, otherInst))
	assert.False(t, nobody.HasPermission(constants.ChecksumRead, otherInst))
	assert.False(t, nobody.HasPermission(constants.StorageRecordRead, otherInst))
	assert.False(t, nobody.HasPermission(constants.UserUpdate, ownInst))
	assert.False(t, nobody.HasPermission(constants.ObjectRequestDelete, ownInst))
	assert.False(t, nobody.HasPermission(constants.ObjectApproveDelete, ownInst))
	assert.False(t, nobody.HasPermission(constants.ObjectRestore, ownInst))
	assert.False(t, nobody.HasPermission(constants.FileRestore, ownInst))
	assert.False(t, nobody.HasPermission(constants.UserUpdate, otherInst))
	assert.False(t, nobody.HasPermission(constants.ObjectRequestDelete, otherInst))
	assert.False(t, nobody.HasPermission(constants.ObjectApproveDelete, otherInst))
	assert.False(t, nobody.HasPermission(constants.ObjectRestore, otherInst))
	assert.False(t, nobody.HasPermission(constants.FileRestore, otherInst))
}

func TestUserDBPerms_SysAdmin(t *testing.T) {
	// These are the various users who will peform test actions.
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)

	// These are the user record we will try to save/update/delete
	inst1User, err := getUser()
	require.Nil(t, err)
	inst1User.InstitutionID = InstOne
	inst1User.Role = constants.RoleInstUser

	inst2User, err := getUser()
	require.Nil(t, err)
	inst2User.InstitutionID = InstTwo
	inst2User.Role = constants.RoleInstUser

	// SysAdmin should be able to perform all actions on any user
	require.Nil(t, models.Save(inst1User, sysAdmin))
	require.Nil(t, models.Save(inst2User, sysAdmin))
	require.Nil(t, models.Delete(inst1User, sysAdmin))
	require.Nil(t, models.Delete(inst2User, sysAdmin))
	require.Nil(t, models.Undelete(inst1User, sysAdmin))
	require.Nil(t, models.Undelete(inst2User, sysAdmin))
}

func TestUserDBPerms_InstAdmin(t *testing.T) {
	instAdmin, err := models.UserFindByEmail(InstAdmin)
	require.Nil(t, err)

	// These are the user record we will try to save/update/delete
	inst1User, err := getUser()
	require.Nil(t, err)
	inst1User.InstitutionID = InstOne
	inst1User.Role = constants.RoleInstUser

	inst2User, err := getUser()
	require.Nil(t, err)
	inst2User.InstitutionID = InstTwo
	inst2User.Role = constants.RoleInstUser

	// Inst Admin can create users at their own institution
	anotherInst1User, err := getUser()
	require.Nil(t, err)
	anotherInst1User.InstitutionID = InstOne
	anotherInst1User.Role = constants.RoleInstUser
	assert.Nil(t, models.Save(anotherInst1User, instAdmin))

	// Inst Admin can edit user at own institution
	assert.Nil(t, models.Save(inst1User, instAdmin))
	assert.Nil(t, models.Delete(inst1User, instAdmin))
	assert.Nil(t, models.Undelete(inst1User, instAdmin))

	// Inst Admin cannot edit user at other institition
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inst2User, instAdmin))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inst2User, instAdmin))
	assert.Equal(t, common.ErrPermissionDenied, models.Undelete(inst2User, instAdmin))

	// Inst Admin can edit self but cannot delete self
	assert.Nil(t, models.Save(instAdmin, instAdmin))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(instAdmin, instAdmin))
}

func TestUserDBPerms_InstUser(t *testing.T) {
	instUser, err := models.UserFindByEmail(InstUser)
	require.Nil(t, err)

	// These are the user record we will try to save/update/delete
	inst1User, err := getUser()
	require.Nil(t, err)
	inst1User.InstitutionID = InstOne
	inst1User.Role = constants.RoleInstUser

	inst2User, err := getUser()
	require.Nil(t, err)
	inst2User.InstitutionID = InstTwo
	inst2User.Role = constants.RoleInstUser

	// Inst User cannot create users
	oneMoreInst1User, err := getUser()
	require.Nil(t, err)
	oneMoreInst1User.InstitutionID = InstOne
	oneMoreInst1User.Role = constants.RoleInstUser
	assert.Equal(t, common.ErrPermissionDenied, models.Save(oneMoreInst1User, instUser))

	// Inst User cannot edit any other users
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inst1User, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inst1User, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Undelete(inst1User, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inst2User, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inst2User, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Undelete(inst2User, instUser))

	// Inst User can edit self but cannot delete self
	assert.Nil(t, models.Save(instUser, instUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(instUser, instUser))
}

func TestUserDBPerms_InactiveUser(t *testing.T) {
	inactiveUser, err := models.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// These are the user record we will try to save/update/delete
	inst1User, err := getUser()
	require.Nil(t, err)
	inst1User.InstitutionID = InstOne
	inst1User.Role = constants.RoleInstUser

	inst2User, err := getUser()
	require.Nil(t, err)
	inst2User.InstitutionID = InstTwo
	inst2User.Role = constants.RoleInstUser

	oneMoreInst1User, err := getUser()
	require.Nil(t, err)
	oneMoreInst1User.InstitutionID = InstOne
	oneMoreInst1User.Role = constants.RoleInstUser

	// Inactive User cannot create or edit any users
	assert.Equal(t, common.ErrPermissionDenied, models.Save(oneMoreInst1User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inst1User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inst1User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Undelete(inst1User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inst2User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inst2User, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Undelete(inst2User, inactiveUser))

	// Inactive User cannot event edit self
	assert.Equal(t, common.ErrPermissionDenied, models.Save(inactiveUser, inactiveUser))
	assert.Equal(t, common.ErrPermissionDenied, models.Delete(inactiveUser, inactiveUser))
}
