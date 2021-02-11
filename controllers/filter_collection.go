package controllers

import (
	"fmt"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/models"
)

// ParamsToQuery converts query string params such as name__eq=Homer to
// a models.Query object that allows us to build a SQL where clause as
// we go.
type FilterCollection struct {
	filters []*ParamFilter
}

// NewFilterCollection returns a ParamFileters object.
func NewFilterCollection() *FilterCollection {
	return &FilterCollection{
		filters: make([]*ParamFilter, 0),
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

// ToQuery returns a query object based on the keys and values passed in.
// The Query's WhereClause() will return the where conditions for the filters
// passed in through Add(), and the Query's Params() method will return the
// params. Conditions and params come back in the order they were added.
func (fc *FilterCollection) ToQuery() (*models.Query, error) {
	query := models.NewQuery()
	for _, filter := range fc.filters {
		if common.ListIsEmpty(filter.Values) {
			continue // no need to apply filter
		}
		err := filter.AddToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return query, nil
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
	// SqlOp is the SQL operator that corresponds to RawOp. For example,
	// "=" or ">".
	SqlOp string
	// Values are the values attached to Key in the query string.
	Values []string
}

// NewParamFilter returns a new ParamFilter object based on the key and
// values submitted in the query string. It will return a custom error
// if it can't parse the key or values. The caller should log the error
// and return a basic common.ErrInvalidParam.
func NewParamFilter(key string, values []string) (*ParamFilter, error) {
	// Parse the column and operator from the key name
	colAndOp := strings.Split(key, "__")
	if len(colAndOp) != 2 {
		return nil, fmt.Errorf("Invalid query string param '%s': missing '__'", key)
	}
	col := colAndOp[0]
	rawOp := colAndOp[1]
	sqlOp, ok := models.QueryOp[rawOp]
	if !ok {
		return nil, fmt.Errorf("Invalid query string param '%s': unknown operator '%s'", key, rawOp)
	}
	return &ParamFilter{
		Key:    key,
		Column: col,
		RawOp:  rawOp,
		SqlOp:  sqlOp,
		Values: values,
	}, nil
}

// AddToQuery adds this ParamFilter to SQL query q. If it can't map the
// RawOp to a known Query method, it returns a custom error. The caller
// should log the error and then return a basic common.ErrInvalidParam.
func (pf *ParamFilter) AddToQuery(q *models.Query) error {
	switch pf.RawOp {
	case "eq":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "ne":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "gt":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "gteq":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "lt":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "lteq":
		q.Where(pf.Column, pf.SqlOp, pf.Values[0])
	case "starts_with":
		q.Where(pf.Column, pf.SqlOp, fmt.Sprintf("%s%%", pf.Values[0]))
	case "contains":
		q.Where(pf.Column, pf.SqlOp, fmt.Sprintf("%%%s%%", pf.Values[0]))
	case "is_null":
		q.IsNull(pf.Column)
	case "not_null":
		q.IsNotNull(pf.Column)
	case "in":
		q.WithAny(pf.Column, pf.InterfaceValues()...)
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
