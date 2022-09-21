package webui

import (
	"net/http"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// PremisEventShowXHR returns the PREMIS event with the specified ID
// as an HTML fragment suitable for loading into a modal dialog or
// other existing HTML element. This does not return a full page with
// header and footer.
//
// GET /events/show_xhr/:id
func PremisEventShowXHR(c *gin.Context) {
	req := NewRequest(c)
	event, err := pgmodels.PremisEventViewByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["event"] = event
	req.TemplateData["showAsModal"] = true
	c.HTML(http.StatusOK, "events/show.html", req.TemplateData)
}

// PremisEventIndex shows list of objects.
// GET /events
func PremisEventIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "events/index.html"
	var events []*pgmodels.PremisEventView
	err := req.LoadResourceList(&events, "date_time", "desc", forms.NewPremisEventFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// PremisEventShow returns the object with the specified id.
// GET /events/show/:id
func PremisEventShow(c *gin.Context) {
	req := NewRequest(c)
	event, err := pgmodels.PremisEventViewByID(req.Auth.ResourceID)
	req.TemplateData["event"] = event
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "events/show.html", req.TemplateData)
}
