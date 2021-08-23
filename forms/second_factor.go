package forms

import (
	"github.com/APTrust/registry/constants"
)

// SecondFactorForm is the form that displays filtering options for
// the alert list page.
type SecondFactorForm struct {
	Form
}

func NewSecondFactorForm() *SecondFactorForm {
	f := &SecondFactorForm{
		Form: NewForm(nil, "users/choose_second_factor.html", "/users"),
	}
	f.init()
	f.SetValues()
	return f
}

func (f *SecondFactorForm) init() {
	f.Fields["second_factor"] = &Field{
		Name:        "second_factor",
		Label:       "Choose Second Factor",
		Placeholder: "",
		Options:     Options(constants.SecondFactorTypes),
		Attrs: map[string]string{
			"required": "",
		},
	}
}

// SetValues sets the form values to match selected values.
// The form must implement this to satisfy the Form interface,
// but in this case, it's a no-op because we don't persist
// this info in any way.
func (f *SecondFactorForm) SetValues() {

}
