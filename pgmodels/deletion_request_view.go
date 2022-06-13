package pgmodels

import (
	"time"
)

var DeletionRequestFilters = []string{
	"institution_id",
	"requested_at__gteq",
	"requested_at__lteq",
	"stage",
	"status",
}

// DeletionRequestsView contains a flattened view of deletion requests
// suitable for the index page.
type DeletionRequestView struct {
	tableName             struct{}  `pg:"deletion_requests_view"`
	ID                    int64     `json:"id"`
	InstitutionID         int64     `json:"institution_id"`
	InstitutionName       string    `json:"institution_name"`
	InstitutionIdentifier string    `json:"institution_identifier"`
	RequestedByID         int64     `json:"requested_by_id"`
	RequestedByName       string    `json:"requested_by_name"`
	RequestedByEmail      string    `json:"requested_by_email"`
	RequestedAt           time.Time `json:"requested_at"`
	ConfirmedByID         int64     `json:"confirmed_by_id"`
	ConfirmedByName       string    `json:"confirmed_by_name"`
	ConfirmedByEmail      string    `json:"confirmed_by_email"`
	ConfirmedAt           time.Time `json:"confirmed_at"`
	CancelledByID         int64     `json:"cancelled_by_id"`
	CancelledByName       string    `json:"cancelled_by_name"`
	CancelledByEmail      string    `json:"cancelled_by_email"`
	CancelledAt           time.Time `json:"cancelled_at"`
	FileCount             int64     `json:"file_count"`
	ObjectCount           int64     `json:"object_count"`
	WorkItemID            int64     `json:"work_item_id"`
	Stage                 string    `json:"stage"`
	Status                string    `json:"status"`
	DateProcessed         time.Time `json:"date_processed"`
	Size                  int64     `json:"size"`
	Note                  string    `json:"note"`
}

// DeletionRequestViewByID returns the DeletionRequestView record
// with the specified id.  Returns pg.ErrNoRows if there is no match.
func DeletionRequestViewByID(id int64) (*DeletionRequestView, error) {
	query := NewQuery().Where("id", "=", id)
	return DeletionRequestViewGet(query)
}

// DeletionRequestViewSelect returns all DeletionRequestView records matching
// the query.
func DeletionRequestViewSelect(query *Query) ([]*DeletionRequestView, error) {
	var requests []*DeletionRequestView
	err := query.Select(&requests)
	return requests, err
}

// DeletionRequestViewGet returns the first user view record matching the query.
func DeletionRequestViewGet(query *Query) (*DeletionRequestView, error) {
	var request DeletionRequestView
	err := query.Select(&request)
	if request.ID == 0 {
		return nil, err
	}
	return &request, err
}

// DisplayStatus returns a string saying whether this deletion request
// has been cancelled, is in progress, or complete, or whatever.
func (request *DeletionRequestView) DisplayStatus() string {
	if request.CancelledByID > 0 {
		return "Rejected"
	}
	if request.Status != "" {
		return request.Status
	}
	return "Awaiting Approval"
}
