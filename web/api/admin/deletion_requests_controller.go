package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// DeletionRequestShow shows the deletion request with the specified id.
//
//
// GET /admin-api/v3/deletions/show/:id
func DeletionRequestShow(c *gin.Context) {
	req := api.NewRequest(c)
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, deletionRequest.ToMin())
}
