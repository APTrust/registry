package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
)

type PremisEvent struct {
	ID                   int64     `json:"id" form:"id" pg:"id"`
	Identifier           string    `json:"identifier" form:"identifier" pg:"identifier"`
	EventType            string    `json:"event_type" form:"event_type" pg:"event_type"`
	DateTime             time.Time `json:"date_time" form:"date_time" pg:"date_time"`
	OutcomeDetail        string    `json:"outcome_detail" form:"outcome_detail" pg:"outcome_detail"`
	Detail               string    `json:"detail" form:"detail" pg:"detail"`
	OutcomeInformation   string    `json:"outcome_information" form:"outcome_information" pg:"outcome_information"`
	Object               string    `json:"object" form:"object" pg:"object"`
	Agent                string    `json:"agent" form:"agent" pg:"agent"`
	IntellectualObjectID int64     `json:"intellectual_object_id" form:"intellectual_object_id" pg:"intellectual_object_id"`
	GenericFileID        int64     `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	CreatedAt            time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	Outcome              string    `json:"outcome" form:"outcome" pg:"outcome"`
	InstitutionID        int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	OldUUID              string    `json:"old_uuid" form:"old_uuid" pg:"old_uuid"`
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
