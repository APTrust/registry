package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

// IntellectualObjectFilters describes the allowed filters for searching
// IntellectualObjects. Some of these are commented out because, while they
// are supported in Pharo's v2 member API, we want to phase them out.
// "Contains" queries, in particular, which use SQL "like" on the backend,
// cause performance problems.
var IntellectualObjectFilters = []string{
	"access",
	"alt_identifier",
	"bag_group_identifier",
	"bag_name",
	"bagit_profile_identifier",
	"created_at__lteq",
	"created_at__gteq",
	"etag",
	"file_count__gteq",
	"file_count__lteq",
	"identifier",
	"institution_id",
	"institution_parent_id",
	"internal_sender_description",
	"internal_sender_identifier",
	"size__gteq",
	"size__lteq",
	"source_organization",
	"state",
	"storage_option",
	// "title",
	// "title__contains",
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
	PayloadFileCount          int64     `json:"payload_file_count"`
	PayloadSize               int64     `json:"payload_size"`
}

// IntellectualObjectViewByID returns the object with the specified id.
// Returns pg.ErrNoRows if there is no match.
func IntellectualObjectViewByID(id int64) (*IntellectualObjectView, error) {
	query := NewQuery().Where(`"intellectual_object_view"."id"`, "=", id)
	return IntellectualObjectViewGet(query)
}

// IntellectualObjectViewByIdentifier returns the object with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func IntellectualObjectViewByIdentifier(identifier string) (*IntellectualObjectView, error) {
	query := NewQuery().Where(`"intellectual_object_view"."identifier"`, "=", identifier)
	return IntellectualObjectViewGet(query)
}

// IntellectualObjectViewGet returns the first object matching the query.
func IntellectualObjectViewGet(query *Query) (*IntellectualObjectView, error) {
	var object IntellectualObjectView
	err := query.Select(&object)
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

// SmallestObjectNotRestoredInXDays returns the smallest intellectual object
// belonging to an institutition that has not been restored in at least X days.
// We use this for restoration spot tests because 1) we don't want to restore
// an item the depositor has recently restored, and 2) we don't want to restore
// massive objects (hundreds of GB) if we can help it.
//
// Param institutionID is the ID of the depositing institution. minSize is the
// minimum size of the object to restore. This should generally be around 1-20 KB.
// days is the number of days since last restoration. This should be 365 or more
// for spot tests, perferably 730 or more.
func SmallestObjectNotRestoredInXDays(institutionID, minSize int64, days int) (*IntellectualObjectView, error) {
	query := `select obj.id 
	from intellectual_objects_view obj
	where obj.institution_id = ?
	and obj."size" >= ?
	and not exists (
		select 1 
		from work_items 
		where intellectual_object_id = obj.id
		and action = 'Restore Object' 
		and status='Success' 
		and updated_at > current_date - ?
	)
	order by obj."size" 
	limit 1
	`
	var objID int64
	_, err := common.Context().DB.Model((*IntellectualObjectView)(nil)).QueryOne(pg.Scan(&objID), query, institutionID, minSize, days)
	if err != nil {
		return nil, err
	}
	return IntellectualObjectViewByID(objID)
}
