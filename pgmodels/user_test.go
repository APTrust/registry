package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/stew/slice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (

	// Test constants for users (from fixture data)
	SysAdmin     = "system@aptrust.org"
	InstAdmin    = "admin@inst1.edu"
	InstUser     = "user@inst1.edu"
	InactiveUser = "inactive@inst1.edu"
	Password     = "password"

	// Institution IDs (from fixture data)
	InstAPTrust = int64(1)
	InstOne     = int64(2)
	InstTwo     = int64(3)
	InstTest    = int64(4)
	InstExample = int64(5)
)

func TestUserHasPermission(t *testing.T) {
	ownInst := int64(2)
	otherInst := int64(3)
	sysAdmin := &pgmodels.User{
		Role:          constants.RoleSysAdmin,
		InstitutionID: int64(1),
	}
	instAdmin := &pgmodels.User{
		Role:          constants.RoleInstAdmin,
		InstitutionID: ownInst,
	}
	instUser := &pgmodels.User{
		Role:          constants.RoleInstUser,
		InstitutionID: ownInst,
	}
	nobody := &pgmodels.User{
		Role:          constants.RoleNone,
		InstitutionID: ownInst,
	}

	// SysAdmin has almost all permissions at all institutions
	sysAdminCannot := []string{
		string(constants.FileRequestDelete),
		string(constants.IntellectualObjectRequestDelete),
		string(constants.DeletionRequestApprove),
	}
	for _, perm := range constants.Permissions {
		if slice.Contains(constants.ForbiddenToAll, perm) || slice.Contains(sysAdminCannot, string(perm)) {
			assert.False(t, sysAdmin.HasPermission(perm, sysAdmin.InstitutionID), perm)
			assert.False(t, sysAdmin.HasPermission(perm, ownInst), perm)
			assert.False(t, sysAdmin.HasPermission(perm, otherInst), perm)
		} else {
			assert.True(t, sysAdmin.HasPermission(perm, sysAdmin.InstitutionID), perm)
			assert.True(t, sysAdmin.HasPermission(perm, ownInst), perm)
			assert.True(t, sysAdmin.HasPermission(perm, otherInst), perm)
		}
	}

	// InstAdmin has extended permissions at own inst only.
	for _, perm := range constants.Permissions {
		if constants.CheckPermission(instAdmin.Role, perm) {
			assert.True(t, instAdmin.HasPermission(perm, ownInst), perm)
		} else {
			assert.False(t, instAdmin.HasPermission(perm, ownInst), perm)
		}
		assert.False(t, instAdmin.HasPermission(perm, otherInst), perm)
	}

	// InstUser has limited permissions at own inst only.
	for _, perm := range constants.Permissions {
		if constants.CheckPermission(instUser.Role, perm) {
			assert.True(t, instUser.HasPermission(perm, ownInst), perm)
		} else {
			assert.False(t, instUser.HasPermission(perm, ownInst), perm)
		}
		assert.False(t, instUser.HasPermission(perm, otherInst), perm)
	}

	// Nobody can't do nothin.
	for _, perm := range constants.Permissions {
		assert.False(t, nobody.HasPermission(perm, sysAdmin.InstitutionID), perm)
		assert.False(t, nobody.HasPermission(perm, ownInst), perm)
		assert.False(t, nobody.HasPermission(perm, otherInst), perm)
	}
}

func TestUserIsAdmin(t *testing.T) {
	user := pgmodels.User{
		Role: constants.RoleInstAdmin,
	}
	assert.False(t, user.IsAdmin())

	user.Role = constants.RoleInstUser
	assert.False(t, user.IsAdmin())

	user.Role = constants.RoleSysAdmin
	assert.True(t, user.IsAdmin())
}

func TestOTPBackupCodes(t *testing.T) {
	db.LoadFixtures()
	// In fixtures, Inactive User with id #4 has some OTP
	// backup codes. We want to make sure these are bound
	// correctly.
	user, err := pgmodels.UserByID(4)
	require.Nil(t, err)

	codes := []string{
		"code1",
		"code2",
		"code3",
	}
	assert.Equal(t, codes, user.OTPBackupCodes)
}

func TestUserByID(t *testing.T) {
	db.LoadFixtures()
	user, err := pgmodels.UserByID(int64(1))
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.NotNil(t, user.Role)
	assert.NotNil(t, user.Institution)
}

func TestUserByEmail(t *testing.T) {
	db.LoadFixtures()
	user, err := pgmodels.UserByEmail(SysAdmin)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, constants.RoleSysAdmin, user.Role)
	assert.Equal(t, "aptrust.org", user.Institution.Identifier)
}

func TestUserByGet(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where(`"user"."name"`, "=", "Inst One Admin")
	user, err := pgmodels.UserGet(query)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Inst One Admin", user.Name)
	assert.NotNil(t, user.Role)
	assert.NotNil(t, user.Institution)
}

func TestUserSignIn_Valid(t *testing.T) {
	db.LoadFixtures()
	// Constants below are defined in models/common_test.go
	users := []string{
		SysAdmin,
		InstAdmin,
		InstUser,
	}
	for _, email := range users {
		user, err := pgmodels.UserSignIn(email, Password, "1.1.1.1")
		require.Nil(t, err, email)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "1.1.1.1", user.CurrentSignInIP)
		assert.True(t, user.SignInCount > 0)
		assert.InDelta(t, time.Now().Unix(), user.CurrentSignInAt.Unix(), 10)
		assert.NotNil(t, user.Role)
		assert.NotNil(t, user.Institution)
		oldSignInTime := user.CurrentSignInAt
		oldSignInCount := user.SignInCount

		user, err = pgmodels.UserSignIn(email, Password, "2.2.2.2")
		require.Nil(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "2.2.2.2", user.CurrentSignInIP)
		assert.True(t, user.SignInCount > oldSignInCount)
		assert.True(t, user.CurrentSignInAt.After(oldSignInTime))
		testUserSignOut(t, user)
	}
}

func testUserSignOut(t *testing.T, user *pgmodels.User) {
	signInIP := user.CurrentSignInIP
	signInTime := user.CurrentSignInAt
	err := user.SignOut()
	require.Nil(t, err)

	savedUser, err := pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, savedUser)

	assert.Empty(t, savedUser.CurrentSignInIP)
	assert.Empty(t, savedUser.CurrentSignInAt)
	assert.Empty(t, savedUser.EncryptedOTPSecret)
	assert.Empty(t, savedUser.EncryptedOTPSentAt)
	assert.Equal(t, signInIP, savedUser.LastSignInIP)

	// On Mac, these times should be identical.
	// On Linux, they can drift by a fraction of a second.
	assert.InDelta(t, signInTime.Unix(), savedUser.LastSignInAt.Unix(), 1)
}

func TestUserSignIn_Invalid(t *testing.T) {
	db.LoadFixtures()

	// User does not exist
	user, err := pgmodels.UserSignIn("noone@example.com", "xyz", "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrInvalidLogin, err)

	// User exists, wrong password
	user, err = pgmodels.UserSignIn(SysAdmin, "xyz", "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrInvalidLogin, err)
}

func TestUserSignIn_Deactivated(t *testing.T) {
	db.LoadFixtures()
	user, err := pgmodels.UserSignIn(InactiveUser, Password, "1.1.1.1")
	require.NotNil(t, err)
	require.Nil(t, user)
	assert.Equal(t, common.ErrAccountDeactivated, err)
}

func TestUserSaveDeleteUndelete(t *testing.T) {
	db.LoadFixtures()
	pwd, err := common.EncryptPassword("Duff Beer")
	require.Nil(t, err)
	user := pgmodels.User{
		Name:              "Barney Gumble",
		Email:             "barney@simpsons.kom",
		InstitutionID:     InstTest,
		Role:              constants.RoleInstUser,
		EncryptedPassword: pwd,
	}
	err = user.Save()
	require.Nil(t, err)
	require.True(t, user.ID > int64(0))
	assert.Empty(t, user.DeactivatedAt)

	err = user.Delete()
	require.Nil(t, err)
	assert.NotEmpty(t, user.DeactivatedAt)

	err = user.Undelete()
	require.Nil(t, err)
	assert.Empty(t, user.DeactivatedAt)
}

func TestUserSelect(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where("institution_id", "=", int64(2)).OrderBy("email", "asc")
	users, err := pgmodels.UserSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, users)

	// These three users are part of the fixture data.
	// There may be more, created by other tests.
	expected := []string{
		"admin@inst1.edu",
		"inactive@inst1.edu",
		"user@inst1.edu",
	}
	assert.True(t, len(users) >= len(expected))
	for _, email := range expected {
		found := false
		for _, user := range users {
			if user.Email == email {
				found = true
			}
		}
		assert.True(t, found, "%s missing from results", email)
	}
}

func TestUserValidate(t *testing.T) {
	user := pgmodels.User{}
	err := user.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrUserName, err.Errors["Name"])
	assert.Equal(t, pgmodels.ErrUserEmail, err.Errors["Email"])
	assert.Equal(t, pgmodels.ErrUserInst, err.Errors["InstitutionID"])
	assert.Equal(t, pgmodels.ErrUserRole, err.Errors["Role"])

	// PhoneNumber can be empty, so no error here
	assert.Empty(t, err.Errors["PhoneNumber"])

	// GracePeriod not required if OTPRequiredForLogin is false
	assert.Empty(t, err.Errors["GracePeriod"])

	// Check Phone and Grace Period
	user.PhoneNumber = "4321-22"
	user.OTPRequiredForLogin = true
	err = user.Validate()
	require.NotNil(t, err)
	assert.Equal(t, pgmodels.ErrUserPhone, err.Errors["PhoneNumber"])
	assert.Equal(t, pgmodels.ErrUserGracePeriod, err.Errors["GracePeriod"])

	user.Email = "not@valid"
	err = user.Validate()
	require.NotNil(t, err)
	assert.Equal(t, pgmodels.ErrUserEmail, err.Errors["Email"])

	user.Role = "Invalid Role"
	err = user.Validate()
	require.NotNil(t, err)
	assert.Equal(t, pgmodels.ErrUserRole, err.Errors["Role"])

	user.InstitutionID = InstTwo
	user.Role = constants.RoleSysAdmin
	err = user.Validate()
	require.NotNil(t, err)
	assert.Equal(t, pgmodels.ErrUserInvalidAdmin, err.Errors["Role"])

	// Now let's make a fully legit user and make sure
	// he passed validation.
	pwd, encErr := common.EncryptPassword("Okaly Dokaly, Homer!")
	require.Nil(t, encErr)

	user.Name = "Ned Flanders"
	user.Email = "flanders@simpson.kom"
	user.PhoneNumber = "+12028089999"
	user.OTPRequiredForLogin = true
	user.GracePeriod = time.Now().UTC().Add(time.Hour * 24)
	user.Role = constants.RoleInstUser
	user.InstitutionID = InstTest
	user.EncryptedPassword = pwd

	err = user.Validate()
	require.Nil(t, err)
}

func TestIsSMSUser(t *testing.T) {
	user := &pgmodels.User{
		EnabledTwoFactor:   true,
		ConfirmedTwoFactor: true,
		AuthyStatus:        constants.TwoFactorSMS,
	}
	assert.True(t, user.IsSMSUser())

	user.AuthyStatus = ""
	assert.True(t, user.IsSMSUser())

	user.AuthyStatus = constants.TwoFactorAuthy
	assert.False(t, user.IsSMSUser())

	user.AuthyStatus = constants.TwoFactorSMS
	user.ConfirmedTwoFactor = false
	assert.False(t, user.IsSMSUser())

	user.EnabledTwoFactor = false
	assert.False(t, user.IsSMSUser())
}

func TestIsAuthyOneTouchUser(t *testing.T) {
	user := &pgmodels.User{
		EnabledTwoFactor:   true,
		ConfirmedTwoFactor: true,
		AuthyStatus:        constants.TwoFactorAuthy,
	}
	assert.True(t, user.IsAuthyOneTouchUser())

	user.AuthyStatus = ""
	assert.False(t, user.IsAuthyOneTouchUser())

	user.AuthyStatus = constants.TwoFactorSMS
	assert.False(t, user.IsAuthyOneTouchUser())

	user.AuthyStatus = constants.TwoFactorAuthy
	user.ConfirmedTwoFactor = false
	assert.False(t, user.IsAuthyOneTouchUser())

	user.EnabledTwoFactor = false
	assert.False(t, user.IsAuthyOneTouchUser())
}

func TestIsTwoFactorUser(t *testing.T) {
	user := &pgmodels.User{
		EnabledTwoFactor:   true,
		ConfirmedTwoFactor: true,
	}
	assert.True(t, user.IsTwoFactorUser())

	user.ConfirmedTwoFactor = false
	assert.False(t, user.IsTwoFactorUser())

	user.EnabledTwoFactor = false
	assert.False(t, user.IsTwoFactorUser())
}

func TestTwoFactorMethod(t *testing.T) {
	user := &pgmodels.User{
		EnabledTwoFactor:   true,
		ConfirmedTwoFactor: true,
		AuthyStatus:        constants.TwoFactorSMS,
	}
	assert.Equal(t, constants.TwoFactorSMS, user.TwoFactorMethod())

	// Holdover from old system. A confirmed two-factor user
	// with empty authy status is an SMS user.
	user.AuthyStatus = ""
	assert.Equal(t, constants.TwoFactorSMS, user.TwoFactorMethod())

	user.AuthyStatus = constants.TwoFactorAuthy
	assert.Equal(t, constants.TwoFactorAuthy, user.TwoFactorMethod())

	user.EnabledTwoFactor = false
	assert.Equal(t, constants.TwoFactorNone, user.TwoFactorMethod())

	user.EnabledTwoFactor = true
	user.ConfirmedTwoFactor = false
	assert.Equal(t, constants.TwoFactorNone, user.TwoFactorMethod())
}

func TestCreateOTPToken(t *testing.T) {
	user, err := pgmodels.UserByEmail(InstUser)
	require.Nil(t, err)
	require.NotNil(t, user)

	defer user.ClearOTPSecret()

	// Start with empty OTP
	user.ClearOTPSecret()

	// Reload user to make sure OTP is empty
	user, err = pgmodels.UserByEmail(InstUser)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Empty(t, user.EncryptedOTPSecret)

	// Now generate a secret
	token, err := user.CreateOTPToken()
	require.Nil(t, err)
	require.Equal(t, 6, len(token))

	// Reload user to make sure OTP is not empty
	user, err = pgmodels.UserByEmail(InstUser)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.True(t, len(user.EncryptedOTPSecret) > 10)
}

func TestUserCountryCodeAndPhone(t *testing.T) {
	user := &pgmodels.User{
		PhoneNumber: "+13135551234",
	}
	countryCode, phone, err := user.CountryCodeAndPhone()
	require.Nil(t, err)
	assert.EqualValues(t, 1, countryCode)
	assert.Equal(t, "3135551234", phone)
}

func TestUserHasUnreadAlerts(t *testing.T) {
	db.LoadFixtures()

	// See db/fixtures/alerts_users.csv
	// In that file, users 1,2 and 3 have
	// unread alerts but other users do not.

	user1, err := pgmodels.UserByID(1)
	require.Nil(t, err)
	require.NotNil(t, user1)
	assert.True(t, user1.HasUnreadAlerts())

	user6, err := pgmodels.UserByID(6)
	require.Nil(t, err)
	require.NotNil(t, user6)
	assert.False(t, user6.HasUnreadAlerts())
}
