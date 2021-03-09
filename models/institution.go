package models

import (
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type Institution struct {
	ID                  int64     `json:"id" form:"id" pg:"id"`
	Name                string    `json:"name" pg:"name" binding:"required,min=2,max=100"`
	Identifier          string    `json:"identifier" pg:"identifier"`
	State               string    `json:"state" pg:"state"`
	Type                string    `json:"type" pg:"type" binding:"oneof=MemberInstitution SubscriptionInstitution"`
	MemberInstitutionID int64     `json:"member_institution_id" pg:"member_institution_id"`
	DeactivatedAt       time.Time `json:"deactivated_at" pg:"deactivated_at"`
	OTPEnabled          bool      `json:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket     string    `json:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket" pg:"restore_bucket"`
	CreatedAt           time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" pg:"updated_at"`

	Users []*User `json:"users" pg:"rel:has-many"`

	// TODO: Add child institutions as an official relation
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
	inst.State = "D"
}

func (inst *Institution) ClearSoftDeleteAttributes() {
	inst.DeactivatedAt = time.Time{}
	inst.State = "A"
}

func (inst *Institution) SetTimestamps() {
	now := time.Now().UTC()
	if inst.CreatedAt.IsZero() {
		inst.CreatedAt = now
	}
	inst.UpdatedAt = now
}

func (inst *Institution) BeforeSave() error {
	if inst.ID == int64(0) {
		ctx := common.Context()
		inst.State = "A"
		inst.ReceivingBucket = fmt.Sprintf("aptrust.receiving%s%s", ctx.Config.BucketQualifier(), inst.Identifier)
		inst.RestoreBucket = fmt.Sprintf("aptrust.restore%s%s", ctx.Config.BucketQualifier(), inst.Identifier)
		ctx.Log.Info().Msgf("Set buckets for new institution '%s' to %s and %s", inst.Name, inst.ReceivingBucket, inst.RestoreBucket)
	}
	return nil
}

// TODO: Struct level validation to ensure a subscribing institution
// has a parent, and to validate bucket names. For example, see:
// https://github.com/go-playground/validator/blob/v9/_examples/struct-level/main.go

// Should we move binding, validation, and error messages from the form
// class to the model? Prolly.

// Should we provide an option to save the model without validating? Hmm.

// Also, it may be a good idea to take advantage of Go Pg's hooks, instead
// of writing our own BeforeSave, etc. See:
// https://pkg.go.dev/github.com/go-pg/pg/v10#AfterDeleteHook (useless)
// https://pg.uptrace.dev/hooks/ (even uselesser)
//
// This helps a little:
// https://github.com/go-pg/pg/issues/1275
//
// Looks like hooks are defined directly on the model, with the
// underscore trick used to force a compiler check.
//
// Consider replacing EVERYTHING in the Model interface with the hooks
// defined at https://pkg.go.dev/github.com/go-pg/pg/v10#AfterDeleteHook
//
// Consider centralizing permission checks, so they're not reimplemented
// in each model.
