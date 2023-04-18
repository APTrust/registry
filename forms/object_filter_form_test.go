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

func getObjectFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("access", []string{constants.AccessConsortia})
	fc.Add("alt_identifier__starts_with", []string{"sap"})
	fc.Add("bag_group_identifier__starts_with", []string{"54321"})
	fc.Add("bag_name", []string{"some_bag"})
	fc.Add("bagit_profile_ideitifier", []string{constants.BTRProfileIdentifier})
	fc.Add("created_at__gteq", []string{"2020-01-01"})
	fc.Add("created_at__lteq", []string{"2024-12-31"})
	fc.Add("etag", []string{"12345"})
	fc.Add("file_count__gteq", []string{"10"})
	fc.Add("file_count__lteq", []string{"100"})
	fc.Add("identifier", []string{"test.edu/some_bag"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("institution__parent_id", []string{"3"})
	fc.Add("internal_sender_identifier", []string{"8675309"})
	fc.Add("size__gteq", []string{"50"})
	fc.Add("size__lteq", []string{"500"})
	fc.Add("source_organization", []string{"test.edu"})
	fc.Add("state", []string{constants.StateActive})
	fc.Add("storage_option", []string{constants.StorageOptionStandard})
	fc.Add("updated_at__gteq", []string{"2020-01-01"})
	fc.Add("updated_at__lteq", []string{"2024-12-31"})
	return fc
}

func getObjectFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getObjectFilters()
	form, err := forms.NewObjectFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestObjectFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getObjectFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testObjectFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestObjectFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getObjectFilterForm(t, user)
		fields := form.GetFields()
		testObjectFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testObjectFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	filters := []string{
		"access",
		"alt_identifier__starts_with",
		"bag_group_identifier__starts_with",
		"bag_name",
		"bagit_profile_identifier",
		"created_at__lteq",
		"created_at__gteq",
		"etag",
		"file_count__lteq",
		"file_count__gteq",
		//"institution_id", // admin only
		"institution_parent_id",
		"internal_sender_identifier",
		"size__lteq",
		"size__gteq",
		"source_organization",
		"state",
		"storage_option",
		"updated_at__lteq",
		"updated_at__gteq",
	}
	for _, filter := range filters {
		require.NotNil(t, fields[filter], filter)
		assert.Equal(t, fc.ValueOf(filter), fields[filter].Value, filter)
	}
}
