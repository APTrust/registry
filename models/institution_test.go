package models_test

import (
	"testing"
	//	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	//	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionGetID(t *testing.T) {
	inst, err := models.InstitutionFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, int64(1), inst.GetID())
}

func TestInstitutionAuthorize(t *testing.T) {
	sysAdmin, err := models.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := models.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := models.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := models.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	inst := &models.Institution{}

	assert.Nil(t, inst.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, inst.Authorize(sysAdmin, constants.ActionRead))
	assert.Nil(t, inst.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Nil(t, inst.Authorize(sysAdmin, constants.ActionDelete))

	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, inst.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, inst.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, inst.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, inst.Authorize(user, constants.ActionDelete))
	}
}

func TestInstitutionIsReadOnly(t *testing.T) {
	inst := &models.Institution{}
	assert.False(t, inst.IsReadOnly())
}

func TestInstitutionSupportsSoftDelete(t *testing.T) {
	inst := &models.Institution{}
	assert.True(t, inst.SupportsSoftDelete())
}

func TestInstitutionSoftDeleteAttributes(t *testing.T) {
	inst := &models.Institution{}
	assert.True(t, inst.DeactivatedAt.IsZero())

	inst.SetSoftDeleteAttributes(&models.User{})
	assert.False(t, inst.DeactivatedAt.IsZero())

	inst.ClearSoftDeleteAttributes()
	assert.True(t, inst.DeactivatedAt.IsZero())
}

func TestInstitutionSetTimestamps(t *testing.T) {
	inst := &models.Institution{}
	assert.True(t, inst.CreatedAt.IsZero())
	assert.True(t, inst.UpdatedAt.IsZero())

	inst.SetTimestamps()
	assert.False(t, inst.CreatedAt.IsZero())
	assert.False(t, inst.UpdatedAt.IsZero())
}

func TestInstitutionFind(t *testing.T) {
	inst, err := models.InstitutionFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, int64(1), inst.ID)
	assert.Equal(t, "aptrust.org", inst.Identifier)
}

func TestInstitutionFindByIdentifier(t *testing.T) {
	inst, err := models.InstitutionFindByIdentifier("aptrust.org")
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, int64(1), inst.ID)
	assert.Equal(t, "aptrust.org", inst.Identifier)
}
