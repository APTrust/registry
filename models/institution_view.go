package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

// InstitutionView contains information about an institution and its
// parent (if it has one). This view simplifies both search and display
// for institution management.
type InstitutionView struct {
	tableName           struct{}  `pg:"users_view"`
	ID                  int64     `json:"id" pg:"id"`
	Name                string    `json:"name" pg:"name"`
	Identifier          string    `json:"identifier" pg:"identifier"`
	State               string    `json:"state" pg:"state"`
	Type                string    `json:"type" pg:"type"`
	DeactivatedAt       time.Time `json:"deactivated_at" pg:"deactivated_at"`
	OTPEnabled          bool      `json:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket     string    `json:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket" pg:"restore_bucket"`
	CreatedAt           time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" pg:"updated_at"`
	ParentId            int64     `json:"parent_id" pg:"parent_id"`
	ParentName          string    `json:"parent_name" pg:"parent_name"`
	ParentIdentifier    string    `json:"parent_identifier" pg:"parent_identifier"`
	ParentState         string    `json:"parent_state" pg:"parent_state"`
	ParentDeactivatedAt time.Time `json:"parent_deactivated_at" pg:"parent_deactivated_at"`
}

func (inst *InstitutionView) GetID() int64 {
	return inst.ID
}

func (inst *InstitutionView) Authorize(actingUser *User, action string) error {
	perm := "Institution" + action
	if !actingUser.HasPermission(constants.Permission(perm), inst.ID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d can't %s on subject institution %d\n", actingUser.ID, perm, inst.ID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (inst *InstitutionView) DeleteIsForbidden() bool {
	return true
}

func (inst *InstitutionView) UpdateIsForbidden() bool {
	return true
}

func (inst *InstitutionView) IsReadOnly() bool {
	return true
}

func (inst *InstitutionView) SupportsSoftDelete() bool {
	return false
}

func (inst *InstitutionView) SetSoftDeleteAttributes(actingUser *User) {
	// No-op. This is a read-only view.
}

func (inst *InstitutionView) ClearSoftDeleteAttributes() {
	// No-op. This is a read-only view.
}

func (inst *InstitutionView) SetTimestamps() {
	// No-op. This is a read-only view.
}

func (inst *InstitutionView) BeforeSave() error {
	return nil
}
