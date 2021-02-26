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

var UserFormErrors = map[string]string{
	"Name":          "Name must contain at least two letters.",
	"Email":         "Valid email address required.",
	"PhoneNumber":   "Please enter a phone number in format 000-000-0000.",
	"InstitutionID": "Please select an institution.",
	"Role":          "Please choose a role for this user.",
}

func NewUserForm(ds *models.DataStore, user *models.User) (*UserForm, error) {
	var err error
	userForm := &UserForm{
		Form: NewForm(ds),
		User: user,
	}
	userForm.instOptions, err = ListInstitutions(ds)
	if err != nil {
		return nil, err
	}
	userForm.init()
	return userForm, err
}

func (f *UserForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		Placeholder: "Name",
		Value:       f.User.Name,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Email"] = &Field{
		Name:        "Email",
		Placeholder: "Email Address",
		Value:       f.User.Email,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["PhoneNumber"] = &Field{
		Name:        "PhoneNumber",
		Placeholder: "Phone in format 212-555-1212",
		Value:       f.User.PhoneNumber,
		Attrs: map[string]string{
			"pattern": "[0-9]{3}-[0-9]{3}-[0-9]{4}",
		},
	}
	f.Fields["OTPRequiredForLogin"] = &Field{
		Name:  "OTPRequiredForLogin",
		Label: "OTP Required For Login",
		Value: f.User.OTPRequiredForLogin,
	}
	f.Fields["GracePeriod"] = &Field{
		Name:        "GracePeriod",
		Label:       "Must enable two-factor auth by",
		Placeholder: "mm/dd/yyyy",
		Value:       f.User.GracePeriod.Format("2006-01-02"),
		Attrs: map[string]string{
			"min": "2021-01-01",
			"max": "2099-12-31",
		},
	}
	f.Fields["InstitutionID"] = &Field{
		Name:    "InstitutionID",
		Label:   "Institution",
		Value:   f.User.InstitutionID,
		Options: f.instOptions,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Role"] = &Field{
		Name:    "Role",
		Label:   "Role",
		Value:   f.User.Role,
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
				f.Fields[fieldErr.Field()].Error = UserFormErrors[fieldErr.Field()]
			}
		}
	}
	f.setValues()
	return err
}

// setValues sets the form values to match the User values.
// This is done in init(), when we first create the form, but it needs
// to be done again when we parse the user-submitted form. If some of
// the user-supplied values are invalid, we need to present the form
// to them again, with validation error messages. The form should have
// all of their saved inputs. Hence this function.
func (f *UserForm) setValues() {
	f.Fields["Name"].Value = f.User.Name
	f.Fields["Email"].Value = f.User.Email
	f.Fields["PhoneNumber"].Value = f.User.PhoneNumber
	f.Fields["OTPRequiredForLogin"].Value = f.User.OTPRequiredForLogin
	f.Fields["GracePeriod"].Value = f.User.GracePeriod.Format("2006-01-02")
	f.Fields["InstitutionID"].Value = f.User.InstitutionID
	f.Fields["Role"].Value = f.User.Role
}
