package pgmodels

import (
	"github.com/APTrust/registry/common"
)

type DepositFormatStats struct {
	FileFormat string  `json:"file_format"`
	FileCount  int64   `json:"file_count"`
	TotalBytes int64   `json:"total_bytes"`
	TotalGB    float64 `json:"total_gb" pg:"total_gb"`
	TotalTB    float64 `json:"total_tb" pg:"total_tb"`
}

func DepositFormatStatsSelect(institutionID, intellectualObjectID int64) ([]*DepositFormatStats, error) {
	var stats []*DepositFormatStats
	_, err := common.Context().DB.Query(&stats, depositFormatQuery,
		institutionID, institutionID,
		intellectualObjectID, intellectualObjectID)
	return stats, err
}

// depositFormatQuery reports on number and total file size by
// institution and/or intellectual object. This is used for donut
// charts on dashboard and/or intellectual object detail pages.
const depositFormatQuery = `
    select
        sum("size") as "total_bytes",
        (sum("size") / 1073741824) as "total_gb",
        (sum("size") / 1099511627776) as "total_tb",
        count(*) as file_count,
        coalesce(file_format, 'Total') as "file_format"
        from generic_files
        where (? = 0 or institution_id = ?)
        and   (? = 0 or intellectual_object_id = ?)
        group by rollup(file_format)
        order by file_format
`
