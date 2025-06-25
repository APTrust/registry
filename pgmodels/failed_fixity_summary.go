package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
)

type FailedFixitySummary struct {
	Failures        int64  `json:"failures"`
	InstitutionID   int64  `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
}

var failedFixityQuery = `select count(id) as "failures", pev.institution_id, pev.institution_name
	from premis_events_view pev
	where pev.event_type = 'fixity check'
	and outcome = 'Failed'
	and date_time > ?
	and date_time < ?
	group by pev.institution_id, pev.institution_name
	order by pev.institution_name;`

func FailedFixitySummarySelect(startDate, endDate time.Time) ([]*FailedFixitySummary, error) {
	var summaries []*FailedFixitySummary
	_, err := common.Context().DB.Query(&summaries, failedFixityQuery, startDate, endDate)
	return summaries, err
}
