package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

	GenericFiles   []*GenericFile   `json:"generic_files" pg:"rel:has-many"`
	PremisEvents   []*PremisEvent   `json:"premis_events" pg:"rel:has-many"`
	Checksums      []*Checksum      `json:"checksumss" pg:"rel:has-many"`
	StorageRecords []*StorageRecord `json:"storage_records" pg:"rel:has-many"`
}

func (gf *GenericFile) GetID() int64 {
	return gf.ID
}

func (gf *GenericFile) Authorize(actingUser *User, action string) error {
	perm := "File" + action
	if !actingUser.HasPermission(constants.Permission(perm), gf.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s file %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, gf.ID, gf.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (gf *GenericFile) DeleteIsForbidden() bool {
	return false
}

func (gf *GenericFile) UpdateIsForbidden() bool {
	return false
}

func (gf *GenericFile) IsReadOnly() bool {
	return false
}

func (gf *GenericFile) SupportsSoftDelete() bool {
	return true
}

func (gf *GenericFile) SetSoftDeleteAttributes(actingUser *User) {
	gf.State = "D"
}

func (gf *GenericFile) ClearSoftDeleteAttributes() {
	gf.State = "A"
}

func (gf *GenericFile) SetTimestamps() {
	now := time.Now().UTC()
	if gf.CreatedAt.IsZero() {
		gf.CreatedAt = now
	}
	gf.UpdatedAt = now
}

func (gf *GenericFile) BeforeSave() error {
	// TODO: Validate
	return nil
}

func GenericFileFind(id int64) (*GenericFile, error) {
	ctx := common.Context()
	gf := &GenericFile{ID: id}
	err := ctx.DB.Model(gf).WherePK().Select()
	return gf, err
}

func GenericFileFindByIdentifier(identifier string) (*GenericFile, error) {
	ctx := common.Context()
	gf := &GenericFile{}
	err := ctx.DB.Model(gf).Where(`"identifier" = ?`, identifier).First()
	return gf, err
}
