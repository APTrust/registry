package models_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
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

	user, err := getUser()
	require.Nil(t, err)
	err = models.Save(user, user)
	require.Nil(t, err)
	assert.True(t, user.ID > int64(0))
	assert.True(t, user.DeactivatedAt.IsZero())

	// We don't hard-delete users. We set a timestamp on
	// User.DeactivatedAt to indicate they're no longer active.
	err = models.Delete(user, user)
	require.Nil(t, err)

	// Reload deleted user. They should exist with a
	// DeactivatedAt timestamp.
	err = models.Find(user, user.ID, user)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.False(t, user.DeactivatedAt.IsZero())

	// Undelete the bastard.
	err = models.Undelete(user, user)
	require.Nil(t, err)

	// His deactivation timestamp should be cleared.
	err = models.Find(user, user.ID, user)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.True(t, user.DeactivatedAt.IsZero())
}
