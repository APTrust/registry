package admin_api_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	admin_api "github.com/APTrust/registry/web/api/admin"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFailedFixityAlerts(t *testing.T) {
	os.Setenv("APT_ENV", "test")
	db.ForceFixtureReload()

	tu.InitHTTPTests(t)
	//defer db.ForceFixtureReload()

	// There should be only one failed fixity alert when we start.
	// This is part of the fixtures in db/fixtures.
	query := pgmodels.NewQuery().
		Where("subject", "=", "Failed Fixity Check").
		Where("institution_id", "=", 1)
	alerts, err := pgmodels.AlertSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 0, len(alerts))

	// Create three failed fixity events at institution 2
	events, err := createFailedFixityEvents(3, 2)
	require.Nil(t, err)
	require.Equal(t, 3, len(events))

	// And three more at institution 3
	events, err = createFailedFixityEvents(3, 3)
	require.Nil(t, err)
	require.Equal(t, 3, len(events))

	resp := tu.SysAdminClient.POST("/admin-api/v3/alerts/generate_failed_fixity_alerts").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusCreated)

	// We should now have three failed fixity alerts.
	// One is for institution one.
	// One for inst two.
	// One is for APTrust admins.
	// The created_at filter filters out the
	// failed fixity alert that was loaded as
	// part of the test filters.
	query = pgmodels.NewQuery().
		Where("subject", "=", "Failed Fixity Check").
		Where("created_at", ">", "2025-06-01")
	alerts, err = pgmodels.AlertSelect(query)
	for _, a := range alerts {
		fmt.Print(a.Users)
	}
	require.Nil(t, err)
	assert.Equal(t, 3, len(alerts))
}

func TestFailedFixityGenForbiddenToNonAdmins(t *testing.T) {
	// Non sys-admin cannot call this endpoint
	tu.Inst1AdminClient.POST("/admin-api/v3/alerts/generate_failed_fixity_alerts").
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/alerts/generate_failed_fixity_alerts").
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.POST("/admin-api/v3/alerts/generate_failed_fixity_alerts").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2UserClient.POST("/admin-api/v3/alerts/generate_failed_fixity_alerts").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
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
	db.ForceFixtureReload()

	// There should be no failed fixity alerts when we start.
	query := pgmodels.NewQuery().
		Where("subject", "=", "Failed Fixity Check").
		Where("institution_id", "=", 3)
	alerts, err := pgmodels.AlertSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 0, len(alerts))

	// Create three failed fixity events at institution 3
	events, err := createFailedFixityEvents(3, 3)
	require.Nil(t, err)
	require.Equal(t, 3, len(events))

	// Get the summaries
	now := time.Now().UTC()
	lastRunDate := now.AddDate(0, -1, 0)
	summaries, err := pgmodels.FailedFixitySummarySelect(lastRunDate, now)
	require.Nil(t, err)
	assert.Equal(t, 1, len(summaries))

	// Generate the alert to institutional admins
	err = admin_api.GenerateFailedFixityAlert(summaries[0], lastRunDate)
	require.Nil(t, err)

	// Check the db to make sure that alert is present.
	// Note that only one alert goes to each admin at
	// each institution. The alert links to the PremisEvents
	// page with pre-set filters, so the admin can see the
	// list of failed items.
	alerts, err = pgmodels.AlertSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 1, len(alerts))
}

func TestAlertAPTrustOfFailedFixities(t *testing.T) {
	db.ForceFixtureReload()

	// There should be no failed fixity alerts for APTrust
	// admins when we start.
	query := pgmodels.NewQuery().
		Where("subject", "=", "Failed Fixity Check").
		Where("institution_id", "=", 1)
	alerts, err := pgmodels.AlertSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 0, len(alerts))

	// Create three failed fixity events at institution 2
	events, err := createFailedFixityEvents(3, 2)
	require.Nil(t, err)
	require.Equal(t, 3, len(events))

	// And three more at institution 3
	events, err = createFailedFixityEvents(3, 3)
	require.Nil(t, err)
	require.Equal(t, 3, len(events))

	// Get the summaries
	now := time.Now().UTC()
	lastRunDate := now.AddDate(0, -1, 0)
	summaries, err := pgmodels.FailedFixitySummarySelect(lastRunDate, now)
	require.Nil(t, err)
	assert.Equal(t, 2, len(summaries))

	// Generate the alert to institutional admins
	err = admin_api.AlertAPTrustOfFailedFixities(summaries, lastRunDate)
	require.Nil(t, err)

	// Check the db to make sure that alerts are present.
	// There should be two, one pertaining to each institution.
	alerts, err = pgmodels.AlertSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 1, len(alerts))

}

func TestGetFailedFixityEvents(t *testing.T) {
	db.ForceFixtureReload()
	lastRunDate, err := time.Parse("2006-01-02", "2020-01-01")
	require.Nil(t, err)
	events, err := admin_api.GetFailedFixityEvents(3, lastRunDate)
	require.Nil(t, err)
	assert.Equal(t, 0, len(events))

	// Add a few failed fixity events and check again
	_, err = createFailedFixityEvents(3, 3)
	require.Nil(t, err)

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

func createFailedFixityEvents(howMany int, institutionID int64) ([]*pgmodels.PremisEvent, error) {
	events := make([]*pgmodels.PremisEvent, howMany)
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	for i := 0; i < 3; i++ {
		event := pgmodels.RandomPremisEvent(constants.EventFixityCheck)
		event.InstitutionID = institutionID
		event.DateTime = yesterday
		event.Outcome = "Failed"
		err := event.Save()
		if err != nil {
			return nil, err
		}
		events[i] = event
	}
	return events, nil
}
