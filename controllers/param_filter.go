package controllers

import (
	"fmt"
	"strings"

	"github.com/APTrust/registry/models"
)

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
	case "ne":
	case "gt":
	case "gteq":
	case "lt":
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
// Golang type suckage.
func (pf *ParamFilter) InterfaceValues() []interface{} {
	iValues := make([]interface{}, len(pf.Values))
	for i, value := range pf.Values {
		iValues[i] = value
	}
	return iValues
}
