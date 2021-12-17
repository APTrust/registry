package helpers_test

import (
	"html/template"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate, _ = time.Parse(time.RFC3339, "2021-04-16T12:24:16Z")
var textString = "The Academic Preservation Trust (APTrust) is committed to the creation and management of a sustainable environment for digital preservation."
var truncatedString = "The Academic Preservation Trust..."

func TestTruncate(t *testing.T) {
	assert.Equal(t, truncatedString, helpers.Truncate(textString, 31))
	assert.Equal(t, "hello", helpers.Truncate("hello", 80))
}

func TestDateUS(t *testing.T) {
	assert.Equal(t, "Apr 16, 2021", helpers.DateUS(testDate))
}

func TestDateISO(t *testing.T) {
	assert.Equal(t, "2021-04-16", helpers.DateISO(testDate))
}

func TestDateTimeISO(t *testing.T) {
	assert.Equal(t, "2021-04-16T12:24:16Z", helpers.DateTimeISO(testDate))
}

func TestStrEq(t *testing.T) {
	assert.True(t, helpers.StrEq("4", int8(4)))
	assert.True(t, helpers.StrEq("200", int16(200)))
	assert.True(t, helpers.StrEq("200", int32(200)))
	assert.True(t, helpers.StrEq("200", int64(200)))

	assert.True(t, helpers.StrEq("true", true))
	assert.True(t, helpers.StrEq("true", "true"))
	assert.True(t, helpers.StrEq(true, true))
	assert.True(t, helpers.StrEq(true, "true"))
	assert.True(t, helpers.StrEq(false, "false"))

	assert.False(t, helpers.StrEq("true", false))
	assert.False(t, helpers.StrEq("200", 909))
}

func TestEscapeAttr(t *testing.T) {
	assert.Equal(t, template.HTMLAttr("O'Blivion's"), helpers.EscapeAttr("O'Blivion's"))
}

func TestEscapeHTML(t *testing.T) {
	assert.Equal(t, template.HTML("<em>escape!</em>"), helpers.EscapeHTML("<em>escape!</em>"))
}

func TestUserCan(t *testing.T) {
	admin := &pgmodels.User{
		Role:          constants.RoleSysAdmin,
		InstitutionID: 1,
	}
	assert.True(t, helpers.UserCan(admin, constants.UserCreate, 1))
	assert.True(t, helpers.UserCan(admin, constants.UserCreate, 2))
	assert.True(t, helpers.UserCan(admin, constants.UserCreate, 100))
	assert.True(t, helpers.UserCan(admin, constants.FileRequestDelete, 1))
	assert.True(t, helpers.UserCan(admin, constants.FileRequestDelete, 2))
	assert.True(t, helpers.UserCan(admin, constants.FileRestore, 1))
	assert.True(t, helpers.UserCan(admin, constants.FileRestore, 2))

	instAdmin := &pgmodels.User{
		Role:          constants.RoleInstAdmin,
		InstitutionID: 1,
	}
	assert.True(t, helpers.UserCan(instAdmin, constants.UserCreate, 1))
	assert.False(t, helpers.UserCan(instAdmin, constants.UserCreate, 2))
	assert.False(t, helpers.UserCan(instAdmin, constants.UserCreate, 100))
	assert.True(t, helpers.UserCan(instAdmin, constants.FileRequestDelete, 1))
	assert.False(t, helpers.UserCan(instAdmin, constants.FileRequestDelete, 2))
	assert.True(t, helpers.UserCan(instAdmin, constants.FileRestore, 1))
	assert.False(t, helpers.UserCan(instAdmin, constants.FileRestore, 2))

	instUser := &pgmodels.User{
		Role:          constants.RoleInstUser,
		InstitutionID: 1,
	}
	assert.False(t, helpers.UserCan(instUser, constants.UserCreate, 1))
	assert.False(t, helpers.UserCan(instUser, constants.UserCreate, 2))
	assert.False(t, helpers.UserCan(instUser, constants.UserCreate, 100))
	assert.False(t, helpers.UserCan(instUser, constants.FileRequestDelete, 1))
	assert.False(t, helpers.UserCan(instUser, constants.FileRequestDelete, 2))
	assert.True(t, helpers.UserCan(instUser, constants.FileRestore, 1))
	assert.False(t, helpers.UserCan(instUser, constants.FileRestore, 2))
}

func TestIconFor(t *testing.T) {
	// Should return item defined in map
	assert.Equal(
		t,
		template.HTML(helpers.IconMap[constants.EventIngestion]),
		helpers.IconFor(constants.EventIngestion))

	// If item is not defined in map, should return IconMissing
	assert.Equal(
		t,
		template.HTML(helpers.IconMissing),
		helpers.IconFor("** missing **"))
}

var longString = "Somewhere in la Mancha, in a place whose name I do not care to remember, a gentleman lived not long ago, one of those who has a lance and ancient shield on a shelf and keeps a skinny nag and a greyhound for racing."

func TestTruncateMiddle(t *testing.T) {
	assert.Equal(t, "Somewher... racing.", helpers.TruncateMiddle(longString, 20))
	assert.Equal(t, "Somewhere in ...d for racing.", helpers.TruncateMiddle(longString, 30))
	assert.Equal(t, longString, helpers.TruncateMiddle(longString, 500))
}

func TestTruncateStart(t *testing.T) {
	assert.Equal(t, "...a greyhound for racing.", helpers.TruncateStart(longString, 20))
	assert.Equal(t, "...y nag and a greyhound for racing.", helpers.TruncateStart(longString, 30))
	assert.Equal(t, longString, helpers.TruncateStart(longString, 500))
}

func TestDict(t *testing.T) {
	expected := map[string]interface{}{
		"key1": 1,
		"key2": "two",
	}
	dict, err := helpers.Dict("key1", 1, "key2", "two")
	require.Nil(t, err)
	assert.Equal(t, expected, dict)

	dict, err = helpers.Dict("key1", 1, "key2")
	assert.Equal(t, common.ErrInvalidParam, err)

	dict, err = helpers.Dict(1, "key2")
	assert.Equal(t, common.ErrWrongDataType, err)
}

func TestDefaultString(t *testing.T) {
	assert.Equal(t, "---", helpers.DefaultString("", "---"))
	assert.Equal(t, "---", helpers.DefaultString("  ", "---"))
	assert.Equal(t, "birdy num num", helpers.DefaultString("birdy num num", "---"))
}

func TestFormatFloat(t *testing.T) {
	f := float64(2137.9786534)
	assert.Equal(t, "2137.98", helpers.FormatFloat(f, 2))
	assert.Equal(t, "2137.979", helpers.FormatFloat(f, 3))
	assert.Equal(t, "2137.9787", helpers.FormatFloat(f, 4))
}

func TestToJSON(t *testing.T) {
	x := struct {
		Names   []string
		Numbers []int
	}{
		Names:   []string{"Homer", "Marge"},
		Numbers: []int{39, 38},
	}
	expected := template.JS(`{"Names":["Homer","Marge"],"Numbers":[39,38]}`)
	assert.Equal(t, expected, helpers.ToJSON(x))
}
