package forms

import (
	"strconv"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// AllRolesList is a list of assignable user roles. Hard-coded instead
// of using Options() function for formatting reasons and because we
// don't want to include the "none" role.
var AllRolesList = []ListOption{
	{constants.RoleInstAdmin, "Institutional Admin"},
	{constants.RoleInstUser, "Institutional User"},
	{constants.RoleSysAdmin, "APTrust System Administrator"},
}

var BagItProfileIdentifiers = []ListOption{
	{constants.DefaultProfileIdentifier, "APTrust"},
	{constants.BTRProfileIdentifier, "BTR"},
}

// DepositChartMetrics come from names of pgmodels.DepositStats properties
var DepositChartMetrics = []ListOption{
	{"object_count", "Object Count"},
	{"file_count", "File Count"},
	{"total_bytes", "Total Bytes"},
	{"total_gb", "Total Gigabytes"},
	{"total_tb", "Total Terabytes"},
	{"monthly_cost", "Monthly Cost"},
}

// InstRolesList is a list of user roles for institutions.
var InstRolesList = []ListOption{
	{constants.RoleInstAdmin, "Institutional Admin"},
	{constants.RoleInstUser, "Institutional User"},
}

var InstTypeList = []ListOption{
	{constants.InstTypeMember, "Member"},
	{constants.InstTypeSubscriber, "Subscriber (Sub-Account)"},
}

var ObjectStateList = []ListOption{
	{constants.StateActive, "Active"},
	{constants.StateDeleted, "Deleted"},
}

var StorageOptionList = []ListOption{
	{constants.StorageOptionGlacierDeepOH, "Glacier Deep - Ohio"},
	{constants.StorageOptionGlacierDeepOR, "Glacier Deep - Oregon"},
	{constants.StorageOptionGlacierDeepVA, "Glacier Deep - Virginia"},
	{constants.StorageOptionGlacierOH, "Glacier - Ohio"},
	{constants.StorageOptionGlacierOR, "Glacier - Oregon"},
	{constants.StorageOptionGlacierVA, "Glacier - Virginia"},
	{constants.StorageOptionStandard, "Standard"},
	{constants.StorageOptionWasabiOR, "Wasabi - Oregon"},
	{constants.StorageOptionWasabiVA, "Wasabi - Virginia"},
}

var YesNoList = []ListOption{
	{"true", "Yes"},
	{"false", "No"},
}

func ListInstitutions(membersOnly bool) ([]ListOption, error) {
	instQuery := pgmodels.NewQuery().Columns("id", "name").OrderBy("name asc").Limit(100).Offset(0)
	if membersOnly {
		instQuery.Where("type", "=", constants.InstTypeMember)
	}
	institutions, err := pgmodels.InstitutionSelect(instQuery)
	if err != nil {
		return nil, err
	}
	options := make([]ListOption, len(institutions))
	for i, inst := range institutions {
		options[i] = ListOption{strconv.FormatInt(inst.ID, 10), inst.Name}
	}
	return options, nil
}

func ListUsers(institutionID int64) ([]ListOption, error) {
	query := pgmodels.NewQuery().Columns("id", "name").OrderBy("name asc").Limit(200).Offset(0)
	if institutionID > 0 {
		query.Where("institution_id", "=", institutionID)
	}
	users, err := pgmodels.UserSelect(query)
	if err != nil {
		return nil, err
	}
	options := make([]ListOption, len(users))
	for i, user := range users {
		options[i] = ListOption{strconv.FormatInt(user.ID, 10), user.Name}
	}
	return options, nil
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
func Options(items []string) []ListOption {
	options := make([]ListOption, len(items))
	for i, item := range items {
		options[i] = ListOption{item, item}
	}
	return options
}
