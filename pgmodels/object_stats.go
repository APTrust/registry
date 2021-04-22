package pgmodels

import (
	"github.com/APTrust/registry/common"
)

type ObjectStats struct {
	FileCount  int64  `json:"file_count"`
	FileSize   int64  `json:"file_size"`
	FileFormat string `json:"file_format"`
}

func GetObjectStats(intellectualObjectID int64) ([]*ObjectStats, error) {
	var stats []*ObjectStats
	_, err := common.Context().DB.Query(&stats,
		`select count(*) as "file_count", sum("size") as "file_size", file_format
         from generic_files where intellectual_object_id = ?
         group by rollup(file_format) order by file_format`,
		intellectualObjectID)
	return stats, err
}

/*
select count(*) as "file_count", sum("size") as "file_size", file_format
from generic_files where intellectual_object_id = 6092 group by rollup(file_format) order by file_format;
*/
