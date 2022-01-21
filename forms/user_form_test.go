package forms_test

import (
	"testing"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserForm(t *testing.T) {
	sysAdmin, err := pgmodels.UserByEmail("system@aptrust.org")
	require.Nil(t, err)
	user, err := pgmodels.UserByID(2)
	require.Nil(t, err)
	form, err := forms.NewUserForm(user, sysAdmin)
	require.Nil(t, err)
	require.NotNil(t, form)
	assert.Equal(t, 7, len(form.Fields))

	assert.Equal(t, user.Name, form.Fields["Name"].Value)
	assert.Equal(t, user.Email, form.Fields["Email"].Value)
	assert.Equal(t, user.PhoneNumber, form.Fields["PhoneNumber"].Value)
	assert.Equal(t, user.OTPRequiredForLogin, form.Fields["OTPRequiredForLogin"].Value)
	assert.Equal(t, user.GracePeriod.Format("2006-01-02"), form.Fields["GracePeriod"].Value)
	assert.Equal(t, user.InstitutionID, form.Fields["InstitutionID"].Value)
	assert.Equal(t, user.Role, form.Fields["Role"].Value)

	assert.Equal(t, "/users/edit/2", form.Action())
	assert.Equal(t, "/users/show/2", form.PostSaveURL())
}
