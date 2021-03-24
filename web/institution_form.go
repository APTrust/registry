package web

import (
	"fmt"
	//"net/http"

	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
)

type InstitutionForm struct {
	Form
	instOptions []ListOption
}

func NewInstitutionForm(request *Request) (*InstitutionForm, error) {
	var err error
	institution := &pgmodels.Institution{}
	if request.ResourceID > 0 {
		institution, err = pgmodels.InstitutionByID(request.ResourceID)
		if err != nil {
			return nil, err
		}
	}
	institutionForm := &InstitutionForm{
		Form: NewForm(request, institution),
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

// func (f *InstitutionForm) Save(c *gin.Context, templateData gin.H) (int, error) {
// 	status := http.StatusCreated
// 	if f.Institution.ID > 0 {
// 		status = http.StatusOK
// 	}
// 	_ = c.ShouldBind(f.Institution)
// 	err := f.Institution.Save()
// 	if err != nil {
// 		status = f.handleError(err, templateData)
// 	}
// 	f.setValues()
// 	return status, err
// }

// func (f *InstitutionForm) handleError(err error, templateData gin.H) int {
// 	status := http.StatusBadRequest
// 	if valErr, ok := err.(*common.ValidationError); ok {
// 		for fieldName, _ := range valErr.Errors {
// 			f.Fields[fieldName].DisplayError = true
// 		}
// 	} else {
// 		templateData["FormError"] = err.Error()
// 		status = http.StatusInternalServerError
// 	}
// 	return status
// }

// setValues sets the form values to match the Institution values.
func (f *InstitutionForm) setValues() {
	fmt.Println("------ CALLED INST setValues() ------")
	institution := f.Model.(*pgmodels.Institution)
	f.Fields["Name"].Value = institution.Name
	f.Fields["Identifier"].Value = institution.Identifier
	f.Fields["Type"].Value = institution.Type
	f.Fields["MemberInstitutionID"].Value = institution.MemberInstitutionID
	f.Fields["OTPEnabled"].Value = institution.OTPEnabled
	f.Fields["ReceivingBucket"].Value = institution.ReceivingBucket
	f.Fields["RestoreBucket"].Value = institution.RestoreBucket
}
