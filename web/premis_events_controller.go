package web

import (
	"net/http"

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
	c.HTML(http.StatusOK, "events/_show.html", req.TemplateData)
}
