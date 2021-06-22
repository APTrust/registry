package pgmodels

import (
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

type GenericFile struct {
	ID                   int64     `json:"id" pg:"id"`
	FileFormat           string    `json:"file_format"`
	Size                 int64     `json:"size"`
	Identifier           string    `json:"identifier"`
	IntellectualObjectID int64     `json:"intellectual_object_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	State                string    `json:"state"`
	LastFixityCheck      time.Time `json:"last_fixity_check"`
	InstitutionID        int64     `json:"institution_id"`
	StorageOption        string    `json:"storage_option"`
	UUID                 string    `json:"uuid" pg:"uuid"`

	Institution        *Institution        `json:"-" pg:"rel:has-one"`
	IntellectualObject *IntellectualObject `json:"-" pg:"rel:has-one"`
	PremisEvents       []*PremisEvent      `json:"premis_events" pg:"rel:has-many"`
	Checksums          []*Checksum         `json:"checksums" pg:"rel:has-many"`
	StorageRecords     []*StorageRecord    `json:"storage_records" pg:"rel:has-many"`
}

// TODO: When selecting relations, order by UpdatedAt asc.

// GenericFileByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func GenericFileByID(id int64) (*GenericFile, error) {
	query := NewQuery().Where(`"generic_file"."id"`, "=", id).Relations("Institution", "IntellectualObject", "PremisEvents", "Checksums", "StorageRecords")
	return GenericFileGet(query)
}

// GenericFileByIdentifier returns the file with the specified
// identifier. Returns pg.ErrNoRows if there is no match.
func GenericFileByIdentifier(identifier string) (*GenericFile, error) {
	query := NewQuery().Where(`"generic_file"."identifier"`, "=", identifier)
	return GenericFileGet(query)
}

// GenericFileGet returns the first file matching the query.
func GenericFileGet(query *Query) (*GenericFile, error) {
	var gf GenericFile
	err := query.Select(&gf)
	return &gf, err
}

// GenericFileSelect returns all files matching the query.
func GenericFileSelect(query *Query) ([]*GenericFile, error) {
	var files []*GenericFile
	err := query.Select(&files)
	return files, err
}

func (gf *GenericFile) GetID() int64 {
	return gf.ID
}

// Save saves this file to the database. This will peform an insert
// if GenericFile.ID is zero. Otherwise, it updates.
func (gf *GenericFile) Save() error {
	if gf.ID == int64(0) {
		return insert(gf)
	}
	return update(gf)
}

// IsGlacierOnly returns true if this file is stored only
// in Glacier.
func (gf *GenericFile) IsGlacierOnly() bool {
	return isGlacierOnly(gf.StorageOption)
}

// ObjectFileCount returns the number of active files with the specified
// Intellectial Object ID.
func ObjectFileCount(objID int64, filter, state string) (int, error) {
	if filter != "" {
		db := common.Context().DB
		likeFilter := fmt.Sprintf("%%%s%%", filter)
		type Result struct {
			Count int
		}
		var result Result
		idQuery := `select count(distinct(gf.id)) from generic_files gf
		left join checksums cs on cs.generic_file_id = gf.id
		where gf.intellectual_object_id = ? and state = ?
		and (gf.identifier like ? or cs.digest = ?)`
		_, err := db.QueryOne(&result, idQuery, objID, state, likeFilter, filter)
		return result.Count, err
	}
	return common.Context().DB.Model((*GenericFile)(nil)).Where(`intellectual_object_id = ? and state = ?`, objID, state).Count()
}

// Object files returns files belonging to an intellectual object.
// This function is used to filter files on the IntellectualObjectShow page.
// objID is the id of the object. filter is an optional file identifier or
// checksum. offset and limit are for paging. state is "A" for active (default)
// or "D" for deleted.
func ObjectFiles(objID int64, filter, state string, offset, limit int) ([]*GenericFile, error) {
	db := common.Context().DB
	var err error
	var files []*GenericFile

	// If we have a string filter, try to match on partial file name
	// or exact checksum value. First get the ids, then return the whole
	// records. Sorting here gets tricky because the sort column must
	// be included in the query while we're using distinct. To sort, we'd
	// need to create a custom struct with id and sort column fields.
	if filter != "" {
		likeFilter := fmt.Sprintf("%%%s%%", filter)
		var fileIds []*int
		idQuery := `select distinct(gf.id) from generic_files gf
		left join checksums cs on cs.generic_file_id = gf.id
		where gf.intellectual_object_id = ? and state = ?
		and (gf.identifier like ? or cs.digest = ?)
		order by gf.id offset ? limit ?`
		_, err = db.Query(&fileIds, idQuery, objID, state, likeFilter, filter, offset, limit)
		if err != nil {
			return nil, err
		}
		if len(fileIds) == 0 {
			return files, err
		}
		err = db.Model(&files).Where("id in (?)", pg.In(fileIds)).Relation("StorageRecords").
			Relation("Checksums", func(q *pg.Query) (*pg.Query, error) {
				return q.Order("datetime desc").Order("algorithm asc"), nil
			}).Select()
		if err != nil {
			return nil, err
		}
	} else {
		// If filter is empty, this query is much simpler.
		// Just get all active files.
		err = db.Model(&files).Where("intellectual_object_id = ? and state = ?", objID, state).
			Relation("StorageRecords").
			Relation("Checksums", func(q *pg.Query) (*pg.Query, error) {
				return q.Order("datetime desc").Order("algorithm asc"), nil
			}).Limit(limit).Offset(offset).Select()
		if err != nil {
			return nil, err
		}
	}
	return files, err
}
