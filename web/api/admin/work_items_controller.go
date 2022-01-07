package admin_api

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// WorkItemCreate creates a new object record.
//
// POST /admin-api/v3/files/create/:institution_id
func WorkItemCreate(c *gin.Context) {
	gf, err := CreateOrUpdateItem(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gf)
}

// WorkItemUpdate updates an existing intellectual
// object record.
//
// PUT /admin-api/v3/files/update/:id
func WorkItemUpdate(c *gin.Context) {
	gf, err := CreateOrUpdateItem(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}

func CreateOrUpdateItem(c *gin.Context) (*pgmodels.WorkItem, error) {
	req := api.NewRequest(c)
	gf, err := WorkItemFromJson(req)
	if err != nil {
		return nil, err
	}
	err = gf.Save()
	return gf, err
}

// WorkItemFromJson returns the WorkItem from the
// JSON in the request body and the existing file record from
// the database (if there is one). It returns an error if the JSON
// can't be parsed, if the existing file can't be found, or if
// changes made to the existing object are not allowed.
func WorkItemFromJson(req *api.Request) (*pgmodels.WorkItem, error) {
	submittedItem := &pgmodels.WorkItem{}
	err := req.GinContext.BindJSON(submittedItem)
	if err != nil {
		return submittedItem, err
	}
	err = req.AssertValidIDs(submittedItem.ID, submittedItem.InstitutionID)
	if err != nil {
		return submittedItem, err
	}
	if req.Auth.ResourceID > 0 {
		existingItem, err := pgmodels.WorkItemByID(req.Auth.ResourceID)
		if err != nil {
			return submittedItem, err
		}
		err = existingItem.ValidateChanges(submittedItem)
	}
	return submittedItem, err
}
