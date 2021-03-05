package forms_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInstitutions(t *testing.T) {
	db.LoadFixtures()
	adminUser := &models.User{
		Role: constants.RoleSysAdmin,
	}
	ds := models.NewDataStore(adminUser)
	options, err := forms.ListInstitutions(ds, false)
	require.Nil(t, err)
	require.NotEmpty(t, options)
	assert.True(t, len(options) >= 4)
	expected := []forms.ListOption{
		{"1", "APTrust"},
		{"5", "Example Institution (for integration tests)"},
		{"2", "Institution One"},
		{"3", "Institution Two"},
		{"4", "Test Institution (for integration tests)"},
		{"6", "Unit Test Institution"},
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
