package helpers

import (
	"fmt"
	"html/template"
	"time"

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
