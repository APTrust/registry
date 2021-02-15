package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPremisEventGetID(t *testing.T) {
	event := &models.PremisEvent{ID: int64(1)}
	assert.Equal(t, int64(1), event.GetID())
}

func TestPremisEventAuthorize(t *testing.T) {
	sysAdmin, err := ds.UserFindByEmail(SysAdmin)
	require.Nil(t, err)
	instAdmin, err := ds.UserFindByEmail(InstAdmin)
	require.Nil(t, err)
	instUser, err := ds.UserFindByEmail(InstUser)
	require.Nil(t, err)
	inactiveUser, err := ds.UserFindByEmail(InactiveUser)
	require.Nil(t, err)

	// This event belongs to same inst as instAdmin, instUser
	// and inactiveUser
	event := &models.PremisEvent{
		InstitutionID: InstOne,
	}

	// SysAdmin can create and read PremisEvents, but no one can
	// update or delete them.
	assert.Nil(t, event.Authorize(sysAdmin, constants.ActionCreate))
	assert.Nil(t, event.Authorize(sysAdmin, constants.ActionRead))
	assert.Equal(t, common.ErrPermissionDenied, event.Authorize(sysAdmin, constants.ActionUpdate))
	assert.Equal(t, common.ErrPermissionDenied, event.Authorize(sysAdmin, constants.ActionDelete))

	// Users other than SysAdmin cannot create, update, or delete events.
	otherUsers := []*models.User{
		instAdmin,
		instUser,
		inactiveUser,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, event.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, event.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, event.Authorize(user, constants.ActionDelete))
	}

	// Inst user and admin can read their own institution's events.
	assert.Nil(t, event.Authorize(instUser, constants.ActionRead))
	assert.Nil(t, event.Authorize(instAdmin, constants.ActionRead))

	// Inactive User cannot read
	assert.Equal(t, common.ErrPermissionDenied, event.Authorize(inactiveUser, constants.ActionRead))

	// No users other than SysAdmin can read events from
	// other institutions.
	foreignEvent := &models.PremisEvent{
		InstitutionID: InstTwo,
	}
	for _, user := range otherUsers {
		assert.Equal(t, common.ErrPermissionDenied, foreignEvent.Authorize(user, constants.ActionCreate))
		assert.Equal(t, common.ErrPermissionDenied, foreignEvent.Authorize(user, constants.ActionUpdate))
		assert.Equal(t, common.ErrPermissionDenied, foreignEvent.Authorize(user, constants.ActionRead))
		assert.Equal(t, common.ErrPermissionDenied, foreignEvent.Authorize(user, constants.ActionDelete))
	}

}

func TestPremisEventIsReadOnly(t *testing.T) {
	event := &models.PremisEvent{}
	assert.False(t, event.IsReadOnly())
}

func TestPremisEventDeleteIsFobidden(t *testing.T) {
	event := &models.PremisEvent{}
	assert.True(t, event.DeleteIsForbidden())
}

func TestPremisEventUpdateIsFobidden(t *testing.T) {
	event := &models.PremisEvent{}
	assert.True(t, event.UpdateIsForbidden())
}

func TestPremisEventSupportsSoftDelete(t *testing.T) {
	event := &models.PremisEvent{}
	assert.False(t, event.SupportsSoftDelete())
}

func TestPremisEventSoftDeleteAttributes(t *testing.T) {
	// No-op
}

func TestPremisEventSetTimestamps(t *testing.T) {
	event := &models.PremisEvent{}
	assert.True(t, event.CreatedAt.IsZero())
	assert.True(t, event.UpdatedAt.IsZero())

	event.SetTimestamps()
	assert.False(t, event.CreatedAt.IsZero())
	assert.False(t, event.UpdatedAt.IsZero())
}
