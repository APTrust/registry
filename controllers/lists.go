package controllers

import (
	"strconv"

	"github.com/APTrust/registry/models"
)

type ListOption struct {
	Value string
	Text  string
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
