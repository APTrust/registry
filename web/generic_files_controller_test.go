package web_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericFileShow(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"institution1.edu/photos/picture1",
		"48771",
		"image/jpeg",
		"Standard",
		"https://localhost:9899/preservation-va/25452f41-1b18-47b7-b334-751dfd5d011e",
		"https://localhost:9899/preservation-or/25452f41-1b18-47b7-b334-751dfd5d011e",
	}

	for _, client := range allClients {
		req := client.GET("/files/show/1")
		require.NotNil(t, req)
		resp := req.Expect()
		require.NotNil(t, resp)
		resp.Status(http.StatusOK)
		html := resp.Body().Raw()
		MatchesAll(t, html, items)
	}

	// This file belongs to institution 2, so sys admin
	// can see it, but inst 1 users cannot.
	sysAdminClient.GET("/files/show/17").Expect().Status(http.StatusOK)
	instAdminClient.GET("/files/show/17").Expect().Status(http.StatusForbidden)
	instUserClient.GET("/files/show/17").Expect().Status(http.StatusForbidden)
}

func TestGenericFileIndex(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"institution1.edu/photos/picture1",
		"institution1.edu/photos/picture2",
		"institution1.edu/photos/picture3",
		"institution1.edu/pdfs/doc1",
		"institution1.edu/pdfs/doc2",
		"institution1.edu/pdfs/doc3",
	}

	commonFilters := []string{
		`type="text" name="identifier"`,
		`select name="state"`,
		`select name="storage_option"`,
		`type="number" name="size__gteq"`,
		`type="number" name="size__lteq"`,
		`type="date" name="created_at__gteq"`,
		`type="date" name="created_at__gteq"`,
		`type="date" name="updated_at__gteq"`,
		`type="date" name="updated_at__gteq"`,
	}

	adminFilters := []string{
		`select name="institution_id"`,
	}

	for _, client := range allClients {
		req := client.GET("/files")
		require.NotNil(t, req)
		resp := req.Expect()
		require.NotNil(t, resp)
		resp.Status(http.StatusOK)
		html := resp.Body().Raw()
		MatchesAll(t, html, items)
		MatchesAll(t, html, commonFilters)
		if client == sysAdminClient {
			MatchesAll(t, html, adminFilters)
			MatchesResultCount(t, html, 49)
		} else {
			MatchesNone(t, html, adminFilters)
			MatchesResultCount(t, html, 11)
		}
	}

	// Test some filters
	for _, client := range allClients {
		req := client.GET("/files").
			WithQuery("size__gteq", 20000).
			WithQuery("size__lteq", 5500000)
		require.NotNil(t, req)
		resp := req.Expect()
		require.NotNil(t, resp)
		resp.Status(http.StatusOK)
		html := resp.Body().Raw()
		if client == sysAdminClient {
			MatchesResultCount(t, html, 34)
		} else {
			MatchesNone(t, html, adminFilters)
			MatchesResultCount(t, html, 10)
		}
	}

	// Sysadmin can see files from inst id 3 (or any inst id).
	// Inst 1 users cannot see files of other inst.
	sysAdminClient.GET("/files").
		WithQuery("institution_id", "3").
		Expect().Status(http.StatusOK)
	instAdminClient.GET("/files").
		WithQuery("institution_id", "3").
		Expect().Status(http.StatusForbidden)
	instUserClient.GET("/files").
		WithQuery("institution_id", "3").
		Expect().Status(http.StatusForbidden)
}

func TestGenericFileRequestDelete(t *testing.T) {
	initHTTPTests(t)
}

func TestGenericFileInitDelete(t *testing.T) {
	initHTTPTests(t)
}

func TestGenericFileRequestRestore(t *testing.T) {
	initHTTPTests(t)
}

func TestGenericFileInitRestore(t *testing.T) {
	initHTTPTests(t)
}
