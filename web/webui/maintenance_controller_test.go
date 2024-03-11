package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMaintenanceIndex(t *testing.T) {
	testutil.InitHTTPTests(t)
	defer func() { common.Context().Config.MaintenanceMode = false }()

	// Make sure we can get there even if we're not logged in.
	client := testutil.GetAnonymousClient(t)
	resp := client.GET("/maintenance").Expect()
	resp.Status(http.StatusOK)
	assert.Contains(t, resp.Body().Raw(), "The APTrust Registry is operating as expected.")
	assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")

	// If Config.MaintenanceMode is false, make sure we get
	// the right HTML and JSON responses.
	// testutil.InitHTTPTests(t)
	for _, client := range testutil.AllClients {
		resp := client.GET("/maintenance").Expect()
		resp.Status(http.StatusOK)
		assert.Contains(t, resp.Body().Raw(), "The APTrust Registry is operating as expected.")
		assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")

		// This should NOT get a redirect
		resp = client.GET("/objects").Expect()
		resp.Status(http.StatusOK)
		assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is operating as expected.")
		assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")

		// Neither should this
		resp = client.GET("/member-api/v3/objects").Expect()
		resp.Status(http.StatusOK)
		assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is operating as expected.")
		assert.NotContains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")
	}

	// If Config.MaintenanceMode is true, make sure we get
	// the right HTML and JSON responses.
	common.Context().Config.MaintenanceMode = true
	for _, client := range testutil.AllClients {
		resp := client.GET("/maintenance").Expect()
		resp.Status(http.StatusServiceUnavailable)
		assert.Contains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")

		// This should get a redirect
		resp = client.GET("/objects").Expect()
		resp.Status(http.StatusServiceUnavailable)
		assert.Contains(t, resp.Body().Raw(), "The APTrust Registry is undergoing system maintenance.")

		// This should get a redirect and a JSON response
		resp = client.GET("/member-api/v3/objects").Expect()
		resp.Status(http.StatusServiceUnavailable)
		assert.Equal(t, `{"Error":"APTrust Registry is currently undergoing maintenance.","Message":"","StatusCode":"503"}`, resp.Body().Raw())
	}

}
