package forms

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type InstitutionForm struct {
	Form
	Institution *pgmodels.Institution
	instOptions []ListOption
}

func NewInstitutionForm(institution *pgmodels.Institution) (*InstitutionForm, error) {
	var err error
	institutionForm := &InstitutionForm{
		Form:        NewForm(),
		Institution: institution,
	}

	// List parent (member) institutions only.
	institutionForm.instOptions, err = ListInstitutions(true)
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

func (f *InstitutionForm) Save(c *gin.Context, templateData gin.H) (int, error) {
	status := http.StatusCreated
	if f.Institution.ID > 0 {
		status = http.StatusOK
	}
	_ = c.ShouldBind(f.Institution)
	err := f.Institution.Save()
	if err != nil {
		status = f.handleError(err, templateData)
	}
	f.setValues()
	return status, err
}

func (f *InstitutionForm) handleError(err error, templateData gin.H) int {
	status := http.StatusBadRequest
	if valErr, ok := err.(*common.ValidationError); ok {
		for fieldName, _ := range valErr.Errors {
			f.Fields[fieldName].DisplayError = true
		}
	} else {
		templateData["FormError"] = err.Error()
		status = http.StatusInternalServerError
	}
	return status
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
