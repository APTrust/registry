package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

/*
   The OTPBackupCodes column in this table is of type _varchar,
   which is a delimited array. The pg driver seems to have a problem
   deserializing the array, so in the User model below, the type
   is set to string instead of []string.
*/

// User is a person who can log in and do stuff.
type User struct {
	ID                     int64        `json:"id" form:"id" pg:"id"`
	Name                   string       `json:"name" pg:"name" binding:"required,min=2,max=100"`
	Email                  string       `json:"email" pg:"email" binding:"required,email"`
	PhoneNumber            string       `json:"phone_number" pg:"phone_number"`
	CreatedAt              time.Time    `json:"created_at" form:"-" pg:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at" form:"-" pg:"updated_at"`
	EncryptedPassword      string       `json:"-" form:"-" pg:"encrypted_password"`
	ResetPasswordToken     string       `json:"-" form:"-" pg:"reset_password_token"`
	ResetPasswordSentAt    time.Time    `json:"reset_password_sent_at" form:"-" pg:"reset_password_sent_at"`
	RememberCreatedAt      time.Time    `json:"-" form:"-" pg:"remember_created_at"`
	SignInCount            int          `json:"sign_in_count" form:"-" pg:"sign_in_count,use_zero"`
	CurrentSignInAt        time.Time    `json:"current_sign_in_at" form:"-" pg:"current_sign_in_at"`
	LastSignInAt           time.Time    `json:"last_sign_in_at" form:"-" pg:"last_sign_in_at"`
	CurrentSignInIP        string       `json:"current_sign_in_ip" form:"-" pg:"current_sign_in_ip"`
	LastSignInIP           string       `json:"last_sign_in_ip" form:"-" pg:"last_sign_in_ip"`
	InstitutionID          int64        `json:"institution_id" pg:"institution_id" binding:"required"`
	EncryptedAPISecretKey  string       `json:"-" form:"-" pg:"encrypted_api_secret_key"`
	PasswordChangedAt      time.Time    `json:"password_changed_at" form:"-" pg:"password_changed_at"`
	EncryptedOTPSecret     string       `json:"-" form:"-" pg:"encrypted_otp_secret"`
	EncryptedOTPSecretIV   string       `json:"-" form:"-" pg:"encrypted_otp_secret_iv"`
	EncryptedOTPSecretSalt string       `json:"-" form:"-" pg:"encrypted_otp_secret_salt"`
	ConsumedTimestep       int          `json:"-" form:"-" pg:"consumed_timestep"`
	OTPRequiredForLogin    bool         `json:"otp_required_for_login" pg:"otp_required_for_login"`
	DeactivatedAt          time.Time    `json:"deactivated_at" form:"-" pg:"deactivated_at"`
	EnabledTwoFactor       bool         `json:"enabled_two_factor" form:"-" pg:"enabled_two_factor"`
	ConfirmedTwoFactor     bool         `json:"confirmed_two_factor" form:"-" pg:"confirmed_two_factor"`
	OTPBackupCodes         string       `json:"-" form:"-" pg:"otp_backup_codes"`
	AuthyID                string       `json:"-" form:"-" pg:"authy_id"`
	LastSignInWithAuthy    time.Time    `json:"last_sign_in_with_authy" form:"-" pg:"last_sign_in_with_authy"`
	AuthyStatus            string       `json:"authy_status" form:"-" pg:"authy_status"`
	EmailVerified          bool         `json:"email_verified" form:"-" pg:"email_verified"`
	InitialPasswordUpdated bool         `json:"initial_password_updated" form:"-" pg:"initial_password_updated"`
	ForcePasswordUpdate    bool         `json:"force_password_update" form:"-" pg:"force_password_update"`
	GracePeriod            time.Time    `json:"grace_period" time_format:"2006-01-02" pg:"grace_period"`
	Role                   string       `json:"role" pg:"role" binding:"required"`
	Institution            *Institution `json:"institution" pg:"rel:has-one"`
}

func (user *User) GetID() int64 {
	return user.ID
}

func (user *User) Authorize(actingUser *User, action string) error {
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

func (user *User) DeleteIsForbidden() bool {
	return false
}

func (user *User) UpdateIsForbidden() bool {
	return false
}

func (user *User) IsReadOnly() bool {
	return false
}

func (user *User) SupportsSoftDelete() bool {
	return true
}

func (user *User) SetSoftDeleteAttributes(*User) {
	user.DeactivatedAt = time.Now().UTC()
}

func (user *User) ClearSoftDeleteAttributes() {
	user.DeactivatedAt = time.Time{}
}

func (user *User) SetTimestamps() {
	now := time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now
}

func (user *User) BeforeSave() error {
	// TODO: Validate
	return nil
}

// HasPermission returns true or false to indicate whether the user has
// sufficient permissions to perform the requested action. Param action
// should be one of the constants from constants/permissions.go. Param
// institutionID should be the ID of the institution that owns the object
// upon which the user is trying to act. In certain cases, such as when a
// user is editing him/herself, this can be zero.
func (user *User) HasPermission(action constants.Permission, institutionID int64) bool {
	// Sys admin's permissions apply across all institutional boundaries.
	if user.IsAdmin() {
		return constants.CheckPermission(user.Role, action)
	}
	// Institutional user and admin permissions apply only within their
	// own institutions.
	return user.InstitutionID == institutionID && constants.CheckPermission(user.Role, action)
}

// IsAdmin returns true if user is a Sys Admin. Returns false for all other
// roles (including institutional admin, which is not a super user).
func (user *User) IsAdmin() bool {
	return user.Role == constants.RoleSysAdmin
}
