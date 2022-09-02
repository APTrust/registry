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
	// Only sys admin should have this permission.
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
	getRedisInfo(req, item)
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

// WorkItemRedisIndex shows a list of WorkItems that have records
// in Redis. This is an admin-only feature.
//
// GET /work_items/redis_list
func WorkItemRedisIndex(c *gin.Context) {
	aptContext := common.Context()
	req := NewRequest(c)

	// Since this is a non-standard query, we have to do most of the
	// work that Request.LoadResourceList usually handles.
	//
	// Start by getting a list of WorkItem ids from Redis.
	// The List function return a max of 500 items, for safety,
	// but in practice, we'll rarely have more than a few dozen.
	ids, err := aptContext.RedisClient.List("*")
	if AbortIfError(c, err) {
		return
	}

	// If there's nothing in Redis, we have to apply this or
	// filter collection will ignore our empty list.
	if len(ids) == 0 {
		ids = []string{"0"}
	}

	// Now get a list of WorkItemView objects matching the ids in Redis.
	filterCollection := req.GetFilterCollection()
	filterCollection.Add("id__in", ids)
	query, err := filterCollection.ToQuery()
	if AbortIfError(c, err) {
		return
	}
	query.OrderBy("date_processed", "desc")
	items, err := pgmodels.WorkItemViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["items"] = items

	// We have to set a pager to avoid a nil pointer exception.
	// We're actually going to show all items at once, so the
	// pager doesn't have much to do.
	pager, err := common.NewPager(req.GinContext, req.PathAndQuery, 500)
	if AbortIfError(c, err) {
		return
	}
	pager.SetCounts(len(ids), len(ids))
	req.TemplateData["pager"] = pager

	// Set up the filter form.
	filterForm, err := forms.NewWorkItemFilterForm(filterCollection, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}
	filterForm.GetFields()["redis_only"].Value = "true"
	req.TemplateData["filterForm"] = filterForm

	c.HTML(http.StatusOK, "work_items/index.html", req.TemplateData)
}

// WorkItemRedisDelete deletes a WorkItem's Redis record.
// This is an admin-only feature.
//
// PUT or POST /work_items/redis_delete/:id
func WorkItemRedisDelete(c *gin.Context) {
	aptContext := common.Context()
	req := NewRequest(c)
	_, err := aptContext.RedisClient.WorkItemDelete(req.Auth.ResourceID)
	if err != nil {
		aptContext.Log.Error().Msgf("Error deleting WorkItem %d from Redis: %v", req.Auth.ResourceID, err)
		AbortIfError(c, err)
		return
	} else {
		aptContext.Log.Info().Msgf("User %s deleted WorkItem %d from Redis.", req.CurrentUser.Email, req.Auth.ResourceID)
	}
	helpers.SetFlashCookie(c, "Redis data for this work item has been deleted.")
	redirectTo := fmt.Sprintf("/work_items/show/%d", req.Auth.ResourceID)
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
	req.TemplateData["item"] = workItem
	return form, req, nil
}

// Show Redis info if user has permission. This should be sys admin only.
func getRedisInfo(req *Request, item *pgmodels.WorkItemView) {
	var err error
	jsonStr := ""
	req.TemplateData["showRedisDelete"] = false
	if !req.CurrentUser.HasPermission(constants.RedisRead, item.InstitutionID) {
		return
	}
	ctx := common.Context()
	if !ctx.RedisClient.KeyExists(item.ID) {
		return
	}
	if item.Action == constants.ActionIngest {
		jsonStr, err = ctx.RedisClient.IngestObjectGet(item.ID, item.GetObjIdentifier())
		if err != nil {
			ctx.Log.Warn().Msgf("Error getting IngestObject from Redis: %v", err)
		}
	} else if item.Action == constants.ActionRestoreFile || item.Action == constants.ActionRestoreObject {
		jsonStr, err = ctx.RedisClient.RestorationObjectGet(item.ID, item.GetObjIdentifier())
		if err != nil {
			ctx.Log.Warn().Msgf("Error getting RestorationObject from Redis: %v", err)
		}
	}
	req.TemplateData["redisInfo"] = jsonStr

	if req.CurrentUser.HasPermission(constants.WorkItemRedisDelete, item.InstitutionID) && (item.HasCompleted() || item.Action == constants.ActionIngest) {
		req.TemplateData["showRedisDelete"] = true
	}
}
