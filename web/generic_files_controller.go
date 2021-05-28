package web

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

// GenericFileDelete deletes a user.
// DELETE /files/delete/:id
func GenericFileDelete(c *gin.Context) {

}

// GenericFileIndex shows list of objects.
// GET /files
func GenericFileIndex(c *gin.Context) {

}

// GenericFileShow returns the object with the specified id.
// GET /files/show/:id
func GenericFileShow(c *gin.Context) {
	req := NewRequest(c)
	file, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	req.TemplateData["file"] = file
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "files/show.html", req.TemplateData)
}

// GenericFileRequestRestore shows a confirmation message asking whether
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
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	// Make sure there are no pending work items...
	pendingWorkItems, err := pgmodels.WorkItemsPendingForFile(gf.ID)
	if AbortIfError(c, err) {
		return
	}
	if len(pendingWorkItems) > 0 {
		AbortIfError(c, common.ErrPendingWorkItems)
		return
	}

	// Create the new restoration work item
	obj, err := pgmodels.IntellectualObjectByID(gf.IntellectualObjectID)
	if AbortIfError(c, err) {
		return
	}
	workItem, err := pgmodels.NewRestorationItem(obj, gf, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}

	// Queue the new work item in NSQ
	topic, err := constants.TopicFor(workItem.Action, workItem.Stage)
	if AbortIfError(c, err) {
		return
	}
	ctx := common.Context()
	err = ctx.NSQClient.Enqueue(topic, workItem.ID)
	if AbortIfError(c, err) {
		return
	}

	workItem.QueuedAt = time.Now().UTC()
	err = workItem.Save()
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
func GenericFileInitDelete(c *gin.Context) {
	req := NewRequest(c)
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	// Make sure there are no pending work items...
	pendingWorkItems, err := pgmodels.WorkItemsPendingForFile(gf.ID)
	if AbortIfError(c, err) {
		return
	}
	if len(pendingWorkItems) > 0 {
		AbortIfError(c, common.ErrPendingWorkItems)
		return
	}

	// Get list of inst admins.
	adminsQuery := pgmodels.NewQuery().Where("institution_id", "=", gf.InstitutionID).Where("role", "=", constants.RoleInstAdmin)
	instAdmins, err := pgmodels.UserSelect(adminsQuery)
	if AbortIfError(c, err) {
		return
	}

	// Create the deletion request and the alert, which will become
	// the deletion confirmation email.
	deleteRequest, err := pgmodels.NewDeletionRequest()
	if AbortIfError(c, err) {
		return
	}
	deleteRequest.InstitutionID = gf.InstitutionID
	deleteRequest.RequestedByID = req.CurrentUser.ID
	deleteRequest.RequestedAt = time.Now().UTC()
	deleteRequest.AddFile(gf)
	err = deleteRequest.Save()
	if AbortIfError(c, err) {
		return
	}

	reviewURL := fmt.Sprintf("%s?token=%s",
		req.BaseURL(),
		deleteRequest.ConfirmationToken)

	alertData := map[string]string{
		"RequesterName":     req.CurrentUser.Name,
		"DeletionReviewURL": reviewURL,
	}
	tmpl := common.TextTemplates["alerts/deletion_requested.txt"]
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, alertData)
	if AbortIfError(c, err) {
		return
	}

	// Put confirmation token into URL
	confirmationAlert := &pgmodels.Alert{
		InstitutionID:     gf.InstitutionID,
		Type:              constants.AlertDeletionRequested,
		DeletionRequestID: deleteRequest.ID,
		Content:           buf.String(),
		CreatedAt:         time.Now().UTC(),
		Users:             instAdmins,
	}
	err = confirmationAlert.Save()
	if AbortIfError(c, err) {
		return
	}

	req.TemplateData["fileIdentifier"] = gf.Identifier
	c.HTML(http.StatusCreated, "files/deletion_requested.html", req.TemplateData)
}
