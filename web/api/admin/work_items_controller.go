package admin_api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// WorkItemCreate creates a new WorkItem.
//
// POST /admin-api/v3/items/create/:institution_id
func WorkItemCreate(c *gin.Context) {
	gf, err := CreateOrUpdateItem(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gf)
}

// WorkItemUpdate updates an existing WorkItem record.
//
// PUT /admin-api/v3/items/update/:id
func WorkItemUpdate(c *gin.Context) {
	gf, err := CreateOrUpdateItem(c)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gf)
}

// WorkItemRequeue requeues a WorkItem to the specified stage.
//
// PUT /admin-api/v3/items/requeue/:id
func WorkItemRequeue(c *gin.Context) {
	stage := c.PostForm("stage")
	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if api.AbortIfError(c, err) {
		return
	}

	item, err := pgmodels.WorkItemByID(itemID)
	if api.AbortIfError(c, err) {
		return
	}

	common.Context().Log.Info().Msgf("Requeueing WorkItem %d to %s", itemID, stage)

	err = item.SetForRequeue(stage)
	if api.AbortIfError(c, err) {
		return
	}

	topic, err := constants.TopicFor(item.Action, stage)
	if api.AbortIfError(c, err) {
		return
	}

	err = common.Context().NSQClient.Enqueue(topic, itemID)
	if api.AbortIfError(c, err) {
		return
	}
	data := map[string]interface{}{
		"StatusCode": http.StatusOK,
		"Message":    fmt.Sprintf("Requeued WorkItem %d to %s", itemID, stage),
	}
	c.JSON(http.StatusOK, data)
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
