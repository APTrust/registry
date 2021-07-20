package constants_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {

	// Spot check a few institutional user permissions
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.AlertRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.IntellectualObjectRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.GenericFileRead))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.IntellectualObjectUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.WorkItemUpdate))

	// Spot check a few institutional admin permissions
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.GenericFileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.DeletionRequestApprove))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.GenericFileRestore))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectRestore))

	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.StorageRecordUpdate))

	// Spot check SysAdmin privileges
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.GenericFileUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.IntellectualObjectUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.InstitutionUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.StorageRecordUpdate))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumDelete))

	// Check these because they were misbehaving in dev
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.GenericFileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.GenericFileRestore))

	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.GenericFileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.GenericFileRestore))

	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.GenericFileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.GenericFileRestore))

}
