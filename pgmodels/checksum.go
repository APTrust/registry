package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/stretchr/stew/slice"
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
	var checksums []*Checksum
	err := query.Select(&checksums)
	return checksums, err
}

// Save saves this file to the database. This will peform an insert
// if Checksum.ID is zero. Otherwise, it returns an error, since updating
// checksums is not allowed. Checksums should only ever change on
// re-ingest, when we get a new version of an existing file. In that case,
// we add a new checksum, so that we have records for all checksums that
// have existed over time.
func (cs *Checksum) Save() error {
	cs.SetTimestamps()
	err := cs.Validate()
	if err != nil {
		return err
	}
	if cs.ID == int64(0) {
		return insert(cs)
	}
	return common.ErrNotSupported
}

func (cs *Checksum) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(cs.Digest) {
		errors["Digest"] = "Checksum Digest cannot be empty"
	}
	if cs.DateTime.IsZero() {
		errors["DateTime"] = "Checksum DateTime is required"
	}
	if !slice.Contains(constants.DigestAlgs, cs.Algorithm) {
		errors["Algorithm"] = "Checksum requires a valid algorithm"
	}
	if cs.GenericFileID <= 0 {
		errors["GenericFileID"] = "Checksum requires a valid file id"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
