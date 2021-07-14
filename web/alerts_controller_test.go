package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/app"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

// TODO: Single setup with admin login, as with loadFixtures, then test.
// TODO: Sysadmin login, inst admin login, inst user login.
//
// TODO: Consider testify's Suite package as described at
// https://pkg.go.dev/github.com/stretchr/testify/suite
// Includes setup, teardown and ordered tests.

func TestGinHandler(t *testing.T) {
	engine := app.InitAppEngine(true)
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://localhost",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(engine),
			Jar:       httpexpect.NewJar(),
			Timeout:   time.Second * 3,
		},
		Reporter: httpexpect.NewAssertReporter(t),
		//Printers: []httpexpect.Printer{
		//	httpexpect.NewDebugPrinter(t, true),
		//},
	})

	signInForm := map[string]string{
		"email":    "system@aptrust.org",
		"password": "password",
	}
	e.POST("/users/sign_in").WithForm(signInForm).Expect().Status(http.StatusOK)

	alert1, err := pgmodels.AlertByID(1)
	require.Nil(t, err)
	require.NotNil(t, alert1)

	// Assert response
	e.GET("/alerts/show/1/1").
		Expect().
		Status(http.StatusOK).Body().Contains(alert1.Content)
}
