package webui

import (
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
	err = del.initObjectDeletionRequest(objID)
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
		Where("role", "=", constants.RoleInstAdmin)
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
func (del *Deletion) initObjectDeletionRequest(objID int64) error {
	obj, err := pgmodels.IntellectualObjectByID(objID)
	if err != nil {
		return err
	}

	deletionRequest, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return err
	}
	deletionRequest.InstitutionID = obj.InstitutionID
	deletionRequest.RequestedByID = del.currentUser.ID
	deletionRequest.RequestedAt = time.Now().UTC()
	deletionRequest.AddObject(obj)
	err = deletionRequest.Save()
	if err != nil {
		return err
	}
	del.DeletionRequest = deletionRequest
	return nil
}

// CreateWorkItem creates a WorkItem describing this deletion. We call
// this only if the admin approves the deletion.
func (del *Deletion) CreateWorkItem() (*pgmodels.WorkItem, error) {
	// Create the deletion WorkItem
	obj := del.DeletionRequest.FirstObject()
	gf := del.DeletionRequest.FirstFile()

	// Deletion may be file only, no object.
	var err error
	if obj == nil && gf != nil {
		obj, err = pgmodels.IntellectualObjectByID(gf.IntellectualObjectID)
		if err != nil {
			return nil, err
		}
	}

	// TODO: WorkItem needs approver.
	workItem, err := pgmodels.NewDeletionItem(obj, gf, del.currentUser)
	if err != nil {
		return nil, err
	}
	del.DeletionRequest.WorkItem = workItem
	err = del.DeletionRequest.Save()

	return workItem, err
}

// QueueWorkItem sends the WorkItem.ID into the appropriate NSQ topic.
// We call this after calling CreateWorkItem, and only if the admin
// approves the deletion.
func (del *Deletion) QueueWorkItem() error {
	workItem := del.DeletionRequest.WorkItem
	if workItem == nil {
		return common.ErrInternal
	}
	topic, err := constants.TopicFor(workItem.Action, workItem.Stage)
	if err != nil {
		return err
	}
	ctx := common.Context()
	ctx.Log.Info().Msgf("Queueing WorkItem %d to topic %s", workItem.ID, topic)
	return ctx.NSQClient.Enqueue(topic, workItem.ID)
}

// CreateAndQueueWorkItem creates and queues a deletion WorkItem.
// We call this only if the admin approves the DeletionRequest.
func (del *Deletion) CreateAndQueueWorkItem() (*pgmodels.WorkItem, error) {
	workItem, err := del.CreateWorkItem()
	if err == nil {
		err = del.QueueWorkItem()
	}
	return workItem, err
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
	workItemURL, err := del.WorkItemURL()
	if err != nil {
		return nil, err
	}
	alertData := map[string]interface{}{
		"deletionRequest":     del.DeletionRequest,
		"workItemURL":         workItemURL,
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

	// createAlert is in account_alert.go
	return createAlert(alert, templateName, alertData)
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

// WorkItemURL returns the URL for the WorkItem for this DeletionRequest.
// If you call this on a cancelled or not-yet-approved request, there is
// no WorkItem and you'll get common.ErrNotSupported.
func (del *Deletion) WorkItemURL() (string, error) {
	if del.DeletionRequest.WorkItemID == 0 {
		return "", common.ErrNotSupported
	}
	return fmt.Sprintf("%s/work_items/show/%d",
		del.baseURL,
		del.DeletionRequest.WorkItemID), nil
}

// ReadOnlyURL returns a URL that displays info about the deletion request
// but does not include buttons to approve or cancel. This view is for
// the depositor's historical/auditing purposes.
func (del *Deletion) ReadOnlyURL() string {
	return fmt.Sprintf("%s/deletions/show/%d", del.baseURL, del.DeletionRequest.ID)
}
