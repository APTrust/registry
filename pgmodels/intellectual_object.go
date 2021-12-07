package pgmodels

import (
	"context"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/stew/slice"
)

// Filters are defined in IntellectualObjectView, since we query on the view.

type IntellectualObject struct {
	ID                        int64     `json:"id"`
	Title                     string    `json:"title"`
	Description               string    `json:"description"`
	Identifier                string    `json:"identifier"`
	AltIdentifier             string    `json:"alt_identifier"`
	Access                    string    `json:"access"`
	BagName                   string    `json:"bag_name"`
	InstitutionID             int64     `json:"institution_id"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	State                     string    `json:"state"`
	ETag                      string    `json:"etag" pg:"etag"`
	BagGroupIdentifier        string    `json:"bag_group_identifier"`
	StorageOption             string    `json:"storage_option"`
	BagItProfileIdentifier    string    `json:"bagit_profile_identifier" pg:"bagit_profile_identifier"`
	SourceOrganization        string    `json:"source_organization"`
	InternalSenderIdentifier  string    `json:"internal_sender_identifier"`
	InternalSenderDescription string    `json:"internal_sender_description"`

	Institution  *Institution   `json:"institution" pg:"rel:has-one"`
	GenericFiles []*GenericFile `json:"generic_files" pg:"rel:has-many"`
	PremisEvents []*PremisEvent `json:"premis_events" pg:"rel:has-many"`
}

// IntellectualObjectByID returns the object with the specified id.
// Returns pg.ErrNoRows if there is no match.
func IntellectualObjectByID(id int64) (*IntellectualObject, error) {
	query := NewQuery().Where(`"intellectual_object"."id"`, "=", id)
	return IntellectualObjectGet(query)
}

// IntellectualObjectByIdentifier returns the object with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func IntellectualObjectByIdentifier(identifier string) (*IntellectualObject, error) {
	query := NewQuery().Where(`"intellectual_object"."identifier"`, "=", identifier)
	return IntellectualObjectGet(query)
}

// IdForFileIdentifier returns the ID of the IntellectualObject
// having the specified identifier.
func IdForObjIdentifier(identifier string) (int64, error) {
	query := NewQuery().Columns("id").Where(`"intellectual_object"."identifier"`, "=", identifier)
	var object IntellectualObject
	err := query.Select(&object)
	return object.ID, err
}

// IntellectualObjectGet returns the first object matching the query.
func IntellectualObjectGet(query *Query) (*IntellectualObject, error) {
	var object IntellectualObject
	err := query.Relations("Institution").Select(&object)
	return &object, err
}

// IntellectualObjectSelect returns all objects matching the query.
func IntellectualObjectSelect(query *Query) ([]*IntellectualObject, error) {
	var objects []*IntellectualObject
	err := query.Select(&objects)
	return objects, err
}

func (obj *IntellectualObject) GetID() int64 {
	return obj.ID
}

// Save saves this object to the database. This will peform an insert
// if IntellectualObject.ID is zero. Otherwise, it updates.
func (obj *IntellectualObject) Save() error {
	if obj.ID == int64(0) {
		return insert(obj)
	}
	return update(obj)
}

// IsGlacierOnly returns true if this object is stored only
// in Glacier.
func (obj *IntellectualObject) IsGlacierOnly() bool {
	return isGlacierOnly(obj.StorageOption)
}

// Delete soft-deletes this object by setting State to 'D' and
// the DeletedAt timestamp to now. You can undo this with Undelete.
func (obj *IntellectualObject) Delete() error {
	hasFiles, err := obj.HasActiveFiles()
	if err != nil {
		return err
	} else if hasFiles {
		return common.ErrActiveFiles
	}
	obj.State = constants.StateDeleted
	obj.UpdatedAt = time.Now().UTC()

	// TODO: Create PremisEvents, update WorkItem. Here? Or Elsewhere?
	// TODO: Wrap obj deletion, event creation into single transaction.
	//       Event will have to get deletion details from WorkItem.

	return update(obj)
}

// HasActiveFiles returns true if this object has any active (non-deleted)
// files. We need to check this before marking an object as deleted.
// Do not mark deleted until all files have been marked deleted.
func (obj *IntellectualObject) HasActiveFiles() (bool, error) {
	db := common.Context().DB
	return db.Model((*GenericFile)(nil)).Where("intellectual_object_id = ? and state = ?", obj.ID, constants.StateActive).Exists()
}

// LastIngestEvent returns the latest ingest event for this object.
// This should never be nil.
func (obj *IntellectualObject) LastIngestEvent() (*PremisEvent, error) {
	return obj.lastEvent(constants.EventIngestion)
}

// LastDeleationEvent returns the latest deletion event for this object,
// which may be nil.
func (obj *IntellectualObject) LastDeletionEvent() (*PremisEvent, error) {
	return obj.lastEvent(constants.EventDeletion)
}

func (obj *IntellectualObject) lastEvent(eventType string) (*PremisEvent, error) {
	query := NewQuery().
		Where("intellectual_object_id", "=", obj.ID).
		Where("event_type", "=", eventType).
		IsNull("generic_file_id").
		OrderBy("created_at", "desc").
		Offset(0).
		Limit(1)
	return PremisEventGet(query)
}

func isGlacierOnly(storageOption string) bool {
	return slice.Contains(constants.GlacierOnlyOptions, storageOption)
}

func (obj *IntellectualObject) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(obj.Title) {
		errors["Title"] = "Title is required"
	}
	if common.IsEmptyString(obj.Identifier) {
		errors["Identifier"] = "Identifier is required"
	}
	if !v.IsIn(obj.State, constants.States...) {
		errors["State"] = ErrInstState
	}
	if !v.IsIn(obj.Access, constants.AccessSettings...) {
		errors["Access"] = "Invalid access value"
	}
	if obj.InstitutionID < 1 {
		errors["InstitutionID"] = "Invalid institution id"
	}
	if !v.IsIn(obj.StorageOption, constants.StorageOptions...) {
		errors["StorageOption"] = "Invalid storage option"
	}
	if common.IsEmptyString(obj.BagItProfileIdentifier) {
		errors["BagItProfileIdentifier"] = "BagItProfileIdentifier is required"
	}
	if common.IsEmptyString(obj.SourceOrganization) {
		errors["SourceOrganization"] = "SourceOrganization is required"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}

func (obj *IntellectualObject) ValidateChanges(updatedObj *IntellectualObject) error {
	if obj.ID != updatedObj.ID {
		return common.ErrIDMismatch
	}
	if obj.InstitutionID != updatedObj.InstitutionID {
		return common.ErrInstIDChange
	}
	if obj.Identifier != updatedObj.Identifier {
		return common.ErrIdentifierChange
	}
	// Caller should force storage option of updated object to
	// match existing object before calling this validation function.
	if obj.State == constants.StateActive && obj.StorageOption != updatedObj.StorageOption {
		return common.ErrStorageOptionChange
	}
	return nil
}

// The following statements have no effect other than to force a compile-time
// check that ensures our IntellectualObject model properly implements these hook
// interfaces.
var (
	_ pg.BeforeInsertHook = (*IntellectualObject)(nil)
	_ pg.BeforeUpdateHook = (*IntellectualObject)(nil)
)

// BeforeInsert sets timestamps and bucket names on creation.
func (obj *IntellectualObject) BeforeInsert(c context.Context) (context.Context, error) {
	now := time.Now().UTC()
	obj.CreatedAt = now
	obj.UpdatedAt = now

	err := obj.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (obj *IntellectualObject) BeforeUpdate(c context.Context) (context.Context, error) {
	err := obj.Validate()
	obj.UpdatedAt = time.Now().UTC()
	if err == nil {
		return c, nil
	}
	return c, err
}

func (obj *IntellectualObject) DeletionItem() (*WorkItem, error) {
	query := NewQuery().
		Where("intellectual_object_id", "=", obj.ID).
		IsNull("generic_file_id").
		Where("action", "=", constants.ActionDelete).
		Where("stage", "=", constants.StatusStarted).
		OrderBy("updated_at", "desc").
		Limit(1)
	return WorkItemGet(query)
}
