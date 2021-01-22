package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StorageRecord #1 belongs to GenericFile #1,
// which belongs to InstOne.

func TestStorageRecordGetID(t *testing.T) {
	sr, err := models.StorageRecordFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, sr)
	assert.Equal(t, int64(1), sr.GetID())
}

func TestStorageRecordAuthorize(t *testing.T) {
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := models.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := models.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := models.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// This StorageRecord's GenericFile belongs to same inst as
	// instAdmin, instUser and inactiveUser
	sr := &models.StorageRecord{
		GenericFileID: int64(1),
	}

	// SysAdmin can create, read, update, delete StorageRecords.
	assert.Nil(t, sr.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, sr.Authorize(sysAdmin, constants.ActionRead))
	assert.Nil(t, sr.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Nil(t, sr.Authorize(sysAdmin, constants.ActionDelete))

	// Users other than SysAdmin cannot create, update or delete
	// StorageRecords.
	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, sr.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, sr.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, sr.Authorize(user, constants.ActionDelete))
	}

	// Inst user and admin can read their own institution's StorageRecords.
	assert.Nil(t, sr.Authorize(instUser, constants.ActionRead))
	assert.Nil(t, sr.Authorize(instAdmin, constants.ActionRead))

	// Inactive User cannot read StorageRecords
	assert.Equal(t, common.ErrPermissionDenied, sr.Authorize(inactiveUser, constants.ActionRead))

	// No users other than SysAdmin can access StorageRecords from
	// other institutions. This StorageRecord's generic file belongs
	// to InstTwo.
	foreignStorageRecord := &models.StorageRecord{
		GenericFileID: int64(11),
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, foreignStorageRecord.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, foreignStorageRecord.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, foreignStorageRecord.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, foreignStorageRecord.Authorize(user, constants.ActionDelete))
	}
}

func TestStorageRecordIsReadOnly(t *testing.T) {
	sr := &models.StorageRecord{}
	assert.False(t, sr.IsReadOnly())
}

func TestStorageRecordDeleteIsFobidden(t *testing.T) {
	sr := &models.StorageRecord{}
	assert.False(t, sr.DeleteIsForbidden())
}

func TestStorageRecordUpdateIsFobidden(t *testing.T) {
	sr := &models.StorageRecord{}
	assert.False(t, sr.UpdateIsForbidden())
}

func TestStorageRecordSupportsSoftDelete(t *testing.T) {
	sr := &models.StorageRecord{}
	assert.False(t, sr.SupportsSoftDelete())
}

func TestStorageRecordSoftDeleteAttributes(t *testing.T) {
	// No-op
}

func TestStorageRecordSetTimestamps(t *testing.T) {
	// No-op
}

func TestStorageRecordFind(t *testing.T) {
	sr, err := models.StorageRecordFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, sr)
	assert.Equal(t, int64(1), sr.ID)
	assert.EqualValues(t, 1, sr.GenericFileID)
	assert.EqualValues(t, "https://localhost:9899/preservation-va/25452f41-1b18-47b7-b334-751dfd5d011e", sr.URL)
}

func TestStorageRecordsForFile(t *testing.T) {
	records, err := models.StorageRecordsForFile(int64(1))
	require.Nil(t, err)
	require.NotEmpty(t, records)
	urls := []string{
		"https://localhost:9899/preservation-or/25452f41-1b18-47b7-b334-751dfd5d011e",
		"https://localhost:9899/preservation-va/25452f41-1b18-47b7-b334-751dfd5d011e",
	}
	for i, sr := range records {
		assert.Equal(t, int64(1), sr.GenericFileID)
		assert.Equal(t, urls[i], records[i].URL)
	}
}
