package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionShow(t *testing.T) {
	testutil.InitHTTPTests(t)

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

	for _, client := range testutil.AllClients {
		html := client.GET("/institutions/show/2").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)

		// Admins should see links to the institution's users,
		// regular users should not see these links.
		if client == testutil.SysAdminClient || client == testutil.Inst1AdminClient {
			testutil.AssertMatchesAll(t, html, userLinks)
		} else {
			testutil.AssertMatchesNone(t, html, userLinks)
		}
	}

	// SysAdmin can see any institution.
	// Others can see only their own institution
	testutil.SysAdminClient.GET("/institutions/show/3").Expect().Status(http.StatusOK)
	testutil.Inst1AdminClient.GET("/institutions/show/3").Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/institutions/show/3").Expect().Status(http.StatusForbidden)

}

func TestInstitutionCreateEditDeleteUndelete(t *testing.T) {
	testutil.InitHTTPTests(t)

	// SysAdmin can get the new institution form. Others cannot.
	testutil.SysAdminClient.GET("/institutions/new").Expect().Status(http.StatusOK)
	testutil.Inst1AdminClient.GET("/institutions/new").Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/institutions/new").Expect().Status(http.StatusForbidden)

	institution := &pgmodels.Institution{
		Name:            "Springfield Yooniversity",
		Identifier:      "springfield.kom",
		Type:            constants.InstTypeMember,
		ReceivingBucket: "aptrust.receiving.springfield.kom",
		RestoreBucket:   "aptrust.restore.springfield.kom",
	}

	// Only Sys Admin can create an institution
	testutil.Inst1UserClient.POST("/institutions/new").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.POST("/institutions/new").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

	// Note that status will be OK instead of created
	// because the controller redirects to the "show"
	// page of the newly created institution.
	testutil.SysAdminClient.POST("/institutions/new").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
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
	testutil.SysAdminClient.PUT("/institutions/edit/{id}", institution.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.SysAdminToken).
		WithForm(institution).
		Expect().Status(http.StatusOK).
		Body().Contains("Springfield University")

	// Make sure non-sysadmins cannot update this other institution.
	testutil.Inst1UserClient.PUT("/institutions/edit/{id}", institution.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.PUT("/institutions/edit/{id}", institution.ID).
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

	// SysAdmin can Delete and Undelete.
	// This should technically use the HTTP delete method,
	// but browsers don't even support that, so we're using GET
	// for the time being. Status is OK because of redirect that
	// eventually displays either the institution list or detail page.
	testutil.SysAdminClient.GET("/institutions/delete/{id}", institution.ID).
		Expect().Status(http.StatusOK)
	testutil.SysAdminClient.GET("/institutions/undelete/{id}", institution.ID).
		Expect().Status(http.StatusOK)

	// Other users cannot delete or undelete, even their own instituion
	testutil.Inst1UserClient.GET("/institutions/delete/{id}", testutil.Inst1User.InstitutionID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.GET("/institutions/delete/{id}", testutil.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/institutions/undelete/{id}", testutil.Inst1User.InstitutionID).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.GET("/institutions/undelete/{id}", testutil.Inst1Admin.InstitutionID).
		WithForm(institution).
		Expect().Status(http.StatusForbidden)

}

func TestInstitutionIndex(t *testing.T) {
	testutil.InitHTTPTests(t)

	instLinks := []string{
		"/institutions/show/1",
		"/institutions/show/2",
		"/institutions/show/3",
		"/institutions/show/4",
	}

	// SysAdmin can see the institutions page
	html := testutil.SysAdminClient.GET("/institutions").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, instLinks)

	// No other roles can see this page.
	testutil.Inst1AdminClient.GET("/institutions").Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/institutions").Expect().Status(http.StatusForbidden)
}

func TestInstitutionEdit(t *testing.T) {
	testutil.InitHTTPTests(t)

	instLinks := []string{
		`name="Name" value="Institution One"`,
		`name="Identifier" value="institution1.edu"`,
		"aptrust.receiving.test.institution1.edu",
	}

	// SysAdmin can see the institutions page
	html := testutil.SysAdminClient.GET("/institutions/edit/2").Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, instLinks)

	// No other roles can edit institutions, even their own institution.
	testutil.Inst1AdminClient.GET("/institutions/edit/2").Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/institutions/edit/2").Expect().Status(http.StatusForbidden)
	testutil.Inst2AdminClient.GET("/institutions/edit/2").Expect().Status(http.StatusForbidden)
	testutil.Inst2UserClient.GET("/institutions/edit/2").Expect().Status(http.StatusForbidden)
}
