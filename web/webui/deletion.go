package webui

import (
	"errors"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// Deletion is a helper object for web requests involving the review,
// approval, and cancellation of deletion requests. This object simply
// encapsulates a lot of the grunt work required by the
// GenericFilesController and the IntellectualObjectsController.
type Deletion struct {
	// DeletionRequest is the DeletionRequest.
	DeletionRequest *pgmodels.DeletionRequest

	// InstAdmins is the list of institutional admins to be alerted
	// about this DeletionRequest. One of these admins has to approve
	// or cancel the request before we move forward.
	InstAdmins []*pgmodels.User

	baseURL     string
	currentUser *pgmodels.User
}

// NewDeletionForFile creates a new DeletionRequest for a GenericFile
// and returns the Deletion object. This constructor is only for initializing
// new DeletionRequests, not for reviewing, approving or cancelling
// existing requests.
func NewDeletionForFile(genericFileID int64, currentUser *pgmodels.User, baseURL string) (*Deletion, error) {
	// Make sure there are no pending work items for this
	// generic file or its parent object.
	pendingWorkItems, err := pgmodels.WorkItemsPendingForFile(genericFileID)
	if err != nil {
		return nil, err
	}
	if len(pendingWorkItems) > 0 {
		return nil, common.ErrPendingWorkItems
	}

	del := &Deletion{
		baseURL:     baseURL,
		currentUser: currentUser,
	}
	err = del.initFileDeletionRequest(genericFileID)
	if err != nil {
		return nil, err
	}
	err = del.loadInstAdmins()
	return del, err
}

// NewDeletionForObject creates a new DeletionRequest for an IntellectualObject
// and returns the Deletion object. This constructor is only for initializing
// new DeletionRequests, not for reviewing, approving or cancelling
// existing requests.
func NewDeletionForObject(objID int64, currentUser *pgmodels.User, baseURL string) (*Deletion, error) {
	obj, err := pgmodels.IntellectualObjectByID(objID)
	if err != nil {
		return nil, err
	}

	// Make sure there are no pending work items for this object.
	pendingWorkItems, err := pgmodels.WorkItemsPendingForObject(obj.ID, obj.BagName)
	if err != nil {
		return nil, err
	}
	if len(pendingWorkItems) > 0 {
		return nil, common.ErrPendingWorkItems
	}

	del := &Deletion{
		baseURL:     baseURL,
		currentUser: currentUser,
	}
	err = del.initObjectDeletionRequest(obj.InstitutionID, []int64{objID})
	if err != nil {
		return nil, err
	}
	err = del.loadInstAdmins()
	return del, err
}

// NewDeletionForObjectBatch creates a new DeletionRequest for a batch of
// IntellectualObjects and returns the Deletion object. This constructor
// is only for initializing new DeletionRequests, not for reviewing, approving
// or cancelling existing requests.
func NewDeletionForObjectBatch(requestorID, institutionID int64, objIDs []int64, baseURL string) (*Deletion, error) {

	requestingUser, err := pgmodels.UserByID(requestorID)
	if err != nil {
		return nil, err
	}
	if requestingUser.InstitutionID != institutionID || requestingUser.Role != constants.RoleInstAdmin {
		common.Context().Log.Error().Msgf("Requesting user %d is not admin at institution %d. Rejecting bulk deletion request.", requestorID, institutionID)
		return nil, common.ErrInvalidRequestorID
	}

	// Make sure that all objects belong to the specified institution.
	validObjectCount, err := pgmodels.CountObjectsThatCanBeDeleted(institutionID, objIDs)
	if err != nil {
		return nil, err
	}
	if validObjectCount != len(objIDs) {
		common.Context().Log.Error().Msgf("Batch deletion requested for %d objects, of which only %d are valid. InstitutionID = %d. Current user = %s. IDs: %v",
			len(objIDs), validObjectCount, institutionID, requestingUser.Email, objIDs)
		return nil, common.ErrInvalidObjectID
	}

	// Make sure there are no pending work items for these objects.
	pendingWorkItems, err := pgmodels.WorkItemsPendingForObjectBatch(objIDs)
	if err != nil {
		return nil, err
	}
	if pendingWorkItems > 0 {
		common.Context().Log.Warn().Msgf("Some objects in batch deletion request have pending work items. Object IDs: %v", objIDs)
		return nil, common.ErrPendingWorkItems
	}

	del := &Deletion{
		baseURL:     baseURL,
		currentUser: requestingUser,
	}
	err = del.initObjectDeletionRequest(institutionID, objIDs)
	if err != nil {
		return nil, err
	}
	err = del.loadInstAdmins()
	return del, err
}

// NewDeletionForReview pulls up information about an existing deletion
// request that an institutional admin will review before deciding whether
// to approve or cancel the request.
func NewDeletionForReview(deletionRequestID int64, currentUser *pgmodels.User, baseURL, token string) (*Deletion, error) {
	del := &Deletion{
		baseURL:     baseURL,
		currentUser: currentUser,
	}
	err := del.loadDeletionRequest(deletionRequestID)
	if err != nil {
		return nil, err
	}

	if !common.ComparePasswords(del.DeletionRequest.EncryptedConfirmationToken, token) {
		return nil, common.ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}

	err = del.loadInstAdmins()
	return del, err
}

// loadDeletionRequest loads an existing request so an admin can
// review it for approval or cancellation.
func (del *Deletion) loadDeletionRequest(deletionRequestID int64) error {
	deletionRequest, err := pgmodels.DeletionRequestByID(deletionRequestID)
	if err != nil {
		return err
	}
	del.DeletionRequest = deletionRequest
	return nil
}

// loadInstAdmins loads the list of institutional admins who should
// receive an alert about this deletion request. The inst admins choose
// whether to approve or deny the request.
func (del *Deletion) loadInstAdmins() error {
	adminsQuery := pgmodels.NewQuery().
		Where("institution_id", "=", del.DeletionRequest.InstitutionID).
		Where("role", "=", constants.RoleInstAdmin).
		IsNull("deactivated_at")
	instAdmins, err := pgmodels.UserSelect(adminsQuery)
	if err != nil {
		return err
	}
	del.InstAdmins = instAdmins
	return nil
}

// initFileDeletionRequest creates a new file DeletionRequest. When this
// request is created, it includes a plaintext token that we add to the
// confirmation URL below. We do not save the plaintext version of the token,
// only the encrypted version. When this new DeletionRequest goes out of
// scope, there's no further access to the token, so get it while you can.
func (del *Deletion) initFileDeletionRequest(genericFileID int64) error {
	gf, err := pgmodels.GenericFileByID(genericFileID)
	if err != nil {
		return err
	}

	deletionRequest, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return err
	}
	deletionRequest.InstitutionID = gf.InstitutionID
	deletionRequest.RequestedByID = del.currentUser.ID
	deletionRequest.RequestedAt = time.Now().UTC()
	deletionRequest.AddFile(gf)
	err = deletionRequest.Save()
	if err != nil {
		return err
	}
	del.DeletionRequest = deletionRequest
	return nil
}

// initObjectDeletionRequest creates a new DeletionRequest for an
// IntellectualObject. When this request is created, it includes a
// plaintext token that we add to the confirmation URL below.
// We do not save the plaintext version of the token,
// only the encrypted version. When this new DeletionRequest goes out of
// scope, there's no further access to the token, so get it while you can.
func (del *Deletion) initObjectDeletionRequest(institutionID int64, objIDs []int64) error {
	deletionRequest, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return err
	}
	deletionRequest.InstitutionID = institutionID
	deletionRequest.RequestedByID = del.currentUser.ID
	deletionRequest.RequestedAt = time.Now().UTC()

	for _, objID := range objIDs {
		obj, err := pgmodels.IntellectualObjectByID(objID)
		if err != nil {
			return err
		}
		deletionRequest.AddObject(obj)
	}
	err = deletionRequest.Save()
	if err != nil {
		return err
	}
	del.DeletionRequest = deletionRequest
	return nil
}

// CreateWorkItem creates a WorkItem describing this deletion. We call
// this only if the admin approves the deletion.
func (del *Deletion) CreateObjDeletionWorkItem(obj *pgmodels.IntellectualObject) error {
	if del.DeletionRequest == nil || del.DeletionRequest.ID == 0 {
		errMsg := "Cannot create deletion work item because deletion request id is zero."
		common.Context().Log.Error().Msgf("%s", errMsg)
		return errors.New(errMsg)
	}
	common.Context().Log.Warn().Msgf("Creating deletion work item for object %d - %s", obj.ID, obj.Identifier)
	workItem, err := pgmodels.NewDeletionItem(obj, nil, del.DeletionRequest.RequestedBy, del.DeletionRequest.ConfirmedBy, del.DeletionRequest.ID)
	if err != nil {
		common.Context().Log.Error().Msgf("%s", err.Error())
		return err
	}
	common.Context().Log.Warn().Msgf("Created deletion work item %d with deletion request id %d", workItem.ID, workItem.DeletionRequestID)
	del.DeletionRequest.WorkItems = append(del.DeletionRequest.WorkItems, workItem)
	return nil
}

func (del *Deletion) CreateFileDeletionWorkItem(gf *pgmodels.GenericFile) error {
	obj, err := pgmodels.IntellectualObjectByID(gf.IntellectualObjectID)
	if err != nil {
		return err
	}

	workItem, err := pgmodels.NewDeletionItem(obj, gf, del.DeletionRequest.RequestedBy, del.DeletionRequest.ConfirmedBy, del.DeletionRequest.ID)
	if err != nil {
		return err
	}
	del.DeletionRequest.WorkItems = append(del.DeletionRequest.WorkItems, workItem)
	return nil
}

// QueueWorkItem sends the WorkItem.ID into the appropriate NSQ topic.
// We call this after calling CreateWorkItem, and only if the admin
// approves the deletion.
func (del *Deletion) QueueWorkItems() error {
	for _, item := range del.DeletionRequest.WorkItems {
		if item == nil {
			return common.ErrInternal
		}
		topic, err := constants.TopicFor(item.Action, item.Stage)
		if err != nil {
			return err
		}
		ctx := common.Context()
		ctx.Log.Info().Msgf("Queueing WorkItem %d to topic %s", item.ID, topic)
		err = ctx.NSQClient.Enqueue(topic, item.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateAndQueueWorkItem creates and queues a deletion WorkItem.
// We call this only if the admin approves the DeletionRequest.
func (del *Deletion) CreateAndQueueWorkItems() error {
	var err error
	for _, gf := range del.DeletionRequest.GenericFiles {
		err = del.CreateFileDeletionWorkItem(gf)
		if err != nil {
			return err
		}
	}
	for _, obj := range del.DeletionRequest.IntellectualObjects {
		err = del.CreateObjDeletionWorkItem(obj)
		if err != nil {
			return err
		}
	}
	err = del.DeletionRequest.Save()
	if err != nil {
		return err
	}
	return del.QueueWorkItems()
}

// CreateRequestAlert creates an alert saying that a user has requested
// a deletion. This alert goes via email to all admins at the institution
// that owns the file or object to be deleted. This method is supported
// only for new deletion requests. If you try to call this on a deletion
// request you retrieved from the DB, you'll get "operation not supported"
// because we don't have access to the plaintext confirmation token for
// the review URL.
func (del *Deletion) CreateRequestAlert() (*pgmodels.Alert, error) {
	templateName := "alerts/deletion_requested.txt"
	alertType := constants.AlertDeletionRequested
	reviewURL, err := del.ReviewURL()
	if err != nil {
		return nil, err
	}
	alertData := map[string]interface{}{
		"requesterName":       del.currentUser.Name,
		"deletionReviewURL":   reviewURL,
		"deletionReadOnlyURL": del.ReadOnlyURL(),
	}
	return del.createDeletionAlert(templateName, alertType, alertData)
}

// CreateApprovalAlert creates an alert saying that an admin has approved
// a deletion. This alert goes via email to all admins at the institution
// that owns the file or object to be deleted.
func (del *Deletion) CreateApprovalAlert() (*pgmodels.Alert, error) {
	templateName := "alerts/deletion_confirmed.txt"
	alertType := constants.AlertDeletionConfirmed
	workItemURLs, err := del.WorkItemURLs()
	if err != nil {
		return nil, err
	}
	alertData := map[string]interface{}{
		"deletionRequest":     del.DeletionRequest,
		"workItemURLs":        workItemURLs,
		"deletionReadOnlyURL": del.ReadOnlyURL(),
	}
	return del.createDeletionAlert(templateName, alertType, alertData)
}

// CreateCancellationAlert creates an alert saying that an admin has
// rejected a deletion request. This alert goes via email to all admins
// at the institution that owns the file or object to be deleted.
func (del *Deletion) CreateCancellationAlert() (*pgmodels.Alert, error) {
	templateName := "alerts/deletion_cancelled.txt"
	alertType := constants.AlertDeletionCancelled
	alertData := map[string]interface{}{
		"deletionRequest":     del.DeletionRequest,
		"deletionReadOnlyURL": del.ReadOnlyURL(),
	}
	return del.createDeletionAlert(templateName, alertType, alertData)
}

// createDeletionAlert does the grunt work for all of the specific
// deletion alert creation methods.
func (del *Deletion) createDeletionAlert(templateName, alertType string, alertData map[string]interface{}) (*pgmodels.Alert, error) {

	alert := &pgmodels.Alert{
		InstitutionID:     del.DeletionRequest.InstitutionID,
		Type:              alertType,
		Subject:           alertType,
		DeletionRequestID: del.DeletionRequest.ID,
		CreatedAt:         time.Now().UTC(),
		Users:             del.InstAdmins,
	}

	return pgmodels.CreateAlert(alert, templateName, alertData)
}

// ReviewURL returns the URL for an institutional admin to review
// this deletion request. Note that the DeletionRequest has a
// plaintext ConfirmationToken ONLY when created, not when we
// retrieve it from the database. If you call this method after
// retrieving a DeletionRequest, you'll get common.ErrNotSupported,
// because we no longer have access to the plaintext ConfirmationToken.
// This works only after calling initFileDeletionRequest or
// initObjectDeletionRequest.
func (del *Deletion) ReviewURL() (string, error) {
	if del.DeletionRequest.ConfirmationToken == "" {
		return "", common.ErrNotSupported
	}
	return fmt.Sprintf("%s/deletions/review/%d?token=%s",
		del.baseURL,
		del.DeletionRequest.ID,
		del.DeletionRequest.ConfirmationToken), nil
}

// WorkItemURLs returns the URL for the WorkItem for this DeletionRequest.
// If you call this on a cancelled or not-yet-approved request, there is
// no WorkItem and you'll get common.ErrNotSupported.
func (del *Deletion) WorkItemURLs() ([]string, error) {
	urls := make([]string, 0)
	if len(del.DeletionRequest.WorkItems) == 0 {
		return urls, common.ErrNotSupported
	}
	for _, item := range del.DeletionRequest.WorkItems {
		itemUrl := fmt.Sprintf("%s/work_items/show/%d",
			del.baseURL,
			item.ID)
		urls = append(urls, itemUrl)
	}
	return urls, nil
}

// ReadOnlyURL returns a URL that displays info about the deletion request
// but does not include buttons to approve or cancel. This view is for
// the depositor's historical/auditing purposes.
func (del *Deletion) ReadOnlyURL() string {
	return fmt.Sprintf("%s/deletions/show/%d", del.baseURL, del.DeletionRequest.ID)
}
