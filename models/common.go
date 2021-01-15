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
func Find(obj Model, id int64, user *User) error {
	ctx := common.Context()
	err := ctx.DB.Model(obj).Where("id = ?", id).Select()
	if err != nil {
		return err
	}
	return obj.Authorize(user, constants.ActionView)
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

func Save(model Model, user *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	if model.GetID() > 0 {
		fmt.Println("Updating model")
		return update(model, user)
	}
	fmt.Println("Inserting model")
	return insert(model, user)
}

func insert(model Model, user *User) error {
	err := model.Authorize(user, constants.ActionCreate)
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

func update(model Model, user *User) error {
	err := model.Authorize(user, constants.ActionEdit)
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
		_, err := db.Model(model).Update()
		return err
	})
}

func Delete(model Model, user *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	err := model.Authorize(user, constants.ActionDelete)
	if err != nil {
		return err
	}
	if model.SupportsSoftDelete() {
		return softDelete(model, user)
	}
	return hardDelete(model, user)
}

func softDelete(model Model, user *User) error {
	model.SetSoftDeleteAttributes(user)
	model.SetTimestamps()
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).Update()
		return err
	})
}

func hardDelete(model Model, user *User) error {
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).Delete()
		return err
	})
}

func Undelete(model Model, user *User) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	err := model.Authorize(user, constants.ActionEdit)
	if err != nil {
		return err
	}
	model.SetTimestamps()
	model.ClearSoftDeleteAttributes()
	db := common.Context().DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).Update()
		return err
	})
}

func Inv64Value(obj interface{}, fieldName string) int64 {
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
