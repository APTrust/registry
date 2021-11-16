package pgmodels

import (
	"fmt"
	"strings"

	"github.com/APTrust/registry/common"
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
	"not_in":      "NOT IN",
	"is_null":     "IS NULL",
	"not_null":    "IS NOT NULL",
}

// Query provides a fluent model for constructing SQL queries,
// though it currently constructs only the where, order by,
// limit and offset clauses, which is all we need to pass to the
// go-pg adapter. The adapter knows how to construct the select
// and from clauses based on the model.
//
// Query's functionality is a subset of go-pg's query builder
// features. It doesn't support joins, for example, though
// it does support relations.
//
// Query exists to help the web front-end construct queries based
// on named inputs in GET requests. The web.FilterCollection functions
// can dynamically translate query string params into complex
// where clauses like "name ilike '%homer%' and age >= 31 and
// city not in ('Shelbyville', 'Atlanta', 'New York')".
//
// We decided to omit complex joins in our dynamic queries and create
// SQL views instead, because after several years in production, we
// already know virtually all of the ways users want to join tables.
// Views allow us to issue simple queries, which this package supports.
type Query struct {
	conditions []string
	params     []interface{}
	columns    []string
	relations  []string
	orderBy    []string
	offset     int
	limit      int
}

func NewQuery() *Query {
	return &Query{
		conditions: make([]string, 0),
		params:     make([]interface{}, 0),
		columns:    make([]string, 0),
		relations:  make([]string, 0),
		orderBy:    make([]string, 0),
		offset:     -1,
		limit:      -1,
	}
}

func (q *Query) Where(col, op string, val interface{}) *Query {
	cond := fmt.Sprintf(`(%s %s ?)`, col, op)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, val)
	return q
}

func (q *Query) IsNull(col string) *Query {
	cond := fmt.Sprintf(`(%s is null)`, col)
	q.conditions = append(q.conditions, cond)
	return q
}

func (q *Query) IsNotNull(col string) *Query {
	cond := fmt.Sprintf(`(%s is not null)`, col)
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
			conditions[i] = fmt.Sprintf(`%s %s ?`, col, ops[i])
			q.params = append(q.params, vals[i])
		}
		cond := fmt.Sprintf("(%s)", strings.Join(conditions, logicOp))
		q.conditions = append(q.conditions, cond)
	}
}

func (q *Query) BetweenInclusive(col string, low, high interface{}) *Query {
	cond := fmt.Sprintf(`(%s >= ? AND %s <= ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
	return q
}

func (q *Query) BetweenExclusive(col string, low, high interface{}) *Query {
	cond := fmt.Sprintf(`(%s > ? AND %s < ?)`, col, col)
	q.conditions = append(q.conditions, cond)
	q.params = append(q.params, low, high)
	return q
}

func (q *Query) WhereIn(col string, vals ...interface{}) *Query {
	return q.inOrNotIn(col, "IN", vals...)
}

func (q *Query) WhereNotIn(col string, vals ...interface{}) *Query {
	return q.inOrNotIn(col, "NOT IN", vals...)
}

func (q *Query) inOrNotIn(col, op string, vals ...interface{}) *Query {
	if len(vals) > 0 {
		paramStr := q.MakePlaceholders(len(q.params), len(vals))
		cond := fmt.Sprintf(`(%s %s (%s))`, col, op, paramStr)
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

func (q *Query) GetRelations() []string {
	return q.relations
}

func (q *Query) Relations(relations ...string) *Query {
	q.relations = make([]string, len(relations))
	for i, rel := range relations {
		q.relations[i] = rel
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

func (q *Query) GetOrderBy() []string {
	return q.orderBy
}

// TODO: limit direction to asc/desc to prevent sql injection.
func (q *Query) OrderBy(colAndDirection string) *Query {
	if q.orderBy == nil {
		q.orderBy = make([]string, 0)
	}
	q.orderBy = append(q.orderBy, colAndDirection)
	return q
}

// Select executes a query and stores the result in structOrSlice,
// which should be either a pointer to a struct (if you want a
// single result) or a slice of pointers if you want multiple results.
// Returns an error if there is one.
//
// Example:
//
// user := User{}
// err := query.Select(&user)
//
// or
//
// var users []*User
// err := query.Select(&users)
//
func (q *Query) Select(structOrSlice interface{}) error {
	orm := common.Context().DB.Model(structOrSlice)
	for _, rel := range q.GetRelations() {
		orm.Relation(rel)
	}
	if !common.ListIsEmpty(q.GetColumns()) {
		orm.Column(q.GetColumns()...)
	}
	// Empty where clause causes orm to generate empty parens -> ()
	// which causes a SQL error. Include where only if non-empty.
	if q.WhereClause() != "" {
		orm.Where(q.WhereClause(), q.Params()...)
	}
	for _, orderBy := range q.GetOrderBy() {
		orm.Order(orderBy)
	}
	if q.GetLimit() > 0 {
		orm.Limit(q.GetLimit())
	}
	if q.GetOffset() >= 0 {
		orm.Offset(q.GetOffset())
	}
	return orm.Select()
}

func (q *Query) Count(model interface{}) (int, error) {
	orm := common.Context().DB.Model(model)
	if q.WhereClause() != "" {
		orm.Where(q.WhereClause(), q.Params()...)
	}
	return orm.Count()
}
