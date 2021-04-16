package web

import (
	"net/http"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// WorkItemIndex shows list of work items.
// GET /work_items
func WorkItemIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "work_items/index.html"
	query := pgmodels.NewQuery().OrderBy("name")
	items, err := pgmodels.WorkItemViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["items"] = items
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// WorkItemShow returns the work item with the specified id.
// GET /work_items/show/:id
func WorkItemShow(c *gin.Context) {
	req := NewRequest(c)
	item, err := pgmodels.WorkItemViewByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["item"] = item

	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "work_items/show.html", req.TemplateData)
}

// WorkItemUpdate saves changes to an exiting work item.
// This is an admin-only feature.
// PUT /work_items/edit/:id
func WorkItemUpdate(c *gin.Context) {
	form, req, err := getFormAndRequest(c)
	if AbortIfError(c, err) {
		return
	}
	if form.Save() {
		c.Redirect(form.Status, form.PostSaveURL())
	} else {
		req.TemplateData["FormError"] = form.Error
		c.HTML(form.Status, form.Template, req.TemplateData)
	}
}

// WorkItemEdit shows a form to edit an exiting work item.
// GET /work_items/edit/:id
func WorkItemEdit(c *gin.Context) {
	form, req, err := getFormAndRequest(c)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

func WorkItemRequeue(c *gin.Context) {
	// TODO: Requeue logic from Pharos.
	// See preservation services code for queueing via HTTPS
}

func getFormAndRequest(c *gin.Context) (*forms.WorkItemForm, *Request, error) {
	req := NewRequest(c)
	workItem, err := pgmodels.WorkItemByID(req.ResourceID)
	if err != nil {
		return nil, nil, err
	}
	c.ShouldBind(workItem)
	form := forms.NewWorkItemForm(workItem)
	req.TemplateData["form"] = form
	return form, req, nil
}
