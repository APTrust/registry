package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
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

func TestInstitutionByID(t *testing.T) {
	db.LoadFixtures()
	inst, err := pgmodels.InstitutionByID(int64(1))
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, int64(1), inst.ID)
}

func TestInstitutionByIdentifier(t *testing.T) {
	db.LoadFixtures()
	inst, err := pgmodels.InstitutionByIdentifier("test.edu")
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, "test.edu", inst.Identifier)
}

func TestInstitutionGet(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where("name", "=", "Institution One")
	inst, err := pgmodels.InstitutionGet(query)
	require.Nil(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, "Institution One", inst.Name)
}

func TestInstitutionSelect(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery()
	query.Where("name", "!=", "Institution One")
	query.Where("name", "!=", "Institution Two")
	query.OrderBy("name asc")
	institutions, err := pgmodels.InstitutionSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, institutions)
	assert.True(t, (len(institutions) > 0 && len(institutions) < 10))
	for _, inst := range institutions {
		assert.NotEqual(t, "Institution One", inst)
		assert.NotEqual(t, "Institution Two", inst)
	}
}

func TestInstitutionSave(t *testing.T) {
	db.LoadFixtures()
	inst := &pgmodels.Institution{
		Name:       "Unit Test Inst #1",
		Identifier: "test1.kom",
		Type:       constants.InstTypeMember,
	}
	err := inst.Save()
	require.Nil(t, err)

	// pg library should set ID, BeforeInsert hook should set other values
	assert.True(t, inst.ID > int64(0))
	assert.Equal(t, constants.StateActive, inst.State)
	assert.Equal(t, "aptrust.receiving.test.test1.kom", inst.ReceivingBucket)
	assert.Equal(t, "aptrust.restore.test.test1.kom", inst.RestoreBucket)
	assert.NotEmpty(t, inst.CreatedAt)
	assert.NotEmpty(t, inst.UpdatedAt)
}

func TestInstitutionDeleteUndelete(t *testing.T) {
	inst := &pgmodels.Institution{
		Name:       "Unit Test Inst #2",
		Identifier: "test2.kom",
		Type:       constants.InstTypeMember,
	}
	err := inst.Save()
	require.Nil(t, err)

	assert.True(t, inst.ID > int64(0))
	assert.Equal(t, constants.StateActive, inst.State)
	assert.Empty(t, inst.DeactivatedAt)

	err = inst.Delete()
	require.Nil(t, err)
	assert.Equal(t, constants.StateDeleted, inst.State)
	assert.NotEmpty(t, inst.DeactivatedAt)

	err = inst.Undelete()
	require.Nil(t, err)
	assert.Equal(t, constants.StateActive, inst.State)
	assert.Empty(t, inst.DeactivatedAt)
}
