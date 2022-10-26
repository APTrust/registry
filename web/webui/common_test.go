package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/web/testutil"
)

func TestDefaultHeaders(t *testing.T) {
	testutil.InitHTTPTests(t)

	secure := []string{
		"/users/sign_in",
		"/dashboard",
		"/objects",
	}

	public := []string{
		"/static/css/registry.css",
		"/static/js/registry.js",
	}

	// Do not cache secure pages
	for _, endpoint := range secure {
		resp := testutil.SysAdminClient.GET(endpoint).Expect()
		resp.Status(http.StatusOK)
		resp.Header("Cache-Control").Equal("no-cache")
		resp.Header("Pragma").Equal("no-store")
	}

	// Do cache public resources
	for _, endpoint := range public {
		resp := testutil.SysAdminClient.GET(endpoint).Expect()
		resp.Header("Cache-Control").Equal("")
		resp.Header("Pragma").Equal("")
	}
}
