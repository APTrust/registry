package pgmodels

import (
	"fmt"
	"strings"

	"github.com/APTrust/registry/common"
)

// FilterCollection converts query string params such as name__eq=Homer to
// a pgmodels.Query object that allows us to build a SQL where clause as
// we go.
type FilterCollection struct {
	filters []*ParamFilter
	sorts   []*SortParam
}

// NewFilterCollection returns a ParamFileters object.
func NewFilterCollection() *FilterCollection {
	return &FilterCollection{
		filters: make([]*ParamFilter, 0),
		sorts:   make([]*SortParam, 0),
	}
}

// Add adds an item to the filter collection. Param key is a filter
// key from the query string. Param values are the values associated
// with that key. For example:
//
// Key: name__in
// Values: ["Bart", "Lisa", "Maggie"]
//
func (fc *FilterCollection) Add(key string, values []string) error {
	filter, err := NewParamFilter(key, values)
	if err != nil {
		return err
	}
	fc.filters = append(fc.filters, filter)
	return nil
}

// AddOrderBy adds sort columns to the filters.
func (fc *FilterCollection) AddOrderBy(values string) {
	sort := NewSortParam(values)
	fc.sorts = append(fc.sorts, sort)
}

// HasExplicitSorting returns true if this object includes explicit
// sort params that will show up as an "order by" clause in the SQL query.
func (fc *FilterCollection) HasExplicitSorting() bool {
	return len(fc.sorts) > 0
}

// ToQuery returns a query object based on the keys and values passed in.
// The Query's WhereClause() will return the where conditions for the filters
// passed in through Add(), and the Query's Params() method will return the
// params. Conditions and params come back in the order they were added.
func (fc *FilterCollection) ToQuery() (*Query, error) {
	query := NewQuery()
	for _, filter := range fc.filters {
		if common.ListIsEmpty(filter.Values) {
			continue // no need to apply filter
		}
		err := filter.AddToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	for _, sort := range fc.sorts {
		query.OrderBy(sort.Column, sort.Direction)
	}
	return query, nil
}

// ValueOf returns the value of the filter with the specified name.
// Returns an empty string if the specified filter is missing or
// has no value.
func (fc *FilterCollection) ValueOf(filterName string) string {
	for _, pf := range fc.filters {
		if pf.Key == filterName && len(pf.Values) > 0 {
			return pf.Values[0]
		}
	}
	return ""
}

// ParamFilter parses query string params into filters that can be added
// to a SQL where clause.
type ParamFilter struct {
	// Key is the name of the query string param.
	Key string
	// Column is derived from Key, and is the name of a database column
	Column string
	// RawOp is the operator in Key. For example, "eq" or "gt".
	RawOp string
	// SQLOp is the SQL operator that corresponds to RawOp. For example,
	// "=" or ">".
	SQLOp string
	// Values are the values attached to Key in the query string.
	Values []string
}

type SortParam struct {
	Column    string
	Direction string // only asc and desc allowed
}

// NewParamFilter returns a new ParamFilter object based on the key and
// values submitted in the query string. It will return a custom error
// if it can't parse the key or values. The caller should log the error
// and return a basic common.ErrInvalidParam.
func NewParamFilter(key string, values []string) (*ParamFilter, error) {
	// Parse the column and operator from the key name
	colAndOp := strings.Split(key, "__")
	if len(colAndOp) != 2 {
		colAndOp = append(colAndOp, "eq") // Assume equals if op not supplied.
	}
	col := colAndOp[0]
	rawOp := colAndOp[1]
	sqlOp, ok := QueryOp[rawOp]
	if !ok {
		return nil, fmt.Errorf("Invalid query string param '%s': unknown operator '%s'", key, rawOp)
	}
	return &ParamFilter{
		Key:    key,
		Column: col,
		RawOp:  rawOp,
		SQLOp:  sqlOp,
		Values: values,
	}, nil
}

// NewSortParam creates a new sort parameter to add to a query.
//
// Sorts will appear on the query string like this:
//
// sort=updated_at__desc&sort=user_id__asc
//
// That means sort first by updated_at descending, then by user_id
// ascending. If sort direction is missing or invalid, it defaults
// to asc.
func NewSortParam(value string) *SortParam {
	colAndDir := strings.Split(value, "__")
	if len(colAndDir) != 2 {
		colAndDir = append(colAndDir, "asc")
	}
	col := colAndDir[0]
	direction := colAndDir[1]
	if direction != "asc" && direction != "desc" {
		direction = "asc"
	}
	return &SortParam{
		Column:    col,
		Direction: direction,
	}
}

// AddToQuery adds this ParamFilter to SQL query q. If it can't map the
// RawOp to a known Query method, it returns a custom error. The caller
// should log the error and then return a basic common.ErrInvalidParam.
func (pf *ParamFilter) AddToQuery(q *Query) error {
	switch pf.RawOp {
	case "eq":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "ne":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "gt":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "gteq":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "lt":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "lteq":
		q.Where(pf.Column, pf.SQLOp, pf.Values[0])
	case "starts_with":
		q.Where(pf.Column, pf.SQLOp, fmt.Sprintf("%s%%", pf.Values[0]))
	case "contains":
		q.Where(pf.Column, pf.SQLOp, fmt.Sprintf("%%%s%%", pf.Values[0]))
	case "is_null":
		q.IsNull(pf.Column)
	case "not_null":
		q.IsNotNull(pf.Column)
	case "in":
		q.WhereIn(pf.Column, pf.InterfaceValues()...)
	default:
		return fmt.Errorf("Invalid query string param '%s': unknown operator '%s'", pf.Key, pf.RawOp)
	}
	return nil
}

// InterfaceValues converts []string Values to []interface{} values.
// We get string values from the HTTP query string, but we need to provide
// interface{} values to the pg library that will query the database.
// Golang type suckage.
func (pf *ParamFilter) InterfaceValues() []interface{} {
	iValues := make([]interface{}, len(pf.Values))
	for i, value := range pf.Values {
		iValues[i] = value
	}
	return iValues
}
