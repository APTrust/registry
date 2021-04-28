package web

import (
	"fmt"
	"net/http"

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
// GET /objects/init_delete/:id
func IntellectualObjectInitDelete(c *gin.Context) {

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
// which is really just a WorkItem that gets queued. Restoration can seconds
// or hours, depending on where the object is stored and how big it is.
// POST /objects/init_restore/:id
func IntellectualObjectInitRestore(c *gin.Context) {

}

// IntellectualObjectIndex shows list of objects.
// GET /objects
func IntellectualObjectIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "objects/index.html"
	query, err := req.GetIndexQuery()
	if AbortIfError(c, err) {
		return
	}
	query.OrderBy("updated_at desc")
	objects, err := pgmodels.IntellectualObjectViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["objects"] = objects
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
	stats, err := pgmodels.GetObjectStats(object.ID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["stats"] = stats
	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "objects/show.html", req.TemplateData)
}

func IntellectualObjectRestore(c *gin.Context) {
	// TODO: Create a restoration WorkItem.
}

// This is called when user pages through events on the
// intellectual object detail page. This returns an HTML
// fragment, not an entire page.
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
	baseURL := fmt.Sprintf("/objects/files/%d", objID)
	pager, err := NewPager(req.GinContext, baseURL, 10)
	if err != nil {
		return err
	}
	fileQuery := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", objID).
		Where("state", "=", "A").
		Relations("StorageRecords", "Checksums").
		OrderBy("identifier").
		Limit(pager.PerPage).
		Offset(pager.QueryOffset)

	//
	// This query won't work as is.
	// See https://pg.uptrace.dev/orm/has-many-relation/
	//

	// fileFilter := req.GinContext.Query("file_filter")
	// if fileFilter != "" {
	// 	cols := []string{"identifier", "checksum.digest"}
	// 	ops := []string{"ilike", "="}
	// 	vals := []interface{}{fmt.Sprintf("%%%s%%", fileFilter), fileFilter}
	// 	fileQuery.Or(cols, ops, vals)
	// }

	files, err := pgmodels.GenericFileSelect(fileQuery)
	if err != nil {
		return err
	}
	totalFileCount, err := pgmodels.ObjectFileCount(objID)
	pager.SetCounts(totalFileCount, len(files))
	req.TemplateData["files"] = files
	req.TemplateData["filePager"] = pager
	return err
}

// Get object-level events only. I.e. those that match our object ID
// but have a null generic file id. Most object will have only a handful
// of object-level events, though they may have thousands or hundreds of
// thousands of file-level events. We'll get the first five, and let the
// user page through from there.
func loadEvents(req *Request, objID int64) error {
	baseURL := fmt.Sprintf("/objects/events/%d", objID)
	pager, err := NewPager(req.GinContext, baseURL, 5)
	if err != nil {
		return err
	}
	eventQuery := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", objID).
		IsNull("generic_file_id").
		OrderBy("created_at desc").
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
