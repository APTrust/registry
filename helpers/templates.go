package helpers

import (
	"fmt"
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
	return date.Format("Jan _2, 2006")
}

// DateISO returns a date in format "2006-01-02"
func DateISO(date time.Time) string {
	return date.Format("2006-01-02")
}
