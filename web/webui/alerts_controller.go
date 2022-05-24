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
	alertView, err := alertLoad(req)
	if AbortIfError(c, err) {
		return
	}
	err = alertMarkAsRead(alertView.ID, alertView.UserID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["alert"] = alertView
	c.HTML(http.StatusOK, "alerts/show.html", req.TemplateData)
}

// AlertIndex shows list of alerts for the logged-in user, unless
// that user is sys admin, who can see all alerts.
// GET /alerts
func AlertIndex(c *gin.Context) {
	req := NewRequest(c)
	var alerts []*pgmodels.AlertView
	err := req.LoadResourceList(&alerts, "created_at", "desc", forms.NewAlertFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "alerts/index.html", req.TemplateData)
}

func alertLoad(req *Request) (*pgmodels.AlertView, error) {
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

func alertMarkAsRead(alertID, userID int64) error {
	alert, err := pgmodels.AlertByID(alertID)
	if err != nil {
		common.Context().Log.Error().Msgf("Error marking alert %d as read for user %d: %v", alertID, userID, err)
		return err
	}
	return alert.MarkAsRead(userID)
}
