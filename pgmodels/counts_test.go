package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanCountFromView(t *testing.T) {
	db.LoadFixtures()

	q := pgmodels.NewQuery()
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFile{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFileView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObject{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObjectView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.PremisEvent{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.PremisEventView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.WorkItem{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.WorkItemView{}))

	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.User{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.UserView{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.Institution{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.InstitutionView{}))

	// Institution ID is ok for all our count views,
	// because it's present in all.
	q.Where("institution_id", "=", 3)
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFile{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFileView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObject{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObjectView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.PremisEvent{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.PremisEventView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.WorkItem{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.WorkItemView{}))

	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.User{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.UserView{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.Institution{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.InstitutionView{}))

	// State is present only in the Object and File count views
	q.Where("state", "=", "A")
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFile{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.GenericFileView{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObject{}))
	assert.True(t, pgmodels.CanCountFromView(q, pgmodels.IntellectualObjectView{}))

	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.PremisEvent{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.PremisEventView{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.WorkItem{}))
	assert.False(t, pgmodels.CanCountFromView(q, pgmodels.WorkItemView{}))

	// Action is present only in the WorkItem count view
	q2 := pgmodels.NewQuery().Where("institution_id", "=", 3)
	q2.Where("action", "=", constants.ActionIngest)
	assert.True(t, pgmodels.CanCountFromView(q2, pgmodels.WorkItem{}))
	assert.True(t, pgmodels.CanCountFromView(q2, pgmodels.WorkItemView{}))

	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.GenericFile{}))
	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.GenericFileView{}))
	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.IntellectualObject{}))
	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.IntellectualObjectView{}))
	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.PremisEvent{}))
	assert.False(t, pgmodels.CanCountFromView(q2, pgmodels.PremisEventView{}))

	// EventType and Outcome are present only in the PremisEvent count view
	q3 := pgmodels.NewQuery().Where("institution_id", "=", 3)
	q3.Where("event_type", "=", constants.EventIngestion)
	q3.Where("outcome", "=", constants.OutcomeSuccess)
	assert.True(t, pgmodels.CanCountFromView(q3, pgmodels.PremisEvent{}))
	assert.True(t, pgmodels.CanCountFromView(q3, pgmodels.PremisEventView{}))

	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.GenericFile{}))
	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.GenericFileView{}))
	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.IntellectualObject{}))
	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.IntellectualObjectView{}))
	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.WorkItem{}))
	assert.False(t, pgmodels.CanCountFromView(q3, pgmodels.WorkItemView{}))

	// No count views support this column.
	q4 := pgmodels.NewQuery().Where("storage_option", "=", constants.StorageOptionStandard)
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.GenericFile{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.GenericFileView{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.IntellectualObject{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.IntellectualObjectView{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.PremisEvent{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.PremisEventView{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.WorkItem{}))
	assert.False(t, pgmodels.CanCountFromView(q4, pgmodels.WorkItemView{}))
}

func TestGetCountFromView(t *testing.T) {
	db.ForceFixtureReload()

	q := pgmodels.NewQuery()
	count, err := pgmodels.GetCountFromView(q, pgmodels.GenericFile{})
	require.Nil(t, err)
	assert.EqualValues(t, 62, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.GenericFileView{})
	require.Nil(t, err)
	assert.EqualValues(t, 62, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObject{})
	require.Nil(t, err)
	assert.EqualValues(t, 14, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObjectView{})
	require.Nil(t, err)
	assert.EqualValues(t, 14, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.PremisEvent{})
	require.Nil(t, err)
	assert.EqualValues(t, 54, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.PremisEventView{})
	require.Nil(t, err)
	assert.EqualValues(t, 54, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.WorkItem{})
	require.Nil(t, err)
	assert.EqualValues(t, 32, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.WorkItemView{})
	require.Nil(t, err)
	assert.EqualValues(t, 32, count)

	q.Where("institution_id", "=", 3)
	count, err = pgmodels.GetCountFromView(q, pgmodels.GenericFile{})
	require.Nil(t, err)
	assert.EqualValues(t, 37, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.GenericFileView{})
	require.Nil(t, err)
	assert.EqualValues(t, 37, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObject{})
	require.Nil(t, err)
	assert.EqualValues(t, 6, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObjectView{})
	require.Nil(t, err)
	assert.EqualValues(t, 6, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.PremisEvent{})
	require.Nil(t, err)
	assert.EqualValues(t, 27, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.PremisEventView{})
	require.Nil(t, err)
	assert.EqualValues(t, 27, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.WorkItem{})
	require.Nil(t, err)
	assert.EqualValues(t, 15, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.WorkItemView{})
	require.Nil(t, err)
	assert.EqualValues(t, 15, count)

	q.Where("state", "=", "A")
	count, err = pgmodels.GetCountFromView(q, pgmodels.GenericFile{})
	require.Nil(t, err)
	assert.EqualValues(t, 35, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.GenericFileView{})
	require.Nil(t, err)
	assert.EqualValues(t, 35, count)

	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObject{})
	require.Nil(t, err)
	assert.EqualValues(t, 5, count)
	count, err = pgmodels.GetCountFromView(q, pgmodels.IntellectualObjectView{})
	require.Nil(t, err)
	assert.EqualValues(t, 5, count)

	// PremisEvents and WorkItems don't have the "state" column.
	_, err = pgmodels.GetCountFromView(q, pgmodels.PremisEvent{})
	assert.NotNil(t, err)
	_, err = pgmodels.GetCountFromView(q, pgmodels.PremisEventView{})
	assert.NotNil(t, err)
	_, err = pgmodels.GetCountFromView(q, pgmodels.WorkItem{})
	assert.NotNil(t, err)
	_, err = pgmodels.GetCountFromView(q, pgmodels.WorkItemView{})
	assert.NotNil(t, err)

	q2 := pgmodels.NewQuery().Where("institution_id", "=", 2).Where("action", "=", constants.ActionIngest)
	count, err = pgmodels.GetCountFromView(q2, pgmodels.WorkItem{})
	require.Nil(t, err)
	assert.EqualValues(t, 14, count)
	count, err = pgmodels.GetCountFromView(q2, pgmodels.WorkItemView{})
	require.Nil(t, err)
	assert.EqualValues(t, 14, count)

	// GenericFile doesn't have the "action" column.
	_, err = pgmodels.GetCountFromView(q2, pgmodels.GenericFile{})
	assert.NotNil(t, err)
	_, err = pgmodels.GetCountFromView(q2, pgmodels.GenericFileView{})
	assert.NotNil(t, err)

	// Check some PremisEvent filters.
	q3 := pgmodels.NewQuery().Where("institution_id", "=", 2).Where("event_type", "=", constants.EventIngestion)
	count, err = pgmodels.GetCountFromView(q3, pgmodels.PremisEvent{})
	require.Nil(t, err)
	assert.EqualValues(t, 15, count)
	count, err = pgmodels.GetCountFromView(q3, pgmodels.PremisEventView{})
	require.Nil(t, err)
	assert.EqualValues(t, 15, count)

	q3.Where("outcome", "=", constants.OutcomeFailure)
	count, err = pgmodels.GetCountFromView(q3, pgmodels.PremisEvent{})
	require.Nil(t, err)
	assert.EqualValues(t, 0, count)
	count, err = pgmodels.GetCountFromView(q3, pgmodels.PremisEventView{})
	require.Nil(t, err)
	assert.EqualValues(t, 0, count)
}
