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
	"starts_with": "ILIKE",
	"contains":    "ILIKE",
	"in":          "IN",
	"is_null":     "IS NULL",
	"not_null":    "IS NOT NULL",
}

type Query struct {
	conditions []string
	params     []interface{}
	columns    []string
	orderBy    string
	offset     int
	limit      int
}

func NewQuery() *Query {
	return &Query{
		conditions: make([]string, 0),
		params:     make([]interface{}, 0),
		columns:    make([]string, 0),
		offset:     -1,
		limit:      -1,
	}
}

func (q *Query) Where(col, op string, val interface{}) *Query {
	cond := fmt.Sprintf(`("%s" %s ?)`, col, op)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, val)
	return q
}

func (q *Query) IsNull(col string) *Query {
	cond := fmt.Sprintf(`("%s" is null)`, col)
	q.conditions = append(q.conditions, cond)
	return q
}

func (q *Query) IsNotNull(col string) *Query {
	cond := fmt.Sprintf(`("%s" is not null)`, col)
	q.conditions = append(q.conditions, cond)
	return q
}

func (q *Query) Or(cols, ops []string, vals []interface{}) *Query {
	q.multi(cols, ops, vals, " OR ")
	return q
}

func (q *Query) And(cols, ops []string, vals []interface{}) *Query {
	q.multi(cols, ops, vals, " AND ")
	return q
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

func (q *Query) BetweenInclusive(col string, low, high interface{}) *Query {
	cond := fmt.Sprintf(`("%s" >= ? AND "%s" <= ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
	return q
}

func (q *Query) BetweenExclusive(col string, low, high interface{}) *Query {
	cond := fmt.Sprintf(`("%s" > ? AND "%s" < ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
	return q
}

func (q *Query) WithAny(col string, vals ...interface{}) *Query {
	if len(vals) > 0 {
		paramStr := q.MakePlaceholders(len(q.params), len(vals))
		cond := fmt.Sprintf(`("%s" IN (%s))`, col, paramStr)
		q.conditions = append(q.conditions, cond)
		q.params = append(q.params, vals...)
	}
	return q
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

func (q *Query) GetColumns() []string {
	return q.columns
}

func (q *Query) Columns(cols ...string) *Query {
	q.columns = make([]string, len(cols))
	for i, col := range cols {
		q.columns[i] = col
	}
	return q
}

func (q *Query) GetOffset() int {
	return q.offset
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

func (q *Query) GetLimit() int {
	return q.limit
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) GetOrderBy() string {
	return q.orderBy
}

func (q *Query) OrderBy(orderBy string) *Query {
	q.orderBy = orderBy
	return q
}
