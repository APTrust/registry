package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// InstitutionIndex shows list of institutions.
//
// GET /admin-api/v3/institutions
func InstitutionIndex(c *gin.Context) {
	req := api.NewRequest(c)
	var institutions []*pgmodels.InstitutionView
	pager, err := req.LoadResourceList(&institutions, "name", "asc")
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, api.NewJsonList(institutions, pager))
}

// InstitutionShow returns the institution with the specified id.
//
// GET /admin-api/v3/objects/show/:id
func InstitutionShow(c *gin.Context) {
	req := api.NewRequest(c)
	inst, err := pgmodels.InstitutionViewByID(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, inst)
}
