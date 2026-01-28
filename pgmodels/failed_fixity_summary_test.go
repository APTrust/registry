package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var June012025 = time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

func TestFailedFixitySummary(t *testing.T) {
	addDummyFailures(t)
	defer deleteDummyFailures(t)

	// Run a query that should return results, and check
	// that results are correct.
	startDate := June012025.AddDate(0, 0, -1)
	endDate := June012025.AddDate(0, 0, 1)
	summaries, err := pgmodels.FailedFixitySummarySelect(startDate, endDate)

	require.Nil(t, err)
	require.NotEmpty(t, summaries)
	assert.Equal(t, 2, len(summaries))

	assert.EqualValues(t, 3, summaries[0].Failures)
	assert.EqualValues(t, 2, summaries[0].InstitutionID)
	assert.Equal(t, "Institution One", summaries[0].InstitutionName)

	assert.EqualValues(t, 3, summaries[1].Failures)
	assert.EqualValues(t, 3, summaries[1].InstitutionID)
	assert.Equal(t, "Institution Two", summaries[1].InstitutionName)

	// Run a query that should produce no results and ensure
	// it really does produce no results.
	startDate = June012025.AddDate(0, 0, 1)
	endDate = June012025.AddDate(0, 0, 3)
	summaries, err = pgmodels.FailedFixitySummarySelect(startDate, endDate)

	require.Nil(t, err)
	require.Empty(t, summaries)
}

func addDummyFailures(t *testing.T) {
	institutions := []int64{2, 3}
	for _, instID := range institutions {
		files, err := pgmodels.GenericFileSelect(pgmodels.NewQuery().Where("institution_id", "=", instID))
		require.Nil(t, err)
		for i := 0; i < 3; i++ {
			file := files[i]
			failure := pgmodels.PremisEvent{
				Agent:                0, // "Registry Unit Test",
				DateTime:             June012025,
				Detail:               "Failed fixity check for unit tests",
				EventType:            constants.EventFixityCheck,
				GenericFileID:        file.ID,
				Identifier:           uuid.NewString(),
				InstitutionID:        instID,
				IntellectualObjectID: file.IntellectualObjectID,
				Object:               0, // "Go language crypto/sha256",
				Outcome:              "Failed",
				OutcomeDetail:        "Yadda yadda",
				OutcomeInformation:   "Not a real failure. This is test data.",
			}
			failure.Save()
		}
	}
}

func deleteDummyFailures(t *testing.T) {
	_, err := common.Context().DB.Exec("delete from premis_events where outcome_information='Not a real failure. This is test data.'")
	assert.Nil(t, err)
}
