package pgmodels

import (
	"context"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
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

var WorkItemFilters = []string{
	"action",
	"bag_date",
	"bucket",
	"date_processed",
	"etag",
	"generic_file_id",
	"instutition_id",
	"intellectual_object_id",
	"name",
	"needs_admin_review",
	"node",
	"pid",
	"queued__is_null",
	"queued__not_null",
	"retry",
	"size__gteq",
	"size__lteq",
	"stage",
	"stage_started_at__is_null",
	"stage_started_at__not_null",
	"status",
	"user",
}

type WorkItem struct {
	ID                   int64     `json:"id" pg:"id"`
	Name                 string    `json:"name" pg:"name"`
	ETag                 string    `json:"etag" pg:"etag"`
	InstitutionID        int64     `json:"institution_id" pg:"institution_id"`
	IntellectualObjectID int64     `json:"intellectual_object_id" pg:"intellectual_object_id"`
	GenericFileID        int64     `json:"generic_file_id" pg:"generic_file_id"`
	Bucket               string    `json:"bucket" pg:"bucket"`
	User                 string    `json:"user" pg:"user"`
	Note                 string    `json:"note" pg:"note"`
	Action               string    `json:"action" pg:"action"`
	Stage                string    `json:"stage" pg:"stage"`
	Status               string    `json:"status" pg:"status"`
	Outcome              string    `json:"outcome" pg:"outcome"`
	BagDate              time.Time `json:"bag_date" pg:"bag_date"`
	DateProcessed        time.Time `json:"date_processed" pg:"date_processed"`
	Retry                bool      `json:"retry" pg:"retry,use_zero"`
	Node                 string    `json:"node" pg:"node"`
	PID                  int       `json:"pid" pg:"pid"`
	NeedsAdminReview     bool      `json:"needs_admin_review" pg:"needs_admin_review,use_zero"`
	QueuedAt             time.Time `json:"queued_at" pg:"queued_at"`
	Size                 int64     `json:"size" pg:"size"`
	StageStartedAt       time.Time `json:"stage_started_at" pg:"stage_started_at"`
	APTrustApprover      string    `json:"aptrust_approver" pg:"aptrust_approver"`
	InstApprover         string    `json:"inst_approver" pg:"inst_approver"`
	CreatedAt            time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" pg:"updated_at"`
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
	topic := constants.TopicFor(item.Action, stage)
	if topic == "" {
		return common.ErrInvalidRequeue
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
