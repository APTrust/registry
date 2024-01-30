package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const (
	ErrDeletionInstitutionID = "Deletion request requires institution id."
	ErrDeletionRequesterID   = "Deletion request requires requester id."
	ErrDeletionWrongInst     = "Deletion request user belongs to wrong institution."
	ErrDeletionWrongRole     = "Deletion confirmer/canceller must be institutional admin."
	ErrDeletionUserNotFound  = "User does not exist."
	ErrDeletionUserInactive  = "User has been deactivated."
	ErrTokenNotEncrypted     = "Token must be encrypted."
	ErrDeletionIllegalObject = "User cannot delete files or objects belonging to other institutions."
	ErrDeletionBadAdmin      = "Admin cannot confirm deletion they requested. This must be approved by a second admin."
	ErrDeletionBadQuery      = "Cannot get admin list for this institution."
)

// init does some setup work so go-pg can recognize many-to-many
// relations. Go automatically calls this function once when package
// is imported.
func init() {
	orm.RegisterTable((*DeletionRequestsGenericFiles)(nil))
	orm.RegisterTable((*DeletionRequestsIntellectualObjects)(nil))
}

type DeletionRequest struct {
	BaseModel
	InstitutionID              int64                 `json:"institution_id"`
	RequestedByID              int64                 `json:"-"`
	RequestedAt                time.Time             `json:"requested_at"`
	ConfirmationToken          string                `json:"-" pg:"-"`
	EncryptedConfirmationToken string                `json:"-"`
	ConfirmedByID              int64                 `json:"-"`
	ConfirmedAt                time.Time             `json:"confirmed_at"`
	CancelledByID              int64                 `json:"-"`
	CancelledAt                time.Time             `json:"cancelled_at"`
	RequestedBy                *User                 `json:"requested_by" pg:"rel:has-one"`
	ConfirmedBy                *User                 `json:"confirmed_by" pg:"rel:has-one"`
	CancelledBy                *User                 `json:"cancelled_by" pg:"rel:has-one"`
	GenericFiles               []*GenericFile        `json:"generic_files" pg:"many2many:deletion_requests_generic_files"`
	IntellectualObjects        []*IntellectualObject `json:"intellectual_objects" pg:"many2many:deletion_requests_intellectual_objects"`
	WorkItems                  []*WorkItem           `json:"work_item" pg:"rel:has-many"`
}

type DeletionRequestsGenericFiles struct {
	tableName         struct{} `pg:"deletion_requests_generic_files"`
	DeletionRequestID int64
	GenericFileID     int64
}

type DeletionRequestsIntellectualObjects struct {
	tableName            struct{} `pg:"deletion_requests_intellectual_objects"`
	DeletionRequestID    int64
	IntellectualObjectID int64
}

func NewDeletionRequest() (*DeletionRequest, error) {
	confToken := common.RandomToken()
	encConfToken, err := common.EncryptPassword(confToken)
	if err != nil {
		return nil, err
	}
	return &DeletionRequest{
		ConfirmationToken:          confToken,
		EncryptedConfirmationToken: encConfToken,
		GenericFiles:               make([]*GenericFile, 0),
		IntellectualObjects:        make([]*IntellectualObject, 0),
	}, nil
}

// DeletionRequestByID returns the institution with the specified id.
// Returns pg.ErrNoRows if there is no match.
func DeletionRequestByID(id int64) (*DeletionRequest, error) {
	query := NewQuery().Relations("RequestedBy", "ConfirmedBy", "CancelledBy", "GenericFiles", "IntellectualObjects", "WorkItems").Where(`"deletion_request"."id"`, "=", id)
	return DeletionRequestGet(query)
}

// DeletionRequestGet returns the first deletion request matching the query.
func DeletionRequestGet(query *Query) (*DeletionRequest, error) {
	var request DeletionRequest
	err := query.Select(&request)
	return &request, err
}

// DeletionRequestSelect returns all deletion requests matching the query.
func DeletionRequestSelect(query *Query) ([]*DeletionRequest, error) {
	var requests []*DeletionRequest
	err := query.Select(&requests)
	return requests, err
}

// DeletionRequestIncludesFile returns true if the deletion request with the
// specified ID includes the generic file with the specified ID.
func DeletionRequestIncludesFile(requestID, gfID int64) (bool, error) {
	db := common.Context().DB
	var count int
	query := `SELECT count(*) FROM deletion_requests_generic_files where deletion_request_id = ? and generic_file_id = ?`
	_, err := db.Model((*DeletionRequestsGenericFiles)(nil)).QueryOne(pg.Scan(&count), query, requestID, gfID)
	return count > 0, err
}

// DeletionRequestIncludesObject returns true if the deletion request with the
// specified ID includes the intellectual object with the specified ID.
func DeletionRequestIncludesObject(requestID, objID int64) (bool, error) {
	db := common.Context().DB
	var count int
	query := `SELECT count(*) FROM deletion_requests_intellectual_objects where deletion_request_id = ? and intellectual_object_id = ?`
	_, err := db.Model((*DeletionRequestsGenericFiles)(nil)).QueryOne(pg.Scan(&count), query, requestID, objID)
	return count > 0, err
}

// Save saves this requestitution to the database. This will peform an insert
// if DeletionRequest.ID is zero. Otherwise, it updates.
func (request *DeletionRequest) Save() error {
	err := request.Validate()
	if err != nil {
		return err
	}
	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		var err error
		if request.ID == 0 {
			_, err = tx.Model(request).Insert()
		} else {
			_, err = tx.Model(request).WherePK().Update()
		}
		if err != nil {
			registryContext.Log.Error().Msgf("Transaction failed. Model: %v. Error: %v", request, err)
			return err
		}
		return request.saveRelations(tx)
	})
}

// Validation enforces business rules, including who can request and
// confirm deletions. Although our general security middleware should
// prevent any of these problems from ever occurring, we want to
// double check everything here because we're a preservation archive
// and deletion is a destructive action. We must be sure deletion is a
// deliberate act initiated and confirmed by authorized individuals.
func (request *DeletionRequest) Validate() *common.ValidationError {
	errors := make(map[string]string)

	if request.InstitutionID < 1 {
		errors["InstitutionID"] = ErrDeletionInstitutionID
	}

	request.validateRequestedBy(errors)
	request.validateConfirmedBy(errors)
	request.validateCancelledBy(errors)

	// Make sure tokens are actually encrypted
	if !common.LooksEncrypted(request.EncryptedConfirmationToken) {
		errors["EncryptedConfirmationToken"] = ErrTokenNotEncrypted
	}

	// Make sure all objects and files belong to the requesting
	// user's institution.
	for _, obj := range request.IntellectualObjects {
		if obj.InstitutionID != request.RequestedBy.InstitutionID {
			errors["IntellectualObjects"] = ErrDeletionIllegalObject
			break
		}
	}
	for _, gf := range request.GenericFiles {
		if gf.InstitutionID != request.RequestedBy.InstitutionID {
			errors["GenericFiles"] = ErrDeletionIllegalObject
			break
		}
	}

	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}

func (request *DeletionRequest) validateRequestedBy(errors map[string]string) {
	var err error
	if request.RequestedByID > 0 && (request.RequestedBy == nil || request.RequestedBy.ID != request.RequestedByID) {
		request.RequestedBy, err = UserByID(request.RequestedByID)
	}
	if request.RequestedByID < 1 {
		errors["RequestedByID"] = ErrDeletionRequesterID
	} else if request.RequestedBy == nil || err != nil {
		errors["RequestedByID"] = ErrDeletionUserNotFound
	} else if !request.RequestedBy.DeactivatedAt.IsZero() {
		errors["RequestedByID"] = ErrDeletionUserInactive
	} else if request.RequestedBy.InstitutionID != request.InstitutionID {
		errors["RequestedByID"] = ErrDeletionWrongInst
	}
}

func (request *DeletionRequest) validateConfirmedBy(errors map[string]string) {
	// Make sure approver has admin role at the right institution
	if request.ConfirmedByID > 0 {
		var err error
		if request.ConfirmedBy == nil || request.ConfirmedBy.ID != request.ConfirmedByID {
			request.ConfirmedBy, err = UserByID(request.ConfirmedByID)
			if err != nil || request.ConfirmedBy == nil || request.ConfirmedBy.ID == 0 {
				errors["ConfirmedByID"] = ErrDeletionUserNotFound
				return
			}
		}
		instAdmins, err := UserSelect(NewQuery().Where("institution_id", "=", request.ConfirmedBy.InstitutionID).Where("role", "=", constants.RoleInstAdmin))
		if err != nil {
			errors["ConfirmedByID"] = ErrDeletionBadQuery
			return
		}
		if request.ConfirmedBy.InstitutionID != request.InstitutionID {
			errors["ConfirmedByID"] = ErrDeletionWrongInst
		} else if request.ConfirmedBy.Role != constants.RoleInstAdmin {
			// fmt.Println("Req Inst:", request.InstitutionID, "Conf Inst:", request.ConfirmedBy.InstitutionID)
			// fmt.Println(request.ConfirmedBy)
			errors["ConfirmedByID"] = ErrDeletionWrongRole
		} else if request.ConfirmedBy.ID == request.RequestedByID && len(instAdmins) > 1 {
			errors["ConfirmedByID"] = ErrDeletionBadAdmin
		}
	}
}

func (request *DeletionRequest) validateCancelledBy(errors map[string]string) {
	// Make sure canceller has admin role at the right institution
	if request.CancelledByID > 0 {
		if request.CancelledBy == nil || request.CancelledByID != request.CancelledBy.ID {
			user, err := UserByID(request.CancelledByID)
			if err != nil || user.ID == 0 {
				errors["CancelledByID"] = ErrDeletionUserNotFound
			} else if user.InstitutionID != request.InstitutionID {
				errors["CancelledByID"] = ErrDeletionWrongInst
			} else if user.Role != constants.RoleInstAdmin {
				errors["CancelledByID"] = ErrDeletionWrongRole
			}
		}
	}
}

func (request *DeletionRequest) AddFile(gf *GenericFile) {
	if request.GenericFiles == nil {
		request.GenericFiles = make([]*GenericFile, 0)
	}
	request.GenericFiles = append(request.GenericFiles, gf)
}

func (request *DeletionRequest) AddObject(obj *IntellectualObject) {
	if request.IntellectualObjects == nil {
		request.IntellectualObjects = make([]*IntellectualObject, 0)
	}
	request.IntellectualObjects = append(request.IntellectualObjects, obj)
}

func (request *DeletionRequest) saveRelations(tx *pg.Tx) error {
	err := request.saveFiles(tx)
	if err != nil {
		return err
	}
	err = request.saveObjects(tx)
	if err != nil {
		return err
	}
	return request.saveWorkItems(tx)
}

func (request *DeletionRequest) saveFiles(tx *pg.Tx) error {
	// Note: on conflict refers to unique index index_drgf_unique
	sql := "insert into deletion_requests_generic_files (deletion_request_id, generic_file_id) values (?, ?) on conflict do nothing"
	for _, gf := range request.GenericFiles {
		_, err := tx.Exec(sql, request.ID, gf.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (request *DeletionRequest) saveObjects(tx *pg.Tx) error {
	// Note: on conflict refers to unique index index_drio_unique
	sql := "insert into deletion_requests_intellectual_objects (deletion_request_id, intellectual_object_id) values (?, ?) on conflict do nothing"
	for _, obj := range request.IntellectualObjects {
		_, err := tx.Exec(sql, request.ID, obj.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (request *DeletionRequest) saveWorkItems(tx *pg.Tx) error {
	for _, item := range request.WorkItems {
		item.DeletionRequestID = request.ID
		err := item.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

// FirstFile returns the first GenericFile associated with this deletion
// request. Use this for simple, single-file deletions.
func (request *DeletionRequest) FirstFile() *GenericFile {
	if len(request.GenericFiles) > 0 {
		return request.GenericFiles[0]
	}
	return nil
}

// FirstObject returns the first IntellectualObject associated with
// this deletion request. Use this for simple, single-object deletions.
func (request *DeletionRequest) FirstObject() *IntellectualObject {
	if len(request.IntellectualObjects) > 0 {
		return request.IntellectualObjects[0]
	}
	return nil
}

// Confirm marks this DeletionRequest as confirmed. It's up to the caller
// to save the request and create an appropriate WorkItem.
func (request *DeletionRequest) Confirm(user *User) {
	request.ConfirmedBy = user
	request.ConfirmedByID = user.ID
	request.ConfirmedAt = time.Now().UTC()
}

// Cancel cancels this DeletionRequest. It's up to the caller to save
// this request after cancelling it.
func (request *DeletionRequest) Cancel(user *User) {
	request.CancelledBy = user
	request.CancelledByID = user.ID
	request.CancelledAt = time.Now().UTC()
}
