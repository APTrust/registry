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