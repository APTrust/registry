package webui

import (
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// AlertShow shows the AlertView with the specified id for the specified
// user. Note that the AlertView table contains multiple entries for each
// Alert ID. But each alert goes only once to each recipient, so the
// combination of AlertID + UserID (recipient) is unique.
// GET /alerts/show/:id
func AlertShow(c *gin.Context) {
	req := NewRequest(c)
	err := alertLoad(req)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "alerts/show.html", req.TemplateData)
}

// AlertIndex shows list of alerts for the logged-in user, unless
// that user is sys admin, who can see all alerts.
// GET /alerts
func AlertIndex(c *gin.Context) {
	req := NewRequest(c)
	var alerts []*pgmodels.AlertView
	err := req.LoadResourceList(&alerts, "created_at desc", forms.NewAlertFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "alerts/index.html", req.TemplateData)
}

func alertLoad(req *Request) error {
	recipientID, err := strconv.ParseInt(req.GinContext.Param("user_id"), 10, 64)
	if err != nil {
		return err
	}
	if !req.CurrentUser.IsAdmin() && recipientID != req.CurrentUser.ID {
		aptContext := common.Context()
		aptContext.Log.Warn().Msgf("User %d illegally tried to access alert %d belonging to user %d. Permission was denied.", req.CurrentUser.ID, req.Auth.ResourceID, recipientID)
		return common.ErrPermissionDenied
	}
	alert, err := pgmodels.AlertViewForUser(req.Auth.ResourceID, recipientID)
	req.TemplateData["alert"] = alert
	return err
}
