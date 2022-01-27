package webui

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// IntellectualObjectRequestDelete shows a message asking if the user
// really wants to delete this object.
// GET /objects/request_delete/:id
func IntellectualObjectRequestDelete(c *gin.Context) {
	req := NewRequest(c)
	obj, err := pgmodels.IntellectualObjectViewByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["object"] = obj
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "objects/_request_delete.html", req.TemplateData)
}

// IntellectualObjectInitDelete creates an object deletion request. This
// request must be approved by an administrator at the depositing institution
// before the deletion will actually be queued.
//
// POST /objects/init_delete/:id
func IntellectualObjectInitDelete(c *gin.Context) {
	req := NewRequest(c)
	del, err := NewDeletionForObject(req.Auth.ResourceID, req.CurrentUser, req.BaseURL())
	if AbortIfError(c, err) {
		return
	}
	_, err = del.CreateRequestAlert()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["objIdentifier"] = del.DeletionRequest.FirstObject().Identifier
	c.HTML(http.StatusCreated, "objects/deletion_requested.html", req.TemplateData)
}

// IntellectualObjectRequestRestore shows a message asking if the user
// really wants to delete this object.
// GET /objects/request_restore/:id
func IntellectualObjectRequestRestore(c *gin.Context) {
	req := NewRequest(c)
	obj, err := pgmodels.IntellectualObjectViewByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["object"] = obj
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "objects/_request_restore.html", req.TemplateData)
}

// IntellectualObjectInitRestore creates an object restoration request,
// which is really just a WorkItem that gets queued. Restoration can take
// seconds or hours, depending on where the object is stored and how big it is.
// POST /objects/init_restore/:id
func IntellectualObjectInitRestore(c *gin.Context) {
	req := NewRequest(c)
	obj, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}

	// Make sure there are no pending work items...
	pendingWorkItems, err := pgmodels.WorkItemsPendingForObject(obj.InstitutionID, obj.BagName)
	if AbortIfError(c, err) {
		return
	}
	if len(pendingWorkItems) > 0 {
		AbortIfError(c, common.ErrPendingWorkItems)
		return
	}

	// Create the new restoration work item
	workItem, err := pgmodels.NewRestorationItem(obj, nil, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}

	// Queue the new work item in NSQ
	topic, err := constants.TopicFor(workItem.Action, workItem.Stage)
	if AbortIfError(c, err) {
		return
	}
	ctx := common.Context()
	err = ctx.NSQClient.Enqueue(topic, workItem.ID)
	if AbortIfError(c, err) {
		return
	}

	workItem.QueuedAt = time.Now().UTC()
	err = workItem.Save()
	if AbortIfError(c, err) {
		return
	}

	// Respond
	helpers.SetFlashCookie(c, "Object has been queued for restoration.")
	redirectUrl := fmt.Sprintf("/objects/show/%d", obj.ID)
	c.Redirect(http.StatusFound, redirectUrl)
}

// IntellectualObjectIndex shows list of objects.
// GET /objects
func IntellectualObjectIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "objects/index.html"
	var objects []*pgmodels.IntellectualObjectView
	err := req.LoadResourceList(&objects, "updated_at", "desc", forms.NewObjectFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// IntellectualObjectShow returns the object with the specified id.
// GET /objects/show/:id
func IntellectualObjectShow(c *gin.Context) {
	req := NewRequest(c)
	object, err := pgmodels.IntellectualObjectViewByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["object"] = object
	err = loadFiles(req, object.ID)
	if AbortIfError(c, err) {
		return
	}
	err = loadEvents(req, object.ID)
	if AbortIfError(c, err) {
		return
	}
	stats, err := pgmodels.DepositFormatStatsSelect(object.InstitutionID, object.ID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["depositFormatStats"] = stats
	c.HTML(http.StatusOK, "objects/show.html", req.TemplateData)
}

// This is called when user pages through events on the
// intellectual object detail page. This returns an HTML
// fragment, not an entire page.
//
// GET /objects/events/:id (intellectual object id)
func IntellectualObjectEvents(c *gin.Context) {
	req := NewRequest(c)
	err := loadEvents(req, req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "objects/_events.html", req.TemplateData)
}

// This is called when user pages through files on the
// intellectual object detail page. This returns an HTML
// fragment, not an entire page.
//
// GET /objects/files/:id (intellectual object id)
func IntellectualObjectFiles(c *gin.Context) {
	req := NewRequest(c)
	err := loadFiles(req, req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "objects/_file_list.html", req.TemplateData)
}

// Select max 20 files to start. Some objects have > 100k files, and
// we definitely don't want that many results. Let the user page through.
func loadFiles(req *Request, objID int64) error {
	baseURL := req.GinContext.Request.URL.Path + "?" + req.GinContext.Request.URL.RawQuery
	pager, err := common.NewPager(req.GinContext, baseURL, 20)
	if err != nil {
		return err
	}

	state := req.GinContext.DefaultQuery("state", "A")
	fileFilter := strings.TrimSpace(req.GinContext.Query("file_filter"))
	files, err := pgmodels.ObjectFiles(objID, fileFilter, state, pager.QueryOffset, pager.PerPage)
	if err != nil {
		return err
	}

	// Depending on where loadFiles is called from, we may or may
	// not have already loaded this. If not, load it now.
	if req.TemplateData["object"] == nil {
		object, err := pgmodels.IntellectualObjectViewByID(objID)
		if err != nil {
			return err
		}
		req.TemplateData["object"] = object
	}

	totalFileCount, err := pgmodels.ObjectFileCount(objID, fileFilter, state)
	pager.SetCounts(totalFileCount, len(files))
	req.TemplateData["fileFilter"] = fileFilter
	req.TemplateData["state"] = state
	req.TemplateData["files"] = files
	req.TemplateData["filePager"] = pager

	if len(files) == 0 {
		pager.ItemFirst = 0
	}

	return err
}

// Get object-level events only. I.e. those that match our object ID
// but have a null generic file id. Most object will have only a handful
// of object-level events, though they may have thousands or hundreds of
// thousands of file-level events. We'll get the first five, and let the
// user page through from there.
func loadEvents(req *Request, objID int64) error {
	baseURL := fmt.Sprintf("/objects/events/%d", objID)
	pager, err := common.NewPager(req.GinContext, baseURL, 5)
	if err != nil {
		return err
	}
	eventQuery := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", objID).
		IsNull("generic_file_id").
		OrderBy("created_at", "desc").
		Limit(pager.PerPage).
		Offset(pager.QueryOffset)
	events, err := pgmodels.PremisEventSelect(eventQuery)
	if err != nil {
		return err
	}
	totalEventCount, err := pgmodels.ObjectEventCount(objID)
	pager.SetCounts(totalEventCount, len(events))
	req.TemplateData["events"] = events
	req.TemplateData["eventPager"] = pager
	return err
}
