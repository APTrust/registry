package pgmodels

import (
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
	InstitutionID   int64   `json:"institution_id"`
	InstitutionName string  `json:"institution_name"`
	StorageOption   string  `json:"storage_option"`
	ObjectCount     int64   `json:"object_count"`
	FileCount       int64   `json:"file_count"`
	TotalBytes      int64   `json:"total_bytes"`
	TotalGB         float64 `json:"total_gb" pg:"total_gb"`
	TotalTB         float64 `json:"total_tb" pg:"total_tb"`
	CostGBPerMonth  float64 `json:"cost_gb_per_month" pg:"cost_gb_per_month"`
	MonthlyCost     float64 `json:"monthly_cost"`
}

func DepositStatsSelect(institutionID int64, storageOption string, updatedBefore time.Time) ([]*DepositStats, error) {
	var stats []*DepositStats
	_, err := common.Context().DB.Query(&stats, depositStatsQuery,
		institutionID, institutionID,
		storageOption, storageOption,
		updatedBefore, updatedBefore)
	return stats, err
}

// Basic depost stats query. Use the "is null / or" trick to deal with
// filters that may or may not be present.
//
// This is used on the deposits report page.
const depositStatsQuery = `
		select
		  coalesce(stats.institution_name, 'Total') as institution_name,
		  i2.id as institution_id,
		  coalesce(stats.storage_option, 'Total') as storage_option,
		  stats.file_count,
		  stats.object_count,
		  stats.total_bytes,
		  (stats.total_bytes / 1073741824) as total_gb,
		  (stats.total_bytes / 1099511627776) as total_tb,
		  so.cost_gb_per_month,
		  ((stats.total_bytes / 1073741824) * so.cost_gb_per_month) as monthly_cost
		from
		  (select
			i."name" as institution_name,
			count(gf.id) as file_count,
			count(distinct(gf.intellectual_object_id)) as object_count,
			sum(gf.size) as total_bytes,
			gf.storage_option
		  from generic_files gf
		  left join institutions i on i.id = gf.institution_id
		  where gf.state = 'A'
		  and (? = 0 or i.id = ?)
		  and (? = '' or gf.storage_option = ?)
		  and (? = '0001-01-01 00:00:00+00:00:00' or gf.updated_at < ?)
		  group by cube (i."name", gf.storage_option)) stats
		left join storage_options so on so."name" = stats.storage_option
		left join institutions i2 on i2."name" = stats.institution_name
		order by stats.institution_name, stats.storage_option
`
