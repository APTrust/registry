package controllers

import (
	"strconv"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
)

type ListOption struct {
	Value string
	Text  string
}

// RolesList is a list of assignable user roles. Hard-coded instead
// of using Options() function for formatting reasons and because we
// don't want to include the "none" role.
var RolesList = []ListOption{
	{constants.RoleInstAdmin, "Institutional Admin"},
	{constants.RoleInstUser, "Institutional User"},
	{constants.RoleSysAdmin, "APTrust System Administrator"},
}

func ListInstitutions(ds *models.DataStore) ([]ListOption, error) {
	instQuery := models.NewQuery().Columns("id", "name").OrderBy("name asc").Limit(100).Offset(0)
	institutions, err := ds.InstitutionList(instQuery)
	if err != nil {
		return nil, err
	}
	options := make([]ListOption, len(institutions))
	for i, inst := range institutions {
		options[i] = ListOption{strconv.FormatInt(inst.ID, 10), inst.Name}
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
