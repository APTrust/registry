package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBillingStats(t *testing.T) {
	addDummyBillingData(t)
	defer deleteDummyBillingData(t)

	startDate := time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)
	stats, err := pgmodels.BillingStatsSelect(999, startDate, endDate, "")

	require.Nil(t, err)
	require.NotEmpty(t, stats)
	assert.Equal(t, 24, len(stats))

	// Check first and last records in the set.
	assert.EqualValues(t, 999, stats[0].InstitutionID)
	assert.Equal(t, "January   2022", stats[0].MonthAndYear)
	assert.Equal(t, "Glacier-Deep-OR", stats[0].StorageOption)
	assert.EqualValues(t, 3814, stats[0].TotalGB)
	assert.EqualValues(t, 3.724609375, stats[0].TotalTB)
	assert.EqualValues(t, 0, stats[0].Overage)

	assert.EqualValues(t, 999, stats[23].InstitutionID)
	assert.Equal(t, "December  2022", stats[23].MonthAndYear)
	assert.Equal(t, "Standard", stats[23].StorageOption)
	assert.EqualValues(t, 24795, stats[23].TotalGB)
	assert.EqualValues(t, 24.2138671875, stats[23].TotalTB)
	assert.EqualValues(t, 14.2138671875, stats[23].Overage)

	// for _, s := range stats {
	// 	fmt.Println(s)
	// }
	// assert.True(t, false)

	// Try again with a storage option filter
	stats, err = pgmodels.BillingStatsSelect(999, startDate, endDate, constants.StorageOptionGlacierDeepOR)
	require.Nil(t, err)
	require.NotEmpty(t, stats)
	assert.Equal(t, 12, len(stats))

	for _, s := range stats {
		assert.Equal(t, constants.StorageOptionGlacierDeepOR, s.StorageOption)
	}

}

func addDummyBillingData(t *testing.T) {
	instID := int64(999)
	instName := "Dummy Inst"
	storageOptions := []string{constants.StorageOptionStandard, constants.StorageOptionGlacierDeepOR}
	insert := `INSERT INTO historical_deposit_stats (institution_id, institution_name, storage_option, object_count, file_count, total_bytes, total_gb, total_tb, cost_gb_per_month, monthly_cost, end_date, member_institution_id, primary_sort, secondary_sort) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	for month := 1; month <= 13; month++ {
		endDate := time.Date(2022, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		if month == 13 {
			endDate = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		}
		//fmt.Println(month, endDate)
		for _, storageOpt := range storageOptions {
			objCount := month * 100
			fileCount := objCount * 10
			totalBytes := int64(fileCount * 2000000 * 1024)
			totalGB := float64(totalBytes / (1024 * 1024 * 1024))
			totalTB := float64(totalGB / 1024)
			_, err := common.Context().DB.Exec(insert, instID, instName, storageOpt, objCount, fileCount, totalBytes, totalGB, totalTB, 0, 0, endDate, 0, "aaa", "bbb")
			require.Nil(t, err)
		}
	}
}

func deleteDummyBillingData(t *testing.T) {
	_, err := common.Context().DB.Exec("delete from historical_deposit_stats where institution_id = 999")
	assert.Nil(t, err)
}
