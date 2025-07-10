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
	"storage_option__in",
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
	PremisEvents         []*PremisEvent      `json:"premis_events" pg:"rel:has-many"`
	Checksums            []*Checksum         `json:"checksums" pg:"rel:has-many"`
	StorageRecords       []*StorageRecord    `json:"storage_records" pg:"rel:has-many"`
	AccessTime           time.Time           `json:"atime"`
	ChangeTime           time.Time           `json:"ctime"`
	ModTime              time.Time           `json:"mtime"`
	Gid                  int64               `json:"gid"`
	Gname                string              `json:"gname"`
	Uid                  int64               `json:"uid"`
	Uname                string              `json:"uname"`
	Mode                 int64               `json:"mode"`
}

// TODO: When selecting relations, order by UpdatedAt asc.

// GenericFileByID returns the file with the specified id.
// Returns pg.ErrNoRows if there is no match.
func GenericFileByID(id int64) (*GenericFile, error) {
	// Use pg's query builder because ours doesn't support ordering relation queries. Oops.
	var gf GenericFile
	db := common.Context().DB
	err := db.Model(&gf).
		Relation("PremisEvents", func(q *pg.Query) (*pg.Query, error) {
			return q.Order("premis_event.date_time desc"), nil
		}).
		Relation("Checksums", func(q *pg.Query) (*pg.Query, error) {
			return q.Order("checksum.created_at desc"), nil
		}).
		Relation("Institution").
		Relation("IntellectualObject").
		Relation("StorageRecords").
		Where("generic_file.id = ?", id).
		First()
	if gf.ID == 0 {
		return nil, err
	}
	return &gf, err
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
// Note that this returns related objects as well, including
// PremisEvents, Checksums, and StorageRecords.
func GenericFileSelect(query *Query) ([]*GenericFile, error) {
	var files []*GenericFile
	err := query.Relations("PremisEvents", "Checksums", "StorageRecords").Select(&files)
	return files, err
}

// Save saves this file to the database. This will peform an insert
// if GenericFile.ID is zero. Otherwise, it updates.
//
// Note that the insert/update also saves all associated records
// (checksums, storage records, and premis events).
func (gf *GenericFile) Save() error {
	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		return gf.saveInTransaction(tx)
	})
}

// GenericFileCreateBatch creates a batch of GenericFiles and
// their dependent records (PremisEvents, Checksums, and StorageRecords)
// in a single transaction. This transaction's single commit is much
// more efficient than doing one commit per insert.
//
// This is used heavily during ingest.
func GenericFileCreateBatch(files []*GenericFile) error {
	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		var err error
		for _, gf := range files {
			err = gf.saveInTransaction(tx)
			if err != nil {
				break
			}
		}
		return err
	})
}

func (gf *GenericFile) saveInTransaction(tx *pg.Tx) error {
	gf.SetTimestamps()
	validationErr := gf.Validate()
	if validationErr != nil {
		common.Context().Log.Error().Msgf("GenericFile save failed on validation of file  (%s). Error: %s", gf.Identifier, validationErr.Error())
		return validationErr
	}
	var err error
	if gf.ID == 0 {
		_, err = tx.Model(gf).Insert()
	} else {
		_, err = tx.Model(gf).WherePK().Update()
	}
	if err == nil {
		err = gf.saveChecksumsTx(tx)
	}
	if err == nil {
		err = gf.saveStorageRecordsTx(tx)
	}
	if err == nil {
		err = gf.saveEventsTx(tx)
	}
	return err
}

func (gf *GenericFile) saveChecksumsTx(tx *pg.Tx) error {
	for _, checksum := range gf.Checksums {
		// Checksums can't be updated, only added.
		if checksum.ID > 0 {
			continue
		}
		checksum.GenericFileID = gf.ID
		checksum.SetTimestamps()
		validationErr := checksum.Validate()
		if validationErr != nil {
			common.Context().Log.Error().Msgf("GenericFile batch insertion failed on validation of checksum (%s) - %s. Error: %s", gf.Identifier, checksum.Digest, validationErr.Error())
			return validationErr
		}
		// Checksums can't be updated, only added.
		_, err := tx.Model(checksum).OnConflict("DO NOTHING").Insert()
		if err != nil {
			common.Context().Log.Error().Msgf("GenericFile save failed on insert of checksum (%s) - %s. Error: %s", gf.Identifier, checksum.Digest, err.Error())
			return err
		}
	}
	return nil
}

func (gf *GenericFile) saveStorageRecordsTx(tx *pg.Tx) error {
	for _, sr := range gf.StorageRecords {
		if sr.ID > 0 {
			continue // already saved and updated aren't allowed
		}
		sr.GenericFileID = gf.ID
		validationErr := sr.Validate()
		if validationErr != nil {
			common.Context().Log.Error().Msgf("GenericFile save failed on validation of storage record (%s) - %s. Error: %s", gf.Identifier, sr.URL, validationErr.Error())
			return validationErr
		}
		// We can have only one record per URL.
		_, err := tx.Model(sr).OnConflict("DO NOTHING").Insert()
		if err != nil {
			common.Context().Log.Error().Msgf("GenericFile save failed on insert of storage record (%s) %s. Error: %s", gf.Identifier, sr.URL, err.Error())
			return err
		}
	}
	return nil
}

func (gf *GenericFile) saveEventsTx(tx *pg.Tx) error {
	for _, event := range gf.PremisEvents {
		if event.ID > 0 {
			continue // already saved and updates aren't allowed
		}
		event.InstitutionID = gf.InstitutionID
		event.IntellectualObjectID = gf.IntellectualObjectID
		event.GenericFileID = gf.ID
		event.SetTimestamps()
		validationErr := event.Validate()
		if validationErr != nil {
			common.Context().Log.Error().Msgf("GenericFile save failed on validation of event (%s) - %s. Error: %s", gf.Identifier, event.EventType, validationErr.Error())
			return validationErr
		}
		// Premis events can only be inserted, not updated.
		_, err := tx.Model(event).Insert()
		if err != nil {
			common.Context().Log.Error().Msgf("GenericFile batch insertion failed on insert of event (%s) - %s. Error: %s", gf.Identifier, event.EventType, err.Error())
			return err
		}
	}
	return nil
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

func (gf *GenericFile) ValidateChanges(updatedFile *GenericFile) error {
	if gf.ID != updatedFile.ID {
		return common.ErrIDMismatch
	}
	if gf.InstitutionID != updatedFile.InstitutionID {
		return common.ErrInstIDChange
	}
	if gf.Identifier != updatedFile.Identifier {
		return common.ErrIdentifierChange
	}
	// Caller should force storage option of updated object to
	// match existing object before calling this validation function.
	if gf.State == constants.StateActive && gf.StorageOption != updatedFile.StorageOption {
		return common.ErrStorageOptionChange
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
		for _, sr := range gf.StorageRecords {
			_, err = tx.Model(sr).Where("id = ?", sr.ID).Delete()
			if err != nil {
				registryContext.Log.Error().Msgf("GenericFile deletion transaction failed on deletion of StorageRecord %d. File: %d (%s). Error: %v", sr.ID, gf.ID, gf.Identifier, err)
			}
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

// ActiveDeletionWorkItem returns the in-progress WorkItem for this file's
// deletion. The WorkItem may be for the deletion of this specific file,
// or of its parent object.
func (gf *GenericFile) ActiveDeletionWorkItem() (*WorkItem, error) {
	cols := []string{
		"generic_file_id",
		"intellectual_object_id",
	}
	ops := []string{
		"=",
		"=",
	}
	vals := []interface{}{
		gf.ID,
		gf.IntellectualObjectID,
	}
	query := NewQuery().
		Or(cols, ops, vals).
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

func (gf *GenericFile) AssertDeletionPreconditions() error {
	if gf.State == constants.StateDeleted {
		return fmt.Errorf("File is already in deleted state")
	}
	if !gf.HasPassedMinimumRetentionPeriod() {
		return fmt.Errorf("File has not passed minimum retention period")
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
	deletionRequest, err := DeletionRequestViewByID(workItem.DeletionRequestID)
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
		Detail:               "All copies of this file have been deleted from preservation storage",
		EventType:            constants.EventDeletion,
		Identifier:           uuid.NewString(),
		InstitutionID:        gf.InstitutionID,
		IntellectualObjectID: gf.IntellectualObjectID,
		GenericFileID:        gf.ID,
		Object:               "Minio S3 library",
		Outcome:              constants.OutcomeSuccess,
		OutcomeDetail:        deletionRequestView.RequestedByEmail,
		OutcomeInformation:   fmt.Sprintf("File deleted at the request of %s. Institutional approver: %s. This event confirms all preservation copies have been deleted.", deletionRequestView.RequestedByEmail, deletionRequestView.ConfirmedByEmail),
	}, nil
}

// EarliestDeletionDate returns the earliest date on which this
// file can be deleted, per retention rules that apply to the
// object's storage option.
//
// The following (rare) case will return a false positive:
// file was ingested five years ago, deleted four years ago,
// and then ingested again yesterday.
//
// We can sort through Premis Events to solve these false positives,
// but that's very expensive and false positives probably are
// less than 0.2% of all cases.
func (gf *GenericFile) EarliestDeletionDate() time.Time {
	minRetentionDays := common.Context().Config.RetentionMinimum.For(gf.StorageOption)
	return gf.CreatedAt.AddDate(0, 0, minRetentionDays)
}

// HasPassedMinimumRetentionPeriod returns true if this object has
// passed the minimum retention period for its storage option.
func (gf *GenericFile) HasPassedMinimumRetentionPeriod() bool {
	return gf.EarliestDeletionDate().Before(time.Now())
}
