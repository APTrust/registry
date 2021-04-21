package pgmodels

import (
	"time"
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
	"institution_parent_id",
	"file_count__gteq",
	"file_count__lteq",
	"size__gteq",
	"size__lteq",
	"created_at__gteq",
	"crated_at__lteq",
	"updated_at__gteq",
	"updated_at__lteq",
}

type IntellectualObjectView struct {
	tableName                 struct{}  `pg:"intellectual_objects_view"`
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
	InstitutionName           string    `json:"institution_name"`
	InstitutionIdentifier     string    `json:"institution_identifier"`
	InstitutionType           string    `json:"institution_type"`
	InstitutionParentID       int64     `json:"institution_parent_id"`
	FileCount                 int64     `json:"file_count"`
	Size                      int64     `json:"size"`
}

// IntellectualObjectViewByID returns the object with the specified id.
// Returns pg.ErrNoRows if there is no match.
func IntellectualObjectViewByID(id int64) (*IntellectualObjectView, error) {
	query := NewQuery().Where(`"intellectual_object"."id"`, "=", id)
	return IntellectualObjectViewGet(query)
}

// IntellectualObjectViewByIdentifier returns the object with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func IntellectualObjectViewByIdentifier(identifier string) (*IntellectualObjectView, error) {
	query := NewQuery().Where(`"intellectual_object"."identifier"`, "=", identifier)
	return IntellectualObjectViewGet(query)
}

// IntellectualObjectViewGet returns the first object matching the query.
func IntellectualObjectViewGet(query *Query) (*IntellectualObjectView, error) {
	var object IntellectualObjectView
	err := query.Relations("Institution").Select(&object)
	return &object, err
}

// IntellectualObjectViewSelect returns all objects matching the query.
func IntellectualObjectViewSelect(query *Query) ([]*IntellectualObjectView, error) {
	var objects []*IntellectualObjectView
	err := query.Select(&objects)
	return objects, err
}

func (obj *IntellectualObjectView) GetID() int64 {
	return obj.ID
}