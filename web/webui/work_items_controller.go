package webui

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// WorkItemIndex shows list of work items.
// GET /work_items
func WorkItemIndex(c *gin.Context) {
	req := NewRequest(c)
	var items []*pgmodels.WorkItemView
	err := req.LoadResourceList(&items, "updated_at", "desc", forms.NewWorkItemFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "work_items/index.html", req.TemplateData)
}

// WorkItemShow returns the work item with the specified id.
// GET /work_items/show/:id
func WorkItemShow(c *gin.Context) {
	req := NewRequest(c)
	item, err := pgmodels.WorkItemViewByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["item"] = item

	// Show requeue options to Admin, if item has not completed.
	userCanRequeue := req.CurrentUser.HasPermission(constants.WorkItemRequeue, item.InstitutionID)
	if userCanRequeue && !item.HasCompleted() {
		workItem, err := pgmodels.WorkItemByID(req.Auth.ResourceID)
		if AbortIfError(c, err) {
			return
		}
		form, err := forms.NewWorkItemRequeueForm(workItem)
		if AbortIfError(c, err) {
			return
		}
		req.TemplateData["form"] = form
	}
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

// WorkItemRequeue requeues a WorkItem. This is an admin-only feature
// typically used to recover from system errors.
//
// PUT or POST /work_items/requeue/:id
func WorkItemRequeue(c *gin.Context) {
	aptContext := common.Context()
	req := NewRequest(c)
	item, err := pgmodels.WorkItemByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	stage := c.Request.PostFormValue("Stage")
	aptContext.Log.Info().Msgf("Requeueing WorkItem %d to %s", item.ID, stage)

	err = item.SetForRequeue(stage)
	if AbortIfError(c, err) {
		return
	}

	topic, err := constants.TopicFor(item.Action, item.Stage)
	if AbortIfError(c, err) {
		return
	}

	err = aptContext.NSQClient.Enqueue(topic, item.ID)
	if AbortIfError(c, err) {
		return
	}
	helpers.SetFlashCookie(c, fmt.Sprintf("Item has been requeued to %s", topic))
	redirectTo := fmt.Sprintf("/work_items/show/%d", item.ID)
	c.Redirect(http.StatusSeeOther, redirectTo)
}

func getFormAndRequest(c *gin.Context) (*forms.WorkItemForm, *Request, error) {
	req := NewRequest(c)
	workItem, err := pgmodels.WorkItemByID(req.Auth.ResourceID)
	if err != nil {
		return nil, nil, err
	}
	c.ShouldBind(workItem)
	form := forms.NewWorkItemForm(workItem)
	req.TemplateData["form"] = form
	return form, req, nil
}