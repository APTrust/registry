package pgmodels_test

import (
	"testing"

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
		Subject:       "Order Up!",
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

	testMarkAsSent(t, alert)
	testMarkAsRead(t, alert)
	testMarkAsUnread(t, alert)
}

func TestAlertValidate(t *testing.T) {
	alert := &pgmodels.Alert{}
	err := alert.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrAlertInstitutionID, err.Errors["InstitutionID"])
	assert.Equal(t, pgmodels.ErrAlertType, err.Errors["Type"])
	assert.Equal(t, pgmodels.ErrAlertContent, err.Errors["Content"])
}

func TestAlertSelect(t *testing.T) {
	query := pgmodels.NewQuery().Where("institution_id", "=", 2)
	alerts, err := pgmodels.AlertSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, alerts)
	assert.Equal(t, 9, len(alerts))
}

func testMarkAsSent(t *testing.T, alert *pgmodels.Alert) {
	userID := alert.Users[0].ID

	// First, make sure it's unsent.
	alertView, err := pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.Empty(t, alertView.SentAt)

	// Mark mark the alert as sent
	err = alert.MarkAsSent(userID)
	require.Nil(t, err)

	// Now make sure it really was marked as sent.
	alertView, err = pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.NotEmpty(t, alertView.SentAt)
}

func testMarkAsRead(t *testing.T, alert *pgmodels.Alert) {
	userID := alert.Users[0].ID

	// First, make sure it's unread.
	alertView, err := pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.Empty(t, alertView.ReadAt)

	// Mark mark the alert as read
	err = alert.MarkAsRead(userID)
	require.Nil(t, err)

	// Now make sure it really was marked as read.
	alertView, err = pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.NotEmpty(t, alertView.ReadAt)
}

func testMarkAsUnread(t *testing.T, alert *pgmodels.Alert) {
	userID := alert.Users[0].ID

	// First, make sure it's read.
	alertView, err := pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.NotEmpty(t, alertView.ReadAt)

	// Mark mark the alert as unread
	err = alert.MarkAsUnread(userID)
	require.Nil(t, err)

	// Now make sure it really was marked as unread.
	alertView, err = pgmodels.AlertViewForUser(alert.ID, userID)
	require.Nil(t, err)
	require.NotNil(t, alertView)
	require.Empty(t, alertView.ReadAt)
}

func TestNewFailedFixityAlert(t *testing.T) {
	instID := int64(4)
	premisEvents := []*pgmodels.PremisEvent{
		pgmodels.RandomPremisEvent(constants.EventFixityCheck),
		pgmodels.RandomPremisEvent(constants.EventFixityCheck),
		pgmodels.RandomPremisEvent(constants.EventFixityCheck),
	}
	recipients := []*pgmodels.User{
		{Email: "homer@simpsons.com"},
		{Email: "duff_man@simpsons.com"},
	}

	alert := pgmodels.NewFailedFixityAlert(instID, premisEvents, recipients)

	assert.Equal(t, int64(0), alert.ID) // Has not been saved
	assert.Equal(t, constants.AlertFailedFixity, alert.Type)
	assert.Equal(t, "Failed Fixity Check", alert.Subject)

	require.Equal(t, 3, len(alert.PremisEvents))
	assert.Equal(t, premisEvents[0], alert.PremisEvents[0])
	assert.Equal(t, premisEvents[1], alert.PremisEvents[1])
	assert.Equal(t, premisEvents[2], alert.PremisEvents[2])

	require.Equal(t, 2, len(alert.Users))
	assert.Equal(t, "homer@simpsons.com", alert.Users[0].Email)
	assert.Equal(t, "duff_man@simpsons.com", alert.Users[1].Email)
}
