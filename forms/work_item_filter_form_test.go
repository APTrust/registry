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

func getWorkItemFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("action__in", []string{constants.ActionGlacierRestore})
	fc.Add("alt_identifier", []string{"blah-blah-blah"})
	fc.Add("bag_date__gteq", []string{"2021-01-01"})
	fc.Add("bag_date__lteq", []string{"2022-01-01"})
	fc.Add("bag_group_identifier", []string{"yadda"})
	fc.Add("bagit_profile_identifier", []string{constants.BTRProfileIdentifier})
	fc.Add("bucket", []string{"some.bucket"})
	fc.Add("date_processed__gteq", []string{"2021-01-01"})
	fc.Add("date_processed__lteq", []string{"2021-10-21"})
	fc.Add("etag", []string{"5556677"})
	fc.Add("generic_file_identifier", []string{"test.edu/bag/data/file.txt"})
	fc.Add("institution_id", []string{"3"})
	fc.Add("name", []string{"barney"})
	fc.Add("needs_admin_review", []string{"true"})
	fc.Add("node__not_null", []string{"true"})
	fc.Add("object_identifier", []string{"test.edu/bag"})
	fc.Add("size__gteq", []string{"800"})
	fc.Add("size__lteq", []string{"1600"})
	fc.Add("stage__in", []string{constants.StageReceive})
	fc.Add("status__in", []string{constants.StatusPending})
	fc.Add("storage_option", []string{constants.StorageOptionGlacierOR})
	fc.Add("user", []string{"barney@example.com"})
	return fc
}

func getWorkItemFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getWorkItemFilters()
	form, err := forms.NewWorkItemFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestWorkItemFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getWorkItemFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testWorkItemFields(t, fc, fields)
	assertListsNotEmpty(t, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestWorkItemFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getWorkItemFilterForm(t, user)
		fields := form.GetFields()
		testWorkItemFields(t, fc, fields)
		assertListsNotEmpty(t, fields)
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func assertListsNotEmpty(t *testing.T, fields map[string]*forms.Field) {
	assert.True(t, len(fields["action__in"].Options) > 1)
	assert.True(t, len(fields["bagit_profile_identifier"].Options) > 1)
	assert.True(t, len(fields["needs_admin_review"].Options) > 1)
	assert.True(t, len(fields["node__not_null"].Options) > 1)
	assert.True(t, len(fields["stage__in"].Options) > 1)
	assert.True(t, len(fields["status__in"].Options) > 1)
	assert.True(t, len(fields["storage_option"].Options) > 1)
}

func testWorkItemFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	fieldNames := []string{
		"action__in",
		"alt_identifier",
		"bag_date__gteq",
		"bag_date__lteq",
		"bag_group_identifier",
		"bagit_profile_identifier",
		"bucket",
		"date_processed__gteq",
		"date_processed__lteq",
		"etag",
		"generic_file_identifier",
		// "institution_id", --> sys admin only
		"name",
		"needs_admin_review",
		"node__not_null",
		"object_identifier",
		"size__gteq",
		"size__lteq",
		"stage__in",
		"status__in",
		"storage_option",
		"user",
	}
	for _, field := range fieldNames {
		assert.Equal(t, fc.ValueOf(field), fields[field].Value, field)
	}
}
