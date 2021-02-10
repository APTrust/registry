package models

import (
	"fmt"
	"strings"
)

var QueryOp = map[string]string{
	"eq":          "=",
	"ne":          "!=",
	"gt":          ">",
	"gteq":        ">=",
	"lt":          "<",
	"lteq":        "<=",
	"starts_with": "LIKE",
	"contains":    "LIKE",
	"in":          "IN",
	"is_null":     "IS NULL",
	"not_null":    "IS NOT NULL",
}

type Query struct {
	conditions []string
	params     []interface{}
	OrderBy    string
	Offset     int
	Limit      int
}

func NewQuery() *Query {
	return &Query{
		conditions: make([]string, 0),
		params:     make([]interface{}, 0),
		Offset:     -1,
		Limit:      -1,
	}
}

func (q *Query) Where(col, op string, val interface{}) {
	cond := fmt.Sprintf(`("%s" %s ?)`, col, op)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, val)
}

func (q *Query) IsNull(col string) {
	cond := fmt.Sprintf(`("%s" is null)`, col)
	q.conditions = append(q.conditions, cond)
}

func (q *Query) IsNotNull(col string) {
	cond := fmt.Sprintf(`("%s" is not null)`, col)
	q.conditions = append(q.conditions, cond)
}

func (q *Query) Or(cols, ops []string, vals []interface{}) {
	q.multi(cols, ops, vals, " OR ")
}

func (q *Query) And(cols, ops []string, vals []interface{}) {
	q.multi(cols, ops, vals, " AND ")
}

func (q *Query) multi(cols, ops []string, vals []interface{}, logicOp string) {
	if len(vals) > 0 && len(cols) == len(vals) {
		conditions := make([]string, len(cols))
		for i, col := range cols {
			conditions[i] = fmt.Sprintf(`"%s" %s ?`, col, ops[i])
			q.params = append(q.params, vals[i])
		}
		cond := fmt.Sprintf("(%s)", strings.Join(conditions, logicOp))
		q.conditions = append(q.conditions, cond)
	}
}

func (q *Query) BetweenInclusive(col string, low, high interface{}) {
	cond := fmt.Sprintf(`("%s" >= ? AND "%s" <= ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
}

func (q *Query) BetweenExclusive(col string, low, high interface{}) {
	cond := fmt.Sprintf(`("%s" > ? AND "%s" < ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
}

func (q *Query) WithAny(col string, vals ...interface{}) {
	if len(vals) > 0 {
		paramStr := q.MakePlaceholders(len(q.params), len(vals))
		cond := fmt.Sprintf(`("%s" IN (%s))`, col, paramStr)
		q.conditions = append(q.conditions, cond)
		q.params = append(q.params, vals...)
	}
}

func (q *Query) MakePlaceholders(start, count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func (q *Query) WhereClause() string {
	if len(q.conditions) > 0 {
		return fmt.Sprintf(strings.Join(q.conditions, " AND "))
	}
	return ""
}

func (q *Query) Params() []interface{} {
	return q.params
}
