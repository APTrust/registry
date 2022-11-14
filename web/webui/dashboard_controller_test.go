package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/web/testutil"
)

// This is a barebones test that ensures we get a 200 response
// and the page has all expected sections.
func TestDashboardShow(t *testing.T) {
	testutil.InitHTTPTests(t)
	sections := []string{
		"Recent Work Items",
		"Notifications",
		"Deposits by Storage Option", // --> Temporarily turned off due to performance on prod
	}
	for _, client := range testutil.AllClients {
		html := client.GET("/dashboard").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, sections)
	}
}
