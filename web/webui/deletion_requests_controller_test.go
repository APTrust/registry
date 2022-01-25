package webui_test

import (
	// "fmt"
	"net/http"
	// "os"
	"testing"
	"time"

	// "github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	// "github.com/APTrust/registry/db"
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
	// START HERE
	// No idea WTF is happening here, but we seem to lose the auth cookie.

	// testutil.InitHTTPTests(t)
	// testutil.ReInitAllClients(t)
	// request := makeDeletionRequest(t)

	// expect := testutil.Inst1AdminClient.GET("deletions/review/{id}", request.ID).
	// 	//WithHeader(constants.CSRFHeaderName, testutil.Inst1AdminToken).
	// 	WithQuery("token", request.ConfirmationToken).Expect()

	// //fmt.Println(expect.Cookie(common.Context().Config.Cookies.SessionCookie))

	// expect.Status(http.StatusOK)
	// html := expect.Body().Raw()

	// expected := []string{
	// 	request.FirstObject().Identifier,
	// 	request.RequestedBy.Email,
	// }
	// testutil.AssertMatchesAll(t, html, expected)
}

func TestDeletionRequestApprove(t *testing.T) {

	// Fix TestDeletionRequestReview before implementing this.

}

func TestDeletionRequestCancel(t *testing.T) {

	// Fix TestDeletionRequestReview before implementing this.

}
