package forms

import (
	//"net/http"

	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	//"github.com/gin-gonic/gin"
)

type UserForm struct {
	Form
	User        *pgmodels.User
	instOptions []ListOption
}

// func NewUserForm(user *pgmodels.User) (*UserForm, error) {
// 	var err error
// 	userForm := &UserForm{
// 		Form: NewForm(),
// 		User: user,
// 	}
// 	userForm.instOptions, err = ListInstitutions(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	userForm.init()
// 	userForm.setValues()
// 	return userForm, err
// }

// func (f *UserForm) init() {
// 	f.Fields["Name"] = &Field{
// 		Name:        "Name",
// 		ErrMsg:      pgmodels.ErrUserName,
// 		Label:       "Name",
// 		Placeholder: "Name",
// 		Attrs: map[string]string{
// 			"required": "",
// 		},
// 	}
// 	f.Fields["Email"] = &Field{
// 		Name:        "Email",
// 		ErrMsg:      pgmodels.ErrUserEmail,
// 		Label:       "Email Address",
// 		Placeholder: "Email Address",
// 		Attrs: map[string]string{
// 			"required": "",
// 		},
// 	}
// 	f.Fields["PhoneNumber"] = &Field{
// 		Name:        "PhoneNumber",
// 		ErrMsg:      pgmodels.ErrUserPhone,
// 		Label:       "Phone",
// 		Placeholder: "Phone in format 212-555-1212",
// 		Attrs:       map[string]string{
// 			//"pattern": "[0-9]{10,11}",
// 		},
// 	}
// 	f.Fields["OTPRequiredForLogin"] = &Field{
// 		Name:    "OTPRequiredForLogin",
// 		Label:   "Require Two-Factor Auth",
// 		ErrMsg:  pgmodels.ErrUser2Factor,
// 		Options: YesNoList,
// 		Attrs: map[string]string{
// 			"required": "",
// 		},
// 	}
// 	f.Fields["GracePeriod"] = &Field{
// 		Name:        "GracePeriod",
// 		Label:       "Must enable two-factor auth by",
// 		Placeholder: "mm/dd/yyyy",
// 		ErrMsg:      pgmodels.ErrUserGracePeriod,
// 		Attrs: map[string]string{
// 			"min": "2021-01-01",
// 			"max": "2099-12-31",
// 		},
// 	}
// 	f.Fields["InstitutionID"] = &Field{
// 		Name:    "InstitutionID",
// 		Label:   "Institution",
// 		ErrMsg:  pgmodels.ErrUserInst,
// 		Options: f.instOptions,
// 		Attrs: map[string]string{
// 			"required": "",
// 		},
// 	}
// 	f.Fields["Role"] = &Field{
// 		Name:    "Role",
// 		ErrMsg:  pgmodels.ErrUserRole,
// 		Label:   "Role",
// 		Options: RolesList,
// 		Attrs: map[string]string{
// 			"required": "",
// 		},
// 	}
// }

// // func (f *UserForm) Save(c *gin.Context, templateData gin.H) (int, error) {
// // 	status := http.StatusCreated
// // 	if f.User.ID > 0 {
// // 		status = http.StatusOK
// // 	}
// // 	_ = c.ShouldBind(f.User)
// // 	err := f.User.Save()
// // 	if err != nil {
// // 		status = f.handleError(err, templateData)
// // 	}
// // 	f.setValues()
// // 	return status, err
// // }

// // func (f *UserForm) handleError(err error, templateData gin.H) int {
// // 	status := http.StatusBadRequest
// // 	if valErr, ok := err.(*common.ValidationError); ok {
// // 		for fieldName, _ := range valErr.Errors {
// // 			f.Fields[fieldName].DisplayError = true
// // 		}
// // 	} else {
// // 		templateData["FormError"] = err.Error()
// // 		status = http.StatusInternalServerError
// // 	}
// // 	return status
// // }

// // setValues sets the form values to match the User values.
// func (f *UserForm) setValues() {
// 	f.Fields["Name"].Value = f.User.Name
// 	f.Fields["Email"].Value = f.User.Email
// 	f.Fields["PhoneNumber"].Value = f.User.PhoneNumber
// 	f.Fields["OTPRequiredForLogin"].Value = f.User.OTPRequiredForLogin
// 	f.Fields["InstitutionID"].Value = f.User.InstitutionID
// 	f.Fields["Role"].Value = f.User.Role

// 	// Don't set date to 0001-01-01, because it makes
// 	// the date picker hard to use. User has to scroll
// 	// ahead 2000 years.
// 	if !f.User.GracePeriod.IsZero() {
// 		f.Fields["GracePeriod"].Value = f.User.GracePeriod.Format("2006-01-02")
// 	}
// }
