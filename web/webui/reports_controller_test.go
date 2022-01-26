package webui_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/web/testutil"
)

func TestDepositReportShow(t *testing.T) {
	testutil.InitHTTPTests(t)
	html := testutil.Inst1AdminClient.GET("/reports/deposits").
		WithQuery("institution_id", testutil.Inst1Admin.InstitutionID).
		WithQuery("storage_option", constants.StorageOptionStandard).
		WithQuery("updated_before", time.Now().UTC().Format(time.RFC3339)).
		Expect().
		Status(http.StatusOK).Body().Raw()

	expected := []string{
		"<td>Institution One</td>",
		"<td>Standard</td>",
		"<td>3</td>",
		"<td>10</td>",
		"<td>72.810</td>",
		"<td>0.071</td>",
		"<td>$1.97</td>",
	}
	testutil.AssertMatchesAll(t, html, expected)

	// User can't get report for other institution.
	// Here, admin from Inst1 is trying to get data for Inst2.
	testutil.Inst1AdminClient.GET("/reports/deposits").
		WithQuery("institution_id", testutil.Inst2Admin.InstitutionID).
		WithQuery("storage_option", constants.StorageOptionStandard).
		WithQuery("updated_before", time.Now().UTC().Format(time.RFC3339)).
		Expect().
		Status(http.StatusForbidden)

	// SysAdmin can get any inst, including 0
	html = testutil.SysAdminClient.GET("/reports/deposits").
		WithQuery("institution_id", testutil.Inst1Admin.InstitutionID).
		WithQuery("storage_option", constants.StorageOptionStandard).
		WithQuery("updated_before", time.Now().UTC().Format(time.RFC3339)).
		Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, expected)

	// SysAdmin can get stats withouth specifying inst id (all institutions).
	// Since we're not specifying storage option either, this should
	// include all options.
	expectedForInst0 := []string{
		"<td>Institution One</td>",
		"<td>Institution Two</td>",
		"<td>Test Institution (for integration tests)</td>",
		"<td>Standard</td>",
		"<td>Wasabi-OR</td>",
		"<td>Wasabi-VA</td>",
		"<td>Glacier-OR</td>",
		"<td>Glacier-Deep-OH</td>",
		"<td>Glacier-Deep-VA</td>",
		"<td>Wasabi-OR</td>",
		"<td>Total</td>",
	}

	html = testutil.SysAdminClient.GET("/reports/deposits").
		WithQuery("institution_id", 0).
		WithQuery("updated_before", time.Now().UTC().Format(time.RFC3339)).
		Expect().
		Status(http.StatusOK).Body().Raw()
	testutil.AssertMatchesAll(t, html, expectedForInst0)
}
