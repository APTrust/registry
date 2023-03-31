package webui_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectShow(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"First Object for Institution One",
		"institution1.edu/photos",
		"Standard",
		"1.5 GB",
		"A bag of photos",
		"/events/show_xhr/37", // link to ingest premis event
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
		"/files/request_restore/1",
	}

	// Only admins see deletion links
	adminOnlyItems := []string{
		"/objects/request_delete/1",
		"/files/request_delete/1",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/objects/show/1").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
		if client == testutil.Inst2AdminClient {
			testutil.AssertMatchesAll(t, html, adminOnlyItems)
		}
	}

	// inst 1 users cannot see objects belonging to inst 2
	testutil.Inst1AdminClient.GET("/objects/show/6").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/objects/show/6").
		Expect().Status(http.StatusForbidden)

}

func TestObjectList(t *testing.T) {
	testutil.InitHTTPTests(t)

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
		`type="text" id="identifier" name="identifier"`,
		`type="text" id="bag_name" name="bag_name"`,
		`type="text" id="alt_identifier" name="alt_identifier"`,
		`type="text" id="bag_group_identifier" name="bag_group_identifier"`,
		`type="text" id="internal_sender_identifier" name="internal_sender_identifier"`,
		`select name="bagit_profile_identifier"`,
		`select name="access"`,
		`type="number" id="size__gteq" name="size__gteq"`,
		`type="number" id="size__lteq" name="size__lteq"`,
		`type="number" id="file_count__gteq" name="file_count__gteq"`,
		`type="number" id="file_count__lteq" name="file_count__lteq"`,
		`type="date" id="created_at__gteq" name="created_at__gteq"`,
		`type="date" id="created_at__lteq" name="created_at__lteq"`,
		`type="date" id="updated_at__gteq" name="updated_at__gteq"`,
		`type="date" id="updated_at__lteq" name="updated_at__lteq"`,
	}

	adminFilters := []string{
		`select name="institution_id"`,
		`select name="institution_parent_id"`,
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/objects").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, inst1Links)
		testutil.AssertMatchesAll(t, html, commonFilters)
		if client == testutil.SysAdminClient {
			testutil.AssertMatchesAll(t, html, adminFilters)
			testutil.AssertMatchesAll(t, html, inst2Links)
			testutil.AssertMatchesResultCount(t, html, 14)
		} else {
			testutil.AssertMatchesNone(t, html, adminFilters)
			testutil.AssertMatchesNone(t, html, inst2Links)
			testutil.AssertMatchesResultCount(t, html, 6)
		}
	}

}

// If we search by identifier and get a single match,
// we should be redirected to the object/show page.
// Ideally, we'd test the location of the redirect using
// resp.Raw().Location(), but that always returns an error,
// so we test the page contents instead.
func TestObjectByIdentifier(t *testing.T) {
	testutil.InitHTTPTests(t)
	items := []string{
		"Object Event History",
		"File Summary",
		"etagforinst1photos",
		"institution1.edu/photos/picture1",
		"institution1.edu/photos/picture2",
		"institution1.edu/photos/picture3",
	}

	// This matches one object, so we should be redirected
	// to the object detail page, which will include the
	// strings above.
	resp := testutil.Inst1UserClient.GET("/objects").WithQuery("identifier", "institution1.edu/photos").Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)
	html := resp.Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	// This should show an empty results page that includes
	// the identifier filter we just applied.
	resp = testutil.Inst1UserClient.GET("/objects").WithQuery("identifier", "institution1.edu/does-not-exist").Expect()
	assert.Equal(t, http.StatusOK, resp.Raw().StatusCode)

	expected := `<input class="input" type="text" id="identifier" name="identifier" value="institution1.edu/does-not-exist" placeholder="Object Identifier"`
	html = resp.Body().Raw()
	assert.Contains(t, html, expected)
}

func TestObjectRequestDelete(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Are you sure you want to delete this object and its files?",
		"institution1.edu/pdfs",
		"pdf_docs_with_lots_of_words",
		"3 files / 58.4 GB",
		"Cancel",
		"Confirm",
	}

	// Sys Admin cannot request deletion
	testutil.SysAdminClient.GET("/objects/request_delete/2").
		Expect().Status(http.StatusForbidden)

	// But Inst Admin can, if item belongs to his own institution
	resp := testutil.Inst1AdminClient.GET("/objects/request_delete/2").
		Expect().Status(http.StatusOK)
	testutil.AssertMatchesAll(t, resp.Body().Raw(), items)

	testutil.Inst1UserClient.GET("/objects/request_delete/2").
		Expect().Status(http.StatusForbidden)

	// No one can request deletions outside their own institution.
	// Object 6 from fixtures belongs to Inst2
	testutil.SysAdminClient.GET("/objects/request_delete/6").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1AdminClient.GET("/objects/request_delete/6").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/objects/request_delete/6").
		Expect().Status(http.StatusForbidden)

}

func TestObjectInitDelete(t *testing.T) {
	// Force fixture reload to prevent "pending work item"
	// error when requesting deletion.
	err := db.ForceFixtureReload()
	require.Nil(t, err)
	testutil.InitHTTPTests(t)

	items := []string{
		"Deletion Requested",
		"institution1.edu/pdfs",
	}

	// Admin at inst 1 can initiate deletion of their own
	// institution's object.
	html := testutil.Inst1AdminClient.POST("/objects/init_delete/2").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusCreated).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	// This should create a deletion request. The fixture data
	// already includes a deletion request for object 2, so we
	// skip over that one with our offset.
	query := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", 2).
		Limit(1).
		Offset(1)
	drio := pgmodels.DeletionRequestsIntellectualObjects{}
	err = query.Select(&drio)
	require.Nil(t, err)

	deletionRequest, err := pgmodels.DeletionRequestByID(drio.DeletionRequestID)
	require.Nil(t, err)
	require.NotNil(t, deletionRequest)

	// Make sure this is the request our test user just made
	require.Equal(t, testutil.Inst1Admin.ID, deletionRequest.RequestedByID)

	// There should also be an alert...
	query = pgmodels.NewQuery().Where("deletion_request_id", "=", drio.DeletionRequestID).Relations("DeletionRequest", "Users")
	alert, err := pgmodels.AlertGet(query)
	require.Nil(t, err)
	assert.NotNil(t, alert)
	assert.NotNil(t, alert.DeletionRequest)
	assert.True(t, alert.DeletionRequest.RequestedAt.After(time.Now().UTC().Add(-5*time.Second)))
	assert.True(t, len(alert.Users) > 0)

	// Alert should include link to review the deletion request.
	assert.True(t, strings.Contains(alert.Content, "token="))

	// The user should NOT be able to initiate deletion of an object
	// that belongs to another institution. In fixture data, object
	// 6 belongs to inst 2.
	testutil.Inst1UserClient.POST("/objects/init_delete/6").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusForbidden)

}

func TestObjectRequestRestore(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"This object will be restored to your institution's receiving bucket",
		"Cancel",
		"Confirm",
	}

	// Users can request deletions of their own objects
	for _, client := range testutil.AllClients {
		html := client.GET("/objects/request_restore/2").
			Expect().Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}

	// Sys Admin can request any deletion, but others cannot
	// request deletions outside their own institution.
	// Object 6 from fixtures belongs to Inst2
	testutil.SysAdminClient.GET("/objects/request_restore/6").
		Expect().Status(http.StatusOK)
	testutil.Inst1AdminClient.GET("/objects/request_restore/6").
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.GET("/objects/request_restore/6").
		Expect().Status(http.StatusForbidden)
}

func TestObjectInitRestore(t *testing.T) {
	// Force fixture reload to prevent "pending work item"
	// error when requesting restoration.
	err := db.ForceFixtureReload()
	require.Nil(t, err)
	testutil.InitHTTPTests(t)

	items := []string{
		"Object <b>institution1.edu/pdfs</b> has been queued for restoration.",
	}

	// User should see flash message saying object is queued for restoration.
	// This means the work item was created and queued.
	html := testutil.Inst1UserClient.POST("/objects/init_restore/2").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusCreated).Body().Raw()
	testutil.AssertMatchesAll(t, html, items)

	query := pgmodels.NewQuery().
		Where("action", "=", constants.ActionRestoreObject).
		Where("intellectual_object_id", "=", 2).
		Limit(1)
	workItem, err := pgmodels.WorkItemGet(query)
	require.Nil(t, err)
	require.NotNil(t, workItem)

	// Users cannot restore objects belonging to other institutions.
	testutil.Inst1AdminClient.POST("/objects/init_restore/6").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1AdminToken).
		Expect().Status(http.StatusForbidden)
	testutil.Inst1UserClient.POST("/objects/init_restore/6").
		WithHeader("Referer", testutil.BaseURL).
		WithFormField(constants.CSRFTokenName, testutil.Inst1UserToken).
		Expect().Status(http.StatusForbidden)

}

func TestIntellectualObjectEvents(t *testing.T) {
	testutil.InitHTTPTests(t)
	expect := testutil.Inst2UserClient.GET("/objects/events/1").Expect()
	html := expect.Body().Raw()
	expect.Status(http.StatusOK)

	expected := []string{
		"Object Event History",
		"ingestion",
		"Aug 26, 2016",
	}
	testutil.AssertMatchesAll(t, html, expected)
}

func TestIntellectualObjectFiles(t *testing.T) {
	testutil.InitHTTPTests(t)
	expect := testutil.Inst1UserClient.GET("/objects/files/1").Expect()
	html := expect.Body().Raw()
	expect.Status(http.StatusOK)

	expected := []string{
		"Active Files",
		"institution1.edu/photos/picture1",
		"institution1.edu/photos/picture2",
		"institution1.edu/photos/picture3",
	}
	testutil.AssertMatchesAll(t, html, expected)
}
