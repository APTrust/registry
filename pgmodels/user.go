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
	ErrUserPwdIncorrect = "Incorrect Password."
)

// User is a person who can log in and do stuff.
type User struct {
	// ID is user's unique ID in the database.
	ID int64 `json:"id" form:"id" pg:"id"`

	// Name is the user's display name.
	Name string `json:"name" pg:"name"`

	// Email is used to login. Must be unique.
	Email string `json:"email" pg:"email"`

	// PhoneNumber should start with + country code (e.g. +1 for US).
	// This number is used for SMS two-factor auth, so it should be
	// a mobile phone.
	PhoneNumber string `json:"phone_number" pg:"phone_number"`

	// CreatedAt is timestamp of user creation.
	CreatedAt time.Time `json:"created_at" form:"-" pg:"created_at"`

	// UpdatedAt is when user was last updated.
	UpdatedAt time.Time `json:"updated_at" form:"-" pg:"updated_at"`

	// EncryptedPassword is the user's password, encrypted.
	EncryptedPassword string `json:"-" form:"-" pg:"encrypted_password"`

	// ResetPasswordToken is an encrypted version of the token sent to a
	// user who wants to reset their password. The plaintext version of
	// this is emailed to the user.
	ResetPasswordToken string `json:"-" form:"-" pg:"reset_password_token"`

	// ResetPasswordSentAt is a timestamp describing when we sent a password
	// reset email. The ResetPasswordToken should be valid for only a few
	// minutes after this timestamp.
	ResetPasswordSentAt time.Time `json:"reset_password_sent_at" form:"-" pg:"reset_password_sent_at"`

	// RememberCreatedAt - ??? - Legacy field from Devise, probably related
	// to time-based OTPs. Not used.
	RememberCreatedAt time.Time `json:"-" form:"-" pg:"remember_created_at"`

	// SignInCount is the number of times this user has successfully signed in.
	SignInCount int `json:"sign_in_count" form:"-" pg:"sign_in_count,use_zero"`

	// CurrentSignInAt is a timestamp describing when this user signed
	// in for their current session.
	CurrentSignInAt time.Time `json:"current_sign_in_at" form:"-" pg:"current_sign_in_at"`

	// LastSignInAt is a timestamp describing when this user signed
	// in for their previous session.
	LastSignInAt time.Time `json:"last_sign_in_at" form:"-" pg:"last_sign_in_at"`

	// CurrentSignInIP is the IP address from which user signed in.
	CurrentSignInIP string `json:"current_sign_in_ip" form:"-" pg:"current_sign_in_ip"`

	// LastSignInIP is the IP address from which user signed in for their
	// prior session.
	LastSignInIP string `json:"last_sign_in_ip" form:"-" pg:"last_sign_in_ip"`

	// InstitituionID is the id of the institution to which this user
	// belongs.
	InstitutionID int64 `json:"institution_id" pg:"institution_id"`

	// EncryptedAPISecretKey is the user's encrypted API key. This may
	// be empty if the user has never requested a key.
	EncryptedAPISecretKey string `json:"-" form:"-" pg:"encrypted_api_secret_key"`

	// PasswordChangedAt is the date and time the user last changed their
	// password.
	PasswordChangedAt time.Time `json:"password_changed_at" form:"-" pg:"password_changed_at"`

	// EncryptedOTPSecret is the encrypted version of a user's one-time
	// password. The plaintext is usually a six-digit code sent via SMS.
	// This will be empty if we didn't send a code, or if the user has
	// correctly entered it and completed two-factor login.
	EncryptedOTPSecret string `json:"-" form:"-" pg:"encrypted_otp_secret"`

	// EncryptedOTPSecretIV is a legacy field from Devise. Not used.
	EncryptedOTPSecretIV string `json:"-" form:"-" pg:"encrypted_otp_secret_iv"`

	// EncryptedOTPSecretSalt is a legacy field from Devise. Not used.
	EncryptedOTPSecretSalt string `json:"-" form:"-" pg:"encrypted_otp_secret_salt"`

	// ConsumedTimestep is a legacy field from Devise, which used it
	// for time-based one-time passwords. Not used.
	ConsumedTimestep int `json:"-" form:"-" pg:"consumed_timestep"`

	// OTPRequiredForLogin indicates whether, as a matter of policy, the
	// user must use some form of OTP to log in. If true, the user should
	// be allowed to log in only with Authy one-touch, six-digit SMS
	// code, or backup code.
	//
	// Use IsTwoFactorUser() to actually determine whether we should
	// force this user to use Authy, SMS, or backup code to get it.
	// This is messy because it has to mirror the legacy Rails implementation
	// for now.
	OTPRequiredForLogin bool `json:"otp_required_for_login" pg:"otp_required_for_login"`

	// DeactivatedAt is a timestamp describing when this user was
	// deactivated. For most users, it will be empty. APTrust admin and
	// Institutional Admin can disable other users. Then can also
	// re-enable a user by clearing this flag.
	DeactivatedAt time.Time `json:"deactivated_at" form:"-" pg:"deactivated_at"`

	// EnabledTwoFactor indicates whether the user's account has
	// enabled two factor authentication. See als ConfirmedTwoFactor.
	EnabledTwoFactor bool `json:"enabled_two_factor" form:"-" pg:"enabled_two_factor"`

	// ConfirmedTwoFactor indicates that the system has confirmed this
	// user's phone number (for SMS) and Authy account (for Authy 2FA).
	ConfirmedTwoFactor bool `json:"confirmed_two_factor" form:"-" pg:"confirmed_two_factor"`

	// OTPBackupCodes is a list of backup codes the user can use to
	// log if they can't get in via SMS or Authy.
	OTPBackupCodes []string `json:"-" form:"-" pg:"otp_backup_codes,array"`

	// AuthyID is the user's Authy ID. We need this to send them push
	// messages to complete one-touch sign-in. This will be empty for
	// those who don't use Authy.
	AuthyID string `json:"-" form:"-" pg:"authy_id"`

	// LastSignInWithAuthy is the timestamp of this user's last
	// successful sign-in with Authy.
	LastSignInWithAuthy time.Time `json:"last_sign_in_with_authy" form:"-" pg:"last_sign_in_with_authy"`

	// AuthyStatus indicates how the user wants to authenticate with
	// Authy. If it's constants.TwoFactorAuthy, we should send them a
	// push, so they can login with one-touch. Anything else means SMS,
	// but call IsTwoFactorUser() to make sure they're actually require
	// two-factor auth before trying to text them.
	AuthyStatus string `json:"authy_status" form:"-" pg:"authy_status"`

	// EmailVerified will be true once the system has verified that the
	// user's email address is correct.
	EmailVerified bool `json:"email_verified" form:"-" pg:"email_verified"`

	// InitialPasswordUpdated will be true once a user updates their
	// initial password. When we create a new user, we generate a random
	// password, then send them an email with a login link. At that
	// point, the user has to set their own password.
	InitialPasswordUpdated bool `json:"initial_password_updated" form:"-" pg:"initial_password_updated"`

	// ForcePasswordUpdate indicated whether the user will be forced to
	// update their password the next time they visit.
	ForcePasswordUpdate bool `json:"force_password_update" form:"-" pg:"force_password_update"`

	// GracePeriod is a legacy field from the old Rails app. It held the
	// date by which a user MUST complete either Authy or SMS two-factor
	// setup. This feature was universally despised, and we won't be
	// using it unless someone complains. For now, consider this as
	// unused.
	GracePeriod time.Time `json:"grace_period" time_format:"2006-01-02" pg:"grace_period"`

	// AwaitingSecondFactor indicates that the user has logged in with
	// email and password, but has not yet completed the second login
	// step (Authy, SMS OTP, or backup code). This flag is only set on
	// users going through the two-factor process. Middleware checks it
	// to prevent partially logged-in users from accessing any pages other
	// than those required to complete the two-factor login process.
	AwaitingSecondFactor bool `json:"-" pg:"awaiting_second_factor,use_zero"`

	// Role is the user's role.
	Role string `json:"role" pg:"role"`

	// Institution is where they lock you up after you've spent too much
	// time trying to figure out the old Rails code.
	Institution *Institution `json:"institution" pg:"rel:has-one"`
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
		return &common.ValidationError{Errors: errors}
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

// IsSMSUser returns true if this user has enabled two-factor authentication
// with SMS/text message.
func (user *User) IsSMSUser() bool {
	return user.IsTwoFactorUser() && (user.AuthyStatus == constants.TwoFactorSMS || user.AuthyStatus == "")
}

// IsAuthyOneTouchUser returns true if the user has enabled Authy one touch
// for two-factor login.
func (user *User) IsAuthyOneTouchUser() bool {
	return user.IsTwoFactorUser() && (user.AuthyStatus == constants.TwoFactorAuthy)
}

// IsTwoFactorUser returns true if this user has enabled and confirmed
// two factor authentication.
//
// Sorry... For now, we have to work with the convoluted logic of the
// old Rails app. Hence the confusion between this an OTPRequiredForLogin.
func (user *User) IsTwoFactorUser() bool {
	return user.EnabledTwoFactor && user.ConfirmedTwoFactor
}

// TwoFactorMethod returns one of the following:
//
// constants.TwoFactorNone if the user does not use two-factor auth.
//
// constants.TwoFactorAuthy if the user uses two-factor auth via Authy.
//
// constants.TwoFactorSMS if the user receives two-factor OTP code via
// text/SMS
func (user *User) TwoFactorMethod() string {
	if !user.IsTwoFactorUser() {
		return constants.TwoFactorNone
	}
	if user.IsSMSUser() {
		return constants.TwoFactorSMS
	}
	return constants.TwoFactorAuthy
}

// CreateOTPToken creates a new one-time password token, typically
// used for SMS-based two-factor authentication. It saves an
// encrypted version of the token to the database and returns
// the plaintext version of the token. For SMS, we use six-digit
// tokens because they're easy for a user to type.
func (user *User) CreateOTPToken() (string, error) {
	token, err := common.NewOTP()
	if err != nil {
		return "", err
	}
	encryptedToken, err := common.EncryptPassword(token)
	if err != nil {
		return "", err
	}
	user.EncryptedOTPSecret = encryptedToken
	err = user.Save()
	if err != nil {
		return "", err
	}
	return token, err
}

// ClearOTPSecret deletes the user's EncryptedOTPSecret.
func (user *User) ClearOTPSecret() error {
	user.EncryptedOTPSecret = ""
	return user.Save()
}
