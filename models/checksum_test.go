package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksumGetID(t *testing.T) {
	cs := &models.Checksum{ID: int64(1)}
	assert.Equal(t, int64(1), cs.GetID())
}

func TestChecksumAuthorize(t *testing.T) {
	sysAdmin, err := ds.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := ds.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := ds.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := ds.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// This checksum's GenericFile belongs to same inst as
	// instAdmin, instUser and inactiveUser
	cs := &models.Checksum{
		GenericFileID: int64(1),
	}

	// SysAdmin can create and read Checksums, but no one can
	// update or delete them.
	assert.Nil(t, cs.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, cs.Authorize(sysAdmin, constants.ActionRead))
	assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(sysAdmin, constants.ActionDelete))

	// Users other than SysAdmin cannot create, update, or delete checksums.
	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(user, constants.ActionDelete))
	}

	// Inst user and admin can read their own institution's checksums.
	assert.Nil(t, cs.Authorize(instUser, constants.ActionRead))
	assert.Nil(t, cs.Authorize(instAdmin, constants.ActionRead))

	// Inactive User cannot read checksums
	assert.Equal(t, common.ErrPermissionDenied, cs.Authorize(inactiveUser, constants.ActionRead))

	// No users other than SysAdmin can events checksums from
	// other institutions. This checksum's generic file belongs
	// to InstTwo.
	foreignChecksum := &models.Checksum{
		GenericFileID: int64(11),
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, foreignChecksum.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, foreignChecksum.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, foreignChecksum.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, foreignChecksum.Authorize(user, constants.ActionDelete))
	}

}

func TestChecksumIsReadOnly(t *testing.T) {
	cs := &models.Checksum{}
	assert.False(t, cs.IsReadOnly())
}

func TestChecksumDeleteIsFobidden(t *testing.T) {
	cs := &models.Checksum{}
	assert.True(t, cs.DeleteIsForbidden())
}

func TestChecksumUpdateIsFobidden(t *testing.T) {
	cs := &models.Checksum{}
	assert.True(t, cs.UpdateIsForbidden())
}

func TestChecksumSupportsSoftDelete(t *testing.T) {
	cs := &models.Checksum{}
	assert.False(t, cs.SupportsSoftDelete())
}

func TestChecksumSoftDeleteAttributes(t *testing.T) {
	// No-op
}

func TestChecksumSetTimestamps(t *testing.T) {
	cs := &models.Checksum{}
	assert.True(t, cs.CreatedAt.IsZero())
	assert.True(t, cs.UpdatedAt.IsZero())

	cs.SetTimestamps()
	assert.False(t, cs.CreatedAt.IsZero())
	assert.False(t, cs.UpdatedAt.IsZero())
}
