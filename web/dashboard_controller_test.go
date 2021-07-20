package web_test

import (
	"net/http"
	"testing"
)

// This is a barebones test that ensures we get a 200 response
// and the page has all expected sections.
func TestDashboardShow(t *testing.T) {
	initHTTPTests(t)
	sections := []string{
		"Recent Work Items",
		"Notifications",
		"Deposits by Storage Option",
	}
	for _, client := range allClients {
		html := client.GET("/dashboard").Expect().
			Status(http.StatusOK).Body().Raw()
		MatchesAll(t, html, sections)
	}
}
