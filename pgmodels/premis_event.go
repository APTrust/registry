package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
)

type PremisEvent struct {
	ID                   int64     `json:"id"`
	Agent                string    `json:"agent"`
	CreatedAt            time.Time `json:"created_at"`
	DateTime             time.Time `json:"date_time"`
	Detail               string    `json:"detail"`
	EventType            string    `json:"event_type"`
	GenericFileID        int64     `json:"generic_file_id"`
	Identifier           string    `json:"identifier"`
	InstitutionID        int64     `json:"institution_id"`
	IntellectualObjectID int64     `json:"intellectual_object_id"`
	Object               string    `json:"object"`
	OldUUID              string    `json:"old_uuid"`
	Outcome              string    `json:"outcome"`
	OutcomeDetail        string    `json:"outcome_detail"`
	OutcomeInformation   string    `json:"outcome_information"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// PremisEventByID returns the event with the specified id.
// Returns pg.ErrNoRows if there is no match.
func PremisEventByID(id int64) (*PremisEvent, error) {
	query := NewQuery().Where(`"premis_event"."id"`, "=", id)
	return PremisEventGet(query)
}

// PremisEventByIdentifier returns the event with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func PremisEventByIdentifier(identifier string) (*PremisEvent, error) {
	query := NewQuery().Where(`"premis_event"."identifier"`, "=", identifier)
	return PremisEventGet(query)
}

// PremisEventGet returns the first event matching the query.
func PremisEventGet(query *Query) (*PremisEvent, error) {
	var event PremisEvent
	err := query.Select(&event)
	return &event, err
}

// PremisEventSelect returns all events matching the query.
func PremisEventSelect(query *Query) ([]*PremisEvent, error) {
	var events []*PremisEvent
	err := query.Select(&events)
	return events, err
}

func (event *PremisEvent) GetID() int64 {
	return event.ID
}

// Save saves this event to the database. This will peform an insert
// if PremisEvent.ID is zero. Otherwise, it updates.
func (event *PremisEvent) Save() error {
	if event.ID == int64(0) {
		return insert(event)
	}
	return update(event)
}

// ObjectEventCount returns the number of object-level PremisEvents for the
// specified IntellectualObject. Note that queries on the premis_events table
// can be potentially quite expensive, since the table has well over 100M rows.
//
// Object-level events include ingest, identifier assignment, access assignment
// and others. They exclude all events related to specific files (such as fixity
// check, etc.).
func ObjectEventCount(intellectualObjectID int64) (int, error) {
	return common.Context().DB.Model((*PremisEvent)(nil)).Where(`intellectual_object_id = ? and generic_file_id is null`, intellectualObjectID).Count()
}
