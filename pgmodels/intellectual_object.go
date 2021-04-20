package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

var IntellectualObjectFilters = []string{
	"title",
	"title__contains",
	"description",
	"description__contains",
	"identifier",
	"alt_identifier",
	"access",
	"bag_name",
	"institution_id",
	"state",
	"etag",
	"bag_group_identifier",
	"storage_option",
	"bagit_profile_identifier",
	"source_organization",
	"internal_sender_identifier",
	"internal_sender_description",
}

type IntellectualObject struct {
	ID                        int64     `json:"id" form:"id" pg:"id"`
	Title                     string    `json:"title" form:"title" pg:"title"`
	Description               string    `json:"description" form:"description" pg:"description"`
	Identifier                string    `json:"identifier" form:"identifier" pg:"identifier"`
	AltIdentifier             string    `json:"alt_identifier" form:"alt_identifier" pg:"alt_identifier"`
	Access                    string    `json:"access" form:"access" pg:"access"`
	BagName                   string    `json:"bag_name" form:"bag_name" pg:"bag_name"`
	InstitutionID             int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	CreatedAt                 time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	State                     string    `json:"state" form:"state" pg:"state"`
	ETag                      string    `json:"etag" form:"etag" pg:"etag"`
	BagGroupIdentifier        string    `json:"bag_group_identifier" form:"bag_group_identifier" pg:"bag_group_identifier"`
	StorageOption             string    `json:"storage_option" form:"storage_option" pg:"storage_option"`
	BagItProfileIdentifier    string    `json:"bagit_profile_identifier" form:"bag_it_profile_identifier" pg:"bagit_profile_identifier"`
	SourceOrganization        string    `json:"source_organization" form:"source_organization" pg:"source_organization"`
	InternalSenderIdentifier  string    `json:"internal_sender_identifier" form:"internal_sender_identifier" pg:"internal_sender_identifier"`
	InternalSenderDescription string    `json:"internal_sender_description" form:"internal_sender_description" pg:"internal_sender_description"`

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

// Delete soft-deletes this object by setting State to 'D' and
// the DeletedAt timestamp to now. You can undo this with Undelete.
func (obj *IntellectualObject) Delete() error {
	obj.State = constants.StateDeleted
	obj.UpdatedAt = time.Now().UTC()

	// TODO: Create PremisEvents, update WorkItem

	return update(obj)
}
