package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
)

type WorkItemView struct {
	tableName                struct{}  `pg:"work_items_view"`
	ID                       int64     `json:"id" pg:"id"`
	Name                     string    `json:"name" pg:"name"`
	ETag                     string    `json:"etag" pg:"etag"`
	InstitutionID            int64     `json:"institution_id" pg:"institution_id"`
	InstitutionName          string    `json:"institution_name" pg:"institution_id"`
	IntellectualObjectID     int64     `json:"intellectual_object_id" pg:"intellectual_object_id"`
	ObjectIdentifier         string    `json:"object_identifier" pg:"object_identifier"`
	AltIdentifier            string    `json:"alt_identifier" pg:"alt_identifier"`
	BagGroupIdentifier       string    `json:"bag_group_identifier" pg:"bag_group_identifier"`
	StorageOption            string    `json:"storage_option" pg:"storage_option"`
	BagItProfileIdentifier   string    `json:"bagit_profile_identifier" pg:"bagit_profile_identifier"`
	SourceOrganization       string    `json:"source_organization" pg:"source_organization"`
	InternalSenderIdentifier string    `json:"internal_sender_identifier" pg:"internal_sender_identifier"`
	GenericFileID            int64     `json:"generic_file_id" pg:"generic_file_id"`
	GenericFileIdentifier    string    `json:"generic_file_identifier" pg:"generic_file_identifier"`
	Bucket                   string    `json:"bucket" pg:"bucket"`
	User                     string    `json:"user" pg:"user"`
	Note                     string    `json:"note" pg:"note"`
	Action                   string    `json:"action" pg:"action"`
	Stage                    string    `json:"stage" pg:"stage"`
	Status                   string    `json:"status" pg:"status"`
	Outcome                  string    `json:"outcome" pg:"outcome"`
	BagDate                  time.Time `json:"bag_date" pg:"bag_date"`
	DateProcessed            time.Time `json:"date_processed" pg:"date_processed"`
	Retry                    bool      `json:"retry" pg:"retry"`
	Node                     string    `json:"node" pg:"node"`
	PID                      int       `json:"pid" pg:"pid"`
	NeedsAdminReview         bool      `json:"needs_admin_review" pg:"needs_admin_review"`
	QueuedAt                 time.Time `json:"queued_at" pg:"queued_at"`
	Size                     int64     `json:"size" pg:"size"`
	StageStartedAt           time.Time `json:"stage_started_at" pg:"stage_started_at"`
	APTrustApprover          string    `json:"aptrust_approver" pg:"aptrust_approver"`
	InstApprover             string    `json:"inst_approver" pg:"inst_approver"`
	CreatedAt                time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" pg:"updated_at"`
}
