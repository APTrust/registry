package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
)

type BillingStats struct {
	InstitutionID   int64     `json:"institution_id"`
	InstitutionName string    `json:"institution_name"`
	EndDate         time.Time `json:"end_date"`
	MonthAndYear    string    `json:"month_and_year"`
	StorageOption   string    `json:"storage_option"`
	TotalGB         float64   `json:"total_gb"`
	TotalTB         float64   `json:"total_tb"`
	Overage         float64   `json:"overage"`
}

var billingStatsQuery = `select
	institution_id,
	institution_name,
	end_date,
	to_char((end_date - interval '1 day'), 'Month YYYY') as month_and_year,
	storage_option,
	total_gb,
	total_tb,
	greatest((total_tb - 10.0), 0.0) as overage
	from historical_deposit_stats
	where institution_id = ?
	and end_date > ?
	and end_date <= ?
	and total_tb > 0
	order by end_date, storage_option`

func BillingStatsSelect(institutionID int64, startDate, endDate time.Time) ([]*BillingStats, error) {
	var stats []*BillingStats
	_, err := common.Context().DB.Query(&stats, billingStatsQuery, institutionID, startDate, endDate)
	return stats, err
}
