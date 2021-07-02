package web

import (
	"fmt"
	"net/http"

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
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deletionRequest"] = deletionRequest

	if deletionRequest.WorkItemID > 0 {
		req.TemplateData["workItemURL"] = fmt.Sprintf("%s/work_items/show/%d",
			req.BaseURL(),
			deletionRequest.WorkItemID)
	}

	c.HTML(http.StatusOK, "deletions/show.html", req.TemplateData)
}

// DeletionRequestIndex shows list of deletion requests.
// GET /deletions
func DeletionRequestIndex(c *gin.Context) {
	req := NewRequest(c)
	var deletions []*pgmodels.DeletionRequestView
	err := req.LoadResourceList(&deletions, "requested_at desc", forms.NewDeletionRequestFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "deletions/index.html", req.TemplateData)
}
