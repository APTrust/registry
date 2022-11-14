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

func getDepositReportFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("updated_at__lteq", []string{"2024-12-31"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("storage_option", []string{constants.StorageOptionGlacierVA, constants.StorageOptionStandard})
	fc.Add("chart_metric", []string{"total_bytes"})
	return fc
}

func getDepositReportFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getDepositReportFilters()
	form, err := forms.NewDepositReportFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestDepositReportFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getDepositReportFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testDepositReportFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestDepositReportFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getDepositReportFilterForm(t, user)
		fields := form.GetFields()
		testDepositReportFields(t, fc, fields)
		// Non sysadmin can see only their own alerts, so
		// there should be no filter options in these lists.
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testDepositReportFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("end_date"), fields["end_date"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("storage_option"), fields["storage_option"].Value)
	assert.Equal(t, fc.ValueOf("chart_metric"), fields["chart_metric"].Value)
}
