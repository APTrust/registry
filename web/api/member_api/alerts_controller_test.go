package memberapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertShow(t *testing.T) {
	tu.InitHTTPTests(t)

	alert1, err := pgmodels.AlertByID(1)
	require.Nil(t, err)
	require.NotNil(t, alert1)

	// Sysadmin can read own alert
	tu.SysAdminClient.GET("/member-api/v3/alerts/show/{id}/{user_id}", alert1.ID, tu.SysAdmin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Sysadmin can read copy of alert sent to inst admin
	tu.SysAdminClient.GET("/member-api/v3/alerts/show/{id}/{user_id}", alert1.ID, tu.Inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin can read own alert
	tu.InstAdminClient.GET("/member-api/v3/alerts/show/{id}/{user_id}", alert1.ID, tu.Inst1Admin.ID).
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)

	// Inst admin CANNOT read sys admin's copy of alert
	tu.InstAdminClient.GET("/member-api/v3/alerts/show/{id}/{user_id}", alert1.ID, tu.SysAdmin.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user CANNOT read inst admin's alert
	tu.InstUserClient.GET("/member-api/v3/alerts/show/{id}/{user_id}", alert1.ID, tu.Inst1Admin.ID).
		Expect().
		Status(http.StatusForbidden)
}

// ---------------------------------------------------------------------
//
// TODO: Find clean way of testing JSON data with httpexpect.
//
// Built-in methods don't work well for nested JSON and when
// tests fail, httpexpect panics and does not report actual vs.
// expected results. assert and require behave much better.
//
// ---------------------------------------------------------------------

func TestAlertIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all alerts and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/alerts").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)
	fmt.Println(resp.Body().Raw())
	resp.JSON().Object().ValueEqual("count", 15)
	resp.JSON().Object().ValueEqual("next", "/member-api/v3/alerts?page=3&per_page=5")
	//resp.JSON().Object().ValueEqual("previous", "/member-api/v1/alerts?page=1&per_page=5")
	assert.Equal(t, "/member-api/v3/alerts?page=1&per_page=5", resp.JSON().Object().Value("previous").String().Raw())

	//assert.Equal(t, "Deletion Confirmed",

	// Make sure filters work. Should be 1 deletion confirmed
	// alerts for the inst 1 admin.
	resp = tu.SysAdminClient.GET("/member-api/v3/alerts").
		WithQuery("user_id", tu.Inst1Admin.ID).
		WithQuery("type", constants.AlertDeletionConfirmed).
		Expect().Status(http.StatusOK)
	//html = resp.Body().Raw()

	// Inst admin should see only his own alerts and the
	// alert type and date filters
	resp = tu.InstAdminClient.GET("/member-api/v3/alerts").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		WithQuery("user_id", tu.Inst1Admin.ID).
		Expect().Status(http.StatusOK)
	//html = resp.Body().Raw()

	// Inst user should see only his own alerts and the
	// alert type and date filters
	resp = tu.InstUserClient.GET("/member-api/v3/alerts").
		WithQuery("institution_id", tu.Inst1User.InstitutionID).
		WithQuery("user_id", tu.Inst1User.ID).
		Expect().Status(http.StatusOK)
	// html = resp.Body().Raw()
}
