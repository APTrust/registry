package web_test

import (
	"net/http"
	"testing"

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
