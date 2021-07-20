package web_test

import (
	"net/http"
	"testing"

	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
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
		html := client.GET("/files/show/1").Expect().
			Status(http.StatusOK).Body().Raw()
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
		html := client.GET("/files").Expect().
			Status(http.StatusOK).Body().Raw()
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
		html := client.GET("/files").
			WithQuery("size__gteq", 20000).
			WithQuery("size__lteq", 5500000).
			Expect().
			Status(http.StatusOK).Body().Raw()
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

	items := []string{
		"Are you sure you want to delete this file?",
		"institution1.edu/photos/picture1",
		"Cancel",
		"Confirm",
	}

	// Users can request deletions of their own files
	for _, client := range allClients {
		html := client.GET("/files/request_delete/1").
			Expect().Status(http.StatusOK).Body().Raw()
		MatchesAll(t, html, items)
	}

	// Sys Admin can request any deletion, but others cannot
	// request deletions outside their own institution.
	// File 18 from fixtures belongs to Inst2
	sysAdminClient.GET("/files/request_delete/18").
		Expect().Status(http.StatusOK)
	instAdminClient.GET("/files/request_delete/18").
		Expect().Status(http.StatusForbidden)
	instUserClient.GET("/files/request_delete/18").
		Expect().Status(http.StatusForbidden)
}

func TestGenericFileInitDelete(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"Deletion Requested",
		"institution1.edu/pdfs/doc1",
	}

	// User at inst 1 can initiate deletion of their own
	// institution's file.
	html := instUserClient.POST("/files/init_delete/4").Expect().Status(http.StatusCreated).Body().Raw()
	MatchesAll(t, html, items)

	// This should create a deletion request...
	query := pgmodels.NewQuery().Where("generic_file_id", "=", 4)
	drgf := pgmodels.DeletionRequestsGenericFiles{}
	err := query.Select(&drgf)
	require.Nil(t, err)
	require.NotEqual(t, int64(0), drgf.DeletionRequestID)

	deletionRequest, err := pgmodels.DeletionRequestByID(drgf.DeletionRequestID)
	require.Nil(t, err)
	require.NotNil(t, deletionRequest)

	// Make sure this is the request our test user just made
	require.Equal(t, inst1User.ID, deletionRequest.RequestedByID)

	// There should also be an alert...
	query = pgmodels.NewQuery().Where("deletion_request_id", "=", drgf.DeletionRequestID).Relations("DeletionRequest", "Users")
	alert, err := pgmodels.AlertGet(query)
	require.Nil(t, err)
	assert.NotNil(t, alert)
	assert.NotNil(t, alert.DeletionRequest)
	assert.True(t, len(alert.Users) > 0)

	// The user should NOT be able to initiate deletion of a file
	// that belongs to another institution. In fixture data, file
	// 34 belongs to inst 2.
	instUserClient.POST("/files/init_delete/34").Expect().Status(http.StatusForbidden)
}

func TestGenericFileRequestRestore(t *testing.T) {
	initHTTPTests(t)
}

func TestGenericFileInitRestore(t *testing.T) {
	initHTTPTests(t)
}
