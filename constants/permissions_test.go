package constants_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {

	// Spot check a few institutional user permissions
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.AlertRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.ObjectRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.FileRead))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.ObjectUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.WorkItemUpdate))

	// Spot check a few institutional admin permissions
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileApproveDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRestore))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ObjectRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ObjectApproveDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ObjectRestore))

	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.StorageRecordUpdate))

	// Spot check SysAdmin privileges
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ObjectUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.InstitutionUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.StorageRecordUpdate))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumDelete))

}
