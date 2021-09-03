package forms

import (
	"github.com/APTrust/registry/pgmodels"
)

type TwoFactorSetupForm struct {
	Form
}

func NewTwoFactorSetupForm(user *pgmodels.User) *TwoFactorSetupForm {
	form := &TwoFactorSetupForm{
		Form: NewForm(user, "users/init_2fa_setup.html", "/users"),
	}
	form.init()
	form.SetValues()
	return form
}

func (f *TwoFactorSetupForm) init() {
	f.Fields["AuthyStatus"] = &Field{
		Name:        "AuthyStatus",
		Label:       "Preferred Method for Two-Factor Auth",
		Placeholder: "",
		ErrMsg:      "Please choose your preferred method.",
		Options:     TwoFactorMethodList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["PhoneNumber"] = &Field{
		Name:        "PhoneNumber",
		Label:       "PhoneNumber",
		Placeholder: "PhoneNumber",
		ErrMsg:      pgmodels.ErrUserPhone,
		Attrs: map[string]string{
			"required": "",
		},
	}
}

// setValues sets the form values
func (f *TwoFactorSetupForm) SetValues() {
	user := f.Model.(*pgmodels.User)
	f.Fields["PhoneNumber"].Value = user.PhoneNumber
	f.Fields["AuthyStatus"].Value = user.AuthyStatus
}
