package web

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
	recipientID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if AbortIfError(c, err) {
		return
	}
	if !req.CurrentUser.IsAdmin() && recipientID != req.CurrentUser.ID {
		aptContext := common.Context()
		aptContext.Log.Warn().Msgf("User %d illegally tried to access alert %d belonging to user %d. Permission was denied.", req.CurrentUser.ID, req.Auth.ResourceID, recipientID)
		AbortIfError(c, common.ErrPermissionDenied)
	}
	alert, err := pgmodels.AlertViewForUser(req.Auth.ResourceID, recipientID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["alert"] = alert
	c.HTML(http.StatusOK, "alerts/show.html", req.TemplateData)
}

// AlertIndex shows list of alerts for the logged-in user, unless
// that user is sys admin, who can see all alerts.
// GET /alerts
func AlertIndex(c *gin.Context) {
	req := NewRequest(c)
	err := alIndexLoadAlerts(req)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "alerts/index.html", req.TemplateData)
}

func alIndexLoadAlerts(req *Request) error {
	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("user_id", "=", req.CurrentUser.ID)
	}
	query.OrderBy("created_at desc")
	baseURL := req.GinContext.Request.URL.Path + "?" + req.GinContext.Request.URL.RawQuery
	pager, err := NewPager(req.GinContext, baseURL, 20)
	if err != nil {
		return err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	alerts, err := pgmodels.AlertViewSelect(query)
	if err != nil {
		return err
	}

	totalRecordCount, err := query.Count(&pgmodels.AlertView{})
	if err != nil {
		return err
	}
	pager.SetCounts(totalRecordCount, len(alerts))

	form, err := forms.NewAlertFilterForm(filterCollection, req.CurrentUser)

	req.TemplateData["alerts"] = alerts
	req.TemplateData["pager"] = pager
	req.TemplateData["filterForm"] = form

	return err
}
