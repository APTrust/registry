package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// AlertFilterForm is the form that displays filtering options for
// the alert list page.
type AlertFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
	userOptions       []ListOption
}

func NewAlertFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &AlertFilterForm{
		Form:             NewForm(nil, "alerts/_filters.html", "/alerts"),
		FilterCollection: fc,
	}
	var err error
	if actingUser.IsAdmin() {
		// SysAdmin can view alerts for all institutions and users.
		f.instOptions, err = ListInstitutions(false)
		if err != nil {
			return nil, err
		}
		f.userOptions, err = ListUsers(0)
		if err != nil {
			return nil, err
		}
	}
	f.init()
	f.SetValues()
	return f, nil
}

func (f *AlertFilterForm) init() {
	f.Fields["created_at__gteq"] = &Field{
		Name:        "created_at__gteq",
		Label:       "Created On or After",
		Placeholder: "Created On or After",
	}
	f.Fields["created_at__lteq"] = &Field{
		Name:        "created_at__lteq",
		Label:       "Created On or Before",
		Placeholder: "Created On or Before",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["type"] = &Field{
		Name:        "type",
		Label:       "Alert Type",
		Placeholder: "Alert Type",
		Options:     Options(constants.AlertTypes),
	}
	f.Fields["user_id"] = &Field{
		Name:        "user_id",
		Label:       "Recipient",
		Placeholder: "Recipient",
		Options:     f.userOptions,
	}
}

// setValues sets the form values to match the Institution values.
func (f *AlertFilterForm) SetValues() {
	for _, fieldName := range pgmodels.AlertFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
