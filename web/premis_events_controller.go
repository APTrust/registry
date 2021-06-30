package web

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
	c.HTML(http.StatusOK, "events/_show.html", req.TemplateData)
}

// PremisEventIndex shows list of objects.
// GET /events
func PremisEventIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "events/index.html"
	err := peIndexLoadEvents(req)
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

func peIndexLoadEvents(req *Request) error {
	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("institution_id", "=", req.CurrentUser.InstitutionID)
	}
	query.OrderBy("updated_at desc")

	baseURL := req.GinContext.Request.URL.Path + "?" + req.GinContext.Request.URL.RawQuery
	pager, err := NewPager(req.GinContext, baseURL, 20)
	if err != nil {
		return err
	}

	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	events, err := pgmodels.PremisEventViewSelect(query)
	if err != nil {
		return err
	}

	count, err := query.Count(&pgmodels.PremisEvent{})
	if err != nil {
		return err
	}
	pager.SetCounts(count, len(events))

	form, err := forms.NewPremisEventFilterForm(filterCollection, req.CurrentUser)

	req.TemplateData["events"] = events
	req.TemplateData["pager"] = pager
	req.TemplateData["filterForm"] = form

	return err
}
