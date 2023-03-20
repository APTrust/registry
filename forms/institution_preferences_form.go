package forms

import (
	"fmt"

	"github.com/APTrust/registry/pgmodels"
)

// InstitutionPreferencesForm allows institutional admins to edit
// a subset of their institution's info. This includes whether to
// require two-factor authentication and how often to run spot tests.
type InstitutionPreferencesForm struct {
	Form
}

func NewInstitutionPreferencesForm(institution *pgmodels.Institution) (*InstitutionPreferencesForm, error) {
	instPrefsForm := &InstitutionPreferencesForm{
		Form: NewForm(institution, "institutions/preferences_form.html", "/institutions"),
	}
	instPrefsForm.init()
	instPrefsForm.SetValues()
	return instPrefsForm, nil
}

// Action returns the html form.action attribute for this form.
func (f *InstitutionPreferencesForm) Action() string {
	return fmt.Sprintf("%s/edit_preferences/%d", f.BaseURL, f.Model.GetID())
}

func (f *InstitutionPreferencesForm) init() {
	f.Fields["OTPEnabled"] = &Field{
		Name:        "OTPEnabled",
		Label:       "Enable two-factor authentication?",
		Placeholder: "Two-Factor Auth Required?",
		ErrMsg:      "Please choose yes or no.",
		Options:     YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["SpotRestoreFrequency"] = &Field{
		Name:        "SpotRestoreFrequency",
		Label:       "Restoration spot test frequency (days)",
		Placeholder: "",
		ErrMsg:      "Please indicate how often to run spot restoration tests. (E.g. 30, 60, 90 days. Use zero to indicate never.)",
		Attrs: map[string]string{
			"required": "",
		},
	}
}

// setValues sets the form values to match the Institution values.
func (f *InstitutionPreferencesForm) SetValues() {
	institution := f.Model.(*pgmodels.Institution)
	f.Fields["OTPEnabled"].Value = institution.OTPEnabled
	f.Fields["SpotRestoreFrequency"].Value = institution.SpotRestoreFrequency
}
