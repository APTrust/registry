package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type Checksum struct {
	ID            int64     `json:"id" form:"id" pg:"id"`
	Algorithm     string    `json:"algorithm" form:"algorithm" pg:"algorithm"`
	DateTime      time.Time `json:"datetime" form:"datetime" pg:"datetime"`
	Digest        string    `json:"digest" form:"digest" pg:"digest"`
	GenericFileID int64     `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	CreatedAt     time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
}

func (cs *Checksum) GetID() int64 {
	return cs.ID
}

func (cs *Checksum) Authorize(actingUser *User, action string) error {
	ctx := common.Context()
	gf, err := GenericFileFind(cs.GenericFileID)
	if err != nil {
		ctx.Log.Error().Msgf("While checking permissions for Checksum %d, could not find parent GenericFile %d", cs.ID, cs.GenericFileID)
		return common.ErrParentRecordNotFound
	}
	perm := "Checksum" + action
	if !actingUser.HasPermission(constants.Permission(perm), gf.InstitutionID) {
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s checksum %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, cs.ID, gf.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

// DeleteIsForbidden returns true because Checksums are our audit trail.
func (cs *Checksum) DeleteIsForbidden() bool {
	return true
}

// UpdateIsForbidden returns true because Checksums are our audit trail.
func (cs *Checksum) UpdateIsForbidden() bool {
	return true
}

func (cs *Checksum) IsReadOnly() bool {
	return false
}

func (cs *Checksum) SupportsSoftDelete() bool {
	return false
}

func (cs *Checksum) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (cs *Checksum) ClearSoftDeleteAttributes() {
	// No-op
}

func (cs *Checksum) SetTimestamps() {
	now := time.Now().UTC()
	if cs.CreatedAt.IsZero() {
		cs.CreatedAt = now
	}
	cs.UpdatedAt = now
}

func (cs *Checksum) BeforeSave() error {
	// TODO: Validate
	return nil
}
