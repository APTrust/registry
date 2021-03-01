package helpers_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/helpers"
	"github.com/stretchr/testify/assert"
)

var testDate, _ = time.Parse(time.RFC3339, "2021-04-16T15:04:05Z")
var textString = "The Academic Preservation Trust (APTrust) is committed to the creation and management of a sustainable environment for digital preservation."
var truncatedString = "The Academic Preservation Trust..."

func TestTruncate(t *testing.T) {
	assert.Equal(t, truncatedString, helpers.Truncate(textString, 31))
}

func TestDateUS(t *testing.T) {
	assert.Equal(t, "Apr 16, 2021", helpers.DateUS(testDate))
}

func TestDateISO(t *testing.T) {
	assert.Equal(t, "2021-04-16", helpers.DateISO(testDate))
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
