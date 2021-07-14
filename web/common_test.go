package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/app"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var appEngine *gin.Engine
var sysAdminClient *httpexpect.Expect
var instAdminClient *httpexpect.Expect
var instUserClient *httpexpect.Expect

var sysAdmin *pgmodels.User
var inst1Admin *pgmodels.User
var inst1User *pgmodels.User

func initHTTPTests(t *testing.T) {
	if appEngine == nil {
		err := db.LoadFixtures()
		require.Nil(t, err)

		appEngine = app.InitAppEngine(true)
		sysAdminClient = initClient(t, "system@aptrust.org", appEngine)
		instAdminClient = initClient(t, "admin@inst1.edu", appEngine)
		instUserClient = initClient(t, "user@inst1.edu", appEngine)

		sysAdmin = initUser(t, "system@aptrust.org")
		inst1Admin = initUser(t, "admin@inst1.edu")
		inst1User = initUser(t, "user@inst1.edu")
	}
}

func initClient(t *testing.T, email string, engine *gin.Engine) *httpexpect.Expect {
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

	// Create sign-in data for the requested user.
	// Remember that in fixtures, the password for all users
	// is "password".
	signInForm := map[string]string{
		"email":    email,
		"password": "password",
	}

	// Sign the user in, and be sure we got on OK.
	// The client cookie jar will store the session
	// cookie for this user.
	e.POST("/users/sign_in").WithForm(signInForm).Expect().Status(http.StatusOK)

	return e
}

func initUser(t *testing.T, email string) *pgmodels.User {
	user, err := pgmodels.UserByEmail(email)
	require.Nil(t, err)
	require.NotNil(t, user)
	return user
}
