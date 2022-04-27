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
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.FileRead))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.IntellectualObjectUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.WorkItemUpdate))

	// Spot check a few institutional admin permissions
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventRead))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.DeletionRequestApprove))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRestore))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectRestore))

	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleInstAdmin, constants.StorageRecordUpdate))

	// Spot check SysAdmin privileges
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.IntellectualObjectUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.InstitutionUpdate))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.StorageRecordUpdate))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventDelete))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.EventUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumUpdate))
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.ChecksumDelete))

	// Check these because they were misbehaving in dev
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstUser, constants.FileRestore))

	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRestore))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileRestore))

}

// Test dangerous permissions to ensure only the right roles have them.
func TestDangerousPermissions(t *testing.T) {

	// Initiate Object Deletion.
	//
	// Admins, yes. Users, no.
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.IntellectualObjectDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.IntellectualObjectDelete))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.IntellectualObjectRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.IntellectualObjectRequestDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.IntellectualObjectRequestDelete))

	// Initiate File Deletion.
	//
	// Admins, yes. Users, no.
	assert.True(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.FileDelete))

	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.FileRequestDelete))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.FileRequestDelete))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.FileRequestDelete))

	// Approve Object Deletion.
	//
	// Inst Admins, yes. Others, no.
	assert.False(t, constants.CheckPermission(constants.RoleSysAdmin, constants.DeletionRequestApprove))
	assert.True(t, constants.CheckPermission(constants.RoleInstAdmin, constants.DeletionRequestApprove))
	assert.False(t, constants.CheckPermission(constants.RoleInstUser, constants.DeletionRequestApprove))

}
