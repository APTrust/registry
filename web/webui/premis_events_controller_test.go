package webui_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/web/testutil"
)

func TestEventShow(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"https://github.com/APTrust/exchange",
		"Aug 26, 2016 18:53:32",
		"Calculated new fixity value",
		"message digest calculation",
		"institution1.edu/photos/picture1",
		"e2b0e887-d54d-4fd2-b4bc-71ea9311afd5",
		"Institution One",
		"SHA-256 thingy",
		"Success",
		"12e6a5fc3c144b31bcf1d781912beb00",
		"New fixididdly",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/events/show/31").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}

	// Inst 1 users cannot see event belonging to inst 2
	testutil.Inst1AdminClient.GET("/events/show/42").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/events/show/42").
		Expect().Status(http.StatusForbidden)
}

func TestEventList(t *testing.T) {
	testutil.InitHTTPTests(t)

	inst1Links := []string{
		"/events/show/37",
		"/events/show/38",
		"/events/show/39",
	}

	inst2Links := []string{
		"/events/show/51",
		"/events/show/52",
		"/events/show/53",
	}

	commonFilters := []string{
		`type="text" id="identifier" name="identifier"`,
		`type="text" id="intellectual_object_identifier" name="intellectual_object_identifier"`,
		`type="text" id="generic_file_identifier" name="generic_file_identifier"`,
		`select name="event_type" id="event_type"`,
		`select name="outcome" id="outcome"`,
		`type="date" id="date_time__gteq" name="date_time__gteq"`,
		`type="date" id="date_time__lteq" name="date_time__lteq"`,
	}

	adminFilters := []string{
		`select name="institution_id" id="institution_id"`,
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/events").
			WithQuery("sort", "date_time__desc").
			WithQuery("sort", "id__desc").
			Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, inst1Links)
		testutil.AssertMatchesAll(t, html, commonFilters)
		if client == testutil.SysAdminClient {
			testutil.AssertMatchesAll(t, html, adminFilters)
			testutil.AssertMatchesAll(t, html, inst2Links)
			testutil.AssertMatchesResultCount(t, html, 54)
		} else {
			testutil.AssertMatchesNone(t, html, adminFilters)
			testutil.AssertMatchesNone(t, html, inst2Links)
			testutil.AssertMatchesResultCount(t, html, 27)
		}
	}
}

func TestEventShowXHR(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"e2b0e887-d54d-4fd2-b4bc-71ea9311afd5",
		"message digest calculation",
		"Aug 26, 2016 18:53:32",
		"Success",
		"12e6a5fc3c144b31bcf1d781912beb00",
		"New fixididdly",
		"SHA-256 thingy",
		"https://github.com/APTrust/exchange",
		"institution1.edu/photos",
		"institution1.edu/photos/picture1",
		"Institution One",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/events/show_xhr/31").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}

	// Inst 1 users cannot see event belonging to inst 2
	testutil.Inst1AdminClient.GET("/events/show_xhr/42").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/events/show_xhr/42").
		Expect().Status(http.StatusForbidden)
}
