package models

import (
	"reflect"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
)

type Model interface {
	GetID() int64
	Authorize(*User, string) error
	DeleteIsForbidden() bool
	UpdateIsForbidden() bool
	IsReadOnly() bool // views are read-only
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
	ctx := common.Context()
	orm := ctx.DB.Model(models)
	// Empty where clause causes orm to generate empty parens -> ()
	// which causes a SQL error. Include where only if non-empty.
	if q.WhereClause() != "" {
		orm = orm.Where(q.WhereClause(), q.Params()...)
	}
	ctx.Log.Debug().Msgf("SELECT PARAMS: %s, %v, order by %s, offset %d, limit %d ", q.WhereClause(), q.Params(), q.OrderBy, q.Offset, q.Limit)
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
		return update(model, actingUser)
	}
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

	ctx := common.Context()
	ctx.Log.Debug().Msgf("Insert %s %v", TypeOf(model), model)
	return ctx.DB.RunInTransaction(ctx.DB.Context(), func(*pg.Tx) error {
		_, err := ctx.DB.Model(model).Insert()
		ctx.Log.Error().Msgf("Transaction on ID %d: %v", model.GetID(), err)
		return err
	})
}

func update(model Model, actingUser *User) error {
	if model.IsReadOnly() || model.UpdateIsForbidden() {
		return common.ErrNotSupported
	}
	err := model.Authorize(actingUser, constants.ActionUpdate)
	if err != nil {
		return err
	}
	model.SetTimestamps()

	err = model.BeforeSave()
	if err != nil {
		return err
	}

	ctx := common.Context()
	ctx.Log.Debug().Msgf("Insert %s %v", TypeOf(model), model)
	return ctx.DB.RunInTransaction(ctx.DB.Context(), func(*pg.Tx) error {
		_, err := ctx.DB.Model(model).WherePK().Update()
		return err
	})
}

func Delete(model Model, actingUser *User) error {
	if model.IsReadOnly() || model.DeleteIsForbidden() {
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
	if model.IsReadOnly() || model.UpdateIsForbidden() {
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

func TypeOf(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

// IsNoRowError returns true if err is pg.ErrNoRows. For some reason,
// err doesn't compare correctly with pg.ErrNoRows, and errors.Is()
// doesn't work either. Probably because pg.ErrNoRows is an alias of
// an error in the pg/internal package, which we cannot access.
func IsNoRowError(err error) bool {
	return err != nil && err.Error() == pg.ErrNoRows.Error()
}
