package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// PremisEventIndex shows list of objects.
//
// GET /member-api/v3/events
// GET /admin-api/v3/events
func PremisEventIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var events []*pgmodels.PremisEventView
	pager, err := req.LoadResourceList(&events, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(events, pager))
}

// PremisEventShow returns the object with the specified id.
//
// GET /member-api/v3/events/show/:id
// GET /admin-api/v3/events/show/:id
func PremisEventShow(c *gin.Context) {
	req := api.NewRequest(c)
	var gf *pgmodels.PremisEventView
	var err error
	if req.Auth.ResourceIdentifier != "" && req.Auth.ResourceID == 0 {
		gf, err = pgmodels.PremisEventViewByIdentifier(req.Auth.ResourceIdentifier)
	} else {
		gf, err = pgmodels.PremisEventViewByID(req.Auth.ResourceID)
	}
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}
