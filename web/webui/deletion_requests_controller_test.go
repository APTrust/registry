package webui_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletionRequestShow(t *testing.T) {
	testutil.InitHTTPTests(t)

	items := []string{
		"Deletion Request",
		"Requested By",
		"Requested At",
		"Work Item",
		"Files",
		"institution1.edu/glass/shard1",
		"institution1.edu/glass/shard2",
		"institution1.edu/glass/shard3",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/deletions/show/1").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, items)
	}

}

func TestDeletionRequestIndex(t *testing.T) {
	testutil.InitHTTPTests(t)

	// All users should see these filters on the index page.
	commonFilters := []string{
		`select name="stage"`,
		`select name="status"`,
		`name="requested_at__gteq"`,
		`name="requested_at__lteq"`,
	}

	// Only sys admin should see these filters.
	sysAdminFilters := []string{
		`select name="institution_id"`,
	}

	deletionLinks := []string{
		"/deletions/show/1",
		"/deletions/show/2",
		"/deletions/show/3",
	}

	for _, client := range testutil.AllClients {
		html := client.GET("/deletions").Expect().
			Status(http.StatusOK).Body().Raw()
		testutil.AssertMatchesAll(t, html, deletionLinks)
		testutil.AssertMatchesAll(t, html, commonFilters)
		if client == testutil.SysAdminClient {
			testutil.AssertMatchesAll(t, html, sysAdminFilters)
		} else {
			testutil.AssertMatchesNone(t, html, sysAdminFilters)
		}
	}
}

func makeDeletionRequest(t *testing.T) *pgmodels.DeletionRequest {
	request, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, request)

	query := pgmodels.NewQuery().
		Where("institution_id", "=", 2).
		Where(`"intellectual_object"."state"`, "=", constants.StateActive).
		OrderBy("id", "asc").
		Limit(1)
	obj, err := pgmodels.IntellectualObjectGet(query)
	require.Nil(t, err)
	require.NotNil(t, obj)

	request.InstitutionID = obj.InstitutionID
	request.RequestedByID = testutil.Inst1Admin.ID
	request.RequestedAt = time.Now().UTC()
	request.AddObject(obj)

	require.Nil(t, request.Save())
	return request
}

func TestDeletionRequestReview(t *testing.T) {
	testutil.InitHTTPTests(t)
	defer db.ForceFixtureReload()
	request := makeDeletionRequest(t)

	expect := testutil.Inst1AdminClient.GET("/deletions/review/{id}", request.ID).
		WithQuery("token", request.ConfirmationToken).Expect()

	expect.Status(http.StatusOK)
	html := expect.Body().Raw()

	// What should be on this page?
	// - Prompt to approve or cancel
	// - URLs for posting confirm & cancel
	// - The identifier of the object to be deleted
	// - The email address of the person who requested the deletion
	// - The confirmation token used to reach this page
	expected := []string{
		"Do you want to approve or cancel this request?",
		fmt.Sprintf(`action="/deletions/approve/%d"`, request.ID),
		fmt.Sprintf(`action="/deletions/cancel/%d"`, request.ID),
		request.FirstObject().Identifier,
		testutil.Inst1Admin.Email, // user who requested deletion
		request.ConfirmationToken,
	}
	testutil.AssertMatchesAll(t, html, expected)
}

func TestDeletionRequestApprove(t *testing.T) {

	// Fix TestDeletionRequestReview before implementing this.

}

func TestDeletionRequestCancel(t *testing.T) {

	// Fix TestDeletionRequestReview before implementing this.

}
