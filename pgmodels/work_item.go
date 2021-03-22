package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
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
