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
	cs, err := models.ChecksumFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.Equal(t, int64(1), cs.GetID())
}

func TestChecksumAuthorize(t *testing.T) {
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := models.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := models.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := models.UserFindByEmail(InactiveUser)
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

func TestChecksumFind(t *testing.T) {
	cs, err := models.ChecksumFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.Equal(t, int64(1), cs.ID)
	assert.EqualValues(t, 1, cs.GenericFileID)
	assert.EqualValues(t, "md5", cs.Algorithm)
	assert.Equal(t, "12345678", cs.Digest)
}

func TestChecksumsForFile(t *testing.T) {
	checksums, err := models.ChecksumsForFile(int64(21))
	require.Nil(t, err)
	require.NotEmpty(t, checksums)
	algs := []string{
		"md5",
		"sha1",
		"sha256",
		"sha512",
	}
	for i, cs := range checksums {
		assert.Equal(t, int64(21), cs.GenericFileID)
		assert.Equal(t, algs[i], checksums[i].Algorithm)
	}
}
