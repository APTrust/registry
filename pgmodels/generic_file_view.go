package pgmodels

import "time"

type GenericFileView struct {
	tableName             struct{}  `pg:"generic_files_view"`
	ID                    int64     `json:"id" pg:"id"`
	FileFormat            string    `json:"file_format"`
	Size                  int64     `json:"size"`
	Identifier            string    `json:"identifier"`
	IntellectualObjectID  int64     `json:"intellection_object_id"`
	ObjectIdentifier      string    `json:"object_identifier"`
	Access                string    `json:"access"`
	State                 string    `json:"state"`
	LastFixityCheck       time.Time `json:"last_fixity_check"`
	InstitutionID         int64     `json:"institution_id"`
	InstitutionName       string    `json:"institution_name"`
	InstitutionIdentifier string    `json:"institution_identifier"`
	StorageOption         string    `json:"storage_option"`
	UUID                  string    `json:"uuid" pg:"uuid"`
	Md5                   string    `json:"md5"`
	Sha1                  string    `json:"sha1"`
	Sha256                string    `json:"sha256"`
	Sha512                string    `json:"sha512"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// GenericFileViewByID returns the GenericFileView record
// with the specified id.  Returns pg.ErrNoRows if there is no match.
func GenericFileViewByID(id int64) (*GenericFileView, error) {
	query := NewQuery().Where("id", "=", id)
	return GenericFileViewGet(query)
}

// GenericFileViewByIdentifier returns the GenericFileView record with the
// specified email address. Returns pg.ErrNoRows if there is no match.
func GenericFileViewByIdentifier(identifier string) (*GenericFileView, error) {
	query := NewQuery().Where("identifier", "=", identifier)
	return GenericFileViewGet(query)
}

// GenericFileViewSelect returns all GenericFileView records matching
// the query.
func GenericFileViewSelect(query *Query) ([]*GenericFileView, error) {
	var requests []*GenericFileView
	err := query.Select(&requests)
	return requests, err
}

// GenericFileViewGet returns the first user view record matching the query.
func GenericFileViewGet(query *Query) (*GenericFileView, error) {
	var request GenericFileView
	err := query.Select(&request)
	return &request, err
}
