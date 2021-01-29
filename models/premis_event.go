package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

func (event *PremisEvent) GetID() int64 {
	return event.ID
}

func (event *PremisEvent) Authorize(actingUser *User, action string) error {
	perm := "Event" + action
	if !actingUser.HasPermission(constants.Permission(perm), event.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s event %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, event.ID, event.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

// DeleteIsForbidden returns true because PremisEvents are our audit trail.
func (event *PremisEvent) DeleteIsForbidden() bool {
	return true
}

// UpdateIsForbidden returns true because PremisEvents are our audit trail.
func (event *PremisEvent) UpdateIsForbidden() bool {
	return true
}

func (event *PremisEvent) IsReadOnly() bool {
	return false
}

func (event *PremisEvent) SupportsSoftDelete() bool {
	return false
}

func (event *PremisEvent) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (event *PremisEvent) ClearSoftDeleteAttributes() {
	// No-op
}

func (event *PremisEvent) SetTimestamps() {
	now := time.Now().UTC()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	event.UpdatedAt = now
}

func (event *PremisEvent) BeforeSave() error {
	// TODO: Validate
	return nil
}

func PremisEventFind(id int64) (*PremisEvent, error) {
	ctx := common.Context()
	event := &PremisEvent{ID: id}
	err := ctx.DB.Model(event).WherePK().Select()
	return event, err
}

func PremisEventFindByIdentifier(identifier string) (*PremisEvent, error) {
	ctx := common.Context()
	event := &PremisEvent{}
	err := ctx.DB.Model(event).Where(`"identifier" = ?`, identifier).First()
	return event, err
}
