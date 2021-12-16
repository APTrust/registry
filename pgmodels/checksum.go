package pgmodels

import (
	"time"
)

type Checksum struct {
	TimestampModel
	Algorithm     string       `json:"algorithm"`
	DateTime      time.Time    `json:"datetime" pg:"datetime"`
	Digest        string       `json:"digest"`
	GenericFileID int64        `json:"generic_file_id" pg:"generic_file_id"`
	GenericFile   *GenericFile `json:"-" pg:"rel:has-one"`
}

// ChecksumByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func ChecksumByID(id int64) (*Checksum, error) {
	query := NewQuery().Where(`"checksum"."id"`, "=", id).Relations()
	return ChecksumGet(query)
}

// ChecksumGet returns the first file matching the query.
func ChecksumGet(query *Query) (*Checksum, error) {
	var cs Checksum
	err := query.Select(&cs)
	return &cs, err
}

// ChecksumSelect returns all files matching the query.
func ChecksumSelect(query *Query) ([]*Checksum, error) {
	var files []*Checksum
	err := query.Select(&files)
	return files, err
}

// Save saves this file to the database. This will peform an insert
// if Checksum.ID is zero. Otherwise, it updates.
func (cs *Checksum) Save() error {
	if cs.ID == int64(0) {
		return insert(cs)
	}
	return update(cs)
}
