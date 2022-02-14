package pgmodels

import (
	"github.com/APTrust/registry/common"
	v "github.com/asaskevich/govalidator"
)

var StorageRecordFilters = []string{
	"generic_file_id",
}

type StorageRecord struct {
	ID            int64        `json:"id" pg:"id"`
	GenericFileID int64        `json:"generic_file_id"`
	URL           string       `json:"url" form:"url" pg:"url"`
	GenericFile   *GenericFile `json:"-" pg:"rel:has-one"`
}

// StorageRecordByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func StorageRecordByID(id int64) (*StorageRecord, error) {
	query := NewQuery().Where("id", "=", id).Relations()
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
	err := sr.Validate()
	if err != nil {
		return err
	}
	if sr.ID == int64(0) {
		return insert(sr)
	}
	return update(sr)
}

func (sr *StorageRecord) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if sr.GenericFileID <= 0 {
		errors["GenericFileID"] = "Storage record requires a valid file id"
	}
	if !v.IsURL(sr.URL) {
		errors["URL"] = "Storage record requires a valid URL"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
