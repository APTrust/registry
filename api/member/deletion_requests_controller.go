package memberapi

import (
	"net/http"

	"github.com/APTrust/registry/api/core"
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
// GET /member-api/v3/deletions/show/:id
func DeletionRequestShow(c *gin.Context) {
	req := core.NewRequest(c)
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if core.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, deletionRequest)
}

// DeletionRequestIndex shows list of deletion requests.
// GET /member-api/v3/deletions
func DeletionRequestIndex(c *gin.Context) {
	req := core.NewRequest(c)
	var deletions []*pgmodels.DeletionRequestView
	pager, err := req.LoadResourceList(&deletions, "requested_at desc")
	if core.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, core.NewJsonList(deletions, pager))
}
