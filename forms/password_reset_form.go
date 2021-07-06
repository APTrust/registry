package forms

import (
	"github.com/APTrust/registry/pgmodels"
)

type PasswordResetForm struct {
	Form
}

func NewPasswordResetForm(userToEdit *pgmodels.User) *PasswordResetForm {
	pwdForm := &PasswordResetForm{
		Form: NewForm(userToEdit, "users/reset_password.html", "/users"),
	}
	pwdForm.init()
	pwdForm.SetValues()
	return pwdForm
}

func (f *PasswordResetForm) init() {
	f.Fields["OldPassword"] = &Field{
		Name:        "OldPassword",
		ErrMsg:      pgmodels.ErrUserPwdIncorrect,
		Label:       "Current Password",
		Placeholder: "Current Password",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["NewPassword"] = &Field{
		Name:        "NewPassword",
		Label:       "New Password",
		Placeholder: "New Password",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["ConfirmNewPassword"] = &Field{
		Name:        "ConfirmNewPassword",
		Label:       "Confirm New Password",
		Placeholder: "Confirm New Password",
		Attrs: map[string]string{
			"required": "",
		},
	}
}

// setValues sets the form values to match the User values.
func (f *PasswordResetForm) SetValues() {
	// No-op
}
