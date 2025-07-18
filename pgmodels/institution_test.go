package pgmodels_test

import (
	"fmt"
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
	inst.Type = constants.InstTypeMember
	inst.SpotRestoreFrequency = 90
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
	query.OrderBy("name", "asc")
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

func TestIdForInstIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForInstIdentifier("institution1.edu")
	require.Nil(t, err)
	assert.Equal(t, int64(2), id)

	id, err = pgmodels.IdForInstIdentifier("institution2.edu")
	require.Nil(t, err)
	assert.Equal(t, int64(3), id)

	id, err = pgmodels.IdForInstIdentifier("bad identifier")
	require.NotNil(t, err)
}

func TestInstitutionGetAssociateMembers(t *testing.T) {
	db.LoadFixtures()
	defer db.ForceFixtureReload()

	inst, err := pgmodels.InstitutionByIdentifier("institution1.edu")
	require.Nil(t, err)
	require.NotNil(t, inst)

	for i := 0; i < 5; i++ {
		inst := pgmodels.Institution{
			Name:                fmt.Sprintf("SubAccount %d", i),
			Identifier:          fmt.Sprintf("subacct%d.edu", i),
			State:               constants.StateActive,
			Type:                constants.InstTypeSubscriber,
			MemberInstitutionID: inst.ID,
			ReceivingBucket:     fmt.Sprintf("aptrust.receiving.text.subacct%d.edu", i),
			RestoreBucket:       fmt.Sprintf("aptrust.restore.text.subacct%d.edu", i),
		}
		require.Nil(t, inst.Save())
	}

	// Should be 5 active sub accounts
	associateAccounts, err := inst.GetAssociateMembers()
	require.Nil(t, err)
	assert.Equal(t, 5, len(associateAccounts))
	for i, acct := range associateAccounts {
		expectedName := fmt.Sprintf("SubAccount %d", i)
		assert.Equal(t, expectedName, acct.Name)
		if i == 2 {
			acct.State = constants.StateDeleted
			require.Nil(t, acct.Save())
		}
	}

	// Should be 4 active sub accounts
	associateAccounts, err = inst.GetAssociateMembers()
	require.Nil(t, err)
	assert.Equal(t, 4, len(associateAccounts))
}

func TestInstitutionHasSubAccounts(t *testing.T) {
	db.LoadFixtures()
	defer db.ForceFixtureReload()

	inst1, err := pgmodels.InstitutionByIdentifier("institution1.edu")
	require.Nil(t, err)
	require.NotNil(t, inst1)
	hasKids, err := inst1.HasSubAccounts()
	require.Nil(t, err)
	assert.False(t, hasKids)

	// Temporarily make Inst2 a subaccount of Inst1
	// We'll reload fixtures at the end of this test
	// so we don't mess up other tests.
	inst2, err := pgmodels.InstitutionByIdentifier("institution2.edu")
	require.Nil(t, err)
	require.NotNil(t, inst2)
	inst2.MemberInstitutionID = inst1.ID
	require.Nil(t, inst2.Save())

	// Now this should be true, because inst2 is a sub-account
	// of inst1.
	hasKids, err = inst1.HasSubAccounts()
	require.Nil(t, err)
	assert.True(t, hasKids)
}

func TestDueForSpotRestore(t *testing.T) {
	db.LoadFixtures()

	// We'll fiddle with this object but the changes
	// won't stick because we're not saving them to
	// the DB.
	inst1, err := pgmodels.InstitutionByIdentifier("institution1.edu")
	require.Nil(t, err)
	require.NotNil(t, inst1)

	// Should be true because this inst is configured for
	// spot tests every 30 days and inst has no last
	// spot test work item id.
	dueForRestore, err := inst1.DueForSpotRestore()
	require.Nil(t, err)
	assert.Equal(t, int64(30), inst1.SpotRestoreFrequency)
	assert.Empty(t, inst1.LastSpotRestoreWorkItemID)
	assert.True(t, dueForRestore)

	// Should be false if we change SpotRestoreFrequency to zero.
	inst1.SpotRestoreFrequency = 0
	dueForRestore, err = inst1.DueForSpotRestore()
	require.Nil(t, err)
	assert.False(t, dueForRestore)

	// Should be true, because frequency is 30 days and
	// WorkItem 31 was a restoration from August, 2016
	// (it's loaded from the fixtures csv file).
	inst1.SpotRestoreFrequency = 30
	inst1.LastSpotRestoreWorkItemID = 31
	dueForRestore, err = inst1.DueForSpotRestore()
	require.Nil(t, err)
	assert.True(t, dueForRestore)

	// Should be false, because frequency is once every
	// 100 years, and the 2016 restoration was less than
	// 100 years ago
	inst1.SpotRestoreFrequency = 36500
	dueForRestore, err = inst1.DueForSpotRestore()
	require.Nil(t, err)
	assert.False(t, dueForRestore)
}

func TestInstDisplayType(t *testing.T) {
	inst := &pgmodels.Institution{
		Type: constants.InstTypeMember,
	}
	assert.Equal(t, "Member", inst.DisplayType())

	inst.Type = constants.InstTypeSubscriber
	assert.Equal(t, "Associate", inst.DisplayType())
}

func TestGetAdmins(t *testing.T) {
	db.LoadFixtures()

	query := pgmodels.NewQuery().Where("name", "=", "Institution One")
	inst, err := pgmodels.InstitutionGet(query)
	require.Nil(t, err)
	require.NotNil(t, inst)

	admins, err := inst.GetAdmins()
	require.Nil(t, err)
	require.Equal(t, 1, len(admins))
	assert.Equal(t, "admin@inst1.edu", admins[0].Email)

	query = pgmodels.NewQuery().Where("name", "=", "Institution Two")
	inst, err = pgmodels.InstitutionGet(query)
	require.Nil(t, err)
	require.NotNil(t, inst)

	admins, err = inst.GetAdmins()
	require.Nil(t, err)
	require.Equal(t, 1, len(admins))
	assert.Equal(t, "admin@inst2.edu", admins[0].Email)
}
