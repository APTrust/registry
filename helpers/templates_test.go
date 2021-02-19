package helpers_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate, _ = time.Parse(time.RFC3339, "2021-04-16T15:04:05Z")
var textString = "The Academic Preservation Trust (APTrust) is committed to the creation and management of a sustainable environment for digital preservation."
var truncatedString = "The Academic Preservation Trust..."

func TestTemplateVars(t *testing.T) {
	c := &gin.Context{}
	c.Set("CurrentUser", &models.User{Name: "John von Neumann"})
	vars := helpers.TemplateVars(c)
	require.NotNil(t, vars)
	require.NotNil(t, vars["CurrentUser"])
	assert.Equal(t, "John von Neumann", vars["CurrentUser"].(*models.User).Name)
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, truncatedString, helpers.Truncate(textString, 31))
}

func TestDateUS(t *testing.T) {
	assert.Equal(t, "Apr 16, 2021", helpers.DateUS(testDate))
}

func TestDateISO(t *testing.T) {
	assert.Equal(t, "2021-04-16", helpers.DateISO(testDate))
}

func TestEqStrInt(t *testing.T) {
	assert.True(t, helpers.EqStrInt64("200", 200))
	assert.False(t, helpers.EqStrInt64("200", 909))
}

func TestEqStrInt64(t *testing.T) {
	assert.True(t, helpers.EqStrInt64("200", int64(200)))
	assert.False(t, helpers.EqStrInt64("200", int64(909)))
}
