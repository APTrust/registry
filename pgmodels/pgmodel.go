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
			registryContext.Log.Error().Msgf("Transaction failed: %v", err)
		}
		return err
	})
}
