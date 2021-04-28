package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
)

type GenericFile struct {
	ID                   int64     `json:"id" pg:"id"`
	FileFormat           string    `json:"file_format" pg:"file_format"`
	Size                 int64     `json:"size" pg:"size"`
	Identifier           string    `json:"identifier" pg:"identifier"`
	IntellectualObjectID int64     `json:"intellectual_object_id" pg:"intellectual_object_id"`
	CreatedAt            time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" pg:"updated_at"`
	State                string    `json:"state" pg:"state"`
	LastFixityCheck      time.Time `json:"last_fixity_check" pg:"last_fixity_check"`
	InstitutionID        int64     `json:"institution_id" pg:"institution_id"`
	StorageOption        string    `json:"storage_option" pg:"storage_option"`
	UUID                 string    `json:"uuid" pg:"uuid"`

	Institution        *Institution        `json:"-" pg:"rel:has-one"`
	IntellectualObject *IntellectualObject `json:"-" pg:"rel:has-one"`
	PremisEvents       []*PremisEvent      `json:"premis_events" pg:"rel:has-many"`
	Checksums          []*Checksum         `json:"checksumss" pg:"rel:has-many"`
	StorageRecords     []*StorageRecord    `json:"storage_records" pg:"rel:has-many"`
}

// GenericFileByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func GenericFileByID(id int64) (*GenericFile, error) {
	query := NewQuery().Where(`"generic_file"."id"`, "=", id)
	return GenericFileGet(query)
}

// GenericFileByIdentifier returns the file with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func GenericFileByIdentifier(identifier string) (*GenericFile, error) {
	query := NewQuery().Where(`"generic_file"."identifier"`, "=", identifier)
	return GenericFileGet(query)
}

// GenericFileGet returns the first file matching the query.
func GenericFileGet(query *Query) (*GenericFile, error) {
	var gf GenericFile
	err := query.Select(&gf)
	return &gf, err
}

// GenericFileSelect returns all files matching the query.
func GenericFileSelect(query *Query) ([]*GenericFile, error) {
	var files []*GenericFile
	err := query.Select(&files)
	return files, err
}

func (gf *GenericFile) GetID() int64 {
	return gf.ID
}

// Save saves this file to the database. This will peform an insert
// if GenericFile.ID is zero. Otherwise, it updates.
func (gf *GenericFile) Save() error {
	if gf.ID == int64(0) {
		return insert(gf)
	}
	return update(gf)
}

// ObjectFileCount returns the number of active files with the specified
// Intellectial Object ID.
func ObjectFileCount(intellectualObjectID int64) (int, error) {
	return common.Context().DB.Model((*GenericFile)(nil)).Where(`intellectual_object_id = ? and state = 'A'`, intellectualObjectID).Count()
}
