package pgmodels

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
)

// Phone: +1234567890 (10 digits)
var rePhone = regexp.MustCompile("^\\+[0-9]{11}$")
var reNumeric = regexp.MustCompile("[^0-9]+")

const (
	ErrUserName         = "Name must contain at least 2 characters."
	ErrUserEmail        = "Email address is required."
	ErrUserPhone        = "Please enter a phone number in format +2125551212."
	ErrUser2Factor      = "Please choose yes or no."
	ErrUserGracePeriod  = "Enter a date specifying when this user must enable two-factor authentication."
	ErrUserInst         = "Please choose an institution."
	ErrUserRole         = "Please choose a role for this user."
	ErrUserInstNotFound = "Internal error: cannot find institution."
	ErrUserInvalidAdmin = "Sys Admin role is not valid for this institution."
	ErrUserPwdMissing   = "Encrypted password is missing."
)

// User is a person who can log in and do stuff.
type User struct {
	ID                     int64        `json:"id" form:"id" pg:"id"`
	Name                   string       `json:"name" pg:"name"`
	Email                  string       `json:"email" pg:"email"`
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
	InstitutionID          int64        `json:"institution_id" pg:"institution_id"`
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
	OTPBackupCodes         []string     `json:"-" form:"-" pg:"otp_backup_codes,array"`
	AuthyID                string       `json:"-" form:"-" pg:"authy_id"`
	LastSignInWithAuthy    time.Time    `json:"last_sign_in_with_authy" form:"-" pg:"last_sign_in_with_authy"`
	AuthyStatus            string       `json:"authy_status" form:"-" pg:"authy_status"`
	EmailVerified          bool         `json:"email_verified" form:"-" pg:"email_verified"`
	InitialPasswordUpdated bool         `json:"initial_password_updated" form:"-" pg:"initial_password_updated"`
	ForcePasswordUpdate    bool         `json:"force_password_update" form:"-" pg:"force_password_update"`
	GracePeriod            time.Time    `json:"grace_period" time_format:"2006-01-02" pg:"grace_period"`
	Role                   string       `json:"role" pg:"role"`
	Institution            *Institution `json:"institution" pg:"rel:has-one"`
}

// UserByID returns the institution with the specified id.
// Returns pg.ErrNoRows if there is no match.
func UserByID(id int64) (*User, error) {
	query := NewQuery().Where(`"user"."id"`, "=", id)
	return UserGet(query)
}

// UserByEmail returns the user with the specified email address.
// Returns pg.ErrNoRows if there is no match.
func UserByEmail(email string) (*User, error) {
	query := NewQuery().Where(`"user"."email"`, "=", email)
	return UserGet(query)
}

// UserGet returns the first user matching the query.
func UserGet(query *Query) (*User, error) {
	var user User
	err := query.Relations("Institution").Select(&user)
	return &user, err
}

// UserSelect returns all user matching the query.
func UserSelect(query *Query) ([]*User, error) {
	var users []*User
	err := query.Select(&users)
	return users, err
}

// UserSignIn signs a user in. If successful, it returns the User
// record with User.Institution properly set. If it fails, check
// the error.
func UserSignIn(email, password, ipAddr string) (*User, error) {
	user, err := UserByEmail(email)
	if IsNoRowError(err) {
		return nil, common.ErrInvalidLogin
	} else if err != nil {
		return nil, err
	}
	if !user.DeactivatedAt.IsZero() {
		return nil, common.ErrAccountDeactivated
	}
	if !common.ComparePasswords(user.EncryptedPassword, password) {
		common.Context().Log.Warn().Msgf("Wrong password for user %s", email)
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
	err = user.Save()
	return user, err
}

func (user *User) GetID() int64 {
	return user.ID
}

// UserSignOut signs a user out.
func (user *User) SignOut() error {
	if user.CurrentSignInIP != "" {
		user.LastSignInIP = user.CurrentSignInIP
	}
	if !user.CurrentSignInAt.IsZero() {
		user.LastSignInAt = user.CurrentSignInAt
	}
	user.CurrentSignInIP = ""
	user.CurrentSignInAt = time.Time{}
	return user.Save()
}

func (user *User) Save() error {
	if user.ID == int64(0) {
		return insert(user)
	}
	return update(user)
}

func (user *User) Delete() error {
	user.DeactivatedAt = time.Now().UTC()
	return update(user)
}

func (user *User) Undelete() error {
	user.DeactivatedAt = time.Time{}
	return update(user)
}

func (user *User) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if !v.IsByteLength(user.Name, 2, 200) {
		errors["Name"] = ErrUserName
	}
	if !v.IsEmail(user.Email) {
		errors["Email"] = ErrUserEmail
	}
	if user.PhoneNumber != "" && !rePhone.MatchString(user.PhoneNumber) {
		errors["PhoneNumber"] = ErrUserPhone
	}
	if user.OTPRequiredForLogin && user.GracePeriod.IsZero() {
		errors["GracePeriod"] = ErrUserGracePeriod
	}
	if user.InstitutionID < int64(1) {
		errors["InstitutionID"] = ErrUserInst
	}
	if !v.IsIn(user.Role, constants.Roles...) {
		errors["Role"] = ErrUserRole
	}
	if user.Role == constants.RoleSysAdmin {
		aptrust, err := InstitutionByIdentifier("aptrust.org")
		if err != nil {
			errors["Role"] = ErrUserInstNotFound
		} else if user.InstitutionID != aptrust.ID {
			errors["Role"] = ErrUserInvalidAdmin
		}
	}
	if user.EncryptedPassword == "" {
		errors["EncryptedPassword"] = ErrUserPwdMissing
	}
	if len(errors) > 0 {
		return &common.ValidationError{errors}
	}
	return nil
}

var (
	_ pg.BeforeInsertHook = (*User)(nil)
	_ pg.BeforeUpdateHook = (*User)(nil)
)

func (user *User) BeforeInsert(c context.Context) (context.Context, error) {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.reformatPhone()
	err := user.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (user *User) BeforeUpdate(c context.Context) (context.Context, error) {
	user.UpdatedAt = time.Now().UTC()
	user.reformatPhone()
	err := user.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

func (user *User) reformatPhone() {
	digitsOnly := reNumeric.ReplaceAllString(user.PhoneNumber, "")
	if len(digitsOnly) > 0 {
		if !strings.HasPrefix(digitsOnly, "1") {
			digitsOnly = "1" + digitsOnly
		}
		user.PhoneNumber = "+" + digitsOnly
	} else {
		user.PhoneNumber = digitsOnly // may be blank string
	}
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
