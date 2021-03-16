package pgmodels_test

import (
	//"fmt"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionValidation(t *testing.T) {
	inst := &pgmodels.Institution{}
	errs := inst.Validate()
	require.NotNil(t, errs)
	valErrors := errs.(v.Errors)

	assert.Equal(t, pgmodels.ErrInstName, valErrors["name"].Error())
	assert.Equal(t, pgmodels.ErrInstIdentifier, valErrors["identifier"].Error())
	assert.Equal(t, pgmodels.ErrInstState, valErrors["state"].Error())
	assert.Equal(t, pgmodels.ErrInstType, valErrors["type"].Error())
	assert.Equal(t, pgmodels.ErrInstReceiving, valErrors["receiving_bucket"].Error())
	assert.Equal(t, pgmodels.ErrInstRestore, valErrors["restore_bucket"].Error())
	// There should be no error here, because only a subscribing
	// institution requires a parent MemberInstitutionID, and our
	// test model doesn't specify any type.
	assert.Nil(t, valErrors["member_institution_id"])

	inst.Name = "x"
	inst.Identifier = "invalid"
	inst.State = "x"
	inst.Type = constants.InstTypeSubscriber
	errs = inst.Validate()
	require.NotNil(t, errs)
	valErrors = errs.(v.Errors)

	assert.Equal(t, pgmodels.ErrInstName, valErrors["name"].Error())
	assert.Equal(t, pgmodels.ErrInstIdentifier, valErrors["identifier"].Error())
	assert.Equal(t, pgmodels.ErrInstState, valErrors["state"].Error())
	assert.Equal(t, pgmodels.ErrInstReceiving, valErrors["receiving_bucket"].Error())
	assert.Equal(t, pgmodels.ErrInstRestore, valErrors["restore_bucket"].Error())
	// This time, because inst.Type is subscriber, we should get
	// the MemberInstitutionID error.
	assert.Nil(t, valErrors["type"])
	assert.Equal(t, pgmodels.ErrInstMemberID, valErrors["member_institution_id"].Error())

	// Now let's make a valid record
	inst.Name = "Valid Institution"
	inst.Identifier = "library.valid.edu"
	inst.State = constants.StateActive
	inst.MemberInstitutionID = int64(33)
	inst.ReceivingBucket = "aptrust.receiving.test.library.valid.edu"
	inst.RestoreBucket = "aptrust.restore.test.library.valid.edu"

	errs = inst.Validate()
	require.Nil(t, errs)
}
