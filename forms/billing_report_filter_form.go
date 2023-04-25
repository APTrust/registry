package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
)

// BillingReportFilterForm is the form that displays filtering options for
// the billing report page.
type BillingReportFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []*ListOption
}

func NewBillingReportFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &BillingReportFilterForm{
		Form:             NewForm(nil, "reports/_billing_filters.html", "/reports/billing"),
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

func (f *BillingReportFilterForm) init() {
	f.Fields["start_date"] = &Field{
		Name:        "start_date",
		Label:       "Deposits from date",
		Placeholder: "Deposits from",
		Options:     ListDepositReportDates(false),
		Attrs:       make(map[string]string),
	}
	f.Fields["end_date"] = &Field{
		Name:        "end_date",
		Label:       "Deposits up to date",
		Placeholder: "Deposits up to",
		Options:     ListDepositReportDates(false),
		Attrs:       make(map[string]string),
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
		Attrs:       make(map[string]string),
	}
}

// setValues sets the form values to match the Institution values.
func (f *BillingReportFilterForm) SetValues() {
	for _, fieldName := range pgmodels.DepositStatsFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for %s", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
		if fieldName == "end_date" && f.FilterCollection.ValueOf(fieldName) == "" {
			f.Fields[fieldName].Value = f.Fields[fieldName].Options[0].Value
		}
	}
}
