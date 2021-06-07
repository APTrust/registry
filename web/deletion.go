package web

import (
	"bytes"
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
	instID      int64
	currentUser *pgmodels.User
}

// NewDeletionForFile creates a new DeletionRequest for a GenericFile
// and returns the Deletion object. This constructor is only for initializing
// new DeletionRequests, not for reviewing, approving or cancelling
// existing requests.
func NewDeletionForFile(req *Request) (*Deletion, error) {
	gf, err := pgmodels.GenericFileByID(req.Auth.ResourceID)
	if err != nil {
		return nil, err
	}

	// Make sure there are no pending work items...
	pendingWorkItems, err := pgmodels.WorkItemsPendingForFile(gf.ID)
	if err != nil {
		return nil, err
	}
	if len(pendingWorkItems) > 0 {
		return nil, common.ErrPendingWorkItems
	}

	// Get list of inst admins.
	adminsQuery := pgmodels.NewQuery().
		Where("institution_id", "=", gf.InstitutionID).
		Where("role", "=", constants.RoleInstAdmin)
	instAdmins, err := pgmodels.UserSelect(adminsQuery)
	if err != nil {
		return nil, err
	}

	// Create the deletion request and the alert, which will become
	// the deletion confirmation email.
	deletionRequest, err := pgmodels.NewDeletionRequest()
	if err != nil {
		return nil, err
	}
	deletionRequest.InstitutionID = gf.InstitutionID
	deletionRequest.RequestedByID = req.CurrentUser.ID
	deletionRequest.RequestedAt = time.Now().UTC()
	deletionRequest.AddFile(gf)
	err = deletionRequest.Save()

	return &Deletion{
		DeletionRequest: deletionRequest,
		InstAdmins:      instAdmins,
		baseURL:         req.BaseURL(),
		currentUser:     req.CurrentUser,
		instID:          gf.InstitutionID,
	}, err
}

// NewDeletionForReview pulls up information about an existing deletion
// request that an institutional admin will review before deciding whether
// to approve or cancel the request.
func NewDeletionForReview(req *Request) (*Deletion, error) {
	// Find the deletion request
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if err != nil {
		return nil, err
	}

	// Make sure the token is valid for that deletion request
	token := req.GinContext.PostForm("token")
	if !common.ComparePasswords(deletionRequest.EncryptedConfirmationToken, token) {
		return nil, common.ErrInvalidToken
	}

	adminsQuery := pgmodels.NewQuery().
		Where("institution_id", "=", deletionRequest.InstitutionID).
		Where("role", "=", constants.RoleInstAdmin)
	instAdmins, err := pgmodels.UserSelect(adminsQuery)
	if err != nil {
		return nil, err
	}

	return &Deletion{
		DeletionRequest: deletionRequest,
		InstAdmins:      instAdmins,
		baseURL:         req.BaseURL(),
		currentUser:     req.CurrentUser,
		instID:          deletionRequest.InstitutionID,
	}, nil
}

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

	workItem, err := pgmodels.NewDeletionItem(obj, gf, del.currentUser)
	if err != nil {
		return nil, err
	}
	del.DeletionRequest.WorkItem = workItem
	err = del.DeletionRequest.Save()

	return workItem, err
}

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

func (del *Deletion) CreateAndQueueWorkItem() (*pgmodels.WorkItem, error) {
	workItem, err := del.CreateWorkItem()
	if err == nil {
		err = del.QueueWorkItem()
	}
	return workItem, err
}

func (del *Deletion) CreateRequestAlert() (*pgmodels.Alert, error) {
	templateName := "alerts/deletion_requested.txt"
	alertType := constants.AlertDeletionRequested
	reviewURL, err := del.ReviewURL()
	if err != nil {
		return nil, err
	}
	alertData := map[string]interface{}{
		"RequesterName":     del.currentUser.Name,
		"DeletionReviewURL": reviewURL,
	}
	return del.createAlert(templateName, alertType, alertData)
}

func (del *Deletion) CreateApprovalAlert() (*pgmodels.Alert, error) {
	templateName := "alerts/deletion_confirmed.txt"
	alertType := constants.AlertDeletionConfirmed
	workItemURL, err := del.WorkItemURL()
	if err != nil {
		return nil, err
	}
	alertData := map[string]interface{}{
		"deletionRequest": del.DeletionRequest,
		"workItemURL":     workItemURL,
	}
	return del.createAlert(templateName, alertType, alertData)
}

func (del *Deletion) CreateCancellationAlert() {

}

func (del *Deletion) createAlert(templateName, alertType string, alertData map[string]interface{}) (*pgmodels.Alert, error) {

	// Create the alert text from the template...
	tmpl := common.TextTemplates[templateName]
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, alertData)
	if err != nil {
		return nil, err
	}

	// Create and save the alert, with our custom text,
	// and make sure it's associated with the right
	// recipients. In this case, it goes to institutional
	// admins at the institution that owns the content.
	alert := &pgmodels.Alert{
		InstitutionID:     del.instID,
		Type:              alertType,
		DeletionRequestID: del.DeletionRequest.ID,
		Content:           buf.String(),
		CreatedAt:         time.Now().UTC(),
		Users:             del.InstAdmins,
	}

	err = alert.Save()
	if err != nil {
		return nil, err
	}

	// Show the alert text in dev and test consoles,
	// so we don't have to look it up in the DB.
	// For dev/test, we need to see the review and
	// confirmation URLS in this alert so we can
	// review and test them.
	envName := common.Context().Config.EnvName
	if envName == "dev" || envName == "test" {
		fmt.Println("***********************")
		fmt.Println(alert.Content)
		fmt.Println("***********************")
	}

	return alert, err
}

// ReviewURL returns the URL for an institutional admin to review
// this deletion request. Note that the DeletionRequest has a
// plaintext ConfirmationToken ONLY when created, not when we
// retrieve it from the database. If you call this method after
// retrieving a DeletionRequest, you'll get common.ErrNotSupported,
// because we no longer have access to the plaintext ConfirmationToken.
func (del *Deletion) ReviewURL() (string, error) {
	if del.DeletionRequest.ConfirmationToken == "" {
		return "", common.ErrNotSupported
	}
	return fmt.Sprintf("%s/files/review_delete/%d?token=%s",
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
