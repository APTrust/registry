package admin_api

import (
	"encoding/json"
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// PremisEventCreate creates a new Premis Event. This function is
// open to sys admin only. Note that Premis Events cannot be updated
// or deleted. Also note that this expects a JSON body, not form values.
//
// POST /admin-api/v3/events/create
func PremisEventCreate(c *gin.Context) {
	jsonBytes, err := c.GetRawData()
	if api.AbortIfError(c, err) {
		return
	}
	event := &pgmodels.PremisEvent{}
	err = json.Unmarshal(jsonBytes, event)
	if api.AbortIfError(c, err) {
		return
	}
	err = event.Save()
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, event)
}
