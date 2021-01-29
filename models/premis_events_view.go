package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type PremisEventsView struct {
	ID                           int64     `json:"id" pg:"id"`
	Identifier                   string    `json:"identifier" pg:"identifier"`
	InstitutionID                int64     `json:"institution_id" pg:"institution_id"`
	InstitutionName              string    `json:"institution_name" pg:"institution_name"`
	IntellectualObjectID         int64     `json:"intellectual_object_id" pg:"intellectual_object_id"`
	IntellectualObjectIdentifier string    `json:"intellectual_object_identifier" pg:"intellectual_object_identifier"`
	GenericFileID                int64     `json:"generic_file_id" pg:"generic_file_id"`
	GenericFileIdentifier        string    `json:"generic_file_identifier" pg:"generic_file_identifier"`
	EventType                    string    `json:"event_type" pg:"event_type"`
	DateTime                     time.Time `json:"date_time" pg:"date_time"`
	Detail                       string    `json:"detail" pg:"detail"`
	Outcome                      string    `json:"outcome" pg:"outcome"`
	OutcomeDetail                string    `json:"outcome_detail" pg:"outcome_detail"`
	OutcomeInformation           string    `json:"outcome_information" pg:"outcome_information"`
	Object                       string    `json:"object" pg:"object"`
	Agent                        string    `json:"agent" pg:"agent"`
	CreatedAt                    time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" pg:"updated_at"`
	OldUUID                      string    `json:"old_uuid" pg:"old_uuid"`
}

func (event *PremisEventsView) GetID() int64 {
	return event.ID
}

func (event *PremisEventsView) Authorize(actingUser *User, action string) error {
	perm := "Event" + action
	if !actingUser.HasPermission(constants.Permission(perm), event.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s event %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, event.ID, event.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

// DeleteIsForbidden returns true because PremisEventsViews are our audit trail.
func (event *PremisEventsView) DeleteIsForbidden() bool {
	return true
}

// UpdateIsForbidden returns true because PremisEventsViews are our audit trail.
func (event *PremisEventsView) UpdateIsForbidden() bool {
	return true
}

// IsReadOnly is true because this is a view.
func (event *PremisEventsView) IsReadOnly() bool {
	return true
}

func (event *PremisEventsView) SupportsSoftDelete() bool {
	return false
}

func (event *PremisEventsView) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (event *PremisEventsView) ClearSoftDeleteAttributes() {
	// No-op
}

func (event *PremisEventsView) SetTimestamps() {
	// No-op. Can't insert or update view.
}

func (event *PremisEventsView) BeforeSave() error {
	// No-op
	return nil
}

func PremisEventsViewFind(id int64) (*PremisEventsView, error) {
	ctx := common.Context()
	event := &PremisEventsView{ID: id}
	err := ctx.DB.Model(event).WherePK().Select()
	return event, err
}

func PremisEventsViewFindByIdentifier(identifier string) (*PremisEventsView, error) {
	ctx := common.Context()
	event := &PremisEventsView{}
	err := ctx.DB.Model(event).Where(`"identifier" = ?`, identifier).First()
	return event, err
}
