package web_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
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
		req := client.GET("/dashboard")
		require.NotNil(t, req)
		resp := req.Expect()
		require.NotNil(t, resp)
		resp.Status(http.StatusOK)
		html := resp.Body().Raw()
		MatchesAll(t, html, sections)
	}
}
