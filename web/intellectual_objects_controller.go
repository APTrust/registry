package web

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// IntellectualObjectDelete deletes a user.
// DELETE /institutions/delete/:id
func IntellectualObjectDelete(c *gin.Context) {
	// req := NewRequest(c)
	// obj, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	// if AbortIfError(c, err) {
	// 	return
	// }

	// // TODO: Create a confirmation email. Check Pharos for existing implementation.

	// err = obj.Delete()
	// if AbortIfError(c, err) {
	// 	return
	// }
	// c.Redirect(http.StatusFound, "/objects")
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

	// Select max 20 files to start. Some objects have > 100k files, and
	// we definitely don't want that many results. Let the user page through.
	fileQuery := pgmodels.NewQuery().Where("intellectual_object_id", "=", object.ID).OrderBy("identifier").Limit(20).Offset(0)
	files, err := pgmodels.GenericFileSelect(fileQuery)
	if AbortIfError(c, err) {
		return
	}

	// Get object-level events only. I.e. those that match our object ID
	// but have a null generic file id. Most object will have only a handful
	// of object-level events, though they may have thousands or hundreds of
	// thousands of file-level events. We'll get the first five, and let the
	// user page through from there.
	eventQuery := pgmodels.NewQuery().Where("intellectual_object_id", "=", object.ID).IsNull("generic_file_id").OrderBy("created_at desc").Limit(5).Offset(0)
	events, err := pgmodels.PremisEventSelect(eventQuery)
	if AbortIfError(c, err) {
		return
	}

	stats, err := pgmodels.GetObjectStats(object.ID)
	if AbortIfError(c, err) {
		return
	}

	eventCount, err := pgmodels.ObjectEventCount(object.ID)
	if AbortIfError(c, err) {
		return
	}

	req.TemplateData["object"] = object
	req.TemplateData["files"] = files

	req.TemplateData["events"] = events
	req.TemplateData["eventsOffsetStart"] = 1
	req.TemplateData["eventsOffsetEnd"] = len(events)
	req.TemplateData["eventsCount"] = eventCount
	req.TemplateData["eventsShowLeftArrow"] = true
	req.TemplateData["eventsShowRightArrow"] = true

	req.TemplateData["stats"] = stats
	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "objects/show.html", req.TemplateData)
}

func IntellectualObjectRestore(c *gin.Context) {
	// TODO: Create a restoration WorkItem.
}
