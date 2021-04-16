package web

import (
	"strconv"

	"github.com/APTrust/registry/pgmodels"
)

type UserForm struct {
	Form
	User              *pgmodels.User
	instOptions       []ListOption
	actingUserIsAdmin bool
}

func NewUserForm(userToEdit *pgmodels.User, actingUser *pgmodels.User) (*UserForm, error) {
	userForm := &UserForm{
		Form:              NewForm(),
		User:              userToEdit,
		actingUserIsAdmin: actingUser.IsAdmin(),
	}

	var err error
	if actingUser.IsAdmin() {
		// SysAdmin can create/edit users at any institutions.
		userForm.instOptions, err = ListInstitutions(false)
		if err != nil {
			return nil, err
		}
	} else {
		// Non-sysadmin (inst admin) can add/edit local users only.
		userForm.instOptions = []ListOption{
			{strconv.FormatInt(actingUser.InstitutionID, 10), actingUser.Institution.Name},
		}
		userToEdit.InstitutionID = actingUser.InstitutionID
	}
	userForm.init()
	userForm.SetValues()
	return userForm, nil
}

func (f *UserForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		ErrMsg:      pgmodels.ErrUserName,
		Label:       "Name",
		Placeholder: "Name",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Email"] = &Field{
		Name:        "Email",
		ErrMsg:      pgmodels.ErrUserEmail,
		Label:       "Email Address",
		Placeholder: "Email Address",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["PhoneNumber"] = &Field{
		Name:        "PhoneNumber",
		ErrMsg:      pgmodels.ErrUserPhone,
		Label:       "Phone",
		Placeholder: "Phone in format 212-555-1212",
		Attrs:       map[string]string{
			//"pattern": "[0-9]{10,11}",
		},
	}
	f.Fields["OTPRequiredForLogin"] = &Field{
		Name:    "OTPRequiredForLogin",
		Label:   "Require Two-Factor Auth",
		ErrMsg:  pgmodels.ErrUser2Factor,
		Options: YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["GracePeriod"] = &Field{
		Name:        "GracePeriod",
		Label:       "Must enable two-factor auth by",
		Placeholder: "mm/dd/yyyy",
		ErrMsg:      pgmodels.ErrUserGracePeriod,
		Attrs: map[string]string{
			"min": "2021-01-01",
			"max": "2099-12-31",
		},
	}
	f.Fields["InstitutionID"] = &Field{
		Name:    "InstitutionID",
		Label:   "Institution",
		ErrMsg:  pgmodels.ErrUserInst,
		Options: f.instOptions,
		Attrs: map[string]string{
			"required": "",
		},
	}
	rolesList := InstRolesList
	if f.actingUserIsAdmin {
		rolesList = AllRolesList
	}
	f.Fields["Role"] = &Field{
		Name:    "Role",
		ErrMsg:  pgmodels.ErrUserRole,
		Label:   "Role",
		Options: rolesList,
		Attrs: map[string]string{
			"required": "",
		},
	}
}

// setValues sets the form values to match the User values.
func (f *UserForm) SetValues() {
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
