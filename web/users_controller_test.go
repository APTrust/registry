package web_test

import (
	"net/http"
	"testing"
	// "github.com/APTrust/registry/constants"
	// "github.com/APTrust/registry/pgmodels"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestUserShow(t *testing.T) {
	initHTTPTests(t)

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
	html := sysAdminClient.GET("/users/show/3").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, items)

	// Inst admin can see users at their own institution
	html = instAdminClient.GET("/users/show/3").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, items)

	// Inst admin cannot view the user belonging to other institution
	instUserClient.GET("/users/show/1").Expect().Status(http.StatusForbidden)

	// Regular user cannot view the user show page, even their own
	instUserClient.GET("/users/show/3").Expect().Status(http.StatusForbidden)

}

func TestUserIndex(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"New",
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
		`type="text" name="name__contains"`,
		`type="text" name="email__contains"`,
		`select name="role"`,
		`select name="institution_id"`,
		"Filter",
	}

	html := sysAdminClient.GET("/users").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, items)
	AssertMatchesAll(t, html, instUserLinks)
	AssertMatchesAll(t, html, nonInst1Links)
	AssertMatchesAll(t, html, adminFilters)
	AssertMatchesResultCount(t, html, 5)

	html = instAdminClient.GET("/users").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, items)
	AssertMatchesAll(t, html, instUserLinks)
	AssertMatchesNone(t, html, nonInst1Links)
	AssertMatchesNone(t, html, adminFilters)
	AssertMatchesResultCount(t, html, 3)

	// Regular user cannot view the user list page
	instUserClient.GET("/users").Expect().Status(http.StatusForbidden)

}

func TestUserCreateEditDeleteUndelete(t *testing.T) {
	initHTTPTests(t)
}

func TestUserSignInSignOut(t *testing.T) {
	initHTTPTests(t)
}

func TestUserChangePassword(t *testing.T) {
	initHTTPTests(t)
}
