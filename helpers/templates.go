package helpers

import (
	"fmt"
	"html/template"
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
