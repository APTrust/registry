package web_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/require"
)

// TODO: Consider testify's Suite package as described at
// https://pkg.go.dev/github.com/stretchr/testify/suite
// Includes setup, teardown and ordered tests.

func TestAlertShow(t *testing.T) {
	initHTTPTests(t)

	alert1, err := pgmodels.AlertByID(1)
	require.Nil(t, err)
	require.NotNil(t, alert1)

	// Sysadmin can read own alert
	sysAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, sysAdmin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Sysadmin can read copy of alert sent to inst admin
	sysAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin can read own alert
	instAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin CANNOT read sys admin's copy of alert
	instAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, sysAdmin.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user CANNOT read inst admin's alert
	instUserClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, inst1Admin.ID).
		Expect().
		Status(http.StatusForbidden)
}

func TestAlertIndex(t *testing.T) {
	initHTTPTests(t)

	// All users should see these filters on the index page.
	commonFilters := []string{
		`select name="type"`,
		`name="created_at__gteq"`,
		`name="created_at__lteq"`,
	}

	// Only sys admin should see these filters.
	sysAdminFilters := []string{
		`select name="user_id"`,
		`select name="institution_id"`,
	}

	// Sys Admin should see all alerts and filters
	resp := sysAdminClient.GET("/alerts").Expect().Status(http.StatusOK)
	html := resp.Body().Raw()
	AssertMatchesAll(t, html, commonFilters)
	AssertMatchesAll(t, html, sysAdminFilters)
	AssertMatchesAll(t, html, constants.AlertTypes)
	AssertMatchesAll(t, html, AllInstitutionNames(t))
	AssertMatchesAll(t, html, AllUserNames(t))
	AssertMatchesResultCount(t, html, 15)

	// Make sure filters work. Should be 1 deletion confirmed
	// alerts for the inst 1 admin.
	resp = sysAdminClient.GET("/alerts").
		WithQuery("user_id", inst1Admin.ID).
		WithQuery("type", constants.AlertDeletionConfirmed).
		Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	AssertMatchesResultCount(t, html, 1)

	// Inst admin should see only his own alerts and the
	// alert type and date filters
	resp = instAdminClient.GET("/alerts").WithQuery("institution_id", inst1Admin.InstitutionID).WithQuery("user_id", inst1Admin.ID).Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	AssertMatchesAll(t, html, commonFilters)
	AssertMatchesNone(t, html, sysAdminFilters)
	AssertMatchesAll(t, html, constants.AlertTypes)
	AssertMatchesResultCount(t, html, 6)

	// Inst user should see only his own alerts and the
	// alert type and date filters
	resp = instUserClient.GET("/alerts").WithQuery("institution_id", inst1User.InstitutionID).WithQuery("user_id", inst1User.ID).Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	AssertMatchesAll(t, html, commonFilters)
	AssertMatchesNone(t, html, sysAdminFilters)
	AssertMatchesAll(t, html, constants.AlertTypes)
	AssertMatchesResultCount(t, html, 2)

}
