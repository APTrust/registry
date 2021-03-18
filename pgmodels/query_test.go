package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TestDate = time.Date(2021, 6, 16, 10, 24, 16, 0, time.UTC)

func TestQueryWhere(t *testing.T) {

	// Empty where clause
	q := pgmodels.NewQuery()
	assert.Equal(t, "", q.WhereClause())
	require.Equal(t, 0, len(q.Params()))

	// Int
	q = pgmodels.NewQuery()
	q.Where("id", "=", 100)
	assert.Equal(t, `(id = ?)`, q.WhereClause())
	require.Equal(t, 1, len(q.Params()))
	assert.Equal(t, 100, q.Params()[0])

	// Int64
	q = pgmodels.NewQuery()
	q.Where("id", "=", int64(100))
	assert.Equal(t, `(id = ?)`, q.WhereClause())
	require.Equal(t, 1, len(q.Params()))
	assert.Equal(t, int64(100), q.Params()[0])

	// String
	q = pgmodels.NewQuery()
	q.Where("name", "=", "Homer Simpson")
	assert.Equal(t, `(name = ?)`, q.WhereClause())
	require.Equal(t, 1, len(q.Params()))
	assert.Equal(t, "Homer Simpson", q.Params()[0])

	// Bool
	q = pgmodels.NewQuery()
	q.Where("active", "=", true)
	assert.Equal(t, `(active = ?)`, q.WhereClause())
	require.Equal(t, 1, len(q.Params()))
	assert.Equal(t, true, q.Params()[0])

	// Time
	q = pgmodels.NewQuery()
	q.Where("created_at", ">=", TestDate)
	assert.Equal(t, `(created_at >= ?)`, q.WhereClause())
	require.Equal(t, 1, len(q.Params()))
	assert.Equal(t, TestDate, q.Params()[0])
}

func TestWithIsNull(t *testing.T) {
	q := pgmodels.NewQuery()
	q.IsNull("email")
	assert.Equal(t, `(email is null)`, q.WhereClause())
	assert.Equal(t, 0, len(q.Params()))
}

func TestWithIsNotNull(t *testing.T) {
	q := pgmodels.NewQuery()
	q.IsNotNull("email")
	assert.Equal(t, `(email is not null)`, q.WhereClause())
	assert.Equal(t, 0, len(q.Params()))
}

func TestOr(t *testing.T) {
	q := pgmodels.NewQuery()
	cols := []string{"col1", "col2", "col3"}
	ops := []string{"=", "<", ">"}
	vals := []interface{}{"val1", "val2", "val3"}
	q.Or(cols, ops, vals)

	assert.Equal(t, `(col1 = ? OR col2 < ? OR col3 > ?)`, q.WhereClause())
	assert.Equal(t, 3, len(q.Params()))
}

func TestAnd(t *testing.T) {
	q := pgmodels.NewQuery()
	cols := []string{"col1", "col2", "col3"}
	ops := []string{"=", "<", ">"}
	vals := []interface{}{"val1", "val2", "val3"}
	q.And(cols, ops, vals)

	assert.Equal(t, `(col1 = ? AND col2 < ? AND col3 > ?)`, q.WhereClause())
	assert.Equal(t, 3, len(q.Params()))
}

func TestBetweenInclusive(t *testing.T) {
	q := pgmodels.NewQuery()
	q.BetweenInclusive("col1", 28, 42)
	assert.Equal(t, `(col1 >= ? AND col1 <= ?)`, q.WhereClause())
	assert.Equal(t, 2, len(q.Params()))
}

func TestBetweenExclusive(t *testing.T) {
	q := pgmodels.NewQuery()
	q.BetweenExclusive("col1", 28, 42)
	assert.Equal(t, `(col1 > ? AND col1 < ?)`, q.WhereClause())
	assert.Equal(t, 2, len(q.Params()))
}

func TestMakePlaceholders(t *testing.T) {
	q := pgmodels.NewQuery()
	assert.Equal(t, "?, ?, ?, ?", q.MakePlaceholders(0, 4))
	assert.Equal(t, "?, ?, ?, ?", q.MakePlaceholders(4, 4))
}

func TestWithAny(t *testing.T) {
	q := pgmodels.NewQuery()
	q.WithAny("col12", []interface{}{1, 2, 3, 4}...)
	assert.Equal(t, `(col12 IN (?, ?, ?, ?))`, q.WhereClause())
	assert.Equal(t, 4, len(q.Params()))
}

func TestWithMultipleConditions(t *testing.T) {
	q := pgmodels.NewQuery()

	q.Where("org_id", "=", int64(100))
	q.Where("name", "=", "Ned Flanders")
	q.Where("active", "=", true)
	q.Where("created_at", ">=", TestDate)
	q.BetweenInclusive("age", 26, 34)

	cols := []string{"col1", "col2", "col3"}
	ops := []string{"=", "<", ">"}
	vals := []interface{}{"val1", "val2", "val3"}
	q.Or(cols, ops, vals)

	q.WithAny("col12", []interface{}{1, 2, 3, 4}...)
	q.IsNull("col99")

	assert.Equal(t, `(org_id = ?) AND (name = ?) AND (active = ?) AND (created_at >= ?) AND (age >= ? AND age <= ?) AND (col1 = ? OR col2 < ? OR col3 > ?) AND (col12 IN (?, ?, ?, ?)) AND (col99 is null)`, q.WhereClause())
	require.Equal(t, 13, len(q.Params()))
	assert.Equal(t, int64(100), q.Params()[0])
	assert.Equal(t, "Ned Flanders", q.Params()[1])
	assert.Equal(t, true, q.Params()[2])
	assert.Equal(t, TestDate, q.Params()[3])
	assert.Equal(t, 26, q.Params()[4])
	assert.Equal(t, 34, q.Params()[5])
	assert.Equal(t, "val1", q.Params()[6])
	assert.Equal(t, "val2", q.Params()[7])
	assert.Equal(t, "val3", q.Params()[8])
	assert.Equal(t, 1, q.Params()[9])
	assert.Equal(t, 2, q.Params()[10])
	assert.Equal(t, 3, q.Params()[11])
	assert.Equal(t, 4, q.Params()[12])
}

func TestSelect(t *testing.T) {
	db.LoadFixtures()
	q := pgmodels.NewQuery().Where("id", "=", 1)
	inst := &pgmodels.Institution{}
	err := q.Select(inst)
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, int64(1), inst.ID)

	q = pgmodels.NewQuery().IsNotNull("created_at")
	var institutions []*pgmodels.Institution
	err = q.Select(&institutions)
	require.Nil(t, err)
	assert.True(t, len(institutions) > 3)
}
