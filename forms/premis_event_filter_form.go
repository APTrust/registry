package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// PremisEventFilterForm is the form that displays filtering options for
// the PREMIS event list page.
type PremisEventFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewPremisEventFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &PremisEventFilterForm{
		Form:             NewForm(nil, "events/_filters.html", "/events"),
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

func (f *PremisEventFilterForm) init() {
	f.Fields["date_time__gteq"] = &Field{
		Name:        "date_time__gteq",
		Label:       "Date on or After",
		Placeholder: "Date on or After",
	}
	f.Fields["date_time__lteq"] = &Field{
		Name:        "date_time__lteq",
		Label:       "Date on or Before",
		Placeholder: "Date on or Before",
	}
	f.Fields["event_type"] = &Field{
		Name:        "event_type",
		Label:       "Event Type",
		Placeholder: "Event Type",
		Options:     Options(constants.EventTypes),
	}
	f.Fields["generic_file_identifier"] = &Field{
		Name:        "generic_file_identifier",
		Label:       "File Identifier",
		Placeholder: "File Identifier",
	}
	f.Fields["identifier"] = &Field{
		Name:        "identifier",
		Label:       "Identifier (UUID)",
		Placeholder: "Identifier (UUID)",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["intellectual_object_identifier"] = &Field{
		Name:        "intellectual_object_identifier",
		Label:       "Object Identifier",
		Placeholder: "Object Identifier",
	}
	f.Fields["outcome"] = &Field{
		Name:        "outcome",
		Label:       "Outcome",
		Placeholder: "Outcome",
		Options:     Options(constants.EventOutcomes),
	}
}

// setValues sets the form values to match the Institution values.
func (f *PremisEventFilterForm) SetValues() {
	for _, fieldName := range pgmodels.PremisEventFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for %s", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
