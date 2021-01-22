package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type Institution struct {
	ID                  int64     `json:"id" form:"id" pg:"id"`
	Name                string    `json:"name" form:"name" pg:"name"`
	Identifier          string    `json:"identifier" form:"identifier" pg:"identifier"`
	CreatedAt           time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" pg:"updated_at"`
	State               string    `json:"state" form:"state" pg:"state"`
	Type                string    `json:"type" form:"type" pg:"type"`
	MemberInstitutionId int64     `json:"member_institution_id" form:"member_institution_id" pg:"member_institution_id"`
	DeactivatedAt       time.Time `json:"deactivated_at" form:"deactivated_at" pg:"deactivated_at"`
	OTPEnabled          bool      `json:"otp_enabled" form:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket     string    `json:"receiving_bucket" form:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket" form:"restore_bucket" pg:"restore_bucket"`

	Users []*User `json:"users" pg:"rel:has-many"`
}

func (inst *Institution) GetID() int64 {
	return inst.ID
}

func (inst *Institution) Authorize(actingUser *User, action string) error {
	perm := "Institution" + action
	if !actingUser.HasPermission(constants.Permission(perm), inst.ID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d can't %s on subject institution %d\n", actingUser.ID, perm, inst.ID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (inst *Institution) DeleteIsForbidden() bool {
	return false
}

func (inst *Institution) UpdateIsForbidden() bool {
	return false
}

func (inst *Institution) IsReadOnly() bool {
	return false
}

func (inst *Institution) SupportsSoftDelete() bool {
	return true
}

func (inst *Institution) SetSoftDeleteAttributes(actingUser *User) {
	inst.DeactivatedAt = time.Now().UTC()
}

func (inst *Institution) ClearSoftDeleteAttributes() {
	inst.DeactivatedAt = time.Time{}
}

func (inst *Institution) SetTimestamps() {
	now := time.Now().UTC()
	if inst.CreatedAt.IsZero() {
		inst.CreatedAt = now
	}
	inst.UpdatedAt = now
}

func (inst *Institution) BeforeSave() error {
	// TODO: Validate
	return nil
}

func InstitutionFind(id int64) (*Institution, error) {
	ctx := common.Context()
	inst := &Institution{ID: id}
	err := ctx.DB.Model(inst).WherePK().Select()
	return inst, err
}

func InstitutionFindByIdentifier(identifier string) (*Institution, error) {
	ctx := common.Context()
	inst := &Institution{}
	err := ctx.DB.Model(inst).Where(`"identifier" = ?`, identifier).First()
	return inst, err
}
