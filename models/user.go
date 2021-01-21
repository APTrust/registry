package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type User struct {
	ID                     int64     `json:"id" form:"id" pg:"id"`
	Name                   string    `json:"name" form:"name" pg:"name"`
	Email                  string    `json:"email" form:"email" pg:"email"`
	PhoneNumber            string    `json:"phone_number" form:"phone_number" pg:"phone_number"`
	CreatedAt              time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	EncryptedPassword      string    `json:"-" form:"-" pg:"encrypted_password"`
	ResetPasswordToken     string    `json:"-" form:"-" pg:"reset_password_token"`
	ResetPasswordSentAt    time.Time `json:"reset_password_sent_at" form:"-" pg:"reset_password_sent_at"`
	RememberCreatedAt      time.Time `json:"-" form:"-" pg:"remember_created_at"`
	SignInCount            int       `json:"sign_in_count" form:"-" pg:"sign_in_count,use_zero"`
	CurrentSignInAt        time.Time `json:"current_sign_in_at" form:"-" pg:"current_sign_in_at"`
	LastSignInAt           time.Time `json:"last_sign_in_at" form:"-" pg:"last_sign_in_at"`
	CurrentSignInIP        string    `json:"current_sign_in_ip" form:"-" pg:"current_sign_in_ip"`
	LastSignInIP           string    `json:"last_sign_in_ip" form:"-" pg:"last_sign_in_ip"`
	InstitutionID          int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	EncryptedAPISecretKey  string    `json:"-" form:"-" pg:"encrypted_api_secret_key"`
	PasswordChangedAt      time.Time `json:"password_changed_at" form:"-" pg:"password_changed_at"`
	EncryptedOTPSecret     string    `json:"-" form:"-" pg:"encrypted_otp_secret"`
	EncryptedOTPSecretIV   string    `json:"-" form:"-" pg:"encrypted_otp_secret_iv"`
	EncryptedOTPSecretSalt string    `json:"-" form:"-" pg:"encrypted_otp_secret_salt"`
	ConsumedTimestep       int       `json:"-" form:"-" pg:"consumed_timestep"`
	OTPRequiredForLogin    bool      `json:"otp_required_for_login" form:"-" pg:"otp_required_for_login"`
	DeactivatedAt          time.Time `json:"deactivated_at" form:"-" pg:"deactivated_at"`
	EnabledTwoFactor       bool      `json:"enabled_two_factor" form:"-" pg:"enabled_two_factor"`
	ConfirmedTwoFactor     bool      `json:"confirmed_two_factor" form:"-" pg:"confirmed_two_factor"`
	OTPBackupCodes         []string  `json:"-" form:"-" pg:"otp_backup_codes"`
	AuthyID                string    `json:"-" form:"-" pg:"authy_id"`
	LastSignInWithAuthy    time.Time `json:"last_sign_in_with_authy" form:"-" pg:"last_sign_in_with_authy"`
	AuthyStatus            string    `json:"authy_status" form:"-" pg:"authy_status"`
	EmailVerified          bool      `json:"email_verified" form:"-" pg:"email_verified"`
	InitialPasswordUpdated bool      `json:"initial_password_updated" form:"-" pg:"initial_password_updated"`
	ForcePasswordUpdate    bool      `json:"force_password_update" form:"-" pg:"force_password_update"`
	GracePeriod            time.Time `json:"grace_period" form:"grace_period" pg:"grace_period"`
	Role                   string    `json:"role" form:"role" pg:"-"`
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

func SignInUser(email, password, ipAddr string) (*User, error) {
	ctx := common.Context()

	user, err := UserFindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.DeactivatedAt.IsZero() {
		return nil, common.ErrAccountDeactivated
	}

	if !common.ComparePasswords(user.EncryptedPassword, password) {
		ctx.Log.Warn().Msgf("Wrong password for user %s", email)
		return nil, common.ErrInvalidLogin
	}

	user.SignInCount = user.SignInCount + 1
	if user.CurrentSignInIP != "" {
		user.LastSignInIP = user.CurrentSignInIP
	}
	if user.CurrentSignInAt.IsZero() {
		user.LastSignInAt = user.CurrentSignInAt
	}
	user.CurrentSignInIP = ipAddr
	user.CurrentSignInAt = time.Now().UTC()
	_, err = ctx.DB.Model(user).WherePK().Update()

	return user, err
}

func UserFind(id int64) (*User, error) {
	ctx := common.Context()
	user := &User{ID: id}
	err := ctx.DB.Model(user).WherePK().Select()
	if err == nil {
		err = user.loadRole()
	}
	return user, err
}

func UserFindByEmail(email string) (*User, error) {
	ctx := common.Context()
	user := &User{}
	err := ctx.DB.Model(user).
		Where("email = ? ", email).
		Offset(0).
		Limit(1).
		Select()
	if IsNoRowError(err) {
		ctx.Log.Error().Msgf("No users matches email %s", email)
		return nil, common.ErrInvalidLogin
	}
	if err != nil {
		return nil, err
	}
	err = user.loadRole()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignOutUser -> copy current sign in time and ip to last, then clear
func (user *User) SignOutUser() error {
	ctx := common.Context()
	if user.CurrentSignInIP != "" {
		user.LastSignInIP = user.CurrentSignInIP
	}
	if !user.CurrentSignInAt.IsZero() {
		user.LastSignInAt = user.CurrentSignInAt
	}
	user.CurrentSignInIP = ""
	user.CurrentSignInAt = time.Time{}
	_, err := ctx.DB.Model(user).WherePK().Update()
	return err
}

func (user *User) loadRole() error {
	// Since users effectively have only one role, role should
	// go into the users table.

	// Initialize to something safe.
	// RoleNone corresponds to the empty permissions list
	// in constants/permissions.go
	user.Role = constants.RoleNone

	// If user is deactivated, code above should prevent them
	// from signing in at all. Just in case, this failsafe leaves
	// them with RoleNone and zero permissions. We do this instead
	// of returning an error because an admin may be loading this
	// user to review their info or to reactivate them.
	if !user.DeactivatedAt.IsZero() {
		return nil
	}

	ctx := common.Context()
	role := &Role{}
	err := ctx.DB.Model(role).
		Join("JOIN roles_users AS ru ON ru.role_id = role.id").
		Where("ru.user_id = ? ", user.ID).
		First()
	if err != nil {
		return err
	}
	user.Role = role.Name
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
	if user.Role == constants.RoleSysAdmin {
		return constants.CheckPermission(user.Role, action)
	}
	// Institutional user and admin permissions apply only within their
	// own institutions.
	return user.InstitutionID == institutionID && constants.CheckPermission(user.Role, action)
}
