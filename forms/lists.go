package forms

import (
	"fmt"
	"strconv"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

var Months = []string{
	"",
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

// AllRolesList is a list of assignable user roles. Hard-coded instead
// of using Options() function for formatting reasons and because we
// don't want to include the "none" role.
var AllRolesList = []*ListOption{
	{constants.RoleInstAdmin, "Institutional Admin", false},
	{constants.RoleInstUser, "Institutional User", false},
	{constants.RoleSysAdmin, "APTrust System Administrator", false},
}

var BagItProfileIdentifiers = []*ListOption{
	{constants.DefaultProfileIdentifier, "APTrust", false},
	{constants.BTRProfileIdentifier, "BTR", false},
}

// DepositChartMetrics come from names of pgmodels.DepositStats properties
var DepositChartMetrics = []*ListOption{
	{"object_count", "Object Count", false},
	{"file_count", "File Count", false},
	{"total_bytes", "Total Bytes", false},
	{"total_gb", "Total Gigabytes", false},
	{"total_tb", "Total Terabytes", false},
	// {"monthly_cost", "Monthly Cost", false}, -- Temporarily removed until we verify costs are correct. https://trello.com/c/pTHtYtK4
}

// InstRolesList is a list of user roles for institutions.
var InstRolesList = []*ListOption{
	{constants.RoleInstAdmin, "Institutional Admin", false},
	{constants.RoleInstUser, "Institutional User", false},
}

var InstTypeList = []*ListOption{
	{constants.InstTypeMember, "Member", false},
	{constants.InstTypeSubscriber, "Associate", false},
}

var ObjectStateList = []*ListOption{
	{constants.StateActive, "Active", false},
	{constants.StateDeleted, "Deleted", false},
}

var StorageOptionList = []*ListOption{
	{constants.StorageOptionGlacierDeepOH, "Glacier Deep - Ohio", false},
	{constants.StorageOptionGlacierDeepOR, "Glacier Deep - Oregon", false},
	{constants.StorageOptionGlacierDeepVA, "Glacier Deep - Virginia", false},
	{constants.StorageOptionGlacierOH, "Glacier - Ohio", false},
	{constants.StorageOptionGlacierOR, "Glacier - Oregon", false},
	{constants.StorageOptionGlacierVA, "Glacier - Virginia", false},
	{constants.StorageOptionStandard, "Standard", false},
	{constants.StorageOptionWasabiOR, "Wasabi - Oregon", false},
	{constants.StorageOptionWasabiTX, "Wasabi - Texas", false},
	{constants.StorageOptionWasabiVA, "Wasabi - Virginia", false},
}

var TwoFactorMethodList = []*ListOption{
	{constants.TwoFactorNone, "None (Turn Off Two-Factor Authentication)", false},
	{constants.TwoFactorTOTP, "Authenticator App (QR Code)", false},
	{constants.TwoFactorAuthy, "Authy OneTouch", false},
	{constants.TwoFactorSMS, "Text Message", false},
}

var YesNoList = []*ListOption{
	{"true", "Yes", false},
	{"false", "No", false},
}

func ListInstitutions(membersOnly bool) ([]*ListOption, error) {
	instQuery := pgmodels.NewQuery().Columns("id", "name").OrderBy("name", "asc").Limit(100).Offset(0)
	if membersOnly {
		instQuery.Where("type", "=", constants.InstTypeMember)
	}
	institutions, err := pgmodels.InstitutionSelect(instQuery)
	if err != nil {
		return nil, err
	}
	options := make([]*ListOption, len(institutions))
	for i, inst := range institutions {
		options[i] = &ListOption{strconv.FormatInt(inst.ID, 10), inst.Name, false}
	}
	return options, nil
}

func ListUsers(institutionID int64) ([]*ListOption, error) {
	query := pgmodels.NewQuery().Columns("id", "name").OrderBy("name", "asc").Limit(200).Offset(0)
	if institutionID > 0 {
		query.Where("institution_id", "=", institutionID)
	}
	users, err := pgmodels.UserSelect(query)
	if err != nil {
		return nil, err
	}
	options := make([]*ListOption, len(users))
	for i, user := range users {
		options[i] = &ListOption{strconv.FormatInt(user.ID, 10), user.Name, false}
	}
	return options, nil
}

// ListDepositReportDates returns a list of dates for deposit reports.
// Note that for each option except "Today", the label is a month and
// year and the value is the first day of the following month. For example,
// if the label is August 2022, the value is 2022-09-01. This would give
// you a report for deposits through the end of August, 2022. In the
// database, the query amounts to "all files created BEFORE 2022-09-01".
//
// Note that past deposit stats are stored in the historical_deposit_stats
// table, and we store month-end reports only. Thus, ALL of the dates in
// that table will be first-of-month dates. Stats go back to December 2014,
// which was APTrust's initial launch (though the system was empty then),
//
// If user chooses the "Today" option, the data will come from the
// current_deposit_stats view, which is updated hourly.
func ListDepositReportDates(includeToday bool) []*ListOption {
	now := time.Now().UTC()
	options := make([]*ListOption, 0)
	if includeToday {
		options = append(options, &ListOption{now.Format("2006-01-02"), "Today", false})
	}
	thisYear, thisMonth, _ := now.Date()
	for year := thisYear; year > 2014; year-- {
		for month := int(time.December); month > 0; month-- {
			if year == thisYear && month > int(thisMonth) {
				continue
			}
			displayYear := year
			displayMonth := Months[month]
			displayDate := fmt.Sprintf("%s %d, %d", displayMonth, 1, displayYear)
			dateValue := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			options = append(options, &ListOption{dateValue.Format("2006-01-02"), displayDate, false})
		}
	}
	return options
}

// Options returns a list of options for the given string list.
// This is intended mainly to provide select list filters
// for the web ui for constants such as:
//
// AccessSettings
// DigestAlgs
// EventTypes
// Stages
// Statuses
// StorageOptions
// WorkItemActions
func Options(items []string) []*ListOption {
	options := make([]*ListOption, len(items))
	for i, item := range items {
		options[i] = &ListOption{item, item, false}
	}
	return options
}
