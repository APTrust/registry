package controllers_test

import (
	"testing"

	"github.com/APTrust/registry/controllers"
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
		SqlOp:  "LIKE",
		Values: []string{"Simpson"},
	},
	&controllers.ParamFilter{
		Key:    "name__contains",
		Column: "name",
		RawOp:  "contains",
		SqlOp:  "LIKE",
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

func TestNewParamFilter(t *testing.T) {
	// Constructor should set all properties correctly based on input params.
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

func TestAddToQuery_Basic(t *testing.T) {
	//q := models.NewQuery()
	//filter, err := controllers.NewParamFilter()
}

func TestAddToQuery_Like(t *testing.T) {

}

func TestAddToQuery_Null(t *testing.T) {

}

func TestAddToQuery_In(t *testing.T) {

}

func TestInterfaceValues(t *testing.T) {

}
