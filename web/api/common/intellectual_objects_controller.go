package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// IntellectualObjectIndex shows list of objects.
//
// GET /member-api/v3/objects
// GET /admin-api/v3/objects
func IntellectualObjectIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var objs []*pgmodels.IntellectualObjectView
	pager, err := req.LoadResourceList(&objs, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(objs, pager))
}

// IntellectualObjectShow returns the object with the specified id.
//
// GET /member-api/v3/objects/show/:id
// GET /admin-api/v3/objects/show/:id
func IntellectualObjectShow(c *gin.Context) {
	req := api.NewRequest(c)
	obj, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, obj)
}
