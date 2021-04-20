package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

	// Show requeue options to Admin, if item has not completed.
	userCanRequeue := req.CurrentUser.HasPermission(constants.WorkItemRequeue, item.InstitutionID)
	if userCanRequeue && !item.HasCompleted() {
		workItem, err := pgmodels.WorkItemByID(req.ResourceID)
		if AbortIfError(c, err) {
			return
		}
		form, err := forms.NewWorkItemRequeueForm(workItem)
		if AbortIfError(c, err) {
			return
		}
		req.TemplateData["form"] = form
	}

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
	aptContext := common.Context()
	req := NewRequest(c)
	item, err := pgmodels.WorkItemByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	stage := c.Request.PostFormValue("Stage")
	aptContext.Log.Info().Msgf("Requeueing WorkItem %d to %s", item.ID, stage)

	err = item.SetForRequeue(stage)
	if AbortIfError(c, err) {
		return
	}

	topic := constants.TopicFor(item.Action, item.Stage)
	err = aptContext.NSQClient.Enqueue(topic, item.ID)
	redirectTo := fmt.Sprintf("/work_items/show/%d?flash=Item+requeued+to+%s", item.ID, url.QueryEscape(topic))
	c.Redirect(http.StatusSeeOther, redirectTo)
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
