package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// DepositReportFilterForm is the form that displays filtering options for
// the deposit report page.
type DepositReportFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewDepositReportFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &DepositReportFilterForm{
		Form:             NewForm(nil, "reports/_deposit_filters.html", "/reports/deposits"),
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

func (f *DepositReportFilterForm) init() {
	f.Fields["updated_at__lteq"] = &Field{
		Name:        "updated_at__lteq",
		Label:       "Updated Before",
		Placeholder: "Updated Before",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["storage_option"] = &Field{
		Name:        "storage_option",
		Label:       "Storage Option",
		Placeholder: "Storage Option",
		Options:     Options(constants.StorageOptions),
	}
	f.Fields["chart_metric"] = &Field{
		Name:        "chart_metric",
		Label:       "Chart Metric",
		Placeholder: "Chart Metric",
		Options:     DepositChartMetrics,
	}
}

// setValues sets the form values to match the Institution values.
func (f *DepositReportFilterForm) SetValues() {
	for _, fieldName := range pgmodels.DepositStatsFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for %s", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
