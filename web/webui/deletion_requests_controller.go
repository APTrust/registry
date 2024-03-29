package webui

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// DeletionRequestShow shows the deletion request with the specified id.
//
// Note that this shows a read-only view of the request. It does not include
// the Approve/Cancel buttons. This read-only view may be available to users
// who do not have permission to initiate, approve, or cancel deletion requests
// but who still need a read-only view of the requests that have been submitted.
//
// Deletions apply to files and/or intellectual objects. The methods for
// initiating, approving and rejecting deletion requests are in the
// Generic Files Controller (for files) and the Intellectual Objects Controller
// (for objects).
//
// GET /deletions/show/:id
func DeletionRequestShow(c *gin.Context) {
	req := NewRequest(c)
	err := deletionRequestLoad(req)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "deletions/show.html", req.TemplateData)
}

// DeletionRequestIndex shows list of deletion requests.
// GET /deletions
func DeletionRequestIndex(c *gin.Context) {
	req := NewRequest(c)
	var deletions []*pgmodels.DeletionRequestView
	err := req.LoadResourceList(&deletions, "requested_at", "desc", forms.NewDeletionRequestFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "deletions/index.html", req.TemplateData)
}

func deletionRequestLoad(req *Request) error {
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if err != nil {
		return err
	}
	req.TemplateData["deletionRequest"] = deletionRequest

	if len(deletionRequest.WorkItems) > 0 {
		urls := make([]string, 0)
		for _, item := range deletionRequest.WorkItems {
			urls = append(urls, fmt.Sprintf("%s/work_items/show/%d", req.BaseURL(), item.ID))
		}
		req.TemplateData["workItemURLs"] = urls
	}
	return nil
}

// DeletionRequestReview displays a page on which an institutional
// admin can review a requested file deletion and choose whether to approve
// or cancel it.
//
// GET /deletions/review/:id?token=<token>
func DeletionRequestReview(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.Query("token"))
	if AbortIfError(c, err) {
		return
	}

	// Present the page describing the request, and if it hasn't
	// already been cancelled or approved, give the user the option
	// to cancel or approve.
	req.TemplateData["deletionRequest"] = del.DeletionRequest
	req.TemplateData["token"] = c.Query("token")

	if len(del.DeletionRequest.IntellectualObjects) == 1 {
		req.TemplateData["itemType"] = "single object"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("object %s", del.DeletionRequest.IntellectualObjects[0].Identifier)
		req.TemplateData["object"] = del.DeletionRequest.IntellectualObjects[0]
	} else if len(del.DeletionRequest.IntellectualObjects) > 1 {
		// Bulk object deletion
		req.TemplateData["itemType"] = "object list"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("%d objects", len(del.DeletionRequest.IntellectualObjects))
		req.TemplateData["objectList"] = del.DeletionRequest.IntellectualObjects
	} else if len(del.DeletionRequest.GenericFiles) > 0 {
		req.TemplateData["itemType"] = "file"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("file %s", del.DeletionRequest.GenericFiles[0].Identifier)
		req.TemplateData["file"] = del.DeletionRequest.GenericFiles[0]
	} else {
		common.Context().Log.Info().Msgf("DeletionRequest with ID %d has no associated files or objects.", req.Auth.ResourceID)
		AbortIfError(c, common.ErrInternal)
		return
	}

	template := "deletions/review.html"
	if del.DeletionRequest.ConfirmedByID > 0 {
		template = "deletions/already_approved.html"
	} else if del.DeletionRequest.CancelledByID > 0 {
		template = "deletions/already_cancelled.html"
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// DeletionRequestApprove handles the case where an institutional
// admin approves a deletion request. Note that the token comes
// in through the post form here, not through the URL.
//
// POST /deletions/approve/:id
func DeletionRequestApprove(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.PostForm("token"))
	if AbortIfError(c, err) {
		return
	}

	if !del.DeletionRequest.CancelledAt.IsZero() {
		err = common.ErrRequestAlreadyCancelled
	} else if !del.DeletionRequest.ConfirmedAt.IsZero() {
		err = common.ErrRequestAlreadyApproved
	}

	if AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("Cannot approve deletion request %d: %v", del.DeletionRequest.ID, err)
		return
	}

	del.DeletionRequest.Confirm(req.CurrentUser)
	err = del.DeletionRequest.Save()
	if AbortIfError(c, err) {
		return
	}
	err = del.CreateAndQueueWorkItems()
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateApprovalAlert()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deletionRequest"] = del.DeletionRequest
	template := "deletions/approved_file.html"
	if len(del.DeletionRequest.IntellectualObjects) > 0 {
		template = "deletions/approved_object.html"
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// DeletionRequestCancel handles the case where an institutional
// admin cancels (rejects) a file deletion request. Token comes in
// through post form, not query string.
//
// POST /deletions/cancel/:id
func DeletionRequestCancel(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.PostForm("token"))
	if AbortIfError(c, err) {
		return
	}

	if !del.DeletionRequest.CancelledAt.IsZero() {
		err = common.ErrRequestAlreadyCancelled
	} else if !del.DeletionRequest.ConfirmedAt.IsZero() {
		err = common.ErrRequestAlreadyApproved
	}

	if AbortIfError(c, err) {
		common.Context().Log.Error().Msgf("Cannot cancel deletion request %d: %v", del.DeletionRequest.ID, err)
		return
	}

	del.DeletionRequest.Cancel(req.CurrentUser)
	err = del.DeletionRequest.Save()
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateCancellationAlert()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deletionRequest"] = del.DeletionRequest
	template := "deletions/cancelled_file.html"
	if len(del.DeletionRequest.IntellectualObjects) > 0 {
		template = "deletions/cancelled_object.html"
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}
