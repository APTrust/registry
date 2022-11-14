package pgmodels

import (
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
)

// Note: chart_metric is ignored by backend. Used only in front-end.
var DepositStatsFilters = []string{
	"chart_metric",
	"updated_at__lteq",
	"institution_id",
	"storage_option",
}

// DepositStats contains info about member deposits and the costs
// of those deposits. This struct does not implement the usual pgmodel
// interface, nor does it map to a single underlying table or view.
// This struct merely represents to the output of a reporting query.
type DepositStats struct {
	InstitutionID   int64     `json:"institution_id"`
	InstitutionName string    `json:"institution_name"`
	StorageOption   string    `json:"storage_option"`
	ObjectCount     int64     `json:"object_count"`
	FileCount       int64     `json:"file_count"`
	TotalBytes      int64     `json:"total_bytes"`
	TotalGB         float64   `json:"total_gb" pg:"total_gb"`
	TotalTB         float64   `json:"total_tb" pg:"total_tb"`
	CostGBPerMonth  float64   `json:"cost_gb_per_month" pg:"cost_gb_per_month"`
	MonthlyCost     float64   `json:"monthly_cost"`
	EndDate         time.Time `json:"end_date"`
}

// DepositStatsSelect returns info about materials a depositor updated
// in our system before a given date. This breaks down deposits by
// storage option and institution. To report on all institutions, use
// zero for institutionID. To report on all storage options, pass an
// empty string for storageOption.
func DepositStatsSelect(institutionID int64, storageOption string, endDate time.Time) ([]*DepositStats, error) {
	var stats []*DepositStats
	statsQuery := getDepositStatsQuery(institutionID, storageOption, endDate)
	_, err := common.Context().DB.Query(&stats, statsQuery,
		institutionID, institutionID,
		storageOption, storageOption,
		endDate, endDate)

	// If we happen to get a query for a date before 2014,
	// we'll get no results. We don't want to return nil, because
	// the caller is likely expected something that can be serialized
	// to JSON. Give the caller an actual answer, saying there was
	// nothing in the system on the date they inquired about.
	if stats == nil {
		stats = make([]*DepositStats, 1)
		stats[0] = &DepositStats{
			InstitutionName: "Total",
			StorageOption:   "Total",
			EndDate:         endDate,
		}
	}
	return stats, err
}

func getDepositStatsQuery(institutionID int64, storageOption string, endDate time.Time) string {
	// Basic depost stats query. Use the "is null / or" trick to deal with
	// filters that may or may not be present.
	q := `select institution_id, 
				institution_name, 
				storage_option, 
				file_count, 
				object_count, 
				total_bytes, 
				total_gb, 
				total_tb, 
				cost_gb_per_month,
				monthly_cost, 
				end_date from %s 
				where (? = 0 or institution_id = ?)
				and (? = '' or storage_option = ?)
				and (? = '0001-01-01 00:00:00+00:00:00' or end_date < ?)`
	tableName := "historical_deposit_stats"

	// Current stats report, which displays on dashboard, passes in
	// time.Now() as end date. In this case, we want to query the
	// current stats table, not historical stats. Realistically,
	// end_data will be within a few milliseconds of time.Now(),
	// let's allow 60 seconds of drift.
	if endDate.After(time.Now().Add(-60 * time.Second)) {
		tableName = "current_deposit_stats"
	}
	return fmt.Sprintf(q, tableName)
}
