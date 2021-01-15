package models

import (
	"time"
)

type User struct {
	ID                     int64     `json:"id" form:"id" pg:"id"`
	Name                   string    `json:"name" form:"name" pg:"name"`
	PhoneNumber            string    `json:"phone_number" form:"phone_number" pg:"phone_number"`
	CreatedAt              time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	EncryptedPassword      string    `json:"-" form:"-" pg:"encrypted_password"`
	ResetPasswordToken     string    `json:"-" form:"-" pg:"reset_password_token"`
	ResetPasswordSentAt    time.Time `json:"reset_password_sent_at" form:"-" pg:"reset_password_sent_at"`
	RememberCreatedAt      time.Time `json:"-" form:"-" pg:"remember_created_at"`
	SignInCount            int       `json:"sign_in_count" form:"-" pg:"sign_in_count"`
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
	ConsumedTimestamp      int       `json:"-" form:"-" pg:"consumed_timestamp"`
	OTPRequiredForLogin    bool      `json:"otp_required_for_login" form:"-" pg:"otp_required_for_login"`
	DeactivatedAt          time.Time `json:"deactivated_at" form:"-" pg:"deactivated_at"`
	EnabledTwoFactor       bool      `json:"enabled_two_factor" form:"-" pg:"enabled_two_factor"`
	ConfirmedTwoFactor     bool      `json:"confirmed_two_factor" form:"-" pg:"confirmed_two_factor"`
	OTPBackupCodes         []string  `json:"-" form:"-" pg:"otp_backup_codes"`
	AuthyID                string    `json:"-" form:"-" pg:"authy_id"`
	LastSignInWithAuthy    time.Time `json:"last_sign_in_with_authy" form:"-" pg:"last_sign_in_with_authy"`
	AuthyStatus            string    `json:"authy_status" form:"-" pg:"authy_status"`
	EmailVerified          bool      `json:"email_verified" form:"-" pg:"email_verified"`
	InitialPasswordUpdated bool      `json:"initial_password_updated" form:"-" pg:"initial_password_updatedd"`
	ForcePasswordUpdate    bool      `json:"force_password_update" form:"-" pg:"force_password_update"`
	GracePeriod            time.Time `json:"grace_period" form:"grace_period" pg:"grace_period"`
}
