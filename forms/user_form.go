package forms

import (
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserForm struct {
	Form
	User        *models.User
	instOptions []ListOption
}

func NewUserForm(ds *models.DataStore, user *models.User) (*UserForm, error) {
	var err error
	userForm := &UserForm{
		Form: NewForm(ds),
		User: user,
	}
	userForm.instOptions, err = ListInstitutions(ds, false)
	if err != nil {
		return nil, err
	}
	userForm.init()
	userForm.setValues()
	return userForm, err
}

func (f *UserForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		ErrMsg:      "Name must contain at least two letters.",
		Label:       "Name",
		Placeholder: "Name",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Email"] = &Field{
		Name:        "Email",
		ErrMsg:      "Valid email address required.",
		Label:       "Email Address",
		Placeholder: "Email Address",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["PhoneNumber"] = &Field{
		Name:        "PhoneNumber",
		ErrMsg:      "Please enter a phone number in format +2125551212.",
		Label:       "Phone",
		Placeholder: "Phone in format +2125551212",
		Attrs: map[string]string{
			"pattern": "\\+[0-9]{10}",
		},
	}
	f.Fields["OTPRequiredForLogin"] = &Field{
		Name:    "OTPRequiredForLogin",
		Label:   "Require Two-Factor Auth",
		ErrMsg:  "Please choose yes or no.",
		Options: YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["GracePeriod"] = &Field{
		Name:        "GracePeriod",
		Label:       "Must enable two-factor auth by",
		Placeholder: "mm/dd/yyyy",
		Attrs: map[string]string{
			"min": "2021-01-01",
			"max": "2099-12-31",
		},
	}
	f.Fields["InstitutionID"] = &Field{
		Name:    "InstitutionID",
		Label:   "Institution",
		Options: f.instOptions,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Role"] = &Field{
		Name:    "Role",
		ErrMsg:  "Please choose a role for this user.",
		Label:   "Role",
		Options: RolesList,
		Attrs: map[string]string{
			"required": "",
		},
	}
}

func (f *UserForm) Bind(c *gin.Context) error {
	err := c.ShouldBind(f.User)
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

// setValues sets the form values to match the User values.
func (f *UserForm) setValues() {
	f.Fields["Name"].Value = f.User.Name
	f.Fields["Email"].Value = f.User.Email
	f.Fields["PhoneNumber"].Value = f.User.PhoneNumber
	f.Fields["OTPRequiredForLogin"].Value = f.User.OTPRequiredForLogin
	f.Fields["InstitutionID"].Value = f.User.InstitutionID
	f.Fields["Role"].Value = f.User.Role

	// Don't set date to 0001-01-01, because it makes
	// the date picker hard to use. User has to scroll
	// ahead 2000 years.
	if !f.User.GracePeriod.IsZero() {
		f.Fields["GracePeriod"].Value = f.User.GracePeriod.Format("2006-01-02")
	}
}