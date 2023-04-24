package pgmodels_test

import (
	"fmt"
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
	stats, err := pgmodels.BillingStatsSelect(999, startDate, endDate)

	require.Nil(t, err)
	require.NotEmpty(t, stats)

	for _, s := range stats {
		fmt.Println(s)
	}
	assert.True(t, false)
}

func addDummyBillingData(t *testing.T) {
	instID := int64(999)
	instName := "Dummy Inst"
	storageOptions := []string{constants.StorageOptionStandard, constants.StorageOptionGlacierDeepOR}
	insert := `INSERT INTO historical_deposit_stats (institution_id, institution_name, storage_option, object_count, file_count, total_bytes, total_gb, total_tb, cost_gb_per_month, monthly_cost, end_date, member_institution_id, primary_sort, secondary_sort) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	for i := 0; i < 12; i++ {
		month := i + 1
		endDate := time.Date(2022, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		for _, storageOpt := range storageOptions {
			objCount := month * 100
			fileCount := objCount * 10
			totalBytes := int64(fileCount * 2000000 * 1024)
			totalGB := int64(totalBytes / (1024 * 1024 * 1024))
			totalTB := int64(totalGB / 1024)
			_, err := common.Context().DB.Exec(insert, instID, instName, storageOpt, objCount, fileCount, totalBytes, totalGB, totalTB, 0, 0, endDate, 0, "aaa", "bbb")
			require.Nil(t, err)
		}
	}
}

func deleteDummyBillingData(t *testing.T) {
	_, err := common.Context().DB.Exec("delete from historical_deposit_stats where institution_id = 999")
	assert.Nil(t, err)
}
