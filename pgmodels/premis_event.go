package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/stretchr/stew/slice"
)

type PremisEvent struct {
	TimestampModel
	Agent                int       `json:"agent"`
	DateTime             time.Time `json:"date_time"`
	Detail               string    `json:"detail"`
	EventType            int       `json:"event_type"`
	GenericFileID        int64     `json:"generic_file_id"`
	Identifier           string    `json:"identifier"`
	InstitutionID        int64     `json:"institution_id"`
	IntellectualObjectID int64     `json:"intellectual_object_id"`
	Object               int       `json:"object"`
	OldUUID              string    `json:"old_uuid"`
	Outcome              string    `json:"outcome"`
	OutcomeDetail        string    `json:"outcome_detail"`
	OutcomeInformation   string    `json:"outcome_information"`
}

type PremisEventType struct {
	EventTypeID int    `json:"event_type_id"`
	EventType   string `json:"event_type"`
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

// IdForEventIdentifier returns the ID of the PremisEvent having
// the specified identifier. These identifiers are UUID strings.
func IdForEventIdentifier(identifier string) (int64, error) {
	query := NewQuery().Columns("id").Where(`"premis_event"."identifier"`, "=", identifier)
	var event PremisEvent
	err := query.Select(&event)
	return event.ID, err
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

// Save saves this event to the database. This will peform an insert
// if PremisEvent.ID is zero. Otherwise, it updates.
func (event *PremisEvent) Save() error {
	if event.ID == int64(0) {
		event.SetTimestamps()
		return insert(event)
	}
	// Premis events cannot be updated
	return common.ErrNotSupported
}

// Validate returns errors if this event isn't valid.
func (event *PremisEvent) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if event.Agent < 0 {
		errors["Agent"] = "Event requires a valid Agent"
	}
	if event.DateTime.IsZero() {
		errors["DateTime"] = "Event DateTime is required"
	}
	if common.IsEmptyString(event.Detail) {
		errors["Detail"] = "Event Detail cannot be empty"
	}
	if !slice.Contains(constants.EventTypes, event.EventType) {
		errors["EventType"] = "Event requires a valid EventType"
	}
	// GenericFileID is optional because not all events pertain
	// to files.
	if !common.LooksLikeUUID(event.Identifier) {
		errors["Identifier"] = "Event identifier should be a UUID"
	}
	if event.InstitutionID <= 0 {
		errors["InstitutionID"] = "Event requires a valid institution id"
	}
	if event.IntellectualObjectID <= 0 {
		errors["IntellectualObjectID"] = "Event requires a valid intellectual object id"
	}
	if event.Object < 0 {
		errors["Object"] = "Event requires a valid Object"
	}
	if !slice.Contains(constants.EventOutcomes, event.Outcome) {
		errors["Outcome"] = "Event requires a valid Outcome value"
	}
	if common.IsEmptyString(event.OutcomeDetail) {
		errors["OutcomeDetail"] = "Event OutcomeDetail cannot be empty"
	}
	if common.IsEmptyString(event.OutcomeInformation) {
		errors["OutcomeInformation"] = "Event OutcomeInformation cannot be empty"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
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

// Provides the full string description of an event type given its int code.
func LookupEventType(eventTypeID int) (string, error) {
	query := NewQuery().Columns("id").Where(`"lookup_event_type"."eventType"`, "=", eventTypeID)
	var premisEventType PremisEventType
	err := query.Select(&premisEventType)
	return premisEventType.EventType, err
}
