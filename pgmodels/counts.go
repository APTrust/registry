package pgmodels

import (
	"fmt"
	"reflect"

	"github.com/APTrust/registry/common"
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
	typeName, allowedFilters, err := typeNameAndFilterColumns(model)
	if err != nil {
		return -1, err
	}

	// Get a copy of the query, minus order by, limit, offset, and relations.
	// We want just the where clause.
	copyOfQuery := query.CopyForCount()
	whereClauseCols := copyOfQuery.GetColumnsInWhereClause()
	for _, col := range allowedFilters {
		if !slice.Contains(whereClauseCols, col) {
			// If a view column is not specified,
			// set it to null to get the rollup value.
			copyOfQuery.IsNull(col)
		}
	}
	var rowCount int64
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

	// NoRowError means our query was valid, but there were
	// no results. This is a legitimate case indicating a
	// count of zero. Our select with rollup in the DB's
	// update_counts() function does not return a row where
	// counts are zero. For example, some depositors may have
	// zero WorkItems where action="Delete", or zero events
	// where event_type="Deletion" and outcome="Failed". For
	// these, we want to return a zero count and no error.
	if IsNoRowError(err) {
		rowCount = 0
		err = nil
	}

	return rowCount, err
}

func CanCountFromView(query *Query, model interface{}) bool {
	_, allowedFilters, err := typeNameAndFilterColumns(model)
	if err != nil {
		return false
	}

	for _, col := range query.GetColumnsInWhereClause() {
		if !slice.Contains(allowedFilters, col) {
			return false
		}
	}
	return true
}

func typeNameAndFilterColumns(model interface{}) (string, []string, error) {
	typeName := reflect.TypeOf(model).Name()
	if typeName == "" {
		typeName = reflect.ValueOf(model).Elem().Type().Name()
	}
	var allowedFilters []string
	var err error
	switch typeName {
	case "GenericFile", "GenericFileView":
		allowedFilters = GenericFileCountFilters
	case "IntellectualObject", "IntellectualObjectView":
		allowedFilters = IntellectualObjectCountFilters
	case "PremisEvent", "PremisEventView":
		allowedFilters = PremisEventCountFilters
	case "WorkItem", "WorkItemView":
		allowedFilters = WorkItemCountFilters
	default:
		err = common.ErrCountTypeNotSupported
	}
	return typeName, allowedFilters, err
}
