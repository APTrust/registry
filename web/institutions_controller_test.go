package web_test

import (
	"net/http"
	"testing"
	//"github.com/APTrust/registry/constants"
	//"github.com/APTrust/registry/pgmodels"
	//"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestInstitutionShow(t *testing.T) {
	initHTTPTests(t)
}

func TestInstitutionCreateEditDeleteUndelete(t *testing.T) {
	initHTTPTests(t)
}

func TestInstitutionIndex(t *testing.T) {
	initHTTPTests(t)

	instLinks := []string{
		"/institutions/show/1",
		"/institutions/show/2",
		"/institutions/show/3",
		"/institutions/show/4",
	}

	// SysAdmin can see the institutions page
	html := sysAdminClient.GET("/institutions").Expect().
		Status(http.StatusOK).Body().Raw()
	AssertMatchesAll(t, html, instLinks)

	// No other roles can see this page.
	instAdminClient.GET("/institutions").Expect().Status(http.StatusForbidden)
	instUserClient.GET("/institutions").Expect().Status(http.StatusForbidden)
}
