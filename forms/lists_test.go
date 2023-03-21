package forms_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/forms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInstitutions(t *testing.T) {
	db.LoadFixtures()
	options, err := forms.ListInstitutions(false)
	require.Nil(t, err)
	require.NotEmpty(t, options)
	assert.True(t, len(options) >= 4)
	expected := []forms.ListOption{
		{"1", "APTrust", false},
		{"5", "Example Institution (for integration tests)", false},
		{"2", "Institution One", false},
		{"3", "Institution Two", false},
		{"4", "Test Institution (for integration tests)", false},
		{"6", "Unit Test Institution", false},
	}
	for i, option := range options {
		assert.Equal(t, expected[i].Value, option.Value)
		assert.Equal(t, expected[i].Text, option.Text)
	}
}

func TestOptions(t *testing.T) {
	options := forms.Options(constants.AccessSettings)
	require.NotEmpty(t, options)
	for i, option := range options {
		assert.Equal(t, constants.AccessSettings[i], option.Value)
		assert.Equal(t, constants.AccessSettings[i], option.Text)
	}
}

func TestListUsers(t *testing.T) {
	db.LoadFixtures()
	options, err := forms.ListUsers(3)
	require.Nil(t, err)
	require.NotEmpty(t, options)
	assert.Equal(t, 2, len(options))
	expected := []forms.ListOption{
		{"5", "Inst Two Admin", false},
		{"7", "Inst Two User", false},
	}
	for i, option := range options {
		assert.Equal(t, expected[i].Value, option.Value)
		assert.Equal(t, expected[i].Text, option.Text)
	}
}

func TestListDepositReportDates(t *testing.T) {
	options := forms.ListDepositReportDates()

	today := time.Now().UTC().Format("2006-01-02")
	assert.Equal(t, "Today", options[0].Text)
	assert.Equal(t, today, options[0].Value)
	assert.True(t, len(options) > 80)
	earliestOption := options[len(options)-1]
	assert.Equal(t, "2015-01-01", earliestOption.Value)
	assert.Equal(t, "January 1, 2015", earliestOption.Text)
}
