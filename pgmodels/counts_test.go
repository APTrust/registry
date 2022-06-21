package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
)

func TestCanCountFromView(t *testing.T) {
	db.ForceFixtureReload()

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
