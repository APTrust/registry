package pgmodels

import "time"

type GenericFileView struct {
	tableName             struct{}  `pg:"generic_files_view"`
	ID                    int64     `json:"id" pg:"id"`
	FileFormat            string    `json:"file_format"`
	Size                  int64     `json:"size"`
	Identifier            string    `json:"identifier"`
	IntellectualObjectID  int64     `json:"intellectual_object_id"`
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

	// This is a late hack. Needed for object restoration.
	// Note that there's no join the DB for this (and there shouldn't
	// be, because it would produce multiple records). The
	// GenericFilesController will hack in the storage records
	// only when specifically requested.
	StorageRecords []*StorageRecord `pg:"-" json:"storage_records,omitempty"`
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
	var files []*GenericFileView
	err := query.Select(&files)
	return files, err
}

// GenericFileViewGet returns the first user view record matching the query.
func GenericFileViewGet(query *Query) (*GenericFileView, error) {
	var gf GenericFileView
	err := query.Select(&gf)
	return &gf, err
}
