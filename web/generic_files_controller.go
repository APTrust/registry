package web

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// GenericFileDelete deletes a user.
// DELETE /files/delete/:id
func GenericFileDelete(c *gin.Context) {

}

// GenericFileIndex shows list of objects.
// GET /files
func GenericFileIndex(c *gin.Context) {

}

// GenericFileShow returns the object with the specified id.
// GET /files/show/:id
func GenericFileShow(c *gin.Context) {
	req := NewRequest(c)
	file, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	req.TemplateData["file"] = file
	req.TemplateData["error"] = err
	c.HTML(http.StatusOK, "files/show.html", req.TemplateData)
}

func GenericFileRestore(c *gin.Context) {
	// TODO: Create a restoration WorkItem.
}
