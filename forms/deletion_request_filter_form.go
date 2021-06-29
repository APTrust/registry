package forms

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// DeletionRequestFilterForm is the form that displays filtering options for
// the deletion request list page.
type DeletionRequestFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewDeletionRequestFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (*DeletionRequestFilterForm, error) {
	f := &DeletionRequestFilterForm{
		Form:             NewForm(nil, "deletions/_filters.html", "/deletions"),
		FilterCollection: fc,
	}
	var err error
	if actingUser.IsAdmin() {
		// SysAdmin can view files at any institutions.
		f.instOptions, err = ListInstitutions(false)
		if err != nil {
			return nil, err
		}
	}
	f.init()
	f.SetValues()
	return f, nil
}

func (f *DeletionRequestFilterForm) init() {
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["requested_at__lteq"] = &Field{
		Name:        "requested_at__lteq",
		Label:       "Requested On or Before",
		Placeholder: "Requested On or Before",
	}
	f.Fields["requested_at__gteq"] = &Field{
		Name:        "requested_at__gteq",
		Label:       "Requested On or After",
		Placeholder: "Requested On or After",
	}
	f.Fields["stage"] = &Field{
		Name:        "stage",
		Label:       "Work Item Stage",
		Placeholder: "Work Item Stage",
		Options:     Options(constants.Stages),
	}
	f.Fields["status"] = &Field{
		Name:        "status",
		Label:       "Work Item Status",
		Placeholder: "Work Item Status",
		Options:     Options(constants.Statuses),
	}
}

// setValues sets the form values to match the Institution values.
func (f *DeletionRequestFilterForm) SetValues() {
	f.Fields["institution_id"].Value = f.FilterCollection.ValueOf("institution_id")
	f.Fields["stage"].Value = f.FilterCollection.ValueOf("stage")
	f.Fields["status"].Value = f.FilterCollection.ValueOf("status")
	f.Fields["requested_at__gteq"].Value = f.FilterCollection.ValueOf("requested_at__gteq")
	f.Fields["requested_at__lteq"].Value = f.FilterCollection.ValueOf("requested_at__lteq")
}
