package admin_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
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

	response := &api.JsonList{}

	// Find out when this process was last run.
	lastRunDate, err := FailedFixityLastRunDate()
	if err != nil {
		ctx.Log.Error().Msgf("Could not get last run date for failed fixity alerts: %v", err)
		reqError := api.RequestError{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		c.JSON(http.StatusInternalServerError, reqError)
		return
	}

	// If this has already run today, don't run it again.
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if lastRunDate.Equal(today) {
		ctx.Log.Info().Msgf("Not generating failed fixity alerts because they were last generated on %v", lastRunDate)
		reqError := api.RequestError{
			StatusCode: http.StatusConflict,
			Error:      "Note: failed fixity alerts have already been generated today. Controller is returning without running them again.",
		}
		c.JSON(http.StatusConflict, reqError)
		return
	}

	ctx.Log.Info().Msgf("Querying for failed fixity checks between %s and %s", lastRunDate.Format(time.RFC3339), now.Format(time.RFC3339))
	summaries, err := pgmodels.FailedFixitySummarySelect(lastRunDate, now)
	if err != nil {
		ctx.Log.Error().Msgf("Error querying for failed fixity alerts: %v", err)
		reqError := api.RequestError{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		c.JSON(http.StatusInternalServerError, reqError)
		return
	}
	ctx.Log.Info().Msgf("Result: failed fixity check query returned %d summaries", len(summaries))
	response.Results = summaries

	// Generate a fixity failure alert for each institution.
	// If there's an error in the alert generation process,
	// add that into the response, but keep processing because
	// we don't want one failed alert to prevent others from
	// being sent.
	errorOccurred := false
	for _, summary := range summaries {
		ctx.Log.Info().Msgf("%s has %d failed fixity checks between %s and %s", summary.InstitutionName, summary.Failures, lastRunDate.Format(time.RFC3339), now.Format(time.RFC3339))
		err = GenerateFailedFixityAlert(c.Request.Host, summary, lastRunDate)
		if err != nil {
			// Log this error, but keep going, so other institutions
			// can get their alerts.
			ctx.Log.Error().Msgf("Error generating failed fixity alert for institution %s: %v", summary.InstitutionName, err)
			errorOccurred = true
		}
		response.Count += int(summary.Failures)
	}

	// Generate only one report for APTrust admins. This will
	// include all failures at all institutions, and the embedded
	// link will show them all.
	err = AlertAPTrustOfFailedFixities(c.Request.Host, summaries, lastRunDate)
	if err != nil {
		ctx.Log.Error().Msgf("Error generating failed fixity alerts for APTrust admins: %v", err)
		reqError := api.RequestError{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		c.JSON(http.StatusInternalServerError, reqError)
		return
	}

	err = SetFailedFixityLastRunDate(now)
	if err != nil {
		ctx.Log.Error().Msgf("Error updating last failed fixity run date in DB: %v", err)
	} else {
		ctx.Log.Info().Msgf("Set last failed fixity run date in DB to %s", now.Format(time.RFC3339))
	}

	if response.Count > 0 {
		if errorOccurred {
			ctx.Log.Info().Msgf("Responding to client with HTTP status 500 because we created alerts for %d institutions, but one or more alerts failed.", len(summaries))
			c.JSON(http.StatusInternalServerError, response)
		} else {
			ctx.Log.Info().Msgf("Responding to client with HTTP status 201 because we created alerts for %d institutions", len(summaries))
			c.JSON(http.StatusCreated, response)
		}
	} else {
		ctx.Log.Info().Msg("Responding to client with HTTP status 200 because we didn't create any alerts.")
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

// SetFailedFixityLastRunDate sets the last run date for
// failed fixity alerts in the database.
func SetFailedFixityLastRunDate(ts time.Time) error {
	metadata, err := pgmodels.InternalMetadataByKey(constants.MetaFixityAlertsLastRun)

	// Now rows in result set means GenerateFailedFixityAlerts has
	// never run before. This should happen only once.
	if pgmodels.IsNoRowError(err) {
		metadata = &pgmodels.InternalMetadata{
			Key:   constants.MetaFixityAlertsLastRun,
			Value: ts.Format(time.RFC3339),
		}
	} else if err != nil {
		return err
	}
	metadata.Value = ts.Format(time.RFC3339)
	return metadata.Save()
}

// GenerateFailedFixityAlert creates fixity failure alerts
// for each admin at the institution in summary.InstitutionID.
func GenerateFailedFixityAlert(hostname string, summary *pgmodels.FailedFixitySummary, lastRunDate time.Time) error {
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
		"AlertURL":  FailedFixityReportURL(hostname, institution.ID, lastRunDate, today),
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

func AlertAPTrustOfFailedFixities(hostname string, summaries []*pgmodels.FailedFixitySummary, lastRunDate time.Time) error {
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

	// Get a list of APTrust admins, except for the system@aptrust.org user,
	// because that's a service account, not a real person.
	query := pgmodels.NewQuery().Where("institution_id", "=", institution.ID).Where("email", "!=", constants.SystemUser)
	ctx.Log.Info().Msgf("Looking for APTrust admins by selecting from users table WHERE %s", query.WhereClause())

	aptrustAdmins, err := pgmodels.UserSelect(query)
	if err != nil {
		ctx.Log.Error().Msgf("Error getting list of APTrust admins for admin/ops email: %v", err)
		return err
	}
	if len(aptrustAdmins) == 0 {
		ctx.Log.Error().Msg("Error getting list of APTrust admins: query returned no users, which should be impossible.")
		return err
	}

	today := time.Now().UTC()
	alert := pgmodels.NewFailedFixityAlert(institution.ID, events, aptrustAdmins)
	alertData := map[string]interface{}{
		"StartDate": lastRunDate.Format("2006-01-02"),
		"EndDate":   today.Format("2006-01-02"),
		"AlertURL":  FailedFixityReportURL(hostname, 0, lastRunDate, today),
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
func FailedFixityReportURL(hostname string, institutionID int64, lastRunDate, currentRunDate time.Time) string {
	startDateString := lastRunDate.Format("2006-01-02")
	endDateString := currentRunDate.Format("2006-01-02")
	if institutionID == 0 {
		return fmt.Sprintf(
			"%s://%s/events?event_type=fixity+check&outcome=Failed&date_time__gteq=%s&date_time__lteq=%s",
			common.Context().Config.HTTPScheme(), hostname, startDateString, endDateString)
	}
	return fmt.Sprintf(
		"%s://%s/events?event_type=fixity+check&outcome=Failed&institution_id=%d&date_time__gteq=%s&date_time__lteq=%s",
		common.Context().Config.HTTPScheme(), hostname, institutionID, startDateString, endDateString)
}
