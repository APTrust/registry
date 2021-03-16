package pgmodels

import (
	"context"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/go-pg/pg/v10"
)

const (
	ErrInstName       = "Name must contain 5-100 characters."
	ErrInstIdentifier = "Identifier must be a domain name."
	ErrInstState      = "State must be 'A' or 'D'."
	ErrInstType       = "Please choose an institution type."
	ErrInstReceiving  = "Receiving bucket name is not valid."
	ErrInstRestore    = "Restoration bucket name is not valid."
	ErrInstMemberID   = "Please choose a parent institution."
)

var ValidStates = common.InterfaceList(constants.States)
var ValidInstTypes = common.InterfaceList(constants.InstTypes)

type Institution struct {
	ID                  int64     `json:"id" form:"id" pg:"id"`
	Name                string    `json:"name" pg:"name" binding:"required,min=2,max=100"`
	Identifier          string    `json:"identifier" pg:"identifier"`
	State               string    `json:"state" pg:"state"`
	Type                string    `json:"type" pg:"type" binding:"oneof=MemberInstitution SubscriptionInstitution"`
	MemberInstitutionID int64     `json:"member_institution_id" pg:"member_institution_id"`
	DeactivatedAt       time.Time `json:"deactivated_at" pg:"deactivated_at,soft_delete"`
	OTPEnabled          bool      `json:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket     string    `json:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket" pg:"restore_bucket"`
	CreatedAt           time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" pg:"updated_at"`

	//Users                   []*User        `json:"users" pg:"rel:has-many"`
	//SubscribingInstitutions []*Institution `json:"subscribing_institutions" pg:"rel:has-many"`
	// TODO: Add child institutions as an official relation
}

// The following statements have no effect other than to force a compile-time
// check that ensures our Institution model properly implements these hook
// interfaces.
var _ pg.BeforeDeleteHook = (*Institution)(nil)
var _ pg.BeforeInsertHook = (*Institution)(nil)
var _ pg.BeforeUpdateHook = (*Institution)(nil)

// BeforeDelete sets Institution.State to "D" before we perform a
// soft delete. Note that DeactivatedAt is the soft delete field,
// which means the pg library sets its timestamp instead of actually
// expunging the record from the DB.
func (inst *Institution) BeforeDelete(c context.Context) (context.Context, error) {
	inst.State = "D"
	return c, inst.Validate()
}

// BeforeInsert sets the CreatedAt and UpdatedAt timestamps on creation.
func (inst *Institution) BeforeInsert(c context.Context) (context.Context, error) {
	ctx := common.Context()
	now := time.Now().UTC()
	inst.CreatedAt = now
	inst.UpdatedAt = now
	inst.ReceivingBucket = fmt.Sprintf("aptrust.receiving%s.%s", ctx.Config.BucketQualifier(), inst.Identifier)
	inst.RestoreBucket = fmt.Sprintf("aptrust.restore%s.%s", ctx.Config.BucketQualifier(), inst.Identifier)
	inst.State = constants.StateActive
	return c, inst.Validate()
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (inst *Institution) BeforeUpdate(c context.Context) (context.Context, error) {
	inst.UpdatedAt = time.Now().UTC()
	return c, inst.Validate()
}

func (inst *Institution) Validate() error {
	return v.ValidateStruct(
		inst,
		v.Field(&inst.Name, v.Required.Error(ErrInstName), v.Length(5, 100).Error(ErrInstName)),
		v.Field(&inst.Identifier, v.Required.Error(ErrInstIdentifier), is.Domain.Error(ErrInstIdentifier)),
		v.Field(&inst.State, v.Required.Error(ErrInstState), v.In(ValidStates...).Error(ErrInstState)),
		v.Field(&inst.Type, v.Required.Error(ErrInstType), v.In(ValidInstTypes...).Error(ErrInstType)),
		v.Field(&inst.ReceivingBucket, v.Required.Error(ErrInstReceiving), v.Length(20, 100).Error(ErrInstReceiving)),
		v.Field(&inst.RestoreBucket, v.Required.Error(ErrInstRestore), v.Length(20, 100).Error(ErrInstRestore)),
		v.Field(&inst.MemberInstitutionID, v.When(inst.Type == constants.InstTypeSubscriber, v.Required.Error(ErrInstMemberID))),
	)
}

// InstitutionGet(db, query) (*Institution, error)
// InstitutionSelect(db, query) ([]*Institution, error)
// func (inst *Institution) Save() error [insert, update]
// func (inst *Institution) Delete() error
// func InstitutionBind(c *gin.Context) (*Institution, bool)
//      - this sets map[string]string of error messages

// Common module:
//
// - error messages
// - valid states, valid inst types, etc.
// - function to convert validation errors to map[string]string
//   - this function should return ["ValidationError"] for internal validation
//     errors
