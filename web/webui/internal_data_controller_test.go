package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/web/testutil"
)

func InternalDataControllerTest(t *testing.T) {
	testutil.InitHTTPTests(t)
	for _, client := range testutil.AllClients {
		if client == testutil.SysAdminClient {
			client.GET("/internal_metadata").Expect().Status(http.StatusOK)
		} else {
			client.GET("/internal_metadata").Expect().Status(http.StatusForbidden)
		}
	}
}
