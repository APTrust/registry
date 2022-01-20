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

func getDeletionRequestFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("requested_at__gteq", []string{"2020-01-01"})
	fc.Add("requested_at__lteq", []string{"2024-12-31"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("stage", []string{constants.StageRequested})
	fc.Add("status", []string{constants.StatusPending})
	return fc
}

func getDeletionRequestFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getDeletionRequestFilters()
	form, err := forms.NewDeletionRequestFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestDeletionRequestFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getDeletionRequestFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testDeletionRequestFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
	assert.True(t, len(fields["stage"].Options) > 1)
	assert.True(t, len(fields["status"].Options) > 1)
}

func TestDeletionRequestFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getDeletionRequestFilterForm(t, user)
		fields := form.GetFields()
		testDeletionRequestFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testDeletionRequestFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("requested_at__gteq"), fields["requested_at__gteq"].Value)
	assert.Equal(t, fc.ValueOf("requested_at__lteq"), fields["requested_at__lteq"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("stage"), fields["stage"].Value)
	assert.Equal(t, fc.ValueOf("status"), fields["status"].Value)
}
