package pgmodels

import (
	"context"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
	"github.com/jinzhu/copier"
	"github.com/stretchr/stew/slice"
)

const (
	ErrItemName          = "Name is required."
	ErrItemETag          = "ETag is required."
	ErrItemBagDate       = "BagDate is required."
	ErrItemBucket        = "Bucket is required."
	ErrItemUser          = "User must be a valid email address."
	ErrItemInstID        = "InstitutionID is required."
	ErrItemDateProcessed = "DateProcessed is required."
	ErrItemNote          = "Note cannot be empty."
	ErrItemAction        = "Action is missing or invalid."
	ErrItemStage         = "Stage is missing or invalid."
	ErrItemStatus        = "Status is missing or invalid."
	ErrItemOutcome       = "Outcome cannot be empty."
)

type WorkItem struct {
	ID                   int64     `json:"id" pg:"id"`
	Name                 string    `json:"name" pg:"name"`
	ETag                 string    `json:"etag" pg:"etag"`
	InstitutionID        int64     `json:"institution_id"`
	IntellectualObjectID int64     `json:"intellectual_object_id"`
	GenericFileID        int64     `json:"generic_file_id"`
	Bucket               string    `json:"bucket"`
	User                 string    `json:"user"`
	Note                 string    `json:"note"`
	Action               string    `json:"action"`
	Stage                string    `json:"stage"`
	Status               string    `json:"status"`
	Outcome              string    `json:"outcome"`
	BagDate              time.Time `json:"bag_date"`
	DateProcessed        time.Time `json:"date_processed"`
	Retry                bool      `json:"retry" pg:",use_zero"`
	Node                 string    `json:"node"`
	PID                  int       `json:"pid"`
	NeedsAdminReview     bool      `json:"needs_admin_review" pg:",use_zero"`
	QueuedAt             time.Time `json:"queued_at"`
	Size                 int64     `json:"size"`
	StageStartedAt       time.Time `json:"stage_started_at"`
	APTrustApprover      string    `json:"aptrust_approver" pg:"aptrust_approver"`
	InstApprover         string    `json:"inst_approver"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// WorkItemByID returns the work item with the specified id.
// Returns pg.ErrNoRows if there is no match.
func WorkItemByID(id int64) (*WorkItem, error) {
	query := NewQuery().Where("id", "=", id)
	return WorkItemGet(query)
}

// WorkItemGet returns the first work item matching the query.
func WorkItemGet(query *Query) (*WorkItem, error) {
	var item WorkItem
	err := query.Select(&item)
	return &item, err
}

// WorkItemSelect returns all work items matching the query.
func WorkItemSelect(query *Query) ([]*WorkItem, error) {
	var items []*WorkItem
	err := query.Select(&items)
	return items, err
}

// WorkItemsPendingForObject returns a list of in-progress WorkItems
// for the IntellectualObject with the specified institution ID and
// bag name. We don't use an IntellectualObjectID here because when
// we're ingesting or re-ingesting an object, the WorkItem won't have
// an ObjectID until the ingest/re-ingest is complete.
//
// This method is called before initializing a new restoration or deletion
// request. We specifically want to avoid the case where a user requests a
// restoration or deletion on an object that is about to be reingested.
// If that were to happen, the delete worker would be deleting files
// that an ingest worker just wrote. Or the ingest worker would be
// overwriting files that the restore worker was trying to restore.
//
// Pharos queried by object id, which was a mistake that would not
// catch re-ingests. This corrects that.
func WorkItemsPendingForObject(instID int64, bagName string) ([]*WorkItem, error) {
	completed := common.InterfaceList(constants.CompletedStatusValues)
	query := NewQuery().Where("institution_id", "=", instID).
		Where("name", "=", bagName).
		WhereNotIn("status", completed...).
		OrderBy(`date_processed desc`)
	return WorkItemSelect(query)
}

// WorkItemsPendingForFile returns a list of in-progress WorkItems
// for the GenericFile with the specified ID.
func WorkItemsPendingForFile(fileID int64) ([]*WorkItem, error) {
	completed := common.InterfaceList(constants.CompletedStatusValues)
	query := NewQuery().Where("generic_file_id", "=", fileID).
		WhereNotIn("status", completed...).
		OrderBy(`date_processed desc`)
	return WorkItemSelect(query)
}

// GetID returns this item's ID.
func (item *WorkItem) GetID() int64 {
	return item.ID
}

// HasCompleted returns true if this item has completed processing.
func (item *WorkItem) HasCompleted() bool {
	return slice.Contains(constants.CompletedStatusValues, item.Status)
}

// Save saves this work item to the database. This will peform an insert
// if WorkItem.ID is zero. Otherwise, it updates.
func (item *WorkItem) Save() error {
	if item.ID == int64(0) {
		return insert(item)
	}
	return update(item)
}

// SetForRequeue sets properies so this item can be requeued.
// Note that it saves the object. It will return common.ErrInvalidRequeue
// if the stage is not valid, and may return validation or pg error
// if the object cannot be saved.
//
// The call is responsible for actually pushing the WorkItem.ID into
// the correct NSQ topic.
func (item *WorkItem) SetForRequeue(stage string) error {
	_, err := constants.TopicFor(item.Action, stage)
	if err != nil {
		return err
	}
	item.Stage = stage
	item.Status = constants.StatusPending
	item.Retry = true
	item.NeedsAdminReview = false
	item.Node = ""
	item.PID = 0
	item.Note = fmt.Sprintf("Requeued for %s", item.Stage)
	return item.Save()
}

// The following statements have no effect other than to force a compile-time
// check that ensures our WorkItem model properly implements these hook
// interfaces.
var (
	_ pg.BeforeInsertHook = (*WorkItem)(nil)
	_ pg.BeforeUpdateHook = (*WorkItem)(nil)
)

// BeforeInsert validates the record and does additional prep work.
func (item *WorkItem) BeforeInsert(c context.Context) (context.Context, error) {
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	err := item.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (item *WorkItem) BeforeUpdate(c context.Context) (context.Context, error) {
	item.UpdatedAt = time.Now().UTC()
	return c, nil
}

func (item *WorkItem) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if !v.IsByteLength(item.Name, 1, 1000) {
		errors["Name"] = ErrItemName
	}
	if !v.IsByteLength(item.ETag, 32, 40) {
		errors["ETag"] = ErrItemETag
	}
	if item.BagDate.IsZero() {
		errors["BagDate"] = ErrItemBagDate
	}
	if !v.IsByteLength(item.Bucket, 1, 1000) {
		errors["Bucket"] = ErrItemBucket
	}
	if !v.IsEmail(item.User) {
		errors["User"] = ErrItemUser
	}
	if item.InstitutionID < 1 {
		errors["InstitutionID"] = ErrItemInstID
	}
	if item.DateProcessed.IsZero() {
		errors["DateProcessed"] = ErrItemDateProcessed
	}
	if !v.IsByteLength(item.Name, 1, 10000) {
		errors["Note"] = ErrItemNote
	}
	if !v.IsIn(item.Action, constants.WorkItemActions...) {
		errors["Action"] = ErrItemAction
	}
	if !v.IsIn(item.Stage, constants.Stages...) {
		errors["Stage"] = ErrItemStage
	}
	if !v.IsIn(item.Status, constants.Statuses...) {
		errors["Status"] = ErrItemStatus
	}
	if !v.IsByteLength(item.Name, 1, 1000) {
		errors["Outcome"] = ErrItemOutcome
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}

// LastSuccessfulIngest returns the last successful
// ingest WorkItem for the specified intellectual object id.
func LastSuccessfulIngest(objID int64) (*WorkItem, error) {
	//db := common.Context().DB
	query := NewQuery().
		Where("intellectual_object_id", "=", objID).
		Where("status", "=", constants.StatusSuccess).
		WhereIn("stage", constants.StageRecord, constants.StageCleanup).
		OrderBy("date_processed desc").
		Limit(1)
	items, err := WorkItemSelect(query)
	if len(items) > 0 {
		return items[0], err
	}
	return nil, err
}

// NewItemFromLastSuccessfulIngest creates a new WorkItem based on
// the last successful ingest WorkItem of the specified object.
// This is used for creating various deletion and restoration WorkItems.
// The returned WorkItem will include the proper object name, object id,
// object identifier and etag. All other fields will be cleared out.
// The caller must set essential fields like Action, User, GenericFileID
// (if appropriate) and the like.
//
// This will return an error if the system can't find the last
// successful ingest record for the specified object.
func NewItemFromLastSuccessfulIngest(objID int64) (*WorkItem, error) {
	item, err := LastSuccessfulIngest(objID)
	if err != nil {
		return nil, err
	}
	newItem := &WorkItem{}
	err = copier.Copy(&newItem, item)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	// Reset essential fields
	newItem.ID = 0
	newItem.CreatedAt = now
	newItem.DateProcessed = now
	newItem.NeedsAdminReview = false
	newItem.Node = ""
	newItem.Note = "Not started"
	newItem.Outcome = "Not started"
	newItem.PID = 0
	newItem.QueuedAt = time.Time{}
	newItem.Retry = true
	newItem.Stage = constants.StageRequested
	newItem.StageStartedAt = time.Time{}
	newItem.Status = constants.StatusPending
	newItem.UpdatedAt = now

	return newItem, err
}

// CreateObjectRestorationItem creates and saves a new WorkItem
// for an object or file restoration.
//
// Param obj (required) is the object to be restored.
// gf is the GenericFile to be restored. This can be zero
// if we're restoring an object instead of a file. Param user is the
// user initiating the restoration.
//
// Before creating a restoration WorkItem, the caller should ensure
// that the object and file have no pending work items. See
// WorkItemsPendingForObject() and WorkItemsPendinForFile().
func NewRestorationItem(obj *IntellectualObject, gf *GenericFile, user *User) (*WorkItem, error) {
	if obj == nil {
		return nil, common.ErrInvalidParam
	}

	restorationItem, err := NewItemFromLastSuccessfulIngest(obj.ID)
	if err != nil {
		return nil, err
	}

	// Figure out the restoration type. This determines which
	// queue it will go into and which worker will handle it.
	if obj.IsGlacierOnly() {
		restorationItem.Action = constants.ActionGlacierRestore
	} else {
		// TODO: https://trello.com/c/GirQ712I
		if gf != nil {
			restorationItem.Action = constants.ActionRestoreFile
		} else {
			restorationItem.Action = constants.ActionRestoreObject
		}
	}

	// If this is a file restoration, we have to set the
	// generic file id.
	if gf != nil {
		restorationItem.GenericFileID = gf.ID
	}
	restorationItem.User = user.Email
	err = restorationItem.Save()
	return restorationItem, err
}

func NewDeletionItem(obj *IntellectualObject, gf *GenericFile, user *User) (*WorkItem, error) {
	if obj == nil {
		return nil, common.ErrInvalidParam
	}

	deletionItem, err := NewItemFromLastSuccessfulIngest(obj.ID)
	if err != nil {
		return nil, err
	}

	// If file deletion, set the file id & override object
	// with file size
	if gf != nil {
		deletionItem.GenericFileID = gf.ID
		deletionItem.Size = gf.Size
	}

	deletionItem.Action = constants.ActionDelete
	deletionItem.User = user.Email
	err = deletionItem.Save()
	return deletionItem, err
}
