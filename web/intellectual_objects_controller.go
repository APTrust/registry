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
	object, err := pgmodels.IntellectualObjectByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["object"] = object
	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "objects/show.html", req.TemplateData)
}

func IntellectualObjectRestore(c *gin.Context) {
	// TODO: Create a restoration WorkItem.
}
