package web

import (
	//"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// GenericFileRequestDelete shows a confirmation message asking
// if user really wants to delete a file.
// DELETE /files/request_delete/:id
func GenericFileRequestDelete(c *gin.Context) {
	req := NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["file"] = gf
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "files/_request_delete.html", req.TemplateData)
}

// GenericFileIndex shows list of objects.
// GET /files
func GenericFileIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "files/index.html"
	var files []*pgmodels.GenericFile
	err := req.LoadResourceList(&files, "updated_at desc", forms.NewFileFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// GenericFileShow returns the object with the specified id.
// GET /files/show/:id
func GenericFileShow(c *gin.Context) {
	req := NewRequest(c)
	file, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	req.TemplateData["file"] = file
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "files/show.html", req.TemplateData)
}

// GenericFileRequestRestore shows a message asking whether
// user really wants to restore a file.
func GenericFileRequestRestore(c *gin.Context) {
	req := NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["file"] = gf
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "files/_request_restore.html", req.TemplateData)
}

// GenericFileInitRestore creates a file restoration request,
// which is really just a WorkItem that gets queued. Restoration can take
// seconds or hours, depending on where the file is stored and how big it is.
// POST /files/init_restore/:id
func GenericFileInitRestore(c *gin.Context) {
	req := NewRequest(c)
	gf, obj, _, err := genericFileInitRestore(req)
	if AbortIfError(c, err) {
		return
	}
	// Respond
	msg := fmt.Sprintf("File %s has been queued for restoration.", gf.Identifier)
	helpers.SetFlashCookie(c, msg)
	redirectUrl := fmt.Sprintf("/objects/show/%d", obj.ID)
	c.Redirect(http.StatusFound, redirectUrl)
}

// GenericFileInitDelete occurs when user clicks the button confirming
// they want to delete a file. This creates a deletion confirmation message,
// which will be emailed to institutional admins for approval.
//
// POST /files/init_delete/:id
func GenericFileInitDelete(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForFile(req.Auth.ResourceID, req.CurrentUser, req.BaseURL())
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateRequestAlert()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["fileIdentifier"] = del.DeletionRequest.FirstFile().Identifier
	c.HTML(http.StatusCreated, "files/deletion_requested.html", req.TemplateData)
}

// GenericFileReviewDelete displays a page on which an institutional
// admin can review a requested file deletion and choose whether to approve
// or cancel it.
// GET /files/review_delete/:id?token=<token>
func GenericFileReviewDelete(c *gin.Context) {
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

	template := "files/review_deletion.html"
	if del.DeletionRequest.ConfirmedByID > 0 {
		template = "files/deletion_already_approved.html"
	} else if del.DeletionRequest.CancelledByID > 0 {
		template = "files/deletion_already_cancelled.html"
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// GenericFileApproveDelete handles the case where an institutional
// admin approves a file deletion request. Note that the token comes
// in through the post form here, not through the URL.
// POST /files/approve_delete/:id
func GenericFileApproveDelete(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.PostForm("token"))
	if AbortIfError(c, err) {
		return
	}
	del.DeletionRequest.Confirm(req.CurrentUser)
	err = del.DeletionRequest.Save()
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateAndQueueWorkItem()
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateApprovalAlert()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deletionRequest"] = del.DeletionRequest
	c.HTML(http.StatusOK, "files/deletion_approved.html", req.TemplateData)
}

// GenericFileCancelDelete handles the case where an institutional
// admin cancels (rejects) a file deletion request. Token comes in
// through post form, not query string.
// POST /files/cancel_delete/:id
func GenericFileCancelDelete(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForReview(req.Auth.ResourceID, req.CurrentUser, req.BaseURL(), c.PostForm("token"))
	if AbortIfError(c, err) {
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
	c.HTML(http.StatusOK, "files/deletion_cancelled.html", req.TemplateData)
}

func genericFileInitRestore(req *Request) (*pgmodels.GenericFile, *pgmodels.IntellectualObject, *pgmodels.WorkItem, error) {
	ctx := common.Context()
	ctx.Log.Info().Msgf("[GenericFileInitRestore] Got restore request for GenericFile %d", req.Auth.ResourceID)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error finding GenericFile %d: %v", req.Auth.ResourceID, err)
		return nil, nil, nil, err
	}

	// Make sure there are no pending work items...
	pendingWorkItems, err := pgmodels.WorkItemsPendingForFile(gf.ID)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error finding pending WorkItems for GenericFile %d: %v", req.Auth.ResourceID, err)
		return gf, nil, nil, err
	}
	if len(pendingWorkItems) > 0 {
		ctx.Log.Warn().Msgf("[GenericFileInitRestore] GenericFile %d can't be restored due to pending work items (%s)", gf.ID, gf.Identifier)
		return gf, nil, nil, common.ErrPendingWorkItems
	}

	// Create the new restoration work item
	obj, err := pgmodels.IntellectualObjectByID(gf.IntellectualObjectID)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error finding parent object of GenericFile %d (IntellectualObjectID = %d): %v", req.Auth.ResourceID, gf.IntellectualObjectID, err)
		return gf, nil, nil, err
	}
	ctx.Log.Info().Msgf("[GenericFileInitRestore] Found Object %d for GenericFile %d", gf.ID, obj.ID)

	workItem, err := pgmodels.NewRestorationItem(obj, gf, req.CurrentUser)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error creating restoration WorkItem for GenericFile %d: %v", req.Auth.ResourceID, err)
		return gf, obj, nil, err
	}
	ctx.Log.Info().Msgf("[GenericFileInitRestore] Created restoration WorkItem %d for GenericFile %d", workItem.ID, gf.ID)

	// Get the name of the NSQ topic for file restorations
	topic, err := constants.TopicFor(workItem.Action, workItem.Stage)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error NSQ topic for GenericFile %d restoration (action=%s, stage=%s): %v", req.Auth.ResourceID, workItem.Action, workItem.Stage, err)
		return gf, obj, workItem, err
	}

	// Queue the new work item in NSQ
	err = ctx.NSQClient.Enqueue(topic, workItem.ID)
	if err != nil {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] NSQ returned error when queueing GenericFile %d. WorkItem %d, topic=%s: %v", req.Auth.ResourceID, workItem.ID, topic, err)
		return gf, obj, workItem, err
	}
	ctx.Log.Info().Msgf("[GenericFileInitRestore] Queued WorkItem %d in topic %s", workItem.ID, topic)

	// Mark the WorkItem as queued
	workItem.QueuedAt = time.Now().UTC()
	err = workItem.Save()

	if err == nil {
		ctx.Log.Info().Msgf("[GenericFileInitRestore] Marked WorkItem %d as queued", workItem.ID)
	} else {
		ctx.Log.Error().Msgf("[GenericFileInitRestore] Error saving WorkItem %d with QueuedAt timestamp for GenericFile %d: %v", workItem.ID, req.Auth.ResourceID, err)
	}

	return gf, obj, workItem, err
}
