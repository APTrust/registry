package forms

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type InstitutionForm struct {
	Form
	instOptions []ListOption
}

func NewInstitutionForm(institution *pgmodels.Institution) (*InstitutionForm, error) {
	institutionForm := &InstitutionForm{
		Form: NewForm(institution, "institutions/form.html", "/institutions"),
	}
	// List parent (member) institutions only.
	var err error
	institutionForm.instOptions, err = ListInstitutions(true)
	if err != nil {
		return nil, err
	}
	institutionForm.init()
	institutionForm.SetValues()
	return institutionForm, err
}

func (f *InstitutionForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		Label:       "Name",
		Placeholder: "Name",
		ErrMsg:      pgmodels.ErrInstName,
		Attrs: map[string]string{
			"required": "",
			"min":      "2",
		},
	}
	f.Fields["Identifier"] = &Field{
		Name:        "Identifier",
		Label:       "Identifier",
		Placeholder: "Identifier",
		ErrMsg:      pgmodels.ErrInstIdentifier,
		Attrs: map[string]string{
			"required": "",
			"pattern":  "[A-Za-z0-9]{2,}\\.[A-Za-z0-9]{2,}",
		},
	}
	f.Fields["Type"] = &Field{
		Name:        "Type",
		Label:       "Institution Type",
		Placeholder: "Institution Type",
		ErrMsg:      pgmodels.ErrInstType,
		Options:     InstTypeList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["MemberInstitutionID"] = &Field{
		Name:        "MemberInstitutionID",
		Label:       "Parent Institution",
		Placeholder: "Parent Institution",
		ErrMsg:      pgmodels.ErrInstMemberID,
		Options:     f.instOptions,
		Attrs:       map[string]string{},
	}
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
	f.Fields["EnableSpotRestore"] = &Field{
		Name:        "EnableSpotRestore",
		Label:       "Enable restoration spot tests?",
		Placeholder: "",
		ErrMsg:      "Please choose yes or no.",
		Options:     YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["ReceivingBucket"] = &Field{
		Name:        "Receiving Bucket",
		Label:       "Receiving Bucket",
		Placeholder: "Receiving Bucket",
		Attrs: map[string]string{
			"disabled": "",
			"readonly": "",
		},
	}
	f.Fields["RestoreBucket"] = &Field{
		Name:        "Restoration Bucket",
		Label:       "Restoration Bucket",
		Placeholder: "Restoration Bucket",
		Attrs: map[string]string{
			"disabled": "",
			"readonly": "",
		},
	}
}

// setValues sets the form values to match the Institution values.
func (f *InstitutionForm) SetValues() {
	institution := f.Model.(*pgmodels.Institution)
	f.Fields["Name"].Value = institution.Name
	f.Fields["Identifier"].Value = institution.Identifier
	f.Fields["Type"].Value = institution.Type
	f.Fields["MemberInstitutionID"].Value = institution.MemberInstitutionID
	f.Fields["OTPEnabled"].Value = institution.OTPEnabled
	f.Fields["EnableSpotRestore"].Value = institution.EnableSpotRestore
	f.Fields["ReceivingBucket"].Value = institution.ReceivingBucket
	f.Fields["RestoreBucket"].Value = institution.RestoreBucket

	if f.Fields["Type"].Value == constants.InstTypeMember {
		f.Fields["MemberInstitutionID"].Attrs["disabled"] = "true"
	}
	if institution.ID > 0 {
		f.Fields["Identifier"].Attrs["readonly"] = "true"
	}
}
