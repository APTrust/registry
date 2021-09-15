package memberapi

import (
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// AlertShow shows the AlertView with the specified id for the specified
// user. Note that the AlertView table contains multiple entries for each
// Alert ID. But each alert goes only once to each recipient, so the
// combination of AlertID + UserID (recipient) is unique.
// GET /member-api/v3/alerts/show/:id/:user_id
func AlertShow(c *gin.Context) {
	req := api.NewRequest(c)
	alertView, err := alertLoad(req)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, alertView)
}

// AlertIndex shows list of alerts for the logged-in user, unless
// that user is sys admin, who can see all alerts.
// GET /member-api/v3/alerts
func AlertIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var alerts []*pgmodels.AlertView
	pager, err := req.LoadResourceList(&alerts, "created_at desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(alerts, pager))
}

func alertLoad(req *api.Request) (*pgmodels.AlertView, error) {
	recipientID, err := strconv.ParseInt(req.GinContext.Param("user_id"), 10, 64)
	if err != nil {
		return nil, err
	}
	if !req.CurrentUser.IsAdmin() && recipientID != req.CurrentUser.ID {
		aptContext := common.Context()
		aptContext.Log.Warn().Msgf("User %d illegally tried to access alert %d belonging to user %d. Permission was denied.", req.CurrentUser.ID, req.Auth.ResourceID, recipientID)
		return nil, common.ErrPermissionDenied
	}
	return pgmodels.AlertViewForUser(req.Auth.ResourceID, recipientID)
}
