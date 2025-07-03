package admin_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// GenerateFailedFixityAlerts generates alerts to institutional
// admins and to APTrust admins describing recent failed fixity
// checks. This is a POST because it may create alert records
// in the database.
//
// POST /admin-api/v3/alerts/generate_failed_fixity_alerts
func GenerateFailedFixityAlerts(c *gin.Context) {
	ctx := common.Context()
	ctx.Log.Info().Msg("Received request to generate failed fixity alerts.")
	type Response struct {
		Error     string `json:"error"`
		Summaries []*pgmodels.FailedFixitySummary
	}
	response := Response{}

	// Find out when this process was last run.
	lastRunDate, err := FailedFixityLastRunDate()
	if err != nil {
		ctx.Log.Error().Msgf("Could not get last run date for failed fixity alerts: %v", err)
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If this has already run today, don't run it again.
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if lastRunDate.Equal(today) {
		ctx.Log.Info().Msgf("Not generating failed fixity alerts because they were last generated on %v", lastRunDate)
		response.Error = "Alerts have already been generated today."
		c.JSON(http.StatusConflict, response)
		return
	}

	summaries, err := pgmodels.FailedFixitySummarySelect(lastRunDate, now)
	if err != nil {
		ctx.Log.Error().Msgf("Error querying for failed fixity alerts: %v", err)
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Summaries = summaries

	// Generate a fixity failure alert for each institution.
	// If there's an error in the alert generation process,
	// add that into the response, but keep processing because
	// we don't want one failed alert to prevent others from
	// being sent.
	for _, summary := range summaries {
		err = GenerateFailedFixityAlert(summary, lastRunDate)
		if err != nil {
			response.Error = fmt.Sprintf("%s; %s", response.Error, err.Error())
		}
	}

	// Generate only one report for APTrust admins. This will
	// include all failures at all institutions, and the embedded
	// link will show them all.
	err = AlertAPTrustOfFailedFixities(summaries, lastRunDate)
	if err != nil {
		ctx.Log.Error().Msgf("Error generating failed fixity alerts for APTrust admins: %v", err)
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if len(response.Summaries) > 0 {
		c.JSON(http.StatusCreated, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// GenerateFailedFixityAlert returns the date on which failed
// fixity checks were last run. The time component of the date
// will always be midnight. If there's no date in the the
// database's ar_internal_metadata table, this returns yesterday's
// date.
func FailedFixityLastRunDate() (time.Time, error) {
	yesterdayTs := time.Now().UTC().AddDate(0, 0, -1)
	yesterday := time.Date(yesterdayTs.Year(), yesterdayTs.Month(), yesterdayTs.Day(), 0, 0, 0, 0, time.UTC)

	metadata, err := pgmodels.InternalMetadataByKey(constants.MetaFixityAlertsLastRun)

	// Now rows in result set means GenerateFailedFixityAlerts has
	// never run before. This should happen only once.
	if pgmodels.IsNoRowError(err) {
		return yesterday, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, metadata.Value)
}

// GenerateFailedFixityAlert creates fixity failure alerts
// for each admin at the institution in summary.InstitutionID.
func GenerateFailedFixityAlert(summary *pgmodels.FailedFixitySummary, lastRunDate time.Time) error {
	ctx := common.Context()
	ctx.Log.Info().Msgf("Generating failed fixity alert for %s with %d failures.", summary.InstitutionName, summary.Failures)

	institution, err := pgmodels.InstitutionByID(summary.InstitutionID)
	if err != nil {
		ctx.Log.Error().Msgf("Error getting institutiton record for institution %d: %v", summary.InstitutionID, err)
		return err
	}

	instAdmins, err := institution.GetAdmins()
	if err != nil {
		ctx.Log.Error().Msgf("Error getting list of admins for %s: %v", institution.Identifier, err)
		return err
	}

	events, err := GetFailedFixityEvents(institution.ID, lastRunDate)
	if err != nil {
		ctx.Log.Error().Msgf("Error collecting list of failed fixity events for %s: %v", institution.Identifier, err)
		return err
	}

	today := time.Now().UTC()
	alert := pgmodels.NewFailedFixityAlert(institution.ID, events, instAdmins)
	alertData := map[string]interface{}{
		"StartDate": lastRunDate.Format("2006-01-02"),
		"EndDate":   today.Format("2006-01-02"),
		"AlertURL":  FailedFixityReportURL(institution.ID, lastRunDate, today),
	}

	alert, err = pgmodels.CreateAlert(alert, "alerts/failed_fixity.txt", alertData)
	if err != nil {
		ctx.Log.Error().Msgf("Error creating failed fixity alert for %s: %v", institution.Identifier, err)
		return err
	}
	if alert != nil {
		ctx.Log.Info().Msgf("Created failed fixity alert for admins at %s.", institution.Identifier)
		return err
	}

	return nil
}

func AlertAPTrustOfFailedFixities(summaries []*pgmodels.FailedFixitySummary, lastRunDate time.Time) error {
	if len(summaries) == 0 {
		return nil
	}
	instCount := 0
	failureCount := 0
	for _, summary := range summaries {
		instCount += 1
		failureCount += int(summary.Failures)
	}
	ctx := common.Context()
	ctx.Log.Info().Msgf("Generating failed fixity alert for APTrust ops/admins showing %d failures at %d institutions.",
		instCount, failureCount)

	events, err := GetFailedFixityEvents(0, lastRunDate)
	if err != nil {
		ctx.Log.Error().Msgf("Error collecting list of failed fixity events for admin/ops email: %v", err)
		return err
	}

	institution, err := pgmodels.InstitutionByIdentifier("aptrust.org")
	if err != nil {
		ctx.Log.Error().Msgf("Error retrieving APTrust institution record for admin/ops email: %v", err)
		return err
	}

	aptrustAdmins, err := institution.GetAdmins()
	if err != nil {
		ctx.Log.Error().Msgf("Error getting list of APTrust admins for admin/ops email: %v", err)
		return err
	}

	today := time.Now().UTC()
	alert := pgmodels.NewFailedFixityAlert(institution.ID, events, aptrustAdmins)
	alertData := map[string]interface{}{
		"StartDate": lastRunDate.Format("2006-01-02"),
		"EndDate":   today.Format("2006-01-02"),
		"AlertURL":  FailedFixityReportURL(0, lastRunDate, today),
	}

	alert, err = pgmodels.CreateAlert(alert, "alerts/failed_fixity.txt", alertData)
	if err != nil {
		ctx.Log.Error().Msgf("Error creating failed fixity alert for admin/ops email: %v", err)
		return err
	}
	if alert != nil {
		ctx.Log.Info().Msg("Created failed fixity alert for APTrust admins.")
		return err
	}

	return nil
}

// GetFailedFixityEvents returns the PremisEvent records for the
// failed fixity events. We want to link these to the Alert records
// in the database, so we know when each admin was alerted to a
// specific failure.
//
// For alerts to APTrust admins, institutionID should be zero,
// because they want to see alerts pertaining to all institutions.
//
// For institutional admins, institutionID should match that of
// the admin's own institution.
func GetFailedFixityEvents(institutionID int64, lastRunDate time.Time) ([]*pgmodels.PremisEvent, error) {
	query := pgmodels.NewQuery().
		Where("event_type", "=", "fixity check").
		Where("outcome", "=", "Failed").
		Where("date_time", ">", lastRunDate)

	// Add inst ID filter for institutional admin...
	if institutionID > 0 {
		query = query.Where("institution_id", "=", institutionID)
	}
	return pgmodels.PremisEventSelect(query)
}

// FailedFixityReportURL returns a link to the Premis Events page
// with pre-applied filters to display failed fixity checks for
// the specified institution between the report's last run date
// and today.
func FailedFixityReportURL(institutionID int64, lastRunDate, currentRunDate time.Time) string {
	var domainAndPort string
	environment := common.Context().Config.EnvName
	switch environment {
	case "prod":
		domainAndPort = "repo.aptrust.org"
	case "demo":
		domainAndPort = "demo.aptrust.org"
	case "staging":
		domainAndPort = "staging.aptrust.org"
	default:
		domainAndPort = "localhost:8080"
	}
	startDateString := lastRunDate.Format("2006-01-02")
	endDateString := currentRunDate.Format("2006-01-02")
	if institutionID == 0 {
		return fmt.Sprintf(
			"%s://%s/events?event_type=fixity+check&outcome=Failed&date_time__gteq=%s&date_time__lteq=%s",
			common.Context().Config.HTTPScheme(), domainAndPort, startDateString, endDateString)
	}
	return fmt.Sprintf(
		"%s://%s/events?event_type=fixity+check&outcome=Failed&institution_id=%d&date_time__gteq=%s&date_time__lteq=%s",
		common.Context().Config.HTTPScheme(), domainAndPort, institutionID, startDateString, endDateString)
}
