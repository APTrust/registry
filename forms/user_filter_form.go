package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
)

// UserFilterForm is the form that displays filtering options for
// the user list page.
type UserFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewUserFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &UserFilterForm{
		Form:             NewForm(nil, "institutions/_filters.html", "/institutions"),
		FilterCollection: fc,
	}
	var err error
	if actingUser.IsAdmin() {
		// SysAdmin can view alerts for all institutions and users.
		f.instOptions, err = ListInstitutions(false)
		if err != nil {
			return nil, err
		}
	}
	f.init()
	f.SetValues()
	return f, nil
}

func (f *UserFilterForm) init() {
	f.Fields["email__contains"] = &Field{
		Name:        "email__contains",
		Label:       "Email Contains",
		Placeholder: "Email Contains",
	}
	f.Fields["name__contains"] = &Field{
		Name:        "name__contains",
		Label:       "Name Contains",
		Placeholder: "Name Contains",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["role"] = &Field{
		Name:        "role",
		Label:       "Role",
		Placeholder: "Role",
		Options:     AllRolesList,
	}
}

// setValues sets the form values to match the User values.
func (f *UserFilterForm) SetValues() {
	for _, fieldName := range pgmodels.UserFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for %s", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
