package pgmodels

import (
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	v "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

var GenericFileFilters = []string{
	"identifier",
	"uuid",
	"intellectual_object_id",
	"institution_id",
	"state",
	"storage_option",
	"size__gteq",
	"size__lteq",
	"created_at__gteq",
	"created_at__lteq",
	"updated_at__gteq",
	"updated_at__lteq",
	"last_fixity_check__gteq",
	"last_fixity_check__lteq",
}

type GenericFile struct {
	TimestampModel
	FileFormat           string              `json:"file_format"`
	Size                 int64               `json:"size"`
	Identifier           string              `json:"identifier"`
	IntellectualObjectID int64               `json:"intellectual_object_id"`
	State                string              `json:"state"`
	LastFixityCheck      time.Time           `json:"last_fixity_check"`
	InstitutionID        int64               `json:"institution_id"`
	StorageOption        string              `json:"storage_option"`
	UUID                 string              `json:"uuid" pg:"uuid"`
	Institution          *Institution        `json:"-" pg:"rel:has-one"`
	IntellectualObject   *IntellectualObject `json:"-" pg:"rel:has-one"`
	PremisEvents         []*PremisEvent      `json:"premis_events,omitempty" pg:"rel:has-many"`
	Checksums            []*Checksum         `json:"checksums,omitempty" pg:"rel:has-many"`
	StorageRecords       []*StorageRecord    `json:"storage_records,omitempty" pg:"rel:has-many"`
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

// IdForFileIdentifier returns the ID of the GenericFile having the
// specified identifier.
func IdForFileIdentifier(identifier string) (int64, error) {
	query := NewQuery().Columns("id").Where(`"generic_file"."identifier"`, "=", identifier)
	var gf GenericFile
	err := query.Select(&gf)
	return gf.ID, err
}

// GenericFileGet returns the first file matching the query.
func GenericFileGet(query *Query) (*GenericFile, error) {
	var gf GenericFile
	err := query.Select(&gf)
	if gf.ID == 0 {
		return nil, err
	}
	return &gf, err
}

// GenericFileSelect returns all files matching the query.
func GenericFileSelect(query *Query) ([]*GenericFile, error) {
	var files []*GenericFile
	err := query.Select(&files)
	return files, err
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

func (gf *GenericFile) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(gf.FileFormat) {
		errors["FileFormat"] = "FileFormat is required"
	}
	if common.IsEmptyString(gf.Identifier) {
		errors["Identifier"] = "Identifier is required"
	}
	if !v.IsIn(gf.State, constants.States...) {
		errors["State"] = ErrInstState
	}
	if gf.Size < 0 {
		errors["Size"] = "Size cannot be negative"
	}
	if gf.InstitutionID < 1 {
		errors["InstitutionID"] = "Invalid institution id"
	}
	if gf.IntellectualObjectID < 1 {
		errors["IntellectualObjectID"] = "Intellectual object ID is required"
	}
	if !v.IsIn(gf.StorageOption, constants.StorageOptions...) {
		errors["StorageOption"] = "Invalid storage option"
	}
	if !common.LooksLikeUUID(gf.UUID) {
		errors["UUID"] = "Valid UUID required"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil

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

// Delete soft-deletes this file by setting State to 'D' and
// the UpdatedAt timestamp to now. You can undo this with Undelete.
// It also creates a deletion PremisEvent. You can't get rid of that.
//
// It is legitimate for a depositor to delete a file, then re-upload
// it later, particularly if they want to change the storage option.
// In that case, the file's state would be set back to "A" after the
// new ingest, and the old deletion event would remain to show that an earlier
// version of the file was once deleted.
//
// We would know the new file is active because state = "A" and it would
// have an ingest event dated after the last deletion event.
func (gf *GenericFile) Delete() error {

	err := gf.AssertDeletionPreconditions()
	if err != nil {
		return err
	}

	gf.State = constants.StateDeleted
	gf.UpdatedAt = time.Now().UTC()

	valErr := gf.Validate()
	if valErr != nil {
		return valErr
	}

	deletionEvent, err := gf.NewDeletionEvent()
	if err != nil {
		return err
	}
	deletionEvent.SetTimestamps()
	valErr = deletionEvent.Validate()
	if valErr != nil {
		return valErr
	}

	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		var err error
		_, err = tx.Model(gf).WherePK().Update()
		if err != nil {
			registryContext.Log.Error().Msgf("GenericFile deletion transaction failed on update of file. File: %d (%s). Error: %v", gf.ID, gf.Identifier, err)
		}
		_, err = tx.Model(deletionEvent).Insert()
		if err != nil {
			registryContext.Log.Error().Msgf("GenericFile deletion transaction failed on insertion of event. File: %d (%s). Error: %v", gf.ID, gf.Identifier, err)
		}
		return err
	})
}

// LastIngestEvent returns the latest ingest event for this file.
// This should never be nil.
func (gf *GenericFile) LastIngestEvent() (*PremisEvent, error) {
	return gf.lastEvent(constants.EventIngestion)
}

// LastDeleationEvent returns the latest deletion event for this file,
// which may be nil.
func (gf *GenericFile) LastDeletionEvent() (*PremisEvent, error) {
	return gf.lastEvent(constants.EventDeletion)
}

func (gf *GenericFile) lastEvent(eventType string) (*PremisEvent, error) {
	query := NewQuery().
		Where("generic_file_id", "=", gf.ID).
		Where("event_type", "=", eventType).
		OrderBy("created_at", "desc").
		Offset(0).
		Limit(1)
	return PremisEventGet(query)
}

func (gf *GenericFile) ActiveDeletionWorkItem() (*WorkItem, error) {
	query := NewQuery().
		Where("generic_file_id", "=", gf.ID).
		Where("action", "=", constants.ActionDelete).
		Where("status", "=", constants.StatusStarted).
		OrderBy("updated_at", "desc").
		Limit(1)
	item, err := WorkItemGet(query)
	if err != nil && err.Error() == pg.ErrNoRows.Error() {
		return nil, nil
	}
	return item, err
}

func (gf *GenericFile) DeletionRequest(workItemID int64) (*DeletionRequestView, error) {
	query := NewQuery().
		Where("work_item_id", "=", workItemID).
		Where("object_count", "=", 0).
		Where("file_count", "=", 1)
	return DeletionRequestViewGet(query)
}

func (gf *GenericFile) AssertDeletionPreconditions() error {
	if gf.State == constants.StateDeleted {
		return fmt.Errorf("File is already in deleted state")
	}
	_, _, err := gf.assertDeletionApproved()
	return err
}

func (gf *GenericFile) assertDeletionApproved() (*WorkItem, *DeletionRequestView, error) {
	workItem, err := gf.ActiveDeletionWorkItem()
	if workItem == nil || IsNoRowError(err) {
		return nil, nil, fmt.Errorf("Missing deletion request work item")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting active deletion work item: %v", err)
	}
	if common.IsEmptyString(workItem.InstApprover) {
		return workItem, nil, fmt.Errorf("Deletion work item is missing institutional approver")
	}
	deletionRequest, err := gf.DeletionRequest(workItem.ID)
	if deletionRequest == nil || IsNoRowError(err) {
		return workItem, nil, fmt.Errorf("No deletion request for work item %d", workItem.ID)
	}
	if err != nil {
		return workItem, nil, fmt.Errorf("Error getting deletion request: %v", err)
	}
	if deletionRequest.RequestedByID == 0 {
		// We should never hit this because RequestedByID has a not-null constraint.
		return workItem, deletionRequest, fmt.Errorf("Deletion request %d has no requestor", deletionRequest.ID)
	}
	if deletionRequest.ConfirmedByID == 0 {
		return workItem, deletionRequest, fmt.Errorf("Deletion request %d has no approver", deletionRequest.ID)
	}
	return workItem, deletionRequest, nil
}

func (gf *GenericFile) NewDeletionEvent() (*PremisEvent, error) {
	_, deletionRequestView, err := gf.assertDeletionApproved()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &PremisEvent{
		Agent:                "APTrust preservation services",
		DateTime:             now,
		Detail:               "File deleted from preservation storage",
		EventType:            constants.EventDeletion,
		Identifier:           uuid.NewString(),
		InstitutionID:        gf.InstitutionID,
		IntellectualObjectID: gf.IntellectualObjectID,
		GenericFileID:        gf.ID,
		Object:               "Minio S3 library",
		Outcome:              constants.OutcomeSuccess,
		OutcomeDetail:        deletionRequestView.RequestedByEmail,
		OutcomeInformation:   fmt.Sprintf("File deleted at the request of %s. Institutional approver: %s.", deletionRequestView.RequestedByEmail, deletionRequestView.ConfirmedByEmail),
	}, nil
}
