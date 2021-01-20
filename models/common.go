package models

import (
	"fmt"
	"reflect"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
)

type Model interface {
	GetID() int64
	Authorize(*User, string) error
	IsReadOnly() bool
	SupportsSoftDelete() bool
	SetSoftDeleteAttributes(*User)
	ClearSoftDeleteAttributes()
	SetTimestamps()
	BeforeSave() error
}

// Find finds an object by ID
func Find(obj Model, id int64, actingUser *User) error {
	ctx := common.Context()
	err := ctx.DB.Model(obj).Where("id = ?", id).Select()
	if err != nil {
		return err
	}
	return obj.Authorize(actingUser, constants.ActionRead)
}

func Select(models interface{}, q *Query) error {
	db := common.Context().DB
	orm := db.Model(models).Where(q.WhereClause(), q.Params()...)
	fmt.Println(q.WhereClause(), q.Params(), q.OrderBy, q.Offset, q.Limit)
	if q.OrderBy != "" {
		orm.Order(q.OrderBy)
	}
	if q.Limit > 0 {
		orm.Limit(q.Limit)
	}
	if q.Offset >= 0 {
		orm.Offset(q.Offset)
	}
	return orm.Select()
}

func Save(model Model, actingUser *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	if model.GetID() > 0 {
		fmt.Println("Updating model")
		return update(model, actingUser)
	}
	fmt.Println("Inserting model")
	return insert(model, actingUser)
}

func insert(model Model, actingUser *User) error {
	err := model.Authorize(actingUser, constants.ActionCreate)
	if err != nil {
		return err
	}
	model.SetTimestamps()

	err = model.BeforeSave()
	if err != nil {
		return err
	}

	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).Insert()
		fmt.Println("Transaction", err, "ID", model.GetID())
		return err
	})
}

func update(model Model, actingUser *User) error {
	err := model.Authorize(actingUser, constants.ActionUpdate)
	if err != nil {
		return err
	}
	model.SetTimestamps()

	err = model.BeforeSave()
	if err != nil {
		return err
	}

	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Update()
		return err
	})
}

func Delete(model Model, actingUser *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	err := model.Authorize(actingUser, constants.ActionDelete)
	if err != nil {
		return err
	}
	if model.SupportsSoftDelete() {
		return softDelete(model, actingUser)
	}
	return hardDelete(model, actingUser)
}

func softDelete(model Model, actingUser *User) error {
	model.SetSoftDeleteAttributes(actingUser)
	model.SetTimestamps()
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Update()
		return err
	})
}

func hardDelete(model Model, actingUser *User) error {
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Delete()
		return err
	})
}

func Undelete(model Model, actingUser *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	err := model.Authorize(actingUser, constants.ActionUpdate)
	if err != nil {
		return err
	}
	model.SetTimestamps()
	model.ClearSoftDeleteAttributes()
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Update()
		return err
	})
}

func Int64Value(obj interface{}, fieldName string) int64 {
	value := reflect.ValueOf(obj)
	if value.Type().Kind() != reflect.Ptr {
		value = reflect.New(reflect.TypeOf(obj))
	}
	field := value.Elem().FieldByName(fieldName)
	if field.IsValid() {
		return field.Int()
	}
	return int64(0)
}

// IsNoRowError returns true if err is pg.ErrNoRows. For some reason,
// err doesn't compare correctly with pg.ErrNoRows, and errors.Is()
// doesn't work either. Probably because pg.ErrNoRows is an alias of
// an error in the pg/internal package, which we cannot access.
func IsNoRowError(err error) bool {
	return err != nil && err.Error() == pg.ErrNoRows.Error()
}
