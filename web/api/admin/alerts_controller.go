package admin_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
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
	if lastRunDate == today {
		ctx.Log.Info().Msgf("Not generating failed fixity alerts because they were last generated on %v", lastRunDate)
		response.Error = "Alerts have already been generated today."
		c.JSON(http.StatusConflict, response)
		return
	}

	summaries, err := pgmodels.FailedFixitySummarySelect(lastRunDate, time.Now().UTC())
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
		err = GenerateFailedFixityAlert(summary)
		if err != nil {
			response.Error = fmt.Sprintf("%s; %s", response.Error, err.Error())
		}
	}

	if len(response.Summaries) > 0 {
		c.JSON(http.StatusCreated, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func FailedFixityLastRunDate() (time.Time, error) {
	yesterdayTs := time.Now().UTC().AddDate(0, 0, -1)
	yesterday := time.Date(yesterdayTs.Year(), yesterdayTs.Month(), yesterdayTs.Day(), 0, 0, 0, 0, time.UTC)

	metadata, err := pgmodels.InternalMetadataByKey("fixity alerts last run")

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

func GenerateFailedFixityAlert(summary *pgmodels.FailedFixitySummary) error {
	ctx := common.Context()
	ctx.Log.Info().Msgf("Generating failed fixity alert for %s with %d failures.", summary.InstitutionName, summary.Failures)

	// TODO: Generate alert here. Log error or success.

	return nil
}
