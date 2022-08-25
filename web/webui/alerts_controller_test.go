package webui_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/APTrust/registry/web/webui"
	"github.com/stretchr/testify/assert"
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

func TestMarkReadAndUnread(t *testing.T) {
	// Get all the alerts for user id 2 - Inst One Admin
	// These come from the fixture data.
	db := common.Context().DB
	userAlerts := db.Model((*pgmodels.AlertsUsers)(nil)).
		ColumnExpr("alert_id").
		Where("user_id = ?", testutil.Inst1Admin.ID)
	var alerts []*pgmodels.Alert
	err := db.Model(&alerts).
		Where("id IN (?)", userAlerts).
		Order("id desc").
		Select(&alerts)

	require.Nil(t, err)
	require.True(t, len(alerts) > 2)

	alertIDs := make([]int64, len(alerts))
	for i, alert := range alerts {
		alertIDs[i] = alert.ID
	}

	resetAlertsToUnread(t, alerts)
	testMarkAlertsRead(t, alerts, alertIDs)
	testMarkAlertsUnread(t, alerts, alertIDs)
	testMarkAllAlertsRead(t, alerts, alertIDs)
}

func resetAlertsToUnread(t *testing.T, alerts []*pgmodels.Alert) {
	for _, alert := range alerts {
		require.Nil(t, alert.MarkAsUnread(testutil.Inst1Admin.ID))
	}
}

func testMarkAlertsRead(t *testing.T, alerts []*pgmodels.Alert, alertIDs []int64) {
	resp := testutil.Inst1AdminClient.PUT("/alerts/mark_as_read").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithFormField("id__in", alertIDs[0]).
		WithFormField("id__in", alertIDs[1]).
		WithFormField("id__in", alertIDs[2]).
		Expect()
	body := resp.Body().Raw()
	resp.Status(http.StatusOK)
	result := &webui.AlertReadResult{}
	err := json.Unmarshal([]byte(body), result)
	require.Nil(t, err)
	assert.Equal(t, 3, len(result.Succeeded))
	assert.Empty(t, result.Failed)
	assert.Empty(t, result.Error)

	// Make sure controller marked these as succeeded
	// and they were really changed in the database.
	for i := 0; i < 3; i++ {
		id := alertIDs[i]
		assert.Contains(t, result.Succeeded, id)
		alertView, err := pgmodels.AlertViewForUser(id, testutil.Inst1Admin.ID)
		require.Nil(t, err)
		assert.NotEmpty(t, alertView.ReadAt)
	}
}

func testMarkAlertsUnread(t *testing.T, alerts []*pgmodels.Alert, alertIDs []int64) {
	resp := testutil.Inst1AdminClient.PUT("/alerts/mark_as_unread").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		WithFormField("id__in", alertIDs[0]).
		WithFormField("id__in", alertIDs[1]).
		WithFormField("id__in", alertIDs[2]).
		Expect()
	body := resp.Body().Raw()
	resp.Status(http.StatusOK)
	result := &webui.AlertReadResult{}
	err := json.Unmarshal([]byte(body), result)
	require.Nil(t, err)
	assert.Equal(t, 3, len(result.Succeeded))
	assert.Empty(t, result.Failed)
	assert.Empty(t, result.Error)

	// Make sure controller marked these as succeeded
	// and they were really changed in the database.
	for i := 0; i < 3; i++ {
		id := alertIDs[i]
		assert.Contains(t, result.Succeeded, id)
		alertView, err := pgmodels.AlertViewForUser(id, testutil.Inst1Admin.ID)
		require.Nil(t, err)
		assert.Empty(t, alertView.ReadAt)
	}
}

func testMarkAllAlertsRead(t *testing.T, alerts []*pgmodels.Alert, alertIDs []int64) {
	// Make sure this user has unread alerts
	resetAlertsToUnread(t, alerts)
	query := pgmodels.NewQuery().
	Columns("id").
	Where("user_id", "=", testutil.Inst1Admin.ID).
	IsNull("read_at")
	alertViews, err := pgmodels.AlertViewSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, alertViews)

	// Hit the endpoint that marks them all as read
	// and verify that we get an OK response.
	resp := testutil.Inst1AdminClient.POST("/alerts/mark_all_as_read").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect()
	resp.Status(http.StatusOK)

	// Make sure they really are marked as read
	alertViews, err = pgmodels.AlertViewSelect(query)
	require.Nil(t, err)
	require.Empty(t, alertViews)
}
