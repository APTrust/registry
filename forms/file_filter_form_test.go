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

func getFileFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("created_at__gteq", []string{"2020-01-01"})
	fc.Add("created_at__lteq", []string{"2024-12-31"})
	fc.Add("updated_at__gteq", []string{"2020-01-01"})
	fc.Add("updated_at__lteq", []string{"2024-12-31"})
	fc.Add("identifier", []string{"test.edu/object/data/file.jpg"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("state", []string{constants.StateActive})
	fc.Add("size__gteq", []string{"100"})
	fc.Add("size__lteq", []string{"500"})
	fc.Add("storage_option", []string{constants.StorageOptionGlacierDeepOR})
	return fc
}

func getFileFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getFileFilters()
	form, err := forms.NewFileFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestFileFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getFileFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testFileFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestFileFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getFileFilterForm(t, user)
		fields := form.GetFields()
		testFileFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testFileFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("created_at__gteq"), fields["created_at__gteq"].Value)
	assert.Equal(t, fc.ValueOf("created_at__lteq"), fields["created_at__lteq"].Value)
	assert.Equal(t, fc.ValueOf("updated_at__gteq"), fields["updated_at__gteq"].Value)
	assert.Equal(t, fc.ValueOf("updated_at__lteq"), fields["updated_at__lteq"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("identifier"), fields["identifier"].Value)
	assert.Equal(t, fc.ValueOf("state"), fields["state"].Value)
	assert.Equal(t, fc.ValueOf("size__gteq"), fields["size__gteq"].Value)
	assert.Equal(t, fc.ValueOf("size__lteq"), fields["size__lteq"].Value)
	assert.Equal(t, fc.ValueOf("storage_option"), fields["storage_option"].Value)
}
