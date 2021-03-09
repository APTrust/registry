package forms

import (
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InstitutionForm struct {
	Form
	Institution *models.Institution
	instOptions []ListOption
}

func NewInstitutionForm(ds *models.DataStore, institution *models.Institution) (*InstitutionForm, error) {
	var err error
	institutionForm := &InstitutionForm{
		Form:        NewForm(ds),
		Institution: institution,
	}

	// List parent (member) institutions only.
	institutionForm.instOptions, err = ListInstitutions(ds, true)
	if err != nil {
		return nil, err
	}
	institutionForm.init()
	institutionForm.setValues()
	return institutionForm, err
}

func (f *InstitutionForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		Label:       "Name",
		Placeholder: "Name",
		ErrMsg:      "Name must have at least two letters.",
		Attrs: map[string]string{
			"required": "",
			"min":      "2",
		},
	}
	f.Fields["Identifier"] = &Field{
		Name:        "Identifier",
		Label:       "Identifier",
		Placeholder: "Identifier",
		ErrMsg:      "Identifier must be a domain name.",
		Attrs: map[string]string{
			"required": "",
			"pattern":  "[A-Za-z0-9]{2,}\\.[A-Za-z0-9]{2,}",
		},
	}
	f.Fields["Type"] = &Field{
		Name:        "Type",
		Label:       "Institution Type",
		Placeholder: "Institution Type",
		ErrMsg:      "Please choose a type.",
		Options:     InstTypeList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["MemberInstitutionID"] = &Field{
		Name:        "Parent Institution",
		Label:       "Parent Institution",
		Placeholder: "Parent Institution",
		ErrMsg:      "You must choose a parent instition if this is a sub-account.",
		Options:     f.instOptions,
	}
	f.Fields["OTPEnabled"] = &Field{
		Name:        "Two-Factor Auth Required?",
		Label:       "Two-Factor Auth Required?",
		Placeholder: "Two-Factor Auth Required?",
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

// TODO: Can this be part of the underlying Form class, and not
// repeated in each form?
func (f *InstitutionForm) Bind(c *gin.Context) error {
	err := c.ShouldBind(f.Institution)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				f.Fields[fieldErr.Field()].DisplayError = true
			}
		}
	}
	f.setValues()
	return err
}

// setValues sets the form values to match the Institution values.
func (f *InstitutionForm) setValues() {
	f.Fields["Name"].Value = f.Institution.Name
	f.Fields["Identifier"].Value = f.Institution.Identifier
	f.Fields["Type"].Value = f.Institution.Type
	f.Fields["MemberInstitutionID"].Value = f.Institution.MemberInstitutionID
	f.Fields["OTPEnabled"].Value = f.Institution.OTPEnabled
	f.Fields["ReceivingBucket"].Value = f.Institution.ReceivingBucket
	f.Fields["RestoreBucket"].Value = f.Institution.RestoreBucket
}
