package web

import (
	"fmt"
	"net/http"

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
	req := NewRequest(c)
	workItem, err := pgmodels.WorkItemByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	c.ShouldBind(workItem)
	form, err := NewWorkItemForm(workItem)
	if AbortIfError(c, err) {
		return
	}
	template := "work_items/form.html"
	form.Action = fmt.Sprintf("/work_items/edit/%d", req.ResourceID)
	req.TemplateData["form"] = form

	status := http.StatusOK
	err = workItem.Save()
	if err != nil {
		status = form.HandleError(err)
		if form.Error != nil {
			req.TemplateData["FormError"] = form.Error
		}
	}
	if err != nil {
		c.HTML(status, template, req.TemplateData)
		return
	}
	location := fmt.Sprintf("/work_items/show/%d?flash=WorkItem+saved", form.Model.GetID())
	c.Redirect(http.StatusSeeOther, location)
}

// WorkItemEdit shows a form to edit an exiting work item.
// GET /work_items/edit/:id
func WorkItemEdit(c *gin.Context) {
	req := NewRequest(c)
	item, err := pgmodels.WorkItemByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	form, err := NewWorkItemForm(item)
	if AbortIfError(c, err) {
		return
	}
	form.Action = fmt.Sprintf("/work_items/edit/%d", item.ID)
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, "work_items/form.html", req.TemplateData)
}

func WorkItemRequeue(c *gin.Context) {
	// TODO: Requeue logic from Pharos.
	// See preservation services code for queueing via HTTPS
}