package pgmodels

import (
	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

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
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		var err error
		if action == TypeInsert {
			_, err = db.Model(model).Insert()
		} else {
			_, err = db.Model(model).WherePK().Update()
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
	case "Checksum":
		cs := &Checksum{}
		err = db.Model(cs).Column("_").Relation("GenericFile.institution_id").Where(`"checksum"."id" = ?`, resourceID).Select()
		if cs != nil && cs.GenericFile != nil {
			id = cs.GenericFile.InstitutionID
		}
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
