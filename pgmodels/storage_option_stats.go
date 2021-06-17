package pgmodels

// StorageOptionStats contains info about how many bytes institutions
// have stored in each storage option.
type StorageOptionStats struct {
	TotalBytes            int64  `json:"total_bytes"`
	InstitutionID         int64  `json:"institution_id"`
	InstitutionName       string `json:"institution_name"`
	InstitutionIdentifier string `json:"institution_identifier"`
	StorageOption         string `json:"storage_option"`
}

// StorageOptionStatsSelect returns StorageOptionStats matching the query.
func StorageOptionStatsSelect(query *Query) ([]*StorageOptionStats, error) {
	var stats []*StorageOptionStats
	err := query.Select(&stats)
	return stats, err
}
