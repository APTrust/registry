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

func getAlertFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("created_at__gteq", []string{"2020-01-01"})
	fc.Add("created_at__lteq", []string{"2024-12-31"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("type", []string{constants.AlertDeletionRequested, constants.AlertDeletionConfirmed})
	fc.Add("user_id", []string{"4"})
	return fc
}

func getAlertFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getAlertFilters()
	form, err := forms.NewAlertFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestAlertFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getAlertFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testAlertFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
	assert.True(t, len(fields["user_id"].Options) > 1)
}

func TestAlertFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getAlertFilterForm(t, user)
		fields := form.GetFields()
		testAlertFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
		assert.Empty(t, fields["user_id"].Options)
	}
}

func testAlertFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("created_at__gteq"), fields["created_at__gteq"].Value)
	assert.Equal(t, fc.ValueOf("created_at__lteq"), fields["created_at__lteq"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("type"), fields["type"].Value)
	assert.Equal(t, fc.ValueOf("user_id"), fields["user_id"].Value)
}
