package pgmodels

import (
	"context"
	//"fmt"
	"time"

	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
)

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
	return c, nil
}

// BeforeInsert sets the CreatedAt and UpdatedAt timestamps on creation.
func (inst *Institution) BeforeInsert(c context.Context) (context.Context, error) {
	now := time.Now().UTC()
	inst.CreatedAt = now
	inst.UpdatedAt = now
	return c, nil
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (inst *Institution) BeforeUpdate(c context.Context) (context.Context, error) {
	inst.UpdatedAt = time.Now().UTC()
	return c, nil
}

// TODO: Struct level validation to ensure a subscribing institution
// has a parent, and to validate bucket names. For example, see:
// https://github.com/go-playground/validator/blob/v9/_examples/struct-level/main.go

// TODO: Check out https://github.com/go-ozzo/ozzo-validation
// as an alternative to gin's baked-in validation package.
// The problems with the baked in are:
//
// 1. Custom field validators are clumsy.
// 2. Custom struct validators are even more clumsy.
// 3. Custom error messages are even worse.
//
//
// Consider validators as separate structs, or consider wrappers
// around the ozzo validators that simplify usage.
//
// Or roll your own...
//
// min, max, regex, datemin, datemax, inList
//
// Also see https://github.com/asaskevich/govalidator,
// which ozzo uses under the hood.

// func InstitutionValidator(sl validator.StructLevel) {
// 	inst := sl.Current().Interface().(Institution)
// 	if inst.Type == constants.InstTypeSubscriber && inst.MemberInstitutionID == int64(0) {
// 		sl.ReportError(inst.MemberInstitutionID, "MemberInstitutionID", "MemberInstitutionID", "fnameorlname", "")
// 	}
// }

// InstitutionGet(db, query) (*Institution, error)
// InstitutionSelect(db, query) ([]*Institution, error)
// func (inst *Institution) Save() error [insert, update]
// func (inst *Institution) Delete() error
// func InstitutionBind(c *gin.Context) (*Institution, bool)
//      - this sets map[string]string of error messages
