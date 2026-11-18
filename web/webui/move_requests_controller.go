package webui

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// MoveRequestShow shows the move request with the specified id.
//
// Note that this shows a read-only view of the request. It does not include
// the Approve/Cancel buttons. This read-only view may be available to users
// who do not have permission to initiate, approve, or cancel move requests
// but who still need a read-only view of the requests that have been submitted.
//
// Moves apply to intellectual objects only. The methods for
// initiating, approving and rejecting move requests are in the
// the Intellectual Objects Controller.
//
// GET /moves/show/:id
func MoveRequestShow(c *gin.Context) {
	req := NewRequest(c)
	err := moveRequestLoad(req)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "moves/show.html", req.TemplateData)
}

// MoveRequestIndex shows list of move requests.
// GET /moves
func MoveRequestIndex(c *gin.Context) {
	req := NewRequest(c)
	var moves []*pgmodels.DeletionRequestView
	err := req.LoadResourceList(&moves, "requested_at", "desc", forms.NewDeletionRequestFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "moves/index.html", req.TemplateData)
}

func moveRequestLoad(req *Request) error {
	moveRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if err != nil {
		return err
	}
	req.TemplateData["moveRequest"] = moveRequest

	if len(moveRequest.WorkItems) > 0 {
		urls := make([]string, 0)
		for _, item := range moveRequest.WorkItems {
			urls = append(urls, fmt.Sprintf("%s/work_items/show/%d", req.BaseURL(), item.ID))
		}
		req.TemplateData["workItemURLs"] = urls
	}
	return nil
}

// MoveRequestReview displays a page on which an institutional
// admin can review a requested move and choose whether to approve
// or cancel it.
//
// GET /moves/review/:id?token=<token>
func MoveRequestReview(c *gin.Context) {
	req := NewRequest(c)
	move, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.Query("token"))
	if AbortIfError(c, err) {
		return
	}

	// Present the page describing the request, and if it hasn't
	// already been cancelled or approved, give the user the option
	// to cancel or approve.
	req.TemplateData["deletionRequest"] = move.DeletionRequest
	req.TemplateData["token"] = c.Query("token")

	if len(move.DeletionRequest.IntellectualObjects) == 1 {
		req.TemplateData["itemType"] = "single object"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("object %s", move.DeletionRequest.IntellectualObjects[0].Identifier)
		req.TemplateData["object"] = move.DeletionRequest.IntellectualObjects[0]
	} else if len(move.DeletionRequest.IntellectualObjects) > 1 {
		// Bulk object deletion
		req.TemplateData["itemType"] = "object list"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("%d objects", len(move.DeletionRequest.IntellectualObjects))
		req.TemplateData["objectList"] = move.DeletionRequest.IntellectualObjects
	} else if len(move.DeletionRequest.GenericFiles) > 0 {
		req.TemplateData["itemType"] = "file"
		req.TemplateData["itemIdentifier"] = fmt.Sprintf("file %s", move.DeletionRequest.GenericFiles[0].Identifier)
		req.TemplateData["file"] = move.DeletionRequest.GenericFiles[0]
	} else {
		common.Context().Log.Info().Msgf("DeletionRequest with ID %d has no associated files or objects.", req.Auth.ResourceID)
		AbortIfError(c, common.ErrInternal)
		return
	}

	template := "deletions/review.html"
	if move.DeletionRequest.ConfirmedByID > 0 {
		template = "deletions/already_approved.html"
	} else if move.DeletionRequest.CancelledByID > 0 {
		template = "deletions/already_cancelled.html"
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// DeletionRequestApprove handles the case where an institutional
// admin approves a deletion request. Note that the token comes
// in through the post form here, not through the URL.
//
// POST /deletions/approve/:id
func MoveRequestApprove(c *gin.Context) {
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
func MoveRequestCancel(c *gin.Context) {
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
