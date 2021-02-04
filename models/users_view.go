package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type UsersView struct {
	ID                     int64     `json:"id" pg:"id"`
	Name                   string    `json:"name" pg:"name"`
	Email                  string    `json:"email" pg:"email"`
	PhoneNumber            string    `json:"phone_number" pg:"phone_number"`
	CreatedAt              time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" pg:"updated_at"`
	ResetPasswordSentAt    time.Time `json:"reset_password_sent_at" pg:"reset_password_sent_at"`
	RememberCreatedAt      time.Time `json:"-" pg:"remember_created_at"`
	SignInCount            int       `json:"sign_in_count" pg:"sign_in_count,use_zero"`
	CurrentSignInAt        time.Time `json:"current_sign_in_at" pg:"current_sign_in_at"`
	LastSignInAt           time.Time `json:"last_sign_in_at" pg:"last_sign_in_at"`
	CurrentSignInIP        string    `json:"current_sign_in_ip" pg:"current_sign_in_ip"`
	LastSignInIP           string    `json:"last_sign_in_ip" pg:"last_sign_in_ip"`
	InstitutionID          int64     `json:"institution_id" pg:"institution_id"`
	PasswordChangedAt      time.Time `json:"password_changed_at" pg:"password_changed_at"`
	ConsumedTimestep       int       `json:"-" pg:"consumed_timestep"`
	OTPRequiredForLogin    bool      `json:"otp_required_for_login" pg:"otp_required_for_login"`
	DeactivatedAt          time.Time `json:"deactivated_at" pg:"deactivated_at"`
	EnabledTwoFactor       bool      `json:"enabled_two_factor" pg:"enabled_two_factor"`
	ConfirmedTwoFactor     bool      `json:"confirmed_two_factor" pg:"confirmed_two_factor"`
	AuthyID                string    `json:"-" pg:"authy_id"`
	LastSignInWithAuthy    time.Time `json:"last_sign_in_with_authy" pg:"last_sign_in_with_authy"`
	AuthyStatus            string    `json:"authy_status" pg:"authy_status"`
	EmailVerified          bool      `json:"email_verified" pg:"email_verified"`
	InitialPasswordUpdated bool      `json:"initial_password_updated" pg:"initial_password_updated"`
	ForcePasswordUpdate    bool      `json:"force_password_update" pg:"force_password_update"`
	GracePeriod            time.Time `json:"grace_period" pg:"grace_period"`
	Role                   string    `json:"role" pg:"role"`
	InstitutionName        string    `json:"institution_name" pg:"institution_name"`
	InstitutionIdentifier  string    `json:"institution_identifier" pg:"institution_identifier"`
	InstitutionState       string    `json:"institution_state" pg:"institution_state"`
	InstitutionType        string    `json:"institution_type" pg:"institution_type"`
	MemberInstitutionId    int64     `json:"member_institution_id" pg:"member_institution_id"`

	MemberInstitutionName       string `json:"member_institution_name" pg:"member_institution_name"`
	MemberInstitutionIdentifier string `json:"member_institution_identifier" pg:"member_institution_identifier"`
	MemberInstitutionState      string `json:"member_institution_state" pg:"member_institution_state"`

	OTPEnabled      bool   `json:"otp_enabled" pg:"otp_enabled"`
	ReceivingBucket string `json:"receiving_bucket" pg:"receiving_bucket"`
	RestoreBucket   string `json:"restore_bucket" pg:"restore_bucket"`
}

func (user *UsersView) GetID() int64 {
	return user.ID
}

func (user *UsersView) Authorize(actingUser *User, action string) error {
	perm := "User" + action
	if actingUser.ID == user.ID && (action == constants.ActionRead || action == constants.ActionUpdate || action == constants.ActionDelete) {
		perm += "Self"
	}
	if !actingUser.HasPermission(constants.Permission(perm), user.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %s at inst %d can't %s on subject user %s at inst %d\n", actingUser.Email, actingUser.InstitutionID, perm, user.Email, user.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (user *UsersView) DeleteIsForbidden() bool {
	return true
}

func (user *UsersView) UpdateIsForbidden() bool {
	return true
}

func (user *UsersView) IsReadOnly() bool {
	return true
}

func (user *UsersView) SupportsSoftDelete() bool {
	return false
}

func (user *UsersView) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (user *UsersView) ClearSoftDeleteAttributes() {
	// No-op
}

func (user *UsersView) SetTimestamps() {
	// No-op, since view is read-only
}

func (user *UsersView) BeforeSave() error {
	// No-op
	return nil
}

func UsersViewFind(id int64) (*UsersView, error) {
	ctx := common.Context()
	user := &UsersView{ID: id}
	err := ctx.DB.Model(user).WherePK().Select()
	return user, err
}
