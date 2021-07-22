package web_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionShow(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"institution1.edu",
		"Member",
		"aptrust.receiving.test.institution1.edu",
		"aptrust.restore.test.institution1.edu",
		"Inst One Admin (Admin)",
		"Inst One User (User)",
	}

	userLinks := []string{
		`href="/users/show/2"`,
		`href="/users/show/3"`,
	}

	for _, client := range allClients {
		html := client.GET("/institutions/show/2").Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, items)

		// Admins should see links to the institution's users,
		// regular users should not see these links.
		if client == sysAdminClient || client == instAdminClient {
			AssertMatchesAll(t, html, userLinks)
		} else {
			AssertMatchesNone(t, html, userLinks)
		}
	}

	// SysAdmin can see any institution.
	// Others can see only their own institution
	sysAdminClient.GET("/institutions/show/3").Expect().Status(http.StatusOK)
	instAdminClient.GET("/institutions/show/3").Expect().Status(http.StatusForbidden)
	instUserClient.GET("/institutions/show/3").Expect().Status(http.StatusForbidden)

}

func TestInstitutionCreateEditDeleteUndelete(t *testing.T) {
	initHTTPTests(t)

	// SysAdmin can get the new institution form. Others cannot.
	sysAdminClient.GET("/institutions/new").Expect().Status(http.StatusOK)
	instAdminClient.GET("/institutions/new").Expect().Status(http.StatusForbidden)
	instUserClient.GET("/institutions/new").Expect().Status(http.StatusForbidden)

	institution := &pgmodels.Institution{
		Name:            "Springfield Yooniversity",
		Identifier:      "springfield.kom",
		Type:            constants.InstTypeMember,
		ReceivingBucket: "aptrust.receiving.springfield.kom",
		RestoreBucket:   "aptrust.restore.springfield.kom",
	}

	// Only Sys Admin can create an institution
	instUserClient.POST("/institutions/new").
		WithForm(institution).
		Expect().Status(http.StatusForbidden)
	instAdminClient.POST("/institutions/new").
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

	// Note that status will be OK instead of created
	// because the controller redirects to the "show"
	// page of the newly created institution.
	sysAdminClient.POST("/institutions/new").
		WithForm(institution).
		Expect().Status(http.StatusOK)

	// Retrieve the id of the new institution
	query := pgmodels.NewQuery().Where("name", "=", institution.Name)
	institution, err := pgmodels.InstitutionGet(query)
	require.Nil(t, err)
	assert.True(t, institution.ID > 0)

	// Fix the spelling of the institution name,
	// save it, and make sure the save works.
	institution.Name = "Springfield University"
	sysAdminClient.PUT("/institutions/edit/{id}", institution.ID).
		WithForm(institution).
		Expect().Status(http.StatusOK).
		Body().Contains("Springfield University")

	// Make sure non-sysadmins cannot update this other institution.
	instUserClient.PUT("/institutions/edit/{id}", institution.ID).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)
	instAdminClient.PUT("/institutions/edit/{id}", institution.ID).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

	// SysAdmin can Delete and Undelete.
	// This should technically use the HTTP delete method,
	// but browsers don't even support that, so we're using GET
	// for the time being. Status is OK because of redirect that
	// eventually displays either the institution list or detail page.
	sysAdminClient.GET("/institutions/delete/{id}", institution.ID).
		Expect().Status(http.StatusOK)
	sysAdminClient.GET("/institutions/undelete/{id}", institution.ID).
		Expect().Status(http.StatusOK)

	// Other users cannot delete or undelete, even their own instituion
	instUserClient.GET("/institutions/delete/{id}", inst1User.InstitutionID).
		Expect().Status(http.StatusForbidden)
	instAdminClient.GET("/institutions/delete/{id}", inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)
	instUserClient.GET("/institutions/undelete/{id}", inst1User.InstitutionID).
		Expect().Status(http.StatusForbidden)
	instAdminClient.GET("/institutions/undelete/{id}", inst1Admin.InstitutionID).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

}

func TestInstitutionIndex(t *testing.T) {
	initHTTPTests(t)

	instLinks := []string{
		"/institutions/show/1",
		"/institutions/show/2",
		"/institutions/show/3",
		"/institutions/show/4",
	}

	// SysAdmin can see the institutions page
	html := sysAdminClient.GET("/institutions").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, instLinks)

	// No other roles can see this page.
	instAdminClient.GET("/institutions").Expect().Status(http.StatusForbidden)
	instUserClient.GET("/institutions").Expect().Status(http.StatusForbidden)
}
