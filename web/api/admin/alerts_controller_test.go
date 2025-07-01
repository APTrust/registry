package admin_api_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	admin_api "github.com/APTrust/registry/web/api/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFailedFixityAlerts(t *testing.T) {

}

func TestFailedFixityLastRunDate(t *testing.T) {
	db.LoadFixtures()

	// This should return an error because ar_internal_metadata
	// doesn't have an entry yet for last fixity alert date.
	_, err := pgmodels.InternalMetadataByKey(constants.MetaFixityAlertsLastRun)
	require.NotNil(t, err, "Last fixity date should not be set yet.")

	// If no timestamp is set in the database, this function
	// should return yesterday's date.
	yesterdayTs := time.Now().UTC().AddDate(0, 0, -1)
	yesterday := time.Date(yesterdayTs.Year(), yesterdayTs.Month(), yesterdayTs.Day(), 0, 0, 0, 0, time.UTC)
	lastRun, err := admin_api.FailedFixityLastRunDate()
	require.Nil(t, err)
	assert.Equal(t, yesterday, lastRun)

	// Now set the value explicitly in the database, and the
	// function should return that date.
	fiveDaysAgo := yesterday.AddDate(0, 0, -4)
	metadata := pgmodels.NewInteralMetadata(constants.MetaFixityAlertsLastRun, fiveDaysAgo.Format(time.RFC3339))
	err = metadata.Save()
	require.Nil(t, err)

	lastRun, err = admin_api.FailedFixityLastRunDate()
	require.Nil(t, err)
	assert.Equal(t, fiveDaysAgo, lastRun)
}

func TestGenerateFailedFixityAlert(t *testing.T) {

}

func TestAlertAPTrustOfFailedFixities(t *testing.T) {

}

func TestGetFailedFixityEvents(t *testing.T) {
	db.LoadFixtures()
	lastRunDate, err := time.Parse("2006-01-02", "2020-01-01")
	require.Nil(t, err)
	events, err := admin_api.GetFailedFixityEvents(3, lastRunDate)
	require.Nil(t, err)
	assert.Equal(t, 0, len(events))

	// Add a few failed fixity events and check again
	tenDaysAgo := time.Now().UTC().AddDate(0, 0, -10)
	for i := 0; i < 3; i++ {
		event := pgmodels.RandomPremisEvent(constants.EventFixityCheck)
		event.InstitutionID = 3
		event.DateTime = tenDaysAgo
		event.Outcome = "Failed"
		event.Save()
	}

	events, err = admin_api.GetFailedFixityEvents(3, lastRunDate)
	require.Nil(t, err)
	assert.Equal(t, 3, len(events))
}

func TestFailedFixityReportURL(t *testing.T) {
	startDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)

	// Expected URL for institutions other than APTrust.
	// These should contain an institution ID filter.
	expectedURL := "http://localhost:8080/events?event_type=fixity+check&outcome=Failed&institution_id=4&date_time__gteq=2025-07-01&date_time__lteq=2025-06-30"
	url := admin_api.FailedFixityReportURL(4, endDate, startDate)
	assert.Equal(t, expectedURL, url)

	// Expected URL for APTrust.
	// This should not have an institution ID filter.
	expectedURL = "http://localhost:8080/events?event_type=fixity+check&outcome=Failed&date_time__gteq=2025-07-01&date_time__lteq=2025-06-30"
	url = admin_api.FailedFixityReportURL(0, endDate, startDate)
	assert.Equal(t, expectedURL, url)
}
