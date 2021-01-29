package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type WorkItemsView struct {
	ID                           int64     `json:"id" form:"id" pg:"id"`
	InstitutionID                int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	InstitutionName              int64     `json:"institution_name" form:"institution_name" pg:"institution_name"`
	IntellectualObjectID         int64     `json:"intellectual_object_id" form:"intellectual_object_id" pg:"intellectual_object_id"`
	IntellectualObjectIdentifier int64     `json:"intellectual_object_identifier" form:"intellectual_object_identifier" pg:"intellectual_object_identifier"`
	GenericFileID                int64     `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	GenericFileIdentifier        int64     `json:"generic_file_identifier" form:"generic_file_identifier" pg:"generic_file_identifier"`
	Name                         string    `json:"name" form:"name" pg:"name"`
	ETag                         string    `json:"etag" form:"etag" pg:"etag"`
	Bucket                       string    `json:"bucket" form:"bucket" pg:"bucket"`
	User                         string    `json:"user" form:"user" pg:"user"`
	Note                         string    `json:"note" form:"note" pg:"note"`
	Action                       string    `json:"action" form:"action" pg:"action"`
	Stage                        string    `json:"stage" form:"stage" pg:"stage"`
	Status                       string    `json:"status" form:"status" pg:"status"`
	Outcome                      string    `json:"outcome" form:"outcome" pg:"outcome"`
	BagDate                      time.Time `json:"bag_date" form:"bag_date" pg:"bag_date"`
	DateProcessed                time.Time `json:"date_processed" form:"date_processed" pg:"date_processed"`
	Retry                        bool      `json:"retry" form:"retry" pg:"retry"`
	Node                         string    `json:"node" form:"node" pg:"node"`
	PID                          int       `json:"pid" form:"pid" pg:"pid"`
	NeedsAdminReview             bool      `json:"needs_admin_review" form:"needs_admin_review" pg:"needs_admin_review"`
	Size                         int64     `json:"size" form:"size" pg:"size"`
	QueuedAt                     time.Time `json:"queued_at" form:"queued_at" pg:"queued_at"`
	StageStartedAt               time.Time `json:"stage_started_at" form:"stage_started_at" pg:"stage_started_at"`
	APTrustApprover              string    `json:"aptrust_approver" form:"aptrust_approver" pg:"aptrust_approver"`
	InstApprover                 string    `json:"inst_approver" form:"inst_approver" pg:"inst_approver"`
	CreatedAt                    time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
}

func (item *WorkItemsView) GetID() int64 {
	return item.ID
}

func (item *WorkItemsView) Authorize(actingUser *User, action string) error {
	perm := "WorkItem" + action
	if !actingUser.HasPermission(constants.Permission(perm), item.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s WorkItem %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, item.ID, item.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (item *WorkItemsView) DeleteIsForbidden() bool {
	return true
}

func (item *WorkItemsView) UpdateIsForbidden() bool {
	return true
}

func (item *WorkItemsView) IsReadOnly() bool {
	return true
}

func (item *WorkItemsView) SupportsSoftDelete() bool {
	return false
}

func (item *WorkItemsView) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (item *WorkItemsView) ClearSoftDeleteAttributes() {
	// No-op
}

func (item *WorkItemsView) SetTimestamps() {
	// No-op, since view is read-only
}

func (item *WorkItemsView) BeforeSave() error {
	// No-op
	return nil
}

func WorkItemsViewFind(id int64) (*WorkItemsView, error) {
	ctx := common.Context()
	item := &WorkItemsView{ID: id}
	err := ctx.DB.Model(item).WherePK().Select()
	return item, err
}
