package pgmodels

import (
	"time"
)

var ChecksumFilters = []string{
	"algorithm",
	"date_time__gteq",
	"date_time__lteq",
	"digest",
	"generic_file_id",
	"generic_file_identifier",
	"institution_id",
	"intellectual_object_id",
	"state",
}

type ChecksumView struct {
	tableName             struct{}  `pg:"checksums_view"`
	ID                    int64     `json:"id"`
	Algorithm             string    `json:"algorithm"`
	DateTime              time.Time `json:"datetime" pg:"datetime"`
	Digest                string    `json:"digest"`
	State                 string    `json:"state"`
	GenericFileIdentifier string    `json:"generic_file_identifier"`
	GenericFileID         int64     `json:"generic_file_id"`
	IntellectualObjectID  int64     `json:"intellectual_object_id"`
	InstitutionID         int64     `json:"institution_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ChecksumViewByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func ChecksumViewByID(id int64) (*ChecksumView, error) {
	query := NewQuery().Where("id", "=", id)
	return ChecksumViewGet(query)
}

// ChecksumViewGet returns the first file matching the query.
func ChecksumViewGet(query *Query) (*ChecksumView, error) {
	var cs ChecksumView
	err := query.Select(&cs)
	return &cs, err
}

// ChecksumViewSelect returns all files matching the query.
func ChecksumViewSelect(query *Query) ([]*ChecksumView, error) {
	var checksums []*ChecksumView
	err := query.Select(&checksums)
	return checksums, err
}
