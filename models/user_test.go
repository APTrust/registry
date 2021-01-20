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

func TestUserCan(t *testing.T) {
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
	assert.True(t, sysAdmin.Can(constants.ObjectUpdate, 0))
	assert.True(t, sysAdmin.Can(constants.ObjectUpdate, 1))
	assert.True(t, sysAdmin.Can(constants.ObjectUpdate, 2))
	assert.True(t, sysAdmin.Can(constants.UserUpdate, 1))
	assert.True(t, sysAdmin.Can(constants.InstitutionUpdate, 1))

	// Inst admin can read and update some things at their own institution.
	assert.True(t, instAdmin.Can(constants.ObjectRead, ownInst))
	assert.True(t, instAdmin.Can(constants.UserUpdate, ownInst))
	assert.True(t, instAdmin.Can(constants.ObjectRequestDelete, ownInst))
	assert.True(t, instAdmin.Can(constants.ObjectApproveDelete, ownInst))
	assert.True(t, instAdmin.Can(constants.ObjectRestore, ownInst))
	assert.True(t, instAdmin.Can(constants.FileRestore, ownInst))

	// Inst admin cannot read and update items at their other institutions.
	assert.False(t, instAdmin.Can(constants.ObjectRead, otherInst))
	assert.False(t, instAdmin.Can(constants.UserUpdate, otherInst))
	assert.False(t, instAdmin.Can(constants.ObjectRequestDelete, otherInst))
	assert.False(t, instAdmin.Can(constants.ObjectApproveDelete, otherInst))
	assert.False(t, instAdmin.Can(constants.ObjectRestore, otherInst))
	assert.False(t, instAdmin.Can(constants.FileRestore, otherInst))

	// Inst user can read records at their own institution...
	assert.True(t, instUser.Can(constants.ObjectRead, ownInst))
	assert.True(t, instUser.Can(constants.FileRead, ownInst))
	assert.True(t, instUser.Can(constants.EventRead, ownInst))
	assert.True(t, instUser.Can(constants.ChecksumRead, ownInst))
	assert.True(t, instUser.Can(constants.StorageRecordRead, ownInst))

	// Inst user cannot read records at other institutions.
	assert.False(t, instUser.Can(constants.ObjectRead, otherInst))
	assert.False(t, instUser.Can(constants.FileRead, otherInst))
	assert.False(t, instUser.Can(constants.EventRead, otherInst))
	assert.False(t, instUser.Can(constants.ChecksumRead, otherInst))
	assert.False(t, instUser.Can(constants.StorageRecordRead, otherInst))

	// Inst user can't edit much at any institution
	assert.False(t, instUser.Can(constants.UserUpdate, ownInst))
	assert.False(t, instUser.Can(constants.ObjectRequestDelete, ownInst))
	assert.False(t, instUser.Can(constants.ObjectApproveDelete, ownInst))
	assert.False(t, instUser.Can(constants.ObjectRestore, ownInst))
	assert.False(t, instUser.Can(constants.FileRestore, ownInst))

	assert.False(t, instUser.Can(constants.UserUpdate, otherInst))
	assert.False(t, instUser.Can(constants.ObjectRequestDelete, otherInst))
	assert.False(t, instUser.Can(constants.ObjectApproveDelete, otherInst))
	assert.False(t, instUser.Can(constants.ObjectRestore, otherInst))
	assert.False(t, instUser.Can(constants.FileRestore, otherInst))

	// User with no role has no permissions.
	// RoleNone is set by default until we can determine from the database
	// what the user's actual role is. See User.loadRole().
	assert.False(t, nobody.Can(constants.ObjectRead, otherInst))
	assert.False(t, nobody.Can(constants.FileRead, otherInst))
	assert.False(t, nobody.Can(constants.EventRead, otherInst))
	assert.False(t, nobody.Can(constants.ChecksumRead, otherInst))
	assert.False(t, nobody.Can(constants.StorageRecordRead, otherInst))
	assert.False(t, nobody.Can(constants.UserUpdate, ownInst))
	assert.False(t, nobody.Can(constants.ObjectRequestDelete, ownInst))
	assert.False(t, nobody.Can(constants.ObjectApproveDelete, ownInst))
	assert.False(t, nobody.Can(constants.ObjectRestore, ownInst))
	assert.False(t, nobody.Can(constants.FileRestore, ownInst))
	assert.False(t, nobody.Can(constants.UserUpdate, otherInst))
	assert.False(t, nobody.Can(constants.ObjectRequestDelete, otherInst))
	assert.False(t, nobody.Can(constants.ObjectApproveDelete, otherInst))
	assert.False(t, nobody.Can(constants.ObjectRestore, otherInst))
	assert.False(t, nobody.Can(constants.FileRestore, otherInst))

}
