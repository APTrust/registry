package forms_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getPremisEventFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("date_time__gteq", []string{"2020-01-01"})
	fc.Add("date_time__lteq", []string{"2024-12-31"})
	fc.Add("event_type", []string{constants.EventAccessAssignmentString, constants.EventIngestionString})
	fc.Add("generic_file_identifier", []string{"test.edu/obj/data/file.json"})
	fc.Add("identifier", []string{"0000-0000-0000"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("intellectual_object_identifier", []string{"test.edu/obj"})
	fc.Add("outcome", []string{constants.OutcomeSuccess})
	return fc
}

func getPremisEventFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getPremisEventFilters()
	form, err := forms.NewPremisEventFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestPremisEventFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getPremisEventFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testPremisEventFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
	assert.True(t, len(fields["event_type"].Options) > 1)
	assert.True(t, len(fields["outcome"].Options) > 1)
}

func TestPremisEventFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getPremisEventFilterForm(t, user)
		fields := form.GetFields()
		testPremisEventFields(t, fc, fields)
		// Non sysadmin can see only their own events, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
		assert.True(t, len(fields["event_type"].Options) > 1)
		assert.True(t, len(fields["outcome"].Options) > 1)

		// Late fix for https://trello.com/c/VooirpKZ
		foundSuccessOption := false
		foundFailedOption := false
		for _, option := range fields["outcome"].Options {
			if option.Value == "Success" {
				foundSuccessOption = true
			} else if option.Value == "Failed" {
				foundFailedOption = true
			}
		}
		assert.True(t, foundSuccessOption)
		assert.True(t, foundFailedOption)
	}
}

func testPremisEventFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("date_time__gteq"), fields["date_time__gteq"].Value)
	assert.Equal(t, fc.ValueOf("date_time__lteq"), fields["date_time__lteq"].Value)
	assert.Equal(t, fc.ValueOf("event_type"), fields["event_type"].Value)
	assert.Equal(t, fc.ValueOf("generic_file_identifier"), fields["generic_file_identifier"].Value)
	assert.Equal(t, fc.ValueOf("identifier"), fields["identifier"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("intellectual_object_identifier"), fields["intellectual_object_identifier"].Value)
	assert.Equal(t, fc.ValueOf("outcome"), fields["outcome"].Value)
}
