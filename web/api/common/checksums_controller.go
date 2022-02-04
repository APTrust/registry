package common_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// ChecksumIndex shows list of objects.
//
// GET /member-api/v3/checksums
// GET /admin-api/v3/checksums
func ChecksumIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var checksums []*pgmodels.Checksum
	pager, err := req.LoadResourceList(&checksums, "updated_at", "desc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(checksums, pager))
}

// ChecksumShow returns the object with the specified id.
//
// GET /member-api/v3/checksums/show/:id
// GET /admin-api/v3/checksums/show/:id
func ChecksumShow(c *gin.Context) {
	req := api.NewRequest(c)
	gf, err := pgmodels.ChecksumByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}
