package forms_test

import (
	"testing"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionForm(t *testing.T) {
	// Test with new/empty institution
	inst := &pgmodels.Institution{}
	form, err := forms.NewInstitutionForm(inst)
	require.Nil(t, err)
	require.NotNil(t, form)
	assert.True(t, len(form.Fields["Type"].Options) > 1)
	assert.Equal(t, "", form.Fields["Type"].Value)
	assert.True(t, len(form.Fields["MemberInstitutionID"].Options) > 2)
	assert.Equal(t, int64(0), form.Fields["MemberInstitutionID"].Value)

	// Test with an existing institution
	inst, err = pgmodels.InstitutionByIdentifier("test.edu")
	require.Nil(t, err)
	require.NotNil(t, inst)
	form, err = forms.NewInstitutionForm(inst)
	require.Nil(t, err)
	require.NotNil(t, form)

	assert.Equal(t, inst.Name, form.Fields["Name"].Value)
	assert.Equal(t, inst.Identifier, form.Fields["Identifier"].Value)
	assert.Equal(t, inst.Type, form.Fields["Type"].Value)
	assert.Equal(t, inst.MemberInstitutionID, form.Fields["MemberInstitutionID"].Value)
	assert.Equal(t, inst.OTPEnabled, form.Fields["OTPEnabled"].Value)
	assert.Equal(t, inst.SpotRestoreFrequency, form.Fields["SpotRestoreFrequency"].Value)
	assert.Equal(t, inst.ReceivingBucket, form.Fields["ReceivingBucket"].Value)
	assert.Equal(t, inst.RestoreBucket, form.Fields["RestoreBucket"].Value)
}
