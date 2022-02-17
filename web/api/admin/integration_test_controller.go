package admin_api

import (
	"net/http"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// PrepareFileDelete sets up preconditions for a file deletion operation.
// This runs only in the test and integration environments.
//
// POST /admin-api/v3/prepare_file_delete/:id
func PrepareFileDelete(c *gin.Context) {
	req := api.NewRequest(c)
	err := prepareFileDeletionPreconditions(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, nil)
}

// PrepareObjectDelete sets up preconditions for an object deletion operation.
// This runs only in the test and integration environments.
//
// POST /admin-api/v3/prepare_object_delete/:id
func PrepareObjectDelete(c *gin.Context) {
	req := api.NewRequest(c)
	err := prepareObjectDeletionPreconditions(req.Auth.ResourceID)
	if api.AbortIfError(c, err) {
		return
	}
	c.JSON(http.StatusOK, nil)
}

func prepareFileDeletionPreconditions(gfID int64) error {
	if !isTestEnv() {
		return common.ErrNotSupported
	}

	gf, err := pgmodels.GenericFileByID(gfID)
	if err != nil {
		return err
	}

	instAdmin, err := getInstAdmin(gf.InstitutionID)
	if err != nil {
		return err
	}

	// Deletion checks for last ingest event on this object.
	event := pgmodels.RandomPremisEvent(constants.EventIngestion)
	event.IntellectualObjectID = gf.IntellectualObjectID
	event.GenericFileID = gf.ID
	event.InstitutionID = gf.InstitutionID
	err = event.Save()
	if err != nil {
		return err
	}

	// Also requires an approved Deletion work item
	item := pgmodels.RandomWorkItem(
		gf.IntellectualObject.BagName,
		constants.ActionDelete,
		gf.IntellectualObjectID,
		gf.ID)
	item.User = instAdmin.Email
	item.InstApprover = instAdmin.Email
	item.Status = constants.StatusStarted
	err = item.Save()
	if err != nil {
		return err
	}

	// Requires approved deletion request
	now := time.Now().UTC()
	request, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return err
	}

	request.GenericFiles = append(request.GenericFiles, gf)
	request.InstitutionID = gf.InstitutionID
	request.RequestedByID = instAdmin.ID
	request.RequestedAt = now
	request.ConfirmedByID = instAdmin.ID
	request.ConfirmedAt = now
	request.WorkItemID = item.ID
	err = request.Save()
	return err
}

func prepareObjectDeletionPreconditions(objID int64) error {
	if !isTestEnv() {
		return common.ErrNotSupported
	}

	obj, err := pgmodels.IntellectualObjectByID(objID)
	if err != nil {
		return err
	}

	instAdmin, err := getInstAdmin(obj.InstitutionID)
	if err != nil {
		return err
	}

	// Deletion checks for last ingest event on this object.
	event := pgmodels.RandomPremisEvent(constants.EventIngestion)
	event.IntellectualObjectID = obj.ID
	event.InstitutionID = obj.InstitutionID
	err = event.Save()
	if err != nil {
		return err
	}

	// Also requires an approved Deletion work item
	item := pgmodels.RandomWorkItem(
		obj.BagName,
		constants.ActionDelete,
		obj.ID,
		0)
	item.User = instAdmin.Email
	item.InstApprover = instAdmin.Email
	item.Status = constants.StatusStarted
	err = item.Save()
	if err != nil {
		return err
	}

	// Requires approved deletion request
	now := time.Now().UTC()
	request, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return err
	}

	request.IntellectualObjects = append(request.IntellectualObjects, obj)
	request.InstitutionID = obj.InstitutionID
	request.RequestedByID = instAdmin.ID
	request.RequestedAt = now
	request.ConfirmedByID = instAdmin.ID
	request.ConfirmedAt = now
	request.WorkItemID = item.ID
	err = request.Save()
	return err
}

// Check some preconditions to ensure it's safe to run the methods above.
// Because those methods circumvent deletion business logic, we want to
// make sure they run only in safe environments.
func isTestEnv() bool {
	config := common.Context().Config
	isTest := config.EnvName == "test" || config.EnvName == "integration"
	isLocalDB := config.DB.Host == "localhost"
	isTestDB := strings.HasSuffix(config.DB.Name, "_test") || strings.HasSuffix(config.DB.Name, "_integration")
	return isTest && isLocalDB && isTestDB
}

func getInstAdmin(instID int64) (*pgmodels.User, error) {
	query := pgmodels.NewQuery().
		Where("institution_id", "=", instID).
		Where("role", "=", constants.RoleInstAdmin).
		IsNull(`"user"."deactivated_at"`).
		Limit(1)
	return pgmodels.UserGet(query)
}
