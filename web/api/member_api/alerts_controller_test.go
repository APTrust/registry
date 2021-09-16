package memberapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
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

func TestAlertIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all alerts and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/alerts").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.AlertViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 15, list.Count)
	assert.Equal(t, "/member-api/v3/alerts?page=3&per_page=5", list.Next)
	assert.Equal(t, "/member-api/v3/alerts?page=1&per_page=5", list.Previous)
	assert.Equal(t, "Deletion Confirmed", list.Results[0].Subject)

	// Make sure filters work. Should be 3 deletion requested
	// alerts for the inst 1 admin.
	resp = tu.SysAdminClient.GET("/member-api/v3/alerts").
		WithQuery("user_id", tu.Inst1Admin.ID).
		WithQuery("type", constants.AlertDeletionRequested).
		Expect().Status(http.StatusOK)

	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 3, len(list.Results))
	assert.Equal(t, "Deletion Requested", list.Results[0].Subject)

	// Inst admin should see only his own alerts.
	resp = tu.InstAdminClient.GET("/member-api/v3/alerts").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 6, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 6, len(list.Results))
	for _, alert := range list.Results {
		assert.Equal(t, tu.Inst1Admin.ID, alert.UserID)
	}

	// Inst admin cannot see results for other insitutions.
	resp = tu.InstAdminClient.GET("/member-api/v3/alerts").
		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own alerts.
	resp = tu.InstUserClient.GET("/member-api/v3/alerts").
		WithQuery("institution_id", tu.Inst1User.InstitutionID).
		WithQuery("user_id", tu.Inst1User.ID).
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, "", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 2, len(list.Results))
	for _, alert := range list.Results {
		assert.Equal(t, tu.Inst1User.ID, alert.UserID)
	}

	// Inst user cannot see other institution's alerts.
	resp = tu.InstUserClient.GET("/member-api/v3/alerts").
		WithQuery("institution_id", tu.Inst2Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
