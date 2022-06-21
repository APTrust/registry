package pgmodels

import (
	"fmt"
	"reflect"

	"github.com/stretchr/stew/slice"
)

var PremisEventCountFilters = []string{
	"institution_id",
	"event_type",
	"outcome",
}

type PremisEventCount struct {
	tableName     struct{} `pg:"premis_event_counts"`
	InstitutionID int64    `json:"institution_id"`
	RowCount      int64    `json:"row_count"`
	EventType     string   `json:"event_type"`
	Outcome       string   `json:"outcome"`
}

var IntellectualObjectCountFilters = []string{
	"institution_id",
	"state",
}

type IntellectualObjectCount struct {
	tableName     struct{} `pg:"intellectual_object_counts"`
	InstitutionID int64    `json:"institution_id"`
	RowCount      int64    `json:"row_count"`
	State         string   `json:"state"`
}

var GenericFileCountFilters = []string{
	"institution_id",
	"state",
}

type GenericFileCount struct {
	tableName     struct{} `pg:"generic_file_counts"`
	InstitutionID int64    `json:"institution_id"`
	RowCount      int64    `json:"row_count"`
	State         string   `json:"state"`
}

var WorkItemCountFilters = []string{
	"institution_id",
	"action",
}

type WorkItemCount struct {
	tableName     struct{} `pg:"work_item_counts"`
	InstitutionID int64    `json:"institution_id"`
	RowCount      int64    `json:"row_count"`
	Action        string   `json:"action"`
}

func GetCountFromView(query *Query, model interface{}) (int64, error) {
	// Get a copy of the query, minus order by, limit, offset, and relations.
	// We want just the where clause.
	copyOfQuery := query.CopyForCount()
	whereClauseCols := query.GetColumnsInWhereClause()
	for _, col := range PremisEventCountFilters {
		if !slice.Contains(whereClauseCols, col) {
			// If a view column is not specified,
			// set it to null to get the rollup value.
			copyOfQuery.IsNull(col)
		}
	}
	var rowCount int64
	var err error
	typeName := reflect.TypeOf(model).Name()
	switch typeName {
	case "GenericFile", "GenericFileView":
		obj := GenericFileCount{}
		err = copyOfQuery.Columns("row_count").Select(&obj)
		rowCount = obj.RowCount
	case "IntellectualObject", "IntellectualObjectView":
		obj := IntellectualObjectCount{}
		err = copyOfQuery.Columns("row_count").Select(&obj)
		rowCount = obj.RowCount
	case "PremisEvent", "PremisEventView":
		obj := PremisEventCount{}
		err = copyOfQuery.Columns("row_count").Select(&obj)
		rowCount = obj.RowCount
	case "WorkItem", "WorkItemView":
		obj := WorkItemCount{}
		err = copyOfQuery.Columns("row_count").Select(&obj)
		rowCount = obj.RowCount
	default:
		err = fmt.Errorf("type not supported for view count")
	}
	return rowCount, err
}

func CanCountFromView(query *Query, model interface{}) bool {
	typeName := reflect.TypeOf(model).Name()
	var allowedFilters []string
	switch typeName {
	case "GenericFile", "GenericFileView":
		allowedFilters = GenericFileCountFilters
	case "IntellectualObject", "IntellectualObjectView":
		allowedFilters = IntellectualObjectCountFilters
	case "PremisEvent", "PremisEventView":
		allowedFilters = PremisEventCountFilters
	case "WorkItem", "WorkItemView":
		allowedFilters = WorkItemCountFilters
	}
	if len(allowedFilters) == 0 {
		return false
	}
	for _, col := range query.GetColumnsInWhereClause() {
		if !slice.Contains(allowedFilters, col) {
			return false
		}
	}
	return true
}
