package pgmodels

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/stretchr/stew/slice"
)

// Filters are defined in IntellectualObjectView, since we query on the view.

type IntellectualObject struct {
	TimestampModel
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Identifier                string         `json:"identifier"`
	AltIdentifier             string         `json:"alt_identifier"`
	Access                    string         `json:"access"`
	BagName                   string         `json:"bag_name"`
	InstitutionID             int64          `json:"institution_id"`
	State                     string         `json:"state"`
	ETag                      string         `json:"etag" pg:"etag"`
	BagGroupIdentifier        string         `json:"bag_group_identifier"`
	StorageOption             string         `json:"storage_option"`
	BagItProfileIdentifier    string         `json:"bagit_profile_identifier" pg:"bagit_profile_identifier"`
	SourceOrganization        string         `json:"source_organization"`
	InternalSenderIdentifier  string         `json:"internal_sender_identifier"`
	InternalSenderDescription string         `json:"internal_sender_description"`
	Institution               *Institution   `json:"institution" pg:"rel:has-one"`
	GenericFiles              []*GenericFile `json:"generic_files" pg:"rel:has-many"`
	PremisEvents              []*PremisEvent `json:"premis_events" pg:"rel:has-many"`
}

// IntellectualObjectByID returns the object with the specified id.
// Returns pg.ErrNoRows if there is no match.
func IntellectualObjectByID(id int64) (*IntellectualObject, error) {
	query := NewQuery().Where(`"intellectual_object"."id"`, "=", id)
	return IntellectualObjectGet(query)
}

// IntellectualObjectByIdentifier returns the object with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func IntellectualObjectByIdentifier(identifier string) (*IntellectualObject, error) {
	query := NewQuery().Where(`"intellectual_object"."identifier"`, "=", identifier)
	return IntellectualObjectGet(query)
}

// IdForFileIdentifier returns the ID of the IntellectualObject
// having the specified identifier.
func IdForObjIdentifier(identifier string) (int64, error) {
	query := NewQuery().Columns("id").Where(`"intellectual_object"."identifier"`, "=", identifier)
	var object IntellectualObject
	err := query.Select(&object)
	return object.ID, err
}

// IntellectualObjectGet returns the first object matching the query.
func IntellectualObjectGet(query *Query) (*IntellectualObject, error) {
	var object IntellectualObject
	err := query.Relations("Institution").Select(&object)
	return &object, err
}

// IntellectualObjectSelect returns all objects matching the query.
func IntellectualObjectSelect(query *Query) ([]*IntellectualObject, error) {
	var objects []*IntellectualObject
	err := query.Select(&objects)
	return objects, err
}

// CountObjectsThatCanBeDeleted returns the number of active objects in the
// list of object IDs that belong to the specified institution. We use this
// when running batch deletions to ensure that no objects belong to an
// institution other than the one requesting the deletion.
//
// If we get a list of 100 ids, the return value should be 100. If it's not
// some object in the ID list was already deleted, or it belongs to someone
// else.
func CountObjectsThatCanBeDeleted(institutionID int64, objIDs []int64) (int, error) {
	return common.Context().DB.Model((*IntellectualObject)(nil)).Where(`institution_id = ? and state = 'A' and id in (?)`, institutionID, pg.In(objIDs)).Count()
}

// Save saves this object to the database. This will peform an insert
// if IntellectualObject.ID is zero. Otherwise, it updates.
func (obj *IntellectualObject) Save() error {
	obj.SetTimestamps()
	err := obj.Validate()
	if err != nil {
		return err
	}
	if obj.ID == int64(0) {
		return insert(obj)
	}

	// Insert relations, as in DeletionRequest?

	return update(obj)
}

// IsGlacierOnly returns true if this object is stored only
// in Glacier.
func (obj *IntellectualObject) IsGlacierOnly() bool {
	return isGlacierOnly(obj.StorageOption)
}

// Delete soft-deletes this object by setting State to 'D' and
// the UpdatedAt timestamp to now. You can undo this with Undelete.
// It also creates a deletion PremisEvent. You can't get rid of that.
//
// It is legitimate for a depositor to delete an object, then re-upload
// it later, particularly if they want to change the storage option.
// In that case, the object's state would be set back to "A" after the
// new ingest, and the old deletion event would remain to show that an earlier
// version of the object was once deleted.
//
// We would know the new object is active because state = "A" and it would
// have an ingest event dated after the last deletion event.
func (obj *IntellectualObject) Delete() error {

	err := obj.AssertDeletionPreconditions()
	if err != nil {
		return err
	}

	obj.State = constants.StateDeleted
	obj.UpdatedAt = time.Now().UTC()

	// We require a BagItProfileIdentifier on the object.
	// However, older objects (before 2022) do not have one.
	// If the BagItProfileIdentifier is not set on the object, we assign it to a default value.
	// In this way, we can validate the object successfully before deleting it.
	if obj.BagItProfileIdentifier == "" {
		obj.BagItProfileIdentifier = constants.DefaultProfileIdentifier
	}

	valErr := obj.Validate()
	if valErr != nil {
		return valErr
	}

	deletionEvent, err := obj.NewDeletionEvent()
	if err != nil {
		return err
	}
	deletionEvent.SetTimestamps()
	valErr = deletionEvent.Validate()
	if valErr != nil {
		return valErr
	}

	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		var err error
		_, err = tx.Model(obj).WherePK().Update()
		if err != nil {
			registryContext.Log.Error().Msgf("Intellectual object deletion transaction failed on update of object. Object: %d (%s). Error: %v", obj.ID, obj.Identifier, err)
			j, _ := json.Marshal(obj)
			registryContext.Log.Error().Msgf("Object: %s", string(j))
		}
		_, err = tx.Model(deletionEvent).Insert()
		if err != nil {
			registryContext.Log.Error().Msgf("Intellectual object deletion transaction failed on insertion of event. Object: %d (%s). Error: %v", obj.ID, obj.Identifier, err)
			j, _ := json.Marshal(deletionEvent)
			registryContext.Log.Error().Msgf("DeletionEvent: %s", string(j))
		}
		return err
	})
}

// HasActiveFiles returns true if this object has any active (non-deleted)
// files. We need to check this before marking an object as deleted.
// Do not mark deleted until all files have been marked deleted.
func (obj *IntellectualObject) HasActiveFiles() (bool, error) {
	db := common.Context().DB
	return db.Model((*GenericFile)(nil)).Where("intellectual_object_id = ? and state = ?", obj.ID, constants.StateActive).Exists()
}

// LastIngestEvent returns the latest ingest event for this object.
// This should never be nil.
func (obj *IntellectualObject) LastIngestEvent() (*PremisEvent, error) {
	return obj.lastEvent(constants.EventIngestion)
}

// LastDeleationEvent returns the latest deletion event for this object,
// which may be nil.
func (obj *IntellectualObject) LastDeletionEvent() (*PremisEvent, error) {
	return obj.lastEvent(constants.EventDeletion)
}

func (obj *IntellectualObject) lastEvent(eventType string) (*PremisEvent, error) {
	query := NewQuery().
		Where("intellectual_object_id", "=", obj.ID).
		Where("event_type", "=", eventType).
		IsNull("generic_file_id").
		OrderBy("created_at", "desc").
		Offset(0).
		Limit(1)
	return PremisEventGet(query)
}

func isGlacierOnly(storageOption string) bool {
	return slice.Contains(constants.GlacierOnlyOptions, storageOption)
}

func (obj *IntellectualObject) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(obj.Title) {
		errors["Title"] = "Title is required"
	}
	if common.IsEmptyString(obj.Identifier) {
		errors["Identifier"] = "Identifier is required"
	}
	if !v.IsIn(obj.State, constants.States...) {
		errors["State"] = ErrInstState
	}
	if !v.IsIn(obj.Access, constants.AccessSettings...) {
		errors["Access"] = "Invalid access value"
	}
	if obj.InstitutionID < 1 {
		errors["InstitutionID"] = "Invalid institution id"
	}
	if !v.IsIn(obj.StorageOption, constants.StorageOptions...) {
		errors["StorageOption"] = "Invalid storage option"
	}
	if common.IsEmptyString(obj.BagItProfileIdentifier) {
		errors["BagItProfileIdentifier"] = "BagItProfileIdentifier is required"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}

func (obj *IntellectualObject) ValidateChanges(updatedObj *IntellectualObject) error {
	if obj.ID != updatedObj.ID {
		return common.ErrIDMismatch
	}
	if obj.InstitutionID != updatedObj.InstitutionID {
		return common.ErrInstIDChange
	}
	if obj.Identifier != updatedObj.Identifier {
		return common.ErrIdentifierChange
	}
	// Caller should force storage option of updated object to
	// match existing object before calling this validation function.
	if obj.State == constants.StateActive && obj.StorageOption != updatedObj.StorageOption {
		return common.ErrStorageOptionChange
	}
	return nil
}

func (obj *IntellectualObject) ActiveDeletionWorkItem() (*WorkItem, error) {
	query := NewQuery().
		Where("intellectual_object_id", "=", obj.ID).
		IsNull("generic_file_id").
		Where("action", "=", constants.ActionDelete).
		Where("status", "=", constants.StatusStarted).
		OrderBy("updated_at", "desc").
		Limit(1)
	item, err := WorkItemGet(query)
	if err != nil && err.Error() == pg.ErrNoRows.Error() {
		// >99% of cases will have no rows. Just return nil.
		return nil, nil
	}
	return item, err
}

func (obj *IntellectualObject) AssertDeletionPreconditions() error {
	err := obj.assertNoActiveFiles()
	if err == nil {
		err = obj.assertNotAlreadyDeleted()
	}
	if err == nil {
		if !obj.HasPassedMinimumRetentionPeriod() {
			err = fmt.Errorf("Object has not passed minimum retention period")
		}
	}
	if err == nil {
		_, _, err = obj.assertDeletionApproved()
	}
	if err != nil {
		common.Context().Log.Error().Msgf(
			"Deletion precondition check failed for object %d (%s): %v",
			obj.ID, obj.Identifier, err)
	}
	return err
}

func (obj *IntellectualObject) assertNoActiveFiles() error {
	hasFiles, err := obj.HasActiveFiles()
	if err != nil {
		return fmt.Errorf("Error checking for active files: %v", err)
	} else if hasFiles {
		return common.ErrActiveFiles
	}
	return nil
}

func (obj *IntellectualObject) assertNotAlreadyDeleted() error {
	var err error
	var lastIngestEvent *PremisEvent

	if obj.State == constants.StateDeleted {
		err = fmt.Errorf("Object is already in deleted state")
	}

	if err == nil {
		lastIngestEvent, err = obj.LastIngestEvent()
		if err != nil {
			err = fmt.Errorf("Error checking for last ingest event: %v", err)
		} else if lastIngestEvent == nil {
			err = fmt.Errorf("Can't find last ingest event")
		}
	}

	if err == nil {
		lastDeletionEvent, err := obj.LastDeletionEvent()
		if err != nil {
			err = fmt.Errorf("Error checking for last deletion event: %v", err)
		}
		if lastDeletionEvent != nil && lastDeletionEvent.CreatedAt.After(lastIngestEvent.CreatedAt) {
			err = fmt.Errorf("Object has already been deleted since last ingest")
		}
	}
	return err
}

func (obj *IntellectualObject) assertDeletionApproved() (*WorkItem, *DeletionRequestView, error) {
	workItem, err := obj.ActiveDeletionWorkItem()
	if workItem == nil || IsNoRowError(err) {
		return nil, nil, fmt.Errorf("Missing deletion request work item")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting active deletion work item: %v", err)
	}
	if common.IsEmptyString(workItem.InstApprover) {
		return workItem, nil, fmt.Errorf("Deletion work item is missing institutional approver")
	}
	deletionRequest, err := DeletionRequestViewByID(workItem.DeletionRequestID)
	if deletionRequest == nil || IsNoRowError(err) {
		fmt.Println(workItem.ID, workItem.DeletionRequestID)
		return workItem, nil, fmt.Errorf("No deletion request for work item %d", workItem.ID)
	}
	if err != nil {
		return workItem, nil, fmt.Errorf("Error getting deletion request: %v", err)
	}
	if deletionRequest.RequestedByID == 0 {
		// We should never hit this because RequestedByID has a not-null constraint.
		return workItem, deletionRequest, fmt.Errorf("Deletion request %d has no requestor", deletionRequest.ID)
	}
	if deletionRequest.ConfirmedByID == 0 {
		return workItem, deletionRequest, fmt.Errorf("Deletion request %d has no approver", deletionRequest.ID)
	}
	return workItem, deletionRequest, nil
}

func (obj *IntellectualObject) NewDeletionEvent() (*PremisEvent, error) {
	_, deletionRequestView, err := obj.assertDeletionApproved()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &PremisEvent{
		Agent:                "APTrust preservation services",
		DateTime:             now,
		Detail:               "Object deleted from preservation storage",
		EventType:            constants.EventDeletion,
		Identifier:           uuid.NewString(),
		InstitutionID:        obj.InstitutionID,
		IntellectualObjectID: obj.ID,
		Object:               "Minio S3 library",
		Outcome:              constants.OutcomeSuccess,
		OutcomeDetail:        deletionRequestView.RequestedByEmail,
		OutcomeInformation:   fmt.Sprintf("Object deleted at the request of %s. Institutional approver: %s.", deletionRequestView.RequestedByEmail, deletionRequestView.ConfirmedByEmail),
	}, nil
}

// EarliestDeletionDate returns the earliest date on which this
// object can be deleted, per retention rules that apply to the
// object's storage option.
//
// This is generally accurate, but can't be 100% accurate, as some
// of the object's files may have been ingested after the object
// creation date.
//
// Also, the following (rare) case will return a false positive:
// object was ingested five years ago, deleted four years ago,
// and then ingested again yesterday.
//
// We can sort through Premis Events to solve these false positives,
// but that's very expensive and false positives probably are
// less than 0.2% of all cases.
func (obj *IntellectualObject) EarliestDeletionDate() time.Time {
	minRetentionDays := common.Context().Config.RetentionMinimum.For(obj.StorageOption)
	return obj.CreatedAt.AddDate(0, 0, minRetentionDays)
}

// HasPassedMinimumRetentionPeriod returns true if this object has
// passed the minimum retention period for its storage option.
func (obj *IntellectualObject) HasPassedMinimumRetentionPeriod() bool {
	return obj.EarliestDeletionDate().Before(time.Now())
}
