package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkItemShow(t *testing.T) {
	initHTTPTests(t)

	items := []string{
		"1f594a4e5bb944e59c74aefe781a3726",
		"institution1.edu/photos",
		"aptrust.receiving.institution1.edu",
		"system@aptrust.org",
		"Delete",
		"Cleanup",
		"Success",
		"Item deleted successfuly",
	}

	adminActions := []string{
		"/work_items/edit/30",
		"/work_items/edit/30",
	}

	for _, client := range allClients {
		html := client.GET("/work_items/show/30").Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, items)
		if client == sysAdminClient {
			AssertMatchesAll(t, html, adminActions)
		} else {
			AssertMatchesNone(t, html, adminActions)
		}
	}
}

func TestWorkItemIndex(t *testing.T) {
	initHTTPTests(t)

	links := []string{
		"/work_items/show/5",
		"/work_items/show/6",
		"/work_items/show/7",
		"/work_items/show/8",
	}

	commonFilters := []string{
		`select name="action"`,
		`select name="stage"`,
		`select name="status"`,
		`type="text" name="name"`,
		`type="text" name="etag"`,
		`type="date" name="date_processed__gteq"`,
		`type="date" name="date_processed__lteq"`,
		`select name="needs_admin_review"`,
		`type="text" name="object_identifier"`,
		`type="text" name="generic_file_identifier"`,
		`select name="storage_option"`,
		`type="text" name="alt_identifier"`,
		`type="text" name="bag_group_identifier"`,
		`type="text" name="user"`,
		`select name="bagit_profile_identifier"`,
		`type="number" name="size__gteq"`,
		`type="number" name="size__lteq"`,
	}

	adminFilters := []string{
		`select name="institution_id"`,
	}

	for _, client := range allClients {
		html := client.GET("/work_items").Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, links)
		AssertMatchesAll(t, html, commonFilters)
		if client == sysAdminClient {
			AssertMatchesAll(t, html, adminFilters)
			AssertMatchesResultCount(t, html, 33)
		} else {
			AssertMatchesNone(t, html, adminFilters)
			AssertMatchesResultCount(t, html, 18)
		}
	}

	// Apply a filter
	objRestorationLinks := []string{
		"/work_items/show/33",
	}
	for _, client := range allClients {
		html := client.GET("/work_items").
			WithQuery("action", constants.ActionRestoreObject).
			Expect().
			Status(http.StatusOK).Body().Raw()
		AssertMatchesAll(t, html, objRestorationLinks)
	}

}

func TestWorkItemEditUpdate(t *testing.T) {
	initHTTPTests(t)

	workItem := createWorkItem(t, "unit_test_bag1.tar")

	// Sys Admin should should be able to see the edit page for this item
	sysAdminClient.GET("/work_items/edit/{id}", workItem.ID).
		Expect().Status(http.StatusOK)

	// InstAdmin and InstUser cannot edit work items
	instAdminClient.GET("/work_items/edit/{id}", workItem.ID).
		Expect().Status(http.StatusForbidden)
	instUserClient.GET("/work_items/edit/{id}", workItem.ID).
		Expect().Status(http.StatusForbidden)

	// Change some values
	workItem.Stage = constants.StageStorageValidation
	workItem.Status = constants.StatusPending
	workItem.Retry = true
	workItem.NeedsAdminReview = false
	workItem.Note = "This has been edited"
	workItem.Node = ""
	workItem.PID = 0

	// SysAdmin should be able to PUT this
	sysAdminClient.PUT("/work_items/edit/{id}", workItem.ID).
		WithForm(workItem).Expect().Status(http.StatusOK)

	// Make sure changes stuck.
	item, err := pgmodels.WorkItemByID(workItem.ID)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, workItem.Stage, item.Stage)
	assert.Equal(t, workItem.Status, item.Status)
	assert.Equal(t, workItem.Retry, item.Retry)
	assert.Equal(t, workItem.NeedsAdminReview, item.NeedsAdminReview)
	assert.Equal(t, workItem.Note, item.Note)
	assert.Equal(t, workItem.Node, item.Node)
	assert.Equal(t, workItem.PID, item.PID)

	// And make sure these roles cannot update work items
	instAdminClient.PUT("/work_items/edit/{id}", workItem.ID).
		WithForm(workItem).Expect().Status(http.StatusForbidden)
	instAdminClient.PUT("/work_items/edit/{id}", workItem.ID).
		WithForm(workItem).Expect().Status(http.StatusForbidden)
}

func TestWorkItemRequeue(t *testing.T) {
	initHTTPTests(t)
}

func createWorkItem(t *testing.T, name string) *pgmodels.WorkItem {
	now := time.Now().UTC()
	workItem := &pgmodels.WorkItem{
		Name:             name,
		ETag:             "54321543215432154321000000000000",
		InstitutionID:    inst1User.InstitutionID,
		Bucket:           "aptrust.receiving.yadda.yadda",
		User:             "system@aptrust.org",
		Note:             "Wheel her in, Homer! I'm not a picky man.",
		Action:           constants.ActionIngest,
		Stage:            constants.StageRecord,
		Status:           constants.StatusStarted,
		Outcome:          "Ourcome? I ain't done yet.",
		BagDate:          now,
		DateProcessed:    now,
		Retry:            false,
		Node:             "oh god, not Node!",
		PID:              3344,
		NeedsAdminReview: true,
		QueuedAt:         now,
		Size:             8900,
		StageStartedAt:   now,
	}
	require.Nil(t, workItem.Save())
	return workItem
}
