package pgmodels_test

import (
	"testing"

	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/stew/slice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHasPermission(t *testing.T) {
	ownInst := int64(2)
	otherInst := int64(3)
	sysAdmin := &pgmodels.User{
		Role:          constants.RoleSysAdmin,
		InstitutionID: int64(1),
	}
	instAdmin := &pgmodels.User{
		Role:          constants.RoleInstAdmin,
		InstitutionID: ownInst,
	}
	instUser := &pgmodels.User{
		Role:          constants.RoleInstUser,
		InstitutionID: ownInst,
	}
	nobody := &pgmodels.User{
		Role:          constants.RoleNone,
		InstitutionID: ownInst,
	}

	// SysAdmin has all permissions at all institutions
	for _, perm := range constants.Permissions {
		if slice.Contains(constants.ForbiddenToAll, perm) {
			assert.False(t, sysAdmin.HasPermission(perm, sysAdmin.InstitutionID), perm)
			assert.False(t, sysAdmin.HasPermission(perm, ownInst), perm)
			assert.False(t, sysAdmin.HasPermission(perm, otherInst), perm)
		} else {
			assert.True(t, sysAdmin.HasPermission(perm, sysAdmin.InstitutionID), perm)
			assert.True(t, sysAdmin.HasPermission(perm, ownInst), perm)
			assert.True(t, sysAdmin.HasPermission(perm, otherInst), perm)
		}
	}

	// InstAdmin has extended permissions at own inst only.
	for _, perm := range constants.Permissions {
		if constants.CheckPermission(instAdmin.Role, perm) {
			assert.True(t, instAdmin.HasPermission(perm, ownInst), perm)
		} else {
			assert.False(t, instAdmin.HasPermission(perm, ownInst), perm)
		}
		assert.False(t, instAdmin.HasPermission(perm, otherInst), perm)
	}

	// InstUser has limited permissions at own inst only.
	for _, perm := range constants.Permissions {
		if constants.CheckPermission(instUser.Role, perm) {
			assert.True(t, instUser.HasPermission(perm, ownInst), perm)
		} else {
			assert.False(t, instUser.HasPermission(perm, ownInst), perm)
		}
		assert.False(t, instUser.HasPermission(perm, otherInst), perm)
	}

	// Nobody can't do nothin.
	for _, perm := range constants.Permissions {
		assert.False(t, nobody.HasPermission(perm, sysAdmin.InstitutionID), perm)
		assert.False(t, nobody.HasPermission(perm, ownInst), perm)
		assert.False(t, nobody.HasPermission(perm, otherInst), perm)
	}
}

func TestUserIsAdmin(t *testing.T) {
	user := pgmodels.User{
		Role: constants.RoleInstAdmin,
	}
	assert.False(t, user.IsAdmin())

	user.Role = constants.RoleInstUser
	assert.False(t, user.IsAdmin())

	user.Role = constants.RoleSysAdmin
	assert.True(t, user.IsAdmin())
}

func TestOTPBackupCodes(t *testing.T) {
	db.LoadFixtures()
	// In fixtures, Inactive User with id #4 has some OTP
	// backup codes. We want to make sure these are bound
	// correctly.
	user, err := pgmodels.UserByID(4)
	require.Nil(t, err)

	codes := []string{
		"code1",
		"code2",
		"code3",
	}
	assert.Equal(t, codes, user.OTPBackupCodes)
}
