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

func getInstitutionFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("name__contains", []string{"test.edu"})
	fc.Add("type", []string{constants.InstTypeMember})
	return fc
}

func getInstitutionFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getInstitutionFilters()
	form, err := forms.NewInstitutionFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestInstitutionFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getInstitutionFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testInstitutionFields(t, fc, fields)
	assert.True(t, len(fields["type"].Options) > 1)
}

func testInstitutionFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("name__contains"), fields["name__contains"].Value)
	assert.Equal(t, fc.ValueOf("type"), fields["type"].Value)
}
