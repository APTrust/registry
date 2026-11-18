package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const (
	ErrTransferInstitutionID     = "Transfer request requires institution id."
	ErrTransferRequesterID       = "Transfer request requires requester id."
	ErrTransferWrongInst         = "Transfer request user belongs to wrong institution."
	ErrTransferWrongRole         = "Transfer confirmer/canceller must be institutional admin."
	ErrTransferUserNotFound      = "User does not exist."
	ErrTransferUserInactive      = "User has been deactivated."
	ErrTransferTokenNotEncrypted = "Token must be encrypted."
	ErrTransferIllegalObject     = "User cannot transfer objects belonging to other institutions."
	ErrTransferBadAdmin          = "Admin cannot confirm transfer they requested. This must be approved by a second admin."
	ErrTransferBadQuery          = "Cannot get admin list for this institution."
)

// init does some setup work so go-pg can recognize many-to-many
// relations. Go automatically calls this function once when package
// is imported.
func init() {
	orm.RegisterTable((*TransferRequestsIntellectualObjects)(nil))
}

type TransferRequest struct {
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
	IntellectualObjects        []*IntellectualObject `json:"intellectual_objects" pg:"many2many:transfer_requests_intellectual_objects"`
	WorkItems                  []*WorkItem           `json:"work_item" pg:"rel:has-many"`
}

type TransferRequestsIntellectualObjects struct {
	tableName            struct{} `pg:"transfer_requests_intellectual_objects"`
	TransferRequestID    int64
	IntellectualObjectID int64
}

func NewTransferRequest() (*TransferRequest, error) {
	confToken := common.RandomToken()
	encConfToken, err := common.EncryptPassword(confToken)
	if err != nil {
		return nil, err
	}
	return &TransferRequest{
		ConfirmationToken:          confToken,
		EncryptedConfirmationToken: encConfToken,
		IntellectualObjects:        make([]*IntellectualObject, 0),
	}, nil
}

// TransferRequestByID returns the institution with the specified id.
// Returns pg.ErrNoRows if there is no match.
func TransferRequestByID(id int64) (*TransferRequest, error) {
	query := NewQuery().Relations("RequestedBy", "ConfirmedBy", "CancelledBy", "IntellectualObjects", "WorkItems").Where(`"transfer_request"."id"`, "=", id)
	return TransferRequestGet(query)
}

// TransferRequestGet returns the first transfer request matching the query.
func TransferRequestGet(query *Query) (*TransferRequest, error) {
	var request TransferRequest
	err := query.Select(&request)
	return &request, err
}

// TransferRequestSelect returns all transfer requests matching the query.
func TransferRequestSelect(query *Query) ([]*TransferRequest, error) {
	var requests []*TransferRequest
	err := query.Select(&requests)
	return requests, err
}

// TransferRequestIncludesObject returns true if the transfer request with the
// specified ID includes the intellectual object with the specified ID.
func TransferRequestIncludesObject(requestID, objID int64) (bool, error) {
	db := common.Context().DB
	var count int
	query := `SELECT count(*) FROM transfer_requests_intellectual_objects where transfer_request_id = ? and intellectual_object_id = ?`
	_, err := db.Model((*TransferRequestsIntellectualObjects)(nil)).QueryOne(pg.Scan(&count), query, requestID, objID)
	return count > 0, err
}

// Save saves this requestitution to the database. This will peform an insert
// if TransferRequest.ID is zero. Otherwise, it updates.
func (request *TransferRequest) Save() error {
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
// confirm transfers. Although our general security middleware should
// prevent any of these problems from ever occurring, we want to
// double check everything here because we're a preservation archive
// and transferring has a large impact on costs. We must be sure transfer is a
// deliberate act initiated and confirmed by authorized individuals.
func (request *TransferRequest) Validate() *common.ValidationError {
	errors := make(map[string]string)

	if request.InstitutionID < 1 {
		errors["InstitutionID"] = ErrTransferInstitutionID
	}

	request.validateRequestedBy(errors)
	request.validateConfirmedBy(errors)
	request.validateCancelledBy(errors)

	// Make sure tokens are actually encrypted
	if !common.LooksEncrypted(request.EncryptedConfirmationToken) {
		errors["EncryptedConfirmationToken"] = ErrTokenNotEncrypted
	}

	// Make sure all objects belong to the requesting
	// user's institution.
	for _, obj := range request.IntellectualObjects {
		if obj.InstitutionID != request.RequestedBy.InstitutionID {
			errors["IntellectualObjects"] = ErrTransferIllegalObject
			break
		}
	}

	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}

func (request *TransferRequest) validateRequestedBy(errors map[string]string) {
	var err error
	if request.RequestedByID > 0 && (request.RequestedBy == nil || request.RequestedBy.ID != request.RequestedByID) {
		request.RequestedBy, err = UserByID(request.RequestedByID)
	}
	if request.RequestedByID < 1 {
		errors["RequestedByID"] = ErrTransferRequesterID
	} else if request.RequestedBy == nil || err != nil {
		errors["RequestedByID"] = ErrTransferUserNotFound
	} else if !request.RequestedBy.DeactivatedAt.IsZero() {
		errors["RequestedByID"] = ErrTransferUserInactive
	} else if request.RequestedBy.InstitutionID != request.InstitutionID {
		errors["RequestedByID"] = ErrTransferWrongInst
	}
}

func (request *TransferRequest) validateConfirmedBy(errors map[string]string) {
	// Make sure approver has admin role at the right institution
	if request.ConfirmedByID > 0 {
		var err error
		if request.ConfirmedBy == nil || request.ConfirmedBy.ID != request.ConfirmedByID {
			request.ConfirmedBy, err = UserByID(request.ConfirmedByID)
			if err != nil || request.ConfirmedBy == nil || request.ConfirmedBy.ID == 0 {
				errors["ConfirmedByID"] = ErrTransferUserNotFound
				return
			}
		}
		query := NewQuery().Where("institution_id", "=", request.ConfirmedBy.InstitutionID).Where("role", "=", constants.RoleInstAdmin).IsNull("deactivated_at")
		instAdmins, err := UserSelect(query)
		if err != nil {
			errors["ConfirmedByID"] = ErrTransferBadQuery
			return
		}
		if request.ConfirmedBy.InstitutionID != request.InstitutionID {
			errors["ConfirmedByID"] = ErrTransferWrongInst
		} else if request.ConfirmedBy.Role != constants.RoleInstAdmin {
			// fmt.Println("Req Inst:", request.InstitutionID, "Conf Inst:", request.ConfirmedBy.InstitutionID)
			// fmt.Println(request.ConfirmedBy)
			errors["ConfirmedByID"] = ErrTransferWrongRole
		} else if request.ConfirmedBy.ID == request.RequestedByID && len(instAdmins) > 1 {
			errors["ConfirmedByID"] = ErrTransferBadAdmin
		}
	}
}

func (request *TransferRequest) validateCancelledBy(errors map[string]string) {
	// Make sure canceller has admin role at the right institution
	if request.CancelledByID > 0 {
		if request.CancelledBy == nil || request.CancelledByID != request.CancelledBy.ID {
			user, err := UserByID(request.CancelledByID)
			if err != nil || user.ID == 0 {
				errors["CancelledByID"] = ErrTransferUserNotFound
			} else if user.InstitutionID != request.InstitutionID {
				errors["CancelledByID"] = ErrTransferWrongInst
			} else if user.Role != constants.RoleInstAdmin {
				errors["CancelledByID"] = ErrTransferWrongRole
			}
		}
	}
}

func (request *TransferRequest) AddObject(obj *IntellectualObject) {
	if request.IntellectualObjects == nil {
		request.IntellectualObjects = make([]*IntellectualObject, 0)
	}
	request.IntellectualObjects = append(request.IntellectualObjects, obj)
}

func (request *TransferRequest) saveRelations(tx *pg.Tx) error {
	err := request.saveObjects(tx)
	if err != nil {
		return err
	}
	return request.saveWorkItems(tx)
}

func (request *TransferRequest) saveObjects(tx *pg.Tx) error {
	// Note: on conflict refers to unique index index_trio_unique
	sql := "insert into transfer_requests_intellectual_objects (transfer_request_id, intellectual_object_id) values (?, ?) on conflict do nothing"
	for _, obj := range request.IntellectualObjects {
		_, err := tx.Exec(sql, request.ID, obj.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (request *TransferRequest) saveWorkItems(tx *pg.Tx) error {
	for _, item := range request.WorkItems {
		// item.TransferRequestID = request.ID
		err := item.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

// FirstObject returns the first IntellectualObject associated with
// this transfer request. Use this for simple, single-object transfers.
func (request *TransferRequest) FirstObject() *IntellectualObject {
	if len(request.IntellectualObjects) > 0 {
		return request.IntellectualObjects[0]
	}
	return nil
}

// Confirm marks this TransferRequest as confirmed. It's up to the caller
// to save the request and create an appropriate WorkItem.
func (request *TransferRequest) Confirm(user *User) {
	request.ConfirmedBy = user
	request.ConfirmedByID = user.ID
	request.ConfirmedAt = time.Now().UTC()
}

// Cancel cancels this TransferRequest. It's up to the caller to save
// this request after cancelling it.
func (request *TransferRequest) Cancel(user *User) {
	request.CancelledBy = user
	request.CancelledByID = user.ID
	request.CancelledAt = time.Now().UTC()
}
