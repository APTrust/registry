package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/require"
)

// TODO: Consider testify's Suite package as described at
// https://pkg.go.dev/github.com/stretchr/testify/suite
// Includes setup, teardown and ordered tests.

func TestAlertShow(t *testing.T) {
	testutil.InitHTTPTests(t)

	alert1, err := pgmodels.AlertByID(1)
	require.Nil(t, err)
	require.NotNil(t, alert1)

	// Sysadmin can read own alert
	testutil.SysAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, testutil.SysAdmin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Sysadmin can read copy of alert sent to inst admin
	testutil.SysAdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, testutil.Inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin can read own alert
	testutil.Inst1AdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, testutil.Inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin CANNOT read sys admin's copy of alert
	testutil.Inst1AdminClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, testutil.SysAdmin.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user CANNOT read inst admin's alert
	testutil.Inst1UserClient.GET("/alerts/show/{id}/{user_id}", alert1.ID, testutil.Inst1Admin.ID).
		Expect().
		Status(http.StatusForbidden)
}

func TestAlertIndex(t *testing.T) {
	testutil.InitHTTPTests(t)

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
	resp := testutil.SysAdminClient.GET("/alerts").Expect().Status(http.StatusOK)
	html := resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, commonFilters)
	testutil.AssertMatchesAll(t, html, sysAdminFilters)
	testutil.AssertMatchesAll(t, html, constants.AlertTypes)
	testutil.AssertMatchesAll(t, html, testutil.AllInstitutionNames(t))
	testutil.AssertMatchesAll(t, html, testutil.AllUserNames(t))
	testutil.AssertMatchesResultCount(t, html, 15)

	// Make sure filters work. Should be 1 deletion confirmed
	// alerts for the inst 1 admin.
	resp = testutil.SysAdminClient.GET("/alerts").
		WithQuery("user_id", testutil.Inst1Admin.ID).
		WithQuery("type", constants.AlertDeletionConfirmed).
		Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	testutil.AssertMatchesResultCount(t, html, 1)

	// Inst admin should see only his own alerts and the
	// alert type and date filters
	resp = testutil.Inst1AdminClient.GET("/alerts").
		WithQuery("institution_id", testutil.Inst1Admin.InstitutionID).
		WithQuery("user_id", testutil.Inst1Admin.ID).
		Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, commonFilters)
	testutil.AssertMatchesNone(t, html, sysAdminFilters)
	testutil.AssertMatchesAll(t, html, constants.AlertTypes)
	testutil.AssertMatchesResultCount(t, html, 6)

	// Inst user should see only his own alerts and the
	// alert type and date filters
	resp = testutil.Inst1UserClient.GET("/alerts").
		WithQuery("institution_id", testutil.Inst1User.InstitutionID).
		WithQuery("user_id", testutil.Inst1User.ID).
		Expect().Status(http.StatusOK)
	html = resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, commonFilters)
	testutil.AssertMatchesNone(t, html, sysAdminFilters)
	testutil.AssertMatchesAll(t, html, constants.AlertTypes)
	testutil.AssertMatchesResultCount(t, html, 2)

}
