package forms_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTwoFactorSetupForm(t *testing.T) {
	user := &pgmodels.User{}
	user.ID = 99999
	user.AuthyStatus = constants.TwoFactorSMS
	user.PhoneNumber = "+1 505-867-5309"
	form := forms.NewTwoFactorSetupForm(user)
	require.NotNil(t, form)
	assert.Equal(t, 2, len(form.Fields))

	assert.NotNil(t, form.Fields["AuthyStatus"])
	assert.Equal(t, user.AuthyStatus, form.Fields["AuthyStatus"].Value)
	assert.NotNil(t, form.Fields["PhoneNumber"])
	assert.Equal(t, user.PhoneNumber, form.Fields["PhoneNumber"].Value)

	assert.Equal(t, "/users/edit/99999", form.Action())
	assert.Equal(t, "/users/show/99999", form.PostSaveURL())
}
