package pgmodels

import (
	"fmt"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
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

var InstitutionFilters = []string{
	"name__contains",
	"type",
}

type Institution struct {
	TimestampModel
	Name                string    `json:"name"`
	Identifier          string    `json:"identifier"`
	State               string    `json:"state"`
	Type                string    `json:"type"`
	MemberInstitutionID int64     `json:"member_institution_id"`
	DeactivatedAt       time.Time `json:"deactivated_at"`
	OTPEnabled          bool      `json:"otp_enabled" pg:",use_zero"`
	EnableSpotRestore   bool      `json:"enable_spot_restore" pg:",use_zero"`
	ReceivingBucket     string    `json:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket"`
}

// InstitutionByID returns the institution with the specified id.
// Returns pg.ErrNoRows if there is no match.
func InstitutionByID(id int64) (*Institution, error) {
	query := NewQuery().Where("id", "=", id)
	return InstitutionGet(query)
}

// InstitutionByIdentifier returns the institution with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func InstitutionByIdentifier(identifier string) (*Institution, error) {
	query := NewQuery().Where("identifier", "=", identifier)
	return InstitutionGet(query)
}

// IdForInstIdentifier returns the id of the insitution with
// the given identifier, or an error if no matching record exists.
func IdForInstIdentifier(identifier string) (int64, error) {
	query := NewQuery().Columns("id").Where("identifier", "=", identifier)
	inst, err := InstitutionGet(query)
	return inst.ID, err
}

// InstitutionGet returns the first institution matching the query.
func InstitutionGet(query *Query) (*Institution, error) {
	var institution Institution
	err := query.Select(&institution)
	return &institution, err
}

// InstitutionSelect returns all institutions matching the query.
func InstitutionSelect(query *Query) ([]*Institution, error) {
	var institutions []*Institution
	err := query.Select(&institutions)
	return institutions, err
}

// Save saves this institution to the database. This will peform an insert
// if Institution.ID is zero. Otherwise, it updates.
func (inst *Institution) Save() error {
	inst.SetTimestamps()
	if inst.ID == 0 {
		inst.ReceivingBucket = inst.bucket("receiving")
		inst.RestoreBucket = inst.bucket("restore")
		inst.State = constants.StateActive
	}
	err := inst.Validate()
	if err != nil {
		return err
	}
	if inst.ID == int64(0) {
		return insert(inst)
	}
	return update(inst)
}

// Delete soft-deletes this institution by setting State to 'D' and
// the DeletedAt timestamp to now. You can undo this with Undelete.
func (inst *Institution) Delete() error {
	inst.State = constants.StateDeleted
	inst.DeactivatedAt = time.Now().UTC()
	return update(inst)
}

// Undelete reactivates this institution by setting State to 'A' and
// clearing the DeletedAt timestamp.
func (inst *Institution) Undelete() error {
	inst.State = constants.StateActive
	inst.DeactivatedAt = time.Time{}
	return update(inst)
}

// bucket returns a valid bucket name for this institution.
// Param name should be "receiving" or "restore"
func (inst *Institution) bucket(name string) string {
	ctx := common.Context()
	return fmt.Sprintf("aptrust.%s%s.%s", name, ctx.Config.BucketQualifier(), inst.Identifier)
}

// Validate validates the model. This is called automatically on insert
// and update.
func (inst *Institution) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if !v.IsByteLength(inst.Name, 5, 200) {
		errors["Name"] = ErrInstName
	}
	// DNS names without dots, such as "localhost" are valid,
	// but we require a DNS name with at least one dot.
	if !v.IsDNSName(inst.Identifier) || !strings.Contains(inst.Identifier, ".") {
		errors["Identifier"] = ErrInstIdentifier
	}
	if !v.IsIn(inst.State, constants.States...) {
		errors["State"] = ErrInstState
	}
	if !v.IsIn(inst.Type, constants.InstTypes...) {
		errors["Type"] = ErrInstType
	}
	if inst.ReceivingBucket != inst.bucket("receiving") {
		errors["ReceivingBucket"] = ErrInstReceiving
	}
	if inst.RestoreBucket != inst.bucket("restore") {
		errors["RestoreBucket"] = ErrInstRestore
	}
	if inst.Type == constants.InstTypeSubscriber && inst.MemberInstitutionID < int64(1) {
		errors["MemberInstitutionID"] = ErrInstMemberID
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
