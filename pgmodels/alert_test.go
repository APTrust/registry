package pgmodels_test

import (
	//"fmt"
	"testing"

	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertRelations(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Limit(3)

	// Get three users
	users, err := pgmodels.UserSelect(query)
	require.Nil(t, err)
	require.NotNil(t, users)

	// And three work items
	items, err := pgmodels.WorkItemSelect(query)
	require.Nil(t, err)
	require.NotNil(t, items)

	// And three premis events
	events, err := pgmodels.PremisEventSelect(query)
	require.Nil(t, err)
	require.NotNil(t, events)

	alert := &pgmodels.Alert{
		InstitutionID: 2,
		Type:          constants.AlertRestorationCompleted,
		Content:       "Your fries are ready.",
		PremisEvents:  events,
		Users:         users,
		WorkItems:     items,
	}

	err = alert.Save()
	require.Nil(t, err)
	assert.True(t, alert.ID > 0)

	// Fetch the alert...
	alert, err = pgmodels.AlertByID(alert.ID)
	require.Nil(t, err)

	assert.Equal(t, 3, len(alert.Users))
	assert.Equal(t, 3, len(alert.PremisEvents))
	assert.Equal(t, 3, len(alert.WorkItems))
}
