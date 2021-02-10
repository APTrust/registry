package controllers_test

import (
	"testing"

	"github.com/APTrust/registry/controllers"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var valid = []*controllers.ParamFilter{
	&controllers.ParamFilter{
		Key:    "name__eq",
		Column: "name",
		RawOp:  "eq",
		SqlOp:  "=",
		Values: []string{"Homer"},
	},
	&controllers.ParamFilter{
		Key:    "name__ne",
		Column: "name",
		RawOp:  "ne",
		SqlOp:  "!=",
		Values: []string{"Homer"},
	},
	&controllers.ParamFilter{
		Key:    "age__gt",
		Column: "age",
		RawOp:  "gt",
		SqlOp:  ">",
		Values: []string{"38"},
	},
	&controllers.ParamFilter{
		Key:    "age__gteq",
		Column: "age",
		RawOp:  "gteq",
		SqlOp:  ">=",
		Values: []string{"38"},
	},
	&controllers.ParamFilter{
		Key:    "age__lt",
		Column: "age",
		RawOp:  "lt",
		SqlOp:  "<",
		Values: []string{"38"},
	},
	&controllers.ParamFilter{
		Key:    "age__lteq",
		Column: "age",
		RawOp:  "lteq",
		SqlOp:  "<=",
		Values: []string{"38"},
	},
	&controllers.ParamFilter{
		Key:    "name__starts_with",
		Column: "name",
		RawOp:  "starts_with",
		SqlOp:  "ILIKE",
		Values: []string{"Simpson"},
	},
	&controllers.ParamFilter{
		Key:    "name__contains",
		Column: "name",
		RawOp:  "contains",
		SqlOp:  "ILIKE",
		Values: []string{"Simpson"},
	},
	&controllers.ParamFilter{
		Key:    "name__is_null",
		Column: "name",
		RawOp:  "is_null",
		SqlOp:  "IS NULL",
		Values: []string{"true"},
	},
	&controllers.ParamFilter{
		Key:    "name__not_null",
		Column: "name",
		RawOp:  "not_null",
		SqlOp:  "IS NOT NULL",
		Values: []string{"true"},
	},
	&controllers.ParamFilter{
		Key:    "name__in",
		Column: "name",
		RawOp:  "in",
		SqlOp:  "IN",
		Values: []string{"Bart", "Lisa", "Maggie"},
	},
}

var invalid = []*controllers.ParamFilter{
	&controllers.ParamFilter{
		Key:    "name__xyz",
		Column: "name",
		RawOp:  "xyz",
		SqlOp:  "",
		Values: []string{},
	},
}

func expectedQuery(index int) *models.Query {
	q := models.NewQuery()
	switch index {
	case 0:
		q.Where("name", "=", "Homer")
	case 1:
		q.Where("name", "!=", "Homer")
	case 2:
		q.Where("age", ">", "38")
	case 3:
		q.Where("age", ">=", "38")
	case 4:
		q.Where("age", "<", "38")
	case 5:
		q.Where("age", "<=", "38")
	case 6:
		q.Where("name", "ILIKE", "Simpson%")
	case 7:
		q.Where("name", "ILIKE", "%Simpson%")
	case 8:
		q.IsNull("name")
	case 9:
		q.IsNotNull("name")
	case 10:
		q.WithAny("name", []interface{}{"Bart", "Lisa", "Maggie"}...)
	}
	return q
}

// TestNewParamFilter checks that Constructor sets all properties correctly
// based on input params.
func TestNewParamFilter(t *testing.T) {
	for _, obj := range valid {
		filter, err := controllers.NewParamFilter(obj.Key, obj.Values)
		assert.Nil(t, err)
		require.NotNil(t, filter)
		assert.Equal(t, obj.Key, filter.Key)
		assert.Equal(t, obj.Column, filter.Column)
		assert.Equal(t, obj.RawOp, filter.RawOp)
		assert.Equal(t, obj.SqlOp, filter.SqlOp)
		assert.Equal(t, obj.Values, filter.Values)
	}
	for _, obj := range invalid {
		filter, err := controllers.NewParamFilter(obj.Key, obj.Values)
		assert.NotNil(t, err)
		require.Nil(t, filter)
	}
}

// TestAddToQuery ensures that ParamFilter correctly translates something
// like this:
//
// name__eq=Homer
//
// into something like this:
//
// WhereClause: "name" = ?
// Params: ["Homer"]
func TestAddToQuery(t *testing.T) {
	for i, obj := range valid {
		q := models.NewQuery()
		filter, err := controllers.NewParamFilter(obj.Key, obj.Values)
		assert.Nil(t, err)
		err = filter.AddToQuery(q)
		require.Nil(t, err)

		expected := expectedQuery(i)
		assert.Equal(t, expected.WhereClause(), q.WhereClause(), "Index = %d", i)
		assert.Equal(t, expected.Params(), q.Params(), "Index = %d", i)
	}
}

// TestInterfaceValues ensures that we can get the filter's string values
// as a slice of []interface{}.
func TestInterfaceValues(t *testing.T) {
	values := []string{
		"val1",
		"val2",
		"val3",
	}
	filter, err := controllers.NewParamFilter("col1__in", values)
	require.Nil(t, err)
	assert.Equal(t, []interface{}{"val1", "val2", "val3"}, filter.InterfaceValues())
}
