package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet

// TemplateVars returns a map to pass into templates.
// The map contains a number of params that are expected to be
// present within most or all templates. We will likely add to
// this as development proceeds.
func TemplateVars(c *gin.Context) gin.H {
	currentUser, _ := c.Get("CurrentUser")
	return gin.H{
		"CurrentUser": currentUser,
	}
}

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

// EqStrInt compares a string to an int and returns true if the int's
// string value matches str.
func EqStrInt(strValue string, intValue int) bool {
	return strValue == strconv.Itoa(intValue)
}

// EqStrInt64 compares a string to an int and returns true if the int's
// string value matches str.
func EqStrInt64(strValue string, int64Value int64) bool {
	return strValue == strconv.FormatInt(int64Value, 10)
}

// Dict takes a list of pairs in the form string, interface{},
// string, interface{}... and returns a map of [string]interface.
// This allows us to pass custom maps as params to nested templates.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("wrong number of params, expected pairs")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
