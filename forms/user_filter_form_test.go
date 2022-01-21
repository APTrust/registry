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

func getUserFilters() *pgmodels.FilterCollection {
	fc := pgmodels.NewFilterCollection()
	fc.Add("email__contains", []string{"simpsons.com"})
	fc.Add("name__contains", []string{"Homer"})
	fc.Add("institution_id", []string{"2"})
	fc.Add("role", []string{constants.RoleInstAdmin})
	return fc
}

func getUserFilterForm(t *testing.T, user *pgmodels.User) (*pgmodels.FilterCollection, forms.FilterForm) {
	fc := getUserFilters()
	form, err := forms.NewUserFilterForm(fc, user)
	require.Nil(t, err)
	require.NotNil(t, form)
	return fc, form
}

func TestUserFilterFormSysAdmin(t *testing.T) {
	sysAdmin := testutil.InitUser(t, "system@aptrust.org")
	fc, form := getUserFilterForm(t, sysAdmin)
	fields := form.GetFields()
	testUserFields(t, fc, fields)
	assert.True(t, len(fields["institution_id"].Options) > 1)
}

func TestUserFilterFormNonAdmin(t *testing.T) {
	nonSysAdmins := []string{
		"admin@inst1.edu",
		"user@inst1.edu",
	}
	for _, email := range nonSysAdmins {
		user := testutil.InitUser(t, email)
		fc, form := getUserFilterForm(t, user)
		fields := form.GetFields()
		testUserFields(t, fc, fields)
		assert.Empty(t, fields["institution_id"].Options)
	}
}

func testUserFields(t *testing.T, fc *pgmodels.FilterCollection, fields map[string]*forms.Field) {
	assert.Equal(t, fc.ValueOf("email__contains"), fields["email__contains"].Value)
	assert.Equal(t, fc.ValueOf("name__contains"), fields["name__contains"].Value)
	assert.Equal(t, fc.ValueOf("institution_id"), fields["institution_id"].Value)
	assert.Equal(t, fc.ValueOf("role"), fields["role"].Value)
}
