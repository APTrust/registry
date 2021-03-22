package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
)

type GenericFile struct {
	ID                   int64     `json:"id" form:"id" pg:"id"`
	FileFormat           string    `json:"file_format" form:"file_format" pg:"file_format"`
	Size                 int64     `json:"size" form:"size" pg:"size"`
	Identifier           string    `json:"identifier" form:"identifier" pg:"identifier"`
	IntellectualObjectID int64     `json:"intellectual_object_id" form:"intellectual_object_id" pg:"intellectual_object_id"`
	CreatedAt            time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	State                string    `json:"state" form:"state" pg:"state"`
	LastFixityCheck      time.Time `json:"last_fixity_check" form:"last_fixity_check" pg:"last_fixity_check"`
	InstitutionID        int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	StorageOption        string    `json:"storage_option" form:"storage_option" pg:"storage_option"`
	UUID                 string    `json:"uuid" form:"uuid" pg:"uuid"`

	Institution        *Institution        `json:"-" pg:"rel:has-one"`
	IntellectualObject *IntellectualObject `json:"-" pg:"rel:has-one"`
	PremisEvents       []*PremisEvent      `json:"premis_events" pg:"rel:has-many"`
	Checksums          []*Checksum         `json:"checksumss" pg:"rel:has-many"`
	StorageRecords     []*StorageRecord    `json:"storage_records" pg:"rel:has-many"`
}
