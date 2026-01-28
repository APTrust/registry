package pgmodels

import (
	"time"
)

var PremisEventFilters = []string{
	"date_time__gteq",
	"date_time__lteq",
	"event_type",
	"generic_file_id",
	"generic_file_id__is_null",
	"generic_file_identifier",
	"identifier",
	"institution_id",
	"intellectual_object_id",
	"intellectual_object_identifier",
	"outcome",
}

type PremisEventView struct {
	tableName                    struct{}  `pg:"premis_events_view"`
	ID                           int64     `json:"id" form:"id"`
	Agent                        string    `json:"agent"`
	DateTime                     time.Time `json:"date_time"`
	Detail                       string    `json:"detail"`
	EventType                    int32     `json:"event_type"`
	GenericFileID                int64     `json:"generic_file_id"`
	GenericFileIdentifier        string    `json:"generic_file_identifier"`
	Identifier                   string    `json:"identifier"`
	InstitutionID                int64     `json:"institution_id"`
	InstitutionName              string    `json:"institution_name"`
	IntellectualObjectID         int64     `json:"intellectual_object_id"`
	IntellectualObjectIdentifier string    `json:"intellectual_object_identifier"`
	Object                       string    `json:"object"`
	Outcome                      string    `json:"outcome"`
	OutcomeDetail                string    `json:"outcome_detail"`
	OutcomeInformation           string    `json:"outcome_information"`
}

// PremisEventViewByID returns the event with the specified id.
// Returns pg.ErrNoRows if there is no match.
func PremisEventViewByID(id int64) (*PremisEventView, error) {
	query := NewQuery().Where("id", "=", id)
	return PremisEventViewGet(query)
}

// PremisEventViewByIdentifier returns the event with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func PremisEventViewByIdentifier(identifier string) (*PremisEventView, error) {
	query := NewQuery().Where("identifier", "=", identifier)
	return PremisEventViewGet(query)
}

// PremisEventViewGet returns the first event matching the query.
func PremisEventViewGet(query *Query) (*PremisEventView, error) {
	var event PremisEventView
	err := query.Select(&event)
	return &event, err
}

// PremisEventViewSelect returns all events matching the query.
func PremisEventViewSelect(query *Query) ([]*PremisEventView, error) {
	var events []*PremisEventView
	err := query.Select(&events)
	return events, err
}
