package pgmodels

import (
//"github.com/APTrust/registry/common"
//"github.com/APTrust/registry/constants"
)

type StorageRecord struct {
	ID            int64        `json:"id" form:"id" pg:"id"`
	GenericFileID int64        `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	URL           string       `json:"url" form:"url" pg:"url"`
	GenericFile   *GenericFile `json:"-" pg:"rel:has-one"`
}

// StorageRecordByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func StorageRecordByID(id int64) (*StorageRecord, error) {
	query := NewQuery().Where(`"checksum"."id"`, "=", id).Relations()
	return StorageRecordGet(query)
}

// StorageRecordGet returns the first file matching the query.
func StorageRecordGet(query *Query) (*StorageRecord, error) {
	var sr StorageRecord
	err := query.Select(&sr)
	return &sr, err
}

// StorageRecordSelect returns all files matching the query.
func StorageRecordSelect(query *Query) ([]*StorageRecord, error) {
	var files []*StorageRecord
	err := query.Select(&files)
	return files, err
}

func (sr *StorageRecord) GetID() int64 {
	return sr.ID
}

// Save saves this file to the database. This will peform an insert
// if StorageRecord.ID is zero. Otherwise, it updates.
func (sr *StorageRecord) Save() error {
	if sr.ID == int64(0) {
		return insert(sr)
	}
	return update(sr)
}
