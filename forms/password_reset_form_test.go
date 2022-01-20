package forms_test

import (
	"testing"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetForm(t *testing.T) {
	user := &pgmodels.User{}
	user.ID = 99999
	form := forms.NewPasswordResetForm(user)
	require.NotNil(t, form)
	assert.Equal(t, 3, len(form.Fields))

	assert.NotNil(t, form.Fields["OldPassword"])
	assert.Nil(t, form.Fields["OldPassword"].Value)
	assert.NotNil(t, form.Fields["NewPassword"])
	assert.Nil(t, form.Fields["NewPassword"].Value)
	assert.NotNil(t, form.Fields["ConfirmNewPassword"])
	assert.Nil(t, form.Fields["ConfirmNewPassword"].Value)

	assert.Equal(t, "/users/edit/99999", form.Action())
	assert.Equal(t, "/users/show/99999", form.PostSaveURL())
}
