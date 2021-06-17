package pgmodels

import (
	"time"
)

type AlertView struct {
	tableName             struct{}  `pg:"alerts_view"`
	ID                    int64     `json:"id"`
	InstitutionID         int64     `json:"institution_id"`
	InstitutionName       string    `json:"institution_name"`
	InstitutionIdentifier string    `json:"institution_identifier"`
	Type                  string    `json:"type"`
	Subject               string    `json:"subject"`
	Content               string    `json:"content"`
	DeletionRequestID     int64     `json:"deletion_request_id"`
	CreatedAt             time.Time `json:"created_at"`
	UserID                int64     `json:"user_id"`
	UserName              string    `json:"user_name"`
	SentAt                time.Time `json:"sent_at"`
	ReadAt                time.Time `json:"read_at"`
}

// AlertViewByID returns the alert with the specified id.
// Returns pg.ErrNoRows if there is no match.
func AlertViewByID(id int64) (*AlertView, error) {
	query := NewQuery().Where(`"alerts_view"."id"`, "=", id)
	return AlertViewGet(query)
}

// AlertViewGet returns the first alert matching the query.
func AlertViewGet(query *Query) (*AlertView, error) {
	var alert AlertView
	err := query.Select(&alert)
	return &alert, err
}

// AlertViewSelect returns all alerts matching the query.
func AlertViewSelect(query *Query) ([]*AlertView, error) {
	var alerts []*AlertView
	err := query.Select(&alerts)
	return alerts, err
}
