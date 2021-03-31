package pgmodels

import (
	"context"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

type WorkItem struct {
	ID                   int64     `json:"id" form:"id" pg:"id"`
	Name                 string    `json:"name" form:"name" pg:"name"`
	ETag                 string    `json:"etag" form:"etag" pg:"etag"`
	InstitutionID        int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	IntellectualObjectID int64     `json:"intellectual_object_id" form:"intellectual_object_id" pg:"intellectual_object_id"`
	GenericFileID        int64     `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	Bucket               string    `json:"bucket" form:"bucket" pg:"bucket"`
	User                 string    `json:"user" form:"user" pg:"user"`
	Note                 string    `json:"note" form:"note" pg:"note"`
	Action               string    `json:"action" form:"action" pg:"action"`
	Stage                string    `json:"stage" form:"stage" pg:"stage"`
	Status               string    `json:"status" form:"status" pg:"status"`
	Outcome              string    `json:"outcome" form:"outcome" pg:"outcome"`
	BagDate              time.Time `json:"bag_date" form:"bag_date" pg:"bag_date"`
	DateProcessed        time.Time `json:"date_processed" form:"date_processed" pg:"date_processed"`
	Retry                bool      `json:"retry" form:"retry" pg:"retry"`
	Node                 string    `json:"node" form:"node" pg:"node"`
	PID                  int       `json:"pid" form:"pid" pg:"pid"`
	NeedsAdminReview     bool      `json:"needs_admin_review" form:"needs_admin_review" pg:"needs_admin_review"`
	QueuedAt             time.Time `json:"queued_at" form:"queued_at" pg:"queued_at"`
	Size                 int64     `json:"size" form:"size" pg:"size"`
	StageStartedAt       time.Time `json:"stage_started_at" form:"stage_started_at" pg:"stage_started_at"`
	APTrustApprover      string    `json:"aptrust_approver" form:"aptrust_approver" pg:"aptrust_approver"`
	InstApprover         string    `json:"inst_approver" form:"inst_approver" pg:"inst_approver"`
	CreatedAt            time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
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

func (item *WorkItem) GetID() int64 {
	return item.ID
}

// Save saves this work item to the database. This will peform an insert
// if WorkItem.ID is zero. Otherwise, it updates.
func (item *WorkItem) Save() error {
	if item.ID == int64(0) {
		return insert(item)
	}
	return update(item)
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
	err := item.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (item *WorkItem) BeforeUpdate(c context.Context) (context.Context, error) {
	return c, nil
}

func (item *WorkItem) Validate() *common.ValidationError {
	// TODO: Validate required. Validate biz rules.
	return nil
}
