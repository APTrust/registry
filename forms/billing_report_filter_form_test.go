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

func getBillingReportFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("start_date", []string{"2024-01-01"})
	fc.Add("end_date", []string{"2024-12-01"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("storage_option", []string{constants.StorageOptionGlacierVA, constants.StorageOptionStandard})
	return fc
}

func getBillingReportFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getBillingReportFilters()
	form, err := forms.NewBillingReportFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestBillingReportFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getBillingReportFilterForm(t, sysAdmin)
	form.SetValues()
	fields := form.GetFields()
	testBillingReportFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestBillingReportFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getBillingReportFilterForm(t, user)
		form.SetValues()
		fields := form.GetFields()
		testBillingReportFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testBillingReportFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, "2024-12-01", fields["end_date"].Value)
	assert.Equal(t, "2024-01-01", fields["start_date"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
}
