package web_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/APTrust/registry/app"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var appEngine *gin.Engine
var sysAdminClient *httpexpect.Expect
var instAdminClient *httpexpect.Expect
var instUserClient *httpexpect.Expect
var allClients []*httpexpect.Expect

var sysAdmin *pgmodels.User
var inst1Admin *pgmodels.User
var inst1User *pgmodels.User

var allInstNames []string
var allUserNames []string

func initHTTPTests(t *testing.T) {
	if appEngine == nil {
		err := db.LoadFixtures()
		require.Nil(t, err)

		appEngine = app.InitAppEngine(true)
		sysAdminClient = initClient(t, "system@aptrust.org", appEngine)
		instAdminClient = initClient(t, "admin@inst1.edu", appEngine)
		instUserClient = initClient(t, "user@inst1.edu", appEngine)
		allClients = []*httpexpect.Expect{
			sysAdminClient,
			instAdminClient,
			instUserClient,
		}

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

		// We can choose NewAssertReporter, which continues after
		// test failures, or NewRequireReporter, which stops on failure.
		// The AssertReporter leads to panics, so let's use require for now.
		Reporter: httpexpect.NewRequireReporter(t),
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

// OptionLabels returns the text labels from a list of HTML options.
func OptionLabels(options []forms.ListOption) []string {
	labels := make([]string, len(options))
	for i, opt := range options {
		labels[i] = opt.Text
	}
	return labels
}

// AllInstitutionNames returns a list of all institution names
// in our test data, in no guaranteed order. We use this to ensure
// that pages containing institution lists do indeed display all
// institutions.
func AllInstitutionNames(t *testing.T) []string {
	if len(allInstNames) == 0 {
		options, err := forms.ListInstitutions(false)
		require.Nil(t, err)
		allInstNames = OptionLabels(options)
	}
	return allInstNames
}

// AllUserNames returns a list of all user names in our test data,
// in no guaranteed order. We use this to ensure that pages containing
// user lists do indeed display all users.
func AllUserNames(t *testing.T) []string {
	if len(allUserNames) == 0 {
		options, err := forms.ListUsers(0)
		require.Nil(t, err)
		allUserNames = OptionLabels(options)
	}
	return allUserNames
}

// InstUserNames returns the names of all users at an institution.
func InstUserNames(t *testing.T, institutionID int64) []string {
	options, err := forms.ListUsers(0)
	require.Nil(t, err)
	return OptionLabels(options)
}

// Note on match functions:
// httpexpect.String includes good matching functions, but they
// don't behave well in loops. We get panics instead of proper
// test failure reports.

// MatchesAll asserts that all strings in items appear in body.
func MatchesAll(t *testing.T, body string, items []string) {
	for _, item := range items {
		assert.True(t, strings.Contains(body, item), "Missing expected string: %s", item)
	}
}

// MatchesNone asserts that no strings in items appear in body.
func MatchesNone(t *testing.T, body string, items []string) {
	for _, item := range items {
		assert.False(t, strings.Contains(body, item), "Found unexpected string: %s", item)
	}
}

// MatchResult count asserts that the result count at the bottom of
// a list/index page matches the expected count. In the HTML pager,
// the result count appears in the format "1 - 20 of 215".
func MatchesResultCount(t *testing.T, body string, count int) {
	countStr := fmt.Sprintf("%d", count)
	re := regexp.MustCompile(`\d+ - \d+ of (\d+)`)
	matches := re.FindAllStringSubmatch(body, 1)
	assert.NotNil(t, matches, "Did not find result count string '1 - N of N'")
	if matches != nil {
		assert.Equal(t, countStr, matches[0][1], "Expected result count %d; got %s. Full string: %s", count, matches[0][1], matches[0][0])
	}
}
