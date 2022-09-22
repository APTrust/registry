package webui_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserShow(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Force Password Reset",
		"Change Password",
		"Deactivate",
		"Edit",
		"Inst One User",
		"User at Institution One",
		"user@inst1.edu",
		"14345551212",
	}

	// Sys Admin can see any user
	html := testutil.SysAdminClient.GET("/users/show/3").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	// Inst admin can see users at their own institution
	html = testutil.Inst1AdminClient.GET("/users/show/3").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	// Inst admin cannot view the user belonging to other institution
	testutil.Inst1UserClient.GET("/users/show/1").Expect().Status(http.StatusForbidden)

	// Regular user cannot view the user show page, even their own
	testutil.Inst1UserClient.GET("/users/show/3").Expect().Status(http.StatusForbidden)

}

func TestUserIndex(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Add new user",
		"Name",
		"Email",
	}

	instUserLinks := []string{
		"/users/show/2",
		"/users/show/3",
		"/users/show/4",
	}

	nonInst1Links := []string{
		"/users/show/1",
		"/users/show/5",
	}

	// Sys Admin sees filters because list of all users is long.
	// Inst admin does not see filters, because most institutions
	// have only 4-6 users.
	adminFilters := []string{
		`type="text" id="name__contains" name="name__contains"`,
		`type="text" id="email__contains" name="email__contains"`,
		`select name="role"`,
		`select name="institution_id"`,
		"Filter",
	}

	html := testutil.SysAdminClient.GET("/users").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)
	testutil.AssertMatchesAll(t, html, instUserLinks)
	testutil.AssertMatchesAll(t, html, nonInst1Links)
	testutil.AssertMatchesAll(t, html, adminFilters)
	testutil.AssertMatchesResultCount(t, html, 9)

	html = testutil.Inst1AdminClient.GET("/users").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)
	testutil.AssertMatchesAll(t, html, instUserLinks)
	testutil.AssertMatchesNone(t, html, nonInst1Links)
	testutil.AssertMatchesNone(t, html, adminFilters)
	testutil.AssertMatchesResultCount(t, html, 4)

	// Regular user cannot view the user list page
	testutil.Inst1UserClient.GET("/users").Expect().Status(http.StatusForbidden)

}

func TestUserCreateEditDeleteUndelete(t *testing.T) {
	testutil.InitHTTPTests(t)

	// Make sure admins can get to this page and regular users cannot.
	testutil.SysAdminClient.GET("/users/new").Expect().Status(http.StatusOK)
	testutil.Inst1AdminClient.GET("/users/new").Expect().Status(http.StatusOK)
	testutil.Inst1UserClient.GET("/users/new").Expect().Status(http.StatusForbidden)

	formData := map[string]interface{}{
		"Name":           "Unit Test User",
		"Email":          "utest-user@inst1.edu",
		"PhoneNumber":    "+12025559815",
		"institution_id": testutil.Inst1Admin.InstitutionID,
		"Role":           constants.RoleInstUser,
	}

	testutil.Inst1AdminClient.POST("/users/new").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the new user exists and has the correct info
	user, err := pgmodels.UserByEmail("utest-user@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, formData["Name"], user.Name)
	assert.Equal(t, formData["Email"], user.Email)
	assert.Equal(t, formData["PhoneNumber"], user.PhoneNumber)
	assert.Equal(t, formData["institution_id"], user.InstitutionID)
	assert.Equal(t, formData["Role"], user.Role)
	assert.NotEmpty(t, user.EncryptedPassword)

	// Make sure we created a new user alert, so this person
	// can get in and choose a password.
	query := pgmodels.NewQuery().
		Where("type", "=", constants.AlertWelcome).
		Where("user_id", "=", user.ID).
		Limit(1)
	alertView, err := pgmodels.AlertViewGet(query)
	require.Nil(t, err)
	require.NotNil(t, alertView)

	reToken := regexp.MustCompile(`this token into the entry box: \w+`)
	assert.True(t, reToken.Match([]byte(alertView.Content)))

	// Get the edit page for the new user
	testutil.Inst1AdminClient.GET("/users/edit/{id}", user.ID).
		Expect().Status(http.StatusOK)

	// Update the user
	formData["Name"] = "Unit Test User (edited)"
	formData["PhoneNumber"] = "+15058981234"
	testutil.Inst1AdminClient.PUT("/users/edit/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByEmail("utest-user@inst1.edu")
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, formData["Name"], user.Name)
	assert.Equal(t, formData["PhoneNumber"], user.PhoneNumber)

	// Test XHR updates
	testUserUpdateXHR(t, user)

	// Delete the user. This winds up with an OK because of redirect.
	// Note that because it's a delete, there's no form, so we have
	// to pass the CSRF token in the header.
	testutil.Inst1AdminClient.DELETE("/users/delete/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithHeader(constants.CSRFHeaderName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusOK)

	// Undelete the user. Again, we get a redirect ending with an OK.
	testutil.Inst1AdminClient.POST("/users/undelete/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithHeader(constants.CSRFHeaderName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusOK)

}

func testUserUpdateXHR(t *testing.T, user *pgmodels.User) {
	// Update the user
	formData := make(map[string]interface{})
	formData["Name"] = "Named edited XHR"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err := pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, formData["Name"], user.Name)
	delete(formData, "Name")

	formData["Email"] = "xhr@example.com"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, formData["Email"], user.Email)
	delete(formData, "Name")

	oldEncryptedPassword := user.EncryptedPassword
	formData["PhoneNumber"] = "+13039998888"
	formData["Password"] = "SuperSekrit1234!"
	formData["Role"] = constants.RoleInstAdmin
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, formData["PhoneNumber"], user.PhoneNumber)
	assert.NotEqual(t, oldEncryptedPassword, user.EncryptedPassword)
	assert.Equal(t, formData["Role"], user.Role)
	delete(formData, "PhoneNumber")
	delete(formData, "Password")
	delete(formData, "Role")

	formData["Status"] = "inactive"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.NotEmpty(t, user.DeactivatedAt)

	formData["Status"] = "active"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Empty(t, user.DeactivatedAt)
	delete(formData, "Status")

	formData["OTPRequiredForLogin"] = "true"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.True(t, user.OTPRequiredForLogin)

	formData["OTPRequiredForLogin"] = "false"
	testutil.Inst1AdminClient.PUT("/users/edit_xhr/{id}", user.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(formData).Expect().Status(http.StatusOK)

	// Make sure the edits were saved
	user, err = pgmodels.UserByID(user.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.False(t, user.OTPRequiredForLogin)
}

func TestUserSignInSignOut(t *testing.T) {
	testutil.InitHTTPTests(t)

	client := testutil.GetAnonymousClient(t)

	// Make sure anonymous client can access the sign-in page
	client.GET("/").Expect().Status(http.StatusOK)

	// Make sure they can sign in and are redirected to dashboard
	html := client.POST("/users/sign_in").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField("email", "user@inst1.edu").
		WithFormField("password", "password").
		Expect().Status(http.StatusOK).Body().Raw()

	dashboardItems := []string{
		"Recent Work Items",
		"Notifications",
		"Deposits by Storage Option",
	}
	testutil.AssertMatchesAll(t, html, dashboardItems)

	// Make sure user can sign out.
	client.GET("/users/sign_out").Expect().Status(http.StatusOK)

	// After signout, attempts to access valid pages should return
	// unauthorized.
	client.GET("/dashboard").Expect().Status(http.StatusUnauthorized)
}

func TestUserSignInBadCredentials(t *testing.T) {
	testutil.InitHTTPTests(t)
	client := testutil.GetAnonymousClient(t)

	// Invalid credentials should get a 400 response
	client.POST("/users/sign_in").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField("email", "bad-email@inst1.edu").
		WithFormField("password", "invalid-password").
		Expect().Status(http.StatusBadRequest)
}

func TestUserChangePassword(t *testing.T) {
	testutil.InitHTTPTests(t)

	// After tests, restore inst1User's password so that we
	// run interactive tests in browser afterward, we can log
	// in with the usual password.
	defer restorePassword(t, testutil.Inst1User)

	originalEncrypedPwd := testutil.Inst1User.EncryptedPassword

	// Make sure user can get to the change password page
	// to change their own password.
	testutil.Inst1UserClient.GET("/users/change_password/{id}", testutil.Inst1User.ID).
		Expect().Status(http.StatusOK)

	// Submit and test the password change.
	// Password requirements: uppercase, lowercase, number, min 8 chars
	testutil.Inst1UserClient.POST("/users/change_password/{id}", testutil.Inst1User.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("NewPassword", "Password1234").
		WithFormField("ConfirmNewPassword", "Password1234").
		Expect().Status(http.StatusOK)

	user, err := pgmodels.UserByID(testutil.Inst1User.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.NotEqual(t, originalEncrypedPwd, user.EncryptedPassword)

	secondPwd := user.EncryptedPassword

	// Institutional admin can change another user's password,
	// as long as that user is at their institution.
	testutil.Inst1AdminClient.GET("/users/change_password/{id}", testutil.Inst1User.ID).
		Expect().Status(http.StatusOK)
	testutil.Inst1AdminClient.POST("/users/change_password/{id}", testutil.Inst1User.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField("NewPassword", "Password5678").
		WithFormField("ConfirmNewPassword", "Password5678").
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusOK)

	user, err = pgmodels.UserByID(testutil.Inst1User.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.NotEqual(t, secondPwd, user.EncryptedPassword)

	// inst1Admin cannot change password for anyone at inst2
	testutil.Inst1AdminClient.GET("/users/change_password/{id}", testutil.Inst2Admin.ID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.POST("/users/change_password/{id}", testutil.Inst2Admin.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField("NewPassword", "Password5678").
		WithFormField("ConfirmNewPassword", "Password5678").
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusForbidden)

	// Regular user cannot access anyone else's change password page
	// In this case, inst1User is trying to change inst1Admin
	testutil.Inst1UserClient.GET("/users/change_password/{id}", testutil.Inst1Admin.ID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.POST("/users/change_password/{id}", testutil.Inst1Admin.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField("NewPassword", "Password5678").
		WithFormField("ConfirmNewPassword", "Password5678").
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusForbidden)
}

func TestUserForcePasswordReset(t *testing.T) {
	testutil.InitHTTPTests(t)

	defer restorePassword(t, testutil.Inst1User)

	testutil.Inst1AdminClient.GET("/users/init_password_reset/{id}", testutil.Inst1User.ID).
		Expect().Status(http.StatusOK)

	// This should create an alert for the user.
	query := pgmodels.NewQuery().
		Where("type", "=", constants.AlertPasswordReset).
		Where("user_id", "=", testutil.Inst1User.ID)
	alertView, err := pgmodels.AlertViewGet(query)
	require.Nil(t, err)
	require.NotNil(t, alertView)

	// It should also set a password reset token
	user, err := pgmodels.UserByID(testutil.Inst1User.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ResetPasswordToken)

	// Extract unencrypted reset token from the URL in the alert message
	re := regexp.MustCompile(`this token into the entry box: ([^\n]+)`)
	m := re.FindAllStringSubmatch(alertView.Content, 1)
	require.True(t, len(m) > 0, "Token is missing from alert")
	require.True(t, len(m[0]) > 1, "Token is missing from alert")
	unencryptedToken := m[0][1]

	// inst1user should be able to use this token to reset
	// their password. Note that the user will arrive without
	// logging in (because they don't have their password).
	client := testutil.GetAnonymousClient(t)

	// User should be able to get to this page to enter their reset token.
	client.GET("/users/complete_password_reset/{id}", testutil.Inst1User.ID).
		Expect().Status(http.StatusOK)

	// Now check the post route with the reset token
	// First: Bad token should result in error
	client.POST("/users/complete_password_reset/{id}", testutil.Inst1User.ID).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("token", "BAD_TOKEN").
		Expect().Status(http.StatusInternalServerError)

	// Good token should succeed.
	client.POST("/users/complete_password_reset/{id}", testutil.Inst1User.ID).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("token", unencryptedToken).
		Expect().Status(http.StatusOK)

	// Need to clear this token for the next two tests, so inst1User
	// isn't being forced to complete their own password reset.
	testutil.Inst1User.ResetPasswordToken = ""
	require.Nil(t, testutil.Inst1User.Save())

	// Regular user cannot reset other user's password.
	testutil.Inst1UserClient.GET("/users/init_password_reset/{id}", testutil.Inst1Admin.ID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/users/init_password_reset/{id}", testutil.Inst2Admin.ID).
		Expect().Status(http.StatusForbidden)

	// Admin cannot reset password of anyone at other institution
	testutil.Inst1AdminClient.GET("/users/init_password_reset/{id}", testutil.Inst2Admin.ID).
		Expect().Status(http.StatusForbidden)
}

func TestUserGetAPIKey(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Your new API key is",
	}

	// Any user can get their own API key
	for _, client := range testutil.AllClients {
		html := client.POST("/users/get_api_key/{id}", testutil.UserFor[client].ID).
			WithHeader("Referer", testutil.BaseURL).
			WithFormField(constants.CSRFTokenName, testutil.TokenFor[client]).
			Expect().Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}

	// No user can get another user's API key
	testutil.Inst1AdminClient.POST("/users/get_api_key/{id}", testutil.Inst1User.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.POST("/users/get_api_key/{id}", testutil.Inst1Admin.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.POST("/users/get_api_key/{id}", testutil.Inst2Admin.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusForbidden)

}

func TestUserMyAccount(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Get API Key",
		"Change Password",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/users/my_account").
			Expect().Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}
}

func TestUserForgotPassword(t *testing.T) {
	testutil.InitHTTPTests(t)

	defer restorePassword(t, testutil.Inst2User)

	// Anonymous client should be able to reach the
	// forgot password form.
	client := testutil.GetAnonymousClient(t)
	resp := client.GET("/users/forgot_password").Expect().Status(http.StatusOK)

	// Get the CSRF token from that page and submit
	// the form with token and email address
	// fmt.Println(resp.Body().Raw())
	csrfToken := testutil.ExtractCSRFToken(t, resp.Body().Raw())

	client.POST("/users/forgot_password").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, csrfToken).
		WithFormField("email", testutil.Inst2User.Email).
		Expect().Status(http.StatusOK)

	// The POST abouve should have created an alert for the user.
	query := pgmodels.NewQuery().
		Where("type", "=", constants.AlertPasswordReset).
		Where("user_id", "=", testutil.Inst2User.ID)
	alertView, err := pgmodels.AlertViewGet(query)
	require.Nil(t, err)
	require.NotNil(t, alertView)

	// It should also set a password reset token
	user, err := pgmodels.UserByID(testutil.Inst2User.ID)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ResetPasswordToken)

	// Extract unencrypted reset token from the URL in the alert message
	re := regexp.MustCompile(`this token into the entry box: ([^\n]+)`)
	m := re.FindAllStringSubmatch(alertView.Content, 1)
	require.True(t, len(m) > 0, "Token is missing from alert")
	require.True(t, len(m[0]) > 1, "Token is missing from alert")
	unencryptedToken := m[0][1]

	// User should be able to get to this page to enter their reset token.
	client.GET("/users/complete_password_reset/{id}", testutil.Inst2User.ID).
		Expect().Status(http.StatusOK)

	// Now check the post route with the reset token
	// First: Bad token should result in error
	client.POST("/users/complete_password_reset/{id}", testutil.Inst2User.ID).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("token", "BAD_TOKEN").
		Expect().Status(http.StatusInternalServerError)

	// Good token should succeed.
	client.POST("/users/complete_password_reset/{id}", testutil.Inst2User.ID).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithFormField("token", unencryptedToken).
		Expect().Status(http.StatusOK)

	// Need to clear this token for the next two tests, so inst1User
	// isn't being forced to complete their own password reset.
	testutil.Inst2User.ResetPasswordToken = ""
	require.Nil(t, testutil.Inst2User.Save())
}

func restorePassword(t *testing.T, user *pgmodels.User) {
	encPwd, err := common.EncryptPassword("password")
	require.Nil(t, err, "After tests, error restoring password for user %s", user.Name)
	user.EncryptedPassword = encPwd
	err = user.Save()
	assert.Nil(t, err, "After tests, error restoring password for user %s", user.Name)
}
