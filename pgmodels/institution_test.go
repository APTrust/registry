package pgmodels_test

import (
	//"fmt"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionValidation(t *testing.T) {
	inst := &pgmodels.Institution{}
	err := inst.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrInstName, err.Errors["Name"])
	assert.Equal(t, pgmodels.ErrInstIdentifier, err.Errors["Identifier"])
	assert.Equal(t, pgmodels.ErrInstState, err.Errors["State"])
	assert.Equal(t, pgmodels.ErrInstType, err.Errors["Type"])
	assert.Equal(t, pgmodels.ErrInstReceiving, err.Errors["ReceivingBucket"])
	assert.Equal(t, pgmodels.ErrInstRestore, err.Errors["RestoreBucket"])
	// There should be no error here, because only a subscribing
	// institution requires a parent MemberInstitutionID, and our
	// test model doesn't specify any type.
	assert.Equal(t, "", err.Errors["MemberInstitutionID"])

	inst.Name = "x"
	inst.Identifier = "invalid"
	inst.State = "x"
	inst.Type = constants.InstTypeSubscriber
	err = inst.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrInstName, err.Errors["Name"])
	assert.Equal(t, pgmodels.ErrInstIdentifier, err.Errors["Identifier"])
	assert.Equal(t, pgmodels.ErrInstState, err.Errors["State"])
	assert.Equal(t, pgmodels.ErrInstReceiving, err.Errors["ReceivingBucket"])
	assert.Equal(t, pgmodels.ErrInstRestore, err.Errors["RestoreBucket"])
	// This time, because inst.Type is subscriber, we should get
	// the MemberInstitutionID error.
	assert.Equal(t, "", err.Errors["Type"])
	assert.Equal(t, pgmodels.ErrInstMemberID, err.Errors["MemberInstitutionID"])

	// Now let's make a valid record
	inst.Name = "Valid Institution"
	inst.Identifier = "library.valid.edu"
	inst.State = constants.StateActive
	inst.MemberInstitutionID = int64(33)
	inst.ReceivingBucket = "aptrust.receiving.test.library.valid.edu"
	inst.RestoreBucket = "aptrust.restore.test.library.valid.edu"

	err = inst.Validate()
	require.Nil(t, err)
}
