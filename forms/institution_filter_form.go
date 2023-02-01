package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
)

// InstitutionFilterForm is the form that displays filtering options for
// the institution list page.
type InstitutionFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []*ListOption
}

func NewInstitutionFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &InstitutionFilterForm{
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

func (f *InstitutionFilterForm) init() {
	f.Fields["name__contains"] = &Field{
		Name:        "name__contains",
		Label:       "Name Contains",
		Placeholder: "Name Contains",
	}
	f.Fields["type"] = &Field{
		Name:        "type",
		Label:       "Type",
		Placeholder: "Type",
		Options:     InstTypeList,
	}
}

// setValues sets the form values to match the Institution values.
func (f *InstitutionFilterForm) SetValues() {
	for _, fieldName := range pgmodels.InstitutionFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for %s", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
