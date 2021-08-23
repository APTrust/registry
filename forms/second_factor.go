package forms

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// SecondFactorForm is the form that displays filtering options for
// the alert list page.
type SecondFactorForm struct {
	Form
	FilterCollection *pgmodels.FilterCollection
}

func NewSecondFactorForm(fc *pgmodels.FilterCollection) *SecondFactorForm {
	f := &SecondFactorForm{
		Form:             NewForm(nil, "users/choose_second_factor.html", "/users"),
		FilterCollection: fc,
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
	}
}

// setValues sets the form values to match the Institution values.
func (f *SecondFactorForm) SetValues() {
	f.Fields["SecondFactor"].Value = f.FilterCollection.ValueOf("SecondFactor")
}
