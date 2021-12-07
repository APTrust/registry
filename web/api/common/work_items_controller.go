package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// WorkItemIndex shows list of objects.
//
// GET /member-api/v3/items
// GET /admin-api/v3/items
func WorkItemIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var items []*pgmodels.WorkItemView
	pager, err := req.LoadResourceList(&items, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(items, pager))
}

// WorkItemShow returns the object with the specified id.
//
// GET /member-api/v3/items/show/:id
// GET /admin-api/v3/items/show/:id
func WorkItemShow(c *gin.Context) {
	req := api.NewRequest(c)
	item, err := pgmodels.WorkItemViewByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, item)
}
