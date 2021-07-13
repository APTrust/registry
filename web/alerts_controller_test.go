package web_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/app"
	"github.com/gavv/httpexpect/v2"
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
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	signInForm := map[string]string{
		"email":    "system@aptrust.org",
		"password": "password",
	}
	e.POST("/users/sign_in").WithForm(signInForm).Expect().Status(http.StatusOK)

	// Assert response
	e.GET("/alerts/show/1/1").
		Expect().
		Status(http.StatusOK)
}
