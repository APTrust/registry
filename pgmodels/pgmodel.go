package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

type Model interface {
	GetID() int64
	Save() error
	SetTimestamps()
	IsJoinModel() bool
	Validate() *common.ValidationError
}

type BaseModel struct {
	ID int64 `pg:"id" form:"id" json:"id"`
}

func (bm *BaseModel) GetID() int64 {
	return bm.ID
}

func (bm *BaseModel) IsJoinModel() bool {
	return false
}

func (bm *BaseModel) SetTimestamps() {
	// No-Op
}

func (bm *BaseModel) Save() error {
	return common.ErrSubclassMustImplement
}

type TimestampModel struct {
	BaseModel
	CreatedAt time.Time `bun:",nullzero" json:",omitempty"`
	UpdatedAt time.Time `bun:",nullzero" json:",omitempty"`
}

func (tsm *TimestampModel) SetTimestamps() {
	now := time.Now().UTC()
	if tsm.CreatedAt.IsZero() {
		tsm.CreatedAt = now
	}
	tsm.UpdatedAt = now
}

type JoinModel struct {
}

func (jm *JoinModel) GetID() int64 {
	return 0
}

func (jm *JoinModel) IsJoinModel() bool {
	return true
}

func (jm *JoinModel) SetTimestamps() {
	// No-Op
}

type xactType int

const (
	TypeInsert xactType = iota
	TypeUpdate
)

func insert(model interface{}) error {
	return transact(model, TypeInsert)
}

func update(model interface{}) error {
	return transact(model, TypeUpdate)
}

func transact(model interface{}, action xactType) error {
	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		var err error
		if action == TypeInsert {
			_, err = tx.Model(model).Insert()
		} else {
			_, err = tx.Model(model).WherePK().Update()
		}
		if err != nil {
			registryContext.Log.Error().Msgf("Transaction failed. Model: %v. Error: %v", model, err)
		}
		return err
	})
}

// IsNoRowError returns true if err is pg.ErrNoRows. For some reason,
// err doesn't compare correctly with pg.ErrNoRows, and errors.Is()
// doesn't work either. Probably because pg.ErrNoRows is an alias of
// an error in the pg/internal package, which we cannot access.
func IsNoRowError(err error) bool {
	return err != nil && err.Error() == pg.ErrNoRows.Error()
}

func InstIDFor(resourceType string, resourceID int64) (id int64, err error) {
	ctx := common.Context()
	db := ctx.DB
	switch resourceType {
	case "Alert":
		alert := &Alert{}
		err = db.Model(alert).Column("institution_id").Where("id = ?", resourceID).Select()
		id = alert.InstitutionID
	case "Checksum":
		cs := &Checksum{}
		err = db.Model(cs).Column("_").Relation("GenericFile.institution_id").Where(`"checksum"."id" = ?`, resourceID).Select()
		if cs != nil && cs.GenericFile != nil {
			id = cs.GenericFile.InstitutionID
		}
	case "DeletionRequest":
		req := &DeletionRequest{}
		err = db.Model(req).Column("institution_id").Where("id = ?", resourceID).Select()
		id = req.InstitutionID
	case "GenericFile":
		gf := &GenericFile{}
		err = db.Model(gf).Column("institution_id").Where("id = ?", resourceID).Select()
		id = gf.InstitutionID
	case "Institution":
		id = resourceID
	case "IntellectualObject":
		obj := &IntellectualObject{}
		err = db.Model(obj).Column("institution_id").Where("id = ?", resourceID).Select()
		id = obj.InstitutionID
	case "PremisEvent":
		pe := &PremisEvent{}
		err = db.Model(pe).Column("institution_id").Where("id = ?", resourceID).Select()
		id = pe.InstitutionID
	case "StorageRecord":
		sr := &StorageRecord{}
		err = db.Model(sr).Column("_").Relation("GenericFile.institution_id").Where(`"storage_record"."id" = ?`, resourceID).Select()
		if sr != nil && sr.GenericFile != nil {
			id = sr.GenericFile.InstitutionID
		}
	case "User":
		user := &User{}
		err = db.Model(user).Column("institution_id").Where("id = ?", resourceID).Select()
		id = user.InstitutionID
	case "WorkItem":
		item := &WorkItem{}
		err = db.Model(item).Column("institution_id").Where("id = ?", resourceID).Select()
		id = item.InstitutionID
	default:
		ctx.Log.Error().Msgf("pgmodels.InstIDFor got unknown type '%s'", resourceType)
		err = common.ErrInvalidParam
	}
	return id, err
}

var filters map[string][]string

func FiltersFor(typeName string) []string {
	if filters == nil || len(filters) == 0 {
		initFilters()
	}
	return filters[typeName]
}

// TODO: Have an init function in each model to take care of this.
// Init function should register info about each model.
func initFilters() {
	filters = make(map[string][]string)
	filters["Alert"] = AlertFilters
	filters["DeletionRequest"] = DeletionRequestFilters
	filters["DepositStats"] = DepositStatsFilters
	filters["GenericFile"] = GenericFileFilters
	filters["IntellectualObject"] = IntellectualObjectFilters
	filters["Institution"] = InstitutionFilters
	filters["PremisEvent"] = PremisEventFilters
	filters["User"] = UserFilters
	filters["WorkItem"] = WorkItemFilters
}
