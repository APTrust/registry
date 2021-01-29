package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

func (item *WorkItem) GetID() int64 {
	return item.ID
}

func (item *WorkItem) Authorize(actingUser *User, action string) error {
	perm := "WorkItem" + action
	if !actingUser.HasPermission(constants.Permission(perm), item.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s WorkItem %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, item.ID, item.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

// DeleteIsForbidden returns true because WorkItems are our audit trail.
func (item *WorkItem) DeleteIsForbidden() bool {
	return true
}

func (item *WorkItem) UpdateIsForbidden() bool {
	return false
}

func (item *WorkItem) IsReadOnly() bool {
	return false
}

func (item *WorkItem) SupportsSoftDelete() bool {
	return false
}

func (item *WorkItem) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (item *WorkItem) ClearSoftDeleteAttributes() {
	// No-op
}

func (item *WorkItem) SetTimestamps() {
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now
}

func (item *WorkItem) BeforeSave() error {
	// TODO: Validate
	return nil
}

func WorkItemFind(id int64) (*WorkItem, error) {
	ctx := common.Context()
	item := &WorkItem{ID: id}
	err := ctx.DB.Model(item).WherePK().Select()
	return item, err
}
