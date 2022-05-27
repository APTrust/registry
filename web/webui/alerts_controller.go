package webui

import (
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type AlertReadResult struct {
	Succeeded []int64 `json:"succeeded"`
	Failed    []int64 `json:"failed"`
	Error     string  `json:"error"`
}

func NewAlertReadResult() *AlertReadResult {
	return &AlertReadResult{
		Succeeded: make([]int64, 0),
		Failed:    make([]int64, 0),
	}
}

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

// AlertsMarkAsReadXHR handles an AJAX request to mark alerts as read.
// The PUT body should contain a csrf token and param called "id__in"
// whose value is a list of alert ids. This will mark each of those
// alerts as read by the current user.
//
// Because this method acts only on alerts belonging to the current
// user, it's impossible for one user to mark another user's alerts
// as read. (Unless the user logs in under someone else's account.)
//
// PUT /alerts/mark_as_read
func AlertMarkAsReadXHR(c *gin.Context) {
	markAlerts(c, "read")
}

// AlertsMarkAsUnreadXHR handles an AJAX request to mark alerts as unread.
// The PUT body should contain a csrf token and param called "id__in"
// whose value is a list of alert ids. This will mark each of those
// alerts as unread by the current user.
//
// Because this method acts only on alerts belonging to the current
// user, it's impossible for one user to mark another user's alerts
// as unread. (Unless the user logs in under someone else's account.)
//
// PUT /alerts/mark_as_unread
func AlertMarkAsUnreadXHR(c *gin.Context) {
	markAlerts(c, "unread")
}

// AlertsMarkAllAsRead marks all of a user's unread alerts as read.
//
//
// Because this method acts only on alerts belonging to the current
// user, it's impossible for one user to mark another user's alerts
// as read. (Unless the user logs in under someone else's account.)
//
// PUT /alerts/mark_all_as_read
func AlertMarkAllAsReadXHR(c *gin.Context) {
	req := NewRequest(c)
	status := http.StatusOK
	result := NewAlertReadResult()
	query := pgmodels.NewQuery().
		Columns("id").
		Where("user_id", "=", req.CurrentUser.ID).
		IsNull("read_at")
	alertViews, err := pgmodels.AlertViewSelect(query)
	if err != nil {
		status = StatusCodeForError(err)
		result.Error = err.Error()
		c.JSON(status, result)
		return
	}
	for _, alertView := range alertViews {
		err = alertMarkAsRead(alertView.ID, req.CurrentUser.ID)
		if err != nil {
			result.Failed = append(result.Failed, alertView.ID)
			status = StatusCodeForError(err)
			result.Error += err.Error() + ", "
		} else {
			result.Succeeded = append(result.Succeeded, alertView.ID)
		}
	}
	c.JSON(status, result)
}

func markAlerts(c *gin.Context, markThemHow string) {
	req := NewRequest(c)
	status := http.StatusOK
	result := NewAlertReadResult()
	ids := getAlertIDsFromForm(c)
	if len(ids) == 0 {
		status = http.StatusBadRequest
		result.Error = "You must specify one or more alerts."
		c.JSON(status, result)
		return
	}
	for _, id := range ids {
		var err error
		if markThemHow == "read" {
			err = alertMarkAsRead(id, req.CurrentUser.ID)
		} else {
			err = alertMarkAsUnread(id, req.CurrentUser.ID)
		}
		if err != nil {
			status = StatusCodeForError(err)
			result.Failed = append(result.Failed, id)
			result.Error += err.Error() + ", "
		} else {
			result.Succeeded = append(result.Succeeded, id)
		}
	}
	c.JSON(status, result)
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

func alertMarkAsUnread(alertID, userID int64) error {
	alert, err := pgmodels.AlertByID(alertID)
	if err != nil {
		common.Context().Log.Error().Msgf("Error marking alert %d as unread for user %d: %v", alertID, userID, err)
		return err
	}
	return alert.MarkAsUnread(userID)
}

func getAlertIDsFromForm(c *gin.Context) []int64 {
	idParams := c.PostFormArray("id__in")
	// Use zero instead of len(idParams) on the longshot
	// chance that some ids fail to parse.
	ids := make([]int64, 0)
	for _, idParam := range idParams {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}
