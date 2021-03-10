package pgmodels

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
	DeactivatedAt       time.Time `json:"deactivated_at" pg:"deactivated_at,soft_delete"`
	OTPEnabled          bool      `json:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket     string    `json:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket" pg:"restore_bucket"`
	CreatedAt           time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" pg:"updated_at"`

	Users                   []*User        `json:"users" pg:"rel:has-many"`
	SubscribingInstitutions []*Institution `json:"subscribing_institutions" pg:"rel:has-many"`
	// TODO: Add child institutions as an official relation
}

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

// Compile-time check to ensure we implement hooks correctly.
// https://pg.uptrace.dev/hooks/
// https://medium.com/@matryer/golang-tip-compile-time-checks-to-ensure-your-type-satisfies-an-interface-c167afed3aae

// InstitutionGet(db, query) (*Institution, error)
// InstitutionSelect(db, query) ([]*Institution, error)
// func (inst *Institution) Save() error [insert, update]
// func (inst *Institution) Delete() error
// func InstitutionBind(c *gin.Context) (*Institution, bool)
//      - this sets map[string]string of error messages
