package helpers_test

import (
	"html/template"
	"net/url"
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

func TestCurrentYear(t *testing.T) {
	assert.Equal(t, time.Now().Year(), helpers.CurrentYear())
}

func TestDateUS(t *testing.T) {
	assert.Equal(t, "Apr 16, 2021", helpers.DateUS(testDate))
	assert.Equal(t, "", helpers.DateUS(time.Time{}))
}

func TestDateTimeUS(t *testing.T) {
	assert.Equal(t, "Apr 16, 2021 12:24:16", helpers.DateTimeUS(testDate))
	assert.Equal(t, "", helpers.DateUS(time.Time{}))
}

func TestDateISO(t *testing.T) {
	assert.Equal(t, "2021-04-16", helpers.DateISO(testDate))
	assert.Equal(t, "", helpers.DateISO(time.Time{}))
}

func TestDateTimeISO(t *testing.T) {
	assert.Equal(t, "2021-04-16T12:24:16Z", helpers.DateTimeISO(testDate))
	assert.Equal(t, "", helpers.DateTimeISO(time.Time{}))
}

func TestRoleName(t *testing.T) {
	assert.Equal(t, "Admin", helpers.RoleName(constants.RoleInstAdmin))
	assert.Equal(t, "User", helpers.RoleName(constants.RoleInstUser))
	assert.Equal(t, "SysAdmin", helpers.RoleName(constants.RoleSysAdmin))
	assert.Equal(t, "not-a-role", helpers.RoleName("not-a-role"))
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
	assert.False(t, helpers.UserCan(admin, constants.FileRequestDelete, 1))
	assert.False(t, helpers.UserCan(admin, constants.FileRequestDelete, 2))
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

func TestHumanSize(t *testing.T) {
	assert.Equal(t, "2.0 kB", helpers.HumanSize(2*1024))
	assert.Equal(t, "2.0 MB", helpers.HumanSize(2*1024*1024))
	assert.Equal(t, "2.0 GB", helpers.HumanSize(2*1024*1024*1024))
	assert.Equal(t, "2.0 TB", helpers.HumanSize(2*1024*1024*1024*1024))
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
	assert.Equal(t, longString, helpers.TruncateStart(longString, 5000))
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
	f := float64(86352137.9786534)
	assert.Equal(t, "86,352,137.98", helpers.FormatFloat(f, 2))
	assert.Equal(t, "86,352,137.979", helpers.FormatFloat(f, 3))
	assert.Equal(t, "86,352,137.9787", helpers.FormatFloat(f, 4))
}

func TestFormatInt(t *testing.T) {
	assert.Equal(t, "1,234,567,890", helpers.FormatInt(1234567890))
}

func TestFormatInt64(t *testing.T) {
	assert.Equal(t, "1,234,567,890", helpers.FormatInt64(1234567890))
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

func TestUnixToISO(t *testing.T) {
	ts := time.Date(2022, 5, 9, 14, 33, 24, 0, time.Local)
	assert.Equal(t, ts.Format(time.RFC3339), helpers.UnixToISO(1652121204))
}

func TestBadgeClass(t *testing.T) {
	assert.Equal(t, template.HTML("is-cancelled"), helpers.BadgeClass(constants.StatusCancelled))
	assert.Equal(t, template.HTML("is-pending"), helpers.BadgeClass(constants.StatusPending))
	assert.Equal(t, template.HTML(""), helpers.BadgeClass("no such class"))
}

func TestSortURL(t *testing.T) {
	currentUrl, err := url.Parse("https://example.com/objects?name=homer&age=39&sort=salary__asc")
	require.Nil(t, err)

	assert.Equal(t, "/objects?age=39&name=homer&sort=salary__desc", helpers.SortUrl(currentUrl, "salary"))
	assert.Equal(t, "/objects?age=39&name=homer&sort=zip_code__asc", helpers.SortUrl(currentUrl, "zip_code"))
}

func TestLinkifyUrl(t *testing.T) {
	text := `Sample alert text.
	This is a local link: http://localhost/alerts/yadda and this
	is an external link: https://example.com/page and nothing else
	should be linked.`
	expected := "Sample alert text.<br/>\tThis is a local link: <a href=\"http://localhost/alerts/yadda\">http://localhost/alerts/yadda</a> and this<br/>\tis an external link: <a href=\"https://example.com/page\" target=\"_blank\">https://example.com/page</a> and nothing else<br/>\tshould be linked."
	assert.Equal(t, template.HTML(expected), helpers.LinkifyUrls(text))
}

func TestSortIcon(t *testing.T) {
	url, _ := url.Parse("https://example.com?sort=")
	assert.Empty(t, helpers.SortIcon(url, "col1"))

	url, _ = url.Parse("https://example.com?sort=col1__asc")
	assert.Equal(t, "keyboard_arrow_up", helpers.SortIcon(url, "col1"))

	url, _ = url.Parse("https://example.com?sort=col1__desc")
	assert.Equal(t, "keyboard_arrow_down", helpers.SortIcon(url, "col1"))

	// No icon, because we're not sorting on this column.
	assert.Empty(t, helpers.SortIcon(url, "col2"))
	assert.Empty(t, helpers.SortIcon(url, "col3"))
}

func TestRevisionURL(t *testing.T) {
	originalCommitID := common.CommitID
	defer func() { common.CommitID = originalCommitID }()

	common.CommitID = ""
	assert.Equal(t, "Missing commit ID", helpers.RevisionURL())

	common.CommitID = "12345678"
	assert.Equal(t, "https://github.com/APTrust/registry/commit/12345678", helpers.RevisionURL())
}

func TestShortCommitHash(t *testing.T) {
	originalCommitID := common.CommitID
	defer func() { common.CommitID = originalCommitID }()

	common.CommitID = ""
	assert.Equal(t, "Missing commit ID", helpers.ShortCommitHash())

	common.CommitID = "0123456789ABCDEF"
	assert.Equal(t, "0123456", helpers.ShortCommitHash())
}

func TestBuildDate(t *testing.T) {
	originalBuildDate := common.BuildDate
	defer func() { common.BuildDate = originalBuildDate }()

	common.BuildDate = ""
	assert.Equal(t, "Missing build date", helpers.BuildDate())

	common.BuildDate = "Tue Jul  8 14:34:31 EDT 2025"
	assert.Equal(t, "Tue Jul  8 14:34:31 EDT 2025", helpers.BuildDate())

}
