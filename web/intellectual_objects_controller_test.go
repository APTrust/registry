package web_test

import (
	"net/http"
	"testing"
	// "github.com/APTrust/registry/constants"
	// "github.com/APTrust/registry/pgmodels"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestObjectShow(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"First Object for Institution One",
		"institution1.edu/photos",
		"Standard",
		"323.6 kB",
		"A bag of photos",
		"/events/show_xhr/37", // link to ingest premis event
		"/objects/request_delete/1",
		"/objects/request_restore/1",
		"File Summary",
		"image/jpeg",
		"Active Files",
		"Show Deleted Files",
		"Filter",
		"institution1.edu/photos/picture1",
		"https://localhost:9899/preservation-va/25452f41-1b18-47b7-b334-751dfd5d011e",
		"https://localhost:9899/preservation-or/25452f41-1b18-47b7-b334-751dfd5d011e",
		"md5",
		"12345678",
		"sha256",
		"9876543210",
		"institution1.edu/photos/picture2",
		"institution1.edu/photos/picture3",
		"/files/request_delete/1",
		"/files/request_restore/1",
	}

	for _, client := range allClients {
		html := client.GET("/objects/show/1").Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, items)
	}
}

func TestObjectList(t *testing.T) {
	initHTTPTests(t)

	inst1Links := []string{
		"objects/show/1",
		"objects/show/2",
		"objects/show/3",
	}

	inst2Links := []string{
		"objects/show/4",
		"objects/show/5",
		"objects/show/6",
	}

	commonFilters := []string{
		`type="text" name="identifier"`,
		`type="text" name="bag_name"`,
		`type="text" name="alt_identifier"`,
		`type="text" name="bag_group_identifier"`,
		`type="text" name="internal_sender_identifier"`,
		`select name="bagit_profile_identifier"`,
		`select name="access"`,
		`type="number" name="size__gteq"`,
		`type="number" name="size__lteq"`,
		`type="number" name="file_count__gteq"`,
		`type="number" name="file_count__lteq"`,
		`type="date" name="created_at__gteq"`,
		`type="date" name="created_at__lteq"`,
		`type="date" name="updated_at__gteq"`,
		`type="date" name="updated_at__lteq"`,
	}

	adminFilters := []string{
		`select name="institution_id"`,
		`select name="institution_parent_id"`,
	}

	for _, client := range allClients {
		html := client.GET("/objects").Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, inst1Links)
		AssertMatchesAll(t, html, commonFilters)
		if client == sysAdminClient {
			AssertMatchesAll(t, html, adminFilters)
			AssertMatchesAll(t, html, inst2Links)
			AssertMatchesResultCount(t, html, 8)
		} else {
			AssertMatchesNone(t, html, adminFilters)
			AssertMatchesNone(t, html, inst2Links)
			AssertMatchesResultCount(t, html, 3)
		}
	}

}
