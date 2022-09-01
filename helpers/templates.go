package helpers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet

// -------------------------------------------------------------------------
// Helper functions to be used inside of templates
//
// These are loaded in t2m.go (main) by r.SetFuncMap()
// -------------------------------------------------------------------------

// Define this here so it's not recompiled on every call to LinkifyUrls.
var reUrl = regexp.MustCompile(`((https?://)[^\s]+)`)

// Truncate truncates the value to the given length, appending
// an ellipses to the end. If value contains HTML tags, they
// will be stripped because truncating HTML can result in unclosed
// tags that will ruin the page layout.
func Truncate(value string, length int) string {
	if len(value) < length {
		return value
	}
	fmtStr := fmt.Sprintf("%%.%ds...", length)
	return fmt.Sprintf(fmtStr, value)
}

// DateUS returns a date in format "Jan 2, 2006"
func DateUS(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("Jan _2, 2006")
}

// DateISO returns a date in format "2006-01-02"
func DateISO(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}

// DateTimeISO returns a date in format "2006-01-02T15:04:05Z"
func DateTimeISO(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(time.RFC3339)
}

// UnixToISO converts a Unix timestamp to ISO format.
func UnixToISO(ts int64) string {
	return time.Unix(ts, 0).Format(time.RFC3339)
}

// RoleName transforms ugly DB role names into more readable ones.
func RoleName(role string) string {
	switch role {
	case "admin":
		return "SysAdmin"
	case "institutional_admin":
		return "Admin"
	case "institutional_user":
		return "User"
	default:
		return role
	}
}

// YesNo returns "Yes" if value is true, "No" if value is false.
func YesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// StrEq compares the string representation of two values and returns
// true if they are equal.
func StrEq(val1, val2 interface{}) bool {
	str1 := fmt.Sprintf("%v", val1)
	str2 := fmt.Sprintf("%v", val2)
	return str1 == str2
}

// EscapeAttr escapes an HTML attribute value.
// This helps avoid the ZgotmplZ problem.
func EscapeAttr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

// EscapeHTML returns an escaped HTML string.
// This helps avoid the ZgotmplZ problem.
func EscapeHTML(s string) template.HTML {
	return template.HTML(s)
}

// UserCan returns true if the user has the specified permission.
func UserCan(user *pgmodels.User, permission constants.Permission, instID int64) bool {
	return user.HasPermission(permission, instID)
}

// HumanSize returns a number of bytes in a human-readable format.
// Note that we use base 1024, not base 1000, because AWS uses 1024
// to calculate the storage size of the objects we're reporting on.
func HumanSize(size int64) string {
	return common.ToHumanSize(size, 1024)
}

// IconFor returns a FontAwesome icon for the specified string, as defined
// in helpers.IconMap. If the IconMap has no entry for the string, this
// returns helpers.IconMissing.
func IconFor(str string) template.HTML {
	icon := IconMap[str]
	if icon == "" {
		icon = IconMissing
	}
	return template.HTML(icon)
}

// BadgeClass returns the css class for the specified string, where
// string is a work item status or other value defined in Constants.
func BadgeClass(str string) template.HTML {
	return template.HTML(BadgeClassMap[str])
}

// TruncateStart trims str to maxLen by removing them from the
// middle of the string. It adds dots to the middle of the string
// to indicate text was trimmed.
func TruncateMiddle(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	half := (maxLen - 3) / 2
	end := len(str) - half
	return str[0:half] + "..." + str[end:len(str)]
}

// TruncateStart trims str to maxLen by removing them from the
// start of the string. It adds leading dots to indicate some
// text was trimmed.
func TruncateStart(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	end := (len(str) - 3) - maxLen
	if end < 0 {
		end = 0
	}
	return "..." + str[end:len(str)]
}

// Dict returns an interface map suitable for passing into
// sub templates.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, common.ErrInvalidParam
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, common.ErrWrongDataType
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// DefaultString returns value if it's non-empty.
// Otherwise, it returns _default.
func DefaultString(value, _default string) string {
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return _default
}

// FormatFloat formats a floating point number to have scale
// digits after the decimal point.
func FormatFloat(value float64, scale int) string {
	fmtString := fmt.Sprintf("%%.%df", scale)
	return fmt.Sprintf(fmtString, value)
}

// ToJSON converts an interface to JSON.
func ToJSON(v interface{}) template.JS {
	jsonString, _ := json.Marshal(v)
	return template.JS(jsonString)
}

// SortUrl returns the url to sort results by the specified column.
// This is used in table column headers on index pages.
// Note that this returns a URL path and query string only. There's
// no hostname, port, or scheme.
func SortUrl(currentUrl *url.URL, colName string) string {
	newSort := fmt.Sprintf("%s__asc", colName)
	vals := currentUrl.Query()
	currentSort := vals.Get("sort")
	if currentSort == fmt.Sprintf("%s__asc", colName) {
		newSort = fmt.Sprintf("%s__desc", colName)
	}
	vals.Set("sort", newSort)
	return fmt.Sprintf("%s?%s", currentUrl.Path, vals.Encode())
}

// SortIcon returns the name of the sort icon to display at the
// top of a table column. This will be either "keyboard_arrow_up"
// or "keyboard_arrow_down"
func SortIcon(currentUrl *url.URL, colName string) string {
	vals := currentUrl.Query()
	currentSort := vals.Get("sort")
	icon := ""
	if currentSort == fmt.Sprintf("%s__desc", colName) {
		icon = "keyboard_arrow_down"
	} else if currentSort == fmt.Sprintf("%s__asc", colName) {
		icon = "keyboard_arrow_up"
	}
	return icon
}

// LinkifyUrls converts urls in text to clickable links. That is,
// it replaces https://example.com with
// <a href="https://example.com">https://example.com</a>
//
// URLs outside the current domain will open in a new tab
// (i.e. will have target="_blank").
//
// This also replaces newlines with <br/> tags.
func LinkifyUrls(text string) template.HTML {
	alreadyReplaced := make(map[string]bool)
	matches := reUrl.FindAllStringSubmatch(text, -1)

	urls := make([]string, len(matches))
	for i, _ := range matches {
		urls[i] = matches[i][0]
	}
	// Sort matches by length, and then reverse so longest is first
	sort.Sort(ByLen(urls))
	reverse(urls)

	// Replace urls with links. Do longest urls first, so we don't
	// double replace items like "https://example.com" and
	// "https://example.com/sub-page"
	for _, u := range urls {
		if strings.HasSuffix(u, ".") || strings.HasSuffix(u, ",") {
			u = u[0 : len(u)-1]
		}
		if alreadyReplaced[u] {
			continue
		}

		// Link should open in new tab, unless it's in our same domain.
		link := fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, u, u)
		parsedUrl, err := url.Parse(u)
		if err == nil && parsedUrl.Hostname() == common.Context().Config.Cookies.Domain {
			link = fmt.Sprintf(`<a href="%s">%s</a>`, u, u)
		}
		text = strings.Replace(text, u, link, -1)
		alreadyReplaced[u] = true
	}

	text = strings.ReplaceAll(text, "\n", "<br/>")

	return template.HTML(text)
}

// ByLen implements sorting by length
type ByLen []string

// Let returns the length of slice a.
func (a ByLen) Len() int {
	return len(a)
}

// Less return true if i is less than j.
func (a ByLen) Less(i, j int) bool {
	return len(a[i]) < len(a[j])
}

// Swap i and j in slice.
func (a ByLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// reverse reverses order of items in slice s
func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
