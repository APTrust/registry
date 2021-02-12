package models

import (
	"reflect"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
)

type DataStore struct {
	ctx        *common.APTContext
	actingUser *User
}

func NewDataStore(actingUser *User) *DataStore {
	return &DataStore{
		ctx:        common.Context(),
		actingUser: actingUser,
	}
}

func (ds *DataStore) InstitutionList(q *Query) ([]*Institution, error) {
	institutions := make([]*Institution, 0)
	err := ds._select(&institutions, q)
	return institutions, err
}

func (ds *DataStore) UserDelete(user *User) error {
	return ds.delete(user)
}

func (ds *DataStore) UserFind(id int64) (*User, error) {
	user := &User{}
	err := ds.ctx.DB.Model(user).Relation("Institution").Where(`"user"."id" = ?`, id).Select()
	if err != nil {
		return nil, err
	}
	err = user.Authorize(ds.actingUser, constants.ActionRead)
	return user, err
}

func (ds *DataStore) UserFindByEmail(email string) (*User, error) {
	user := &User{}
	err := ds.ctx.DB.Model(user).Relation("Institution").Where(`"user"."email" = ?`, email).Select()
	if err != nil {
		return nil, err
	}
	err = user.Authorize(ds.actingUser, constants.ActionRead)
	return user, err
}

func (ds *DataStore) UserList(q *Query) ([]*User, error) {
	users := make([]*User, 0)
	err := ds._select(&users, q)
	return users, err
}

func (ds *DataStore) UserSave(user *User) error {
	return ds.save(user)
}

func (ds *DataStore) UserSignIn(email, password, ipAddr string) (*User, error) {
	user, err := ds.UserFindByEmail(email)
	if IsNoRowError(err) {
		return nil, common.ErrInvalidLogin
	} else if err != nil {
		return nil, err
	}
	if !user.DeactivatedAt.IsZero() {
		return nil, common.ErrAccountDeactivated
	}
	if !common.ComparePasswords(user.EncryptedPassword, password) {
		ds.ctx.Log.Warn().Msgf("Wrong password for user %s", email)
		return nil, common.ErrInvalidLogin
	}
	user.SignInCount = user.SignInCount + 1
	if user.CurrentSignInIP != "" {
		user.LastSignInIP = user.CurrentSignInIP
	}
	if user.CurrentSignInAt.IsZero() {
		user.LastSignInAt = user.CurrentSignInAt
	}
	user.CurrentSignInIP = ipAddr
	user.CurrentSignInAt = time.Now().UTC()
	err = ds.save(user)
	return user, err
}

func (ds *DataStore) UserSignOut(user *User) error {
	if user.CurrentSignInIP != "" {
		user.LastSignInIP = user.CurrentSignInIP
	}
	if !user.CurrentSignInAt.IsZero() {
		user.LastSignInAt = user.CurrentSignInAt
	}
	user.CurrentSignInIP = ""
	user.CurrentSignInAt = time.Time{}
	return ds.save(user)
}

func (ds *DataStore) UserUndelete(user *User) error {
	return ds.undelete(user)
}

func (ds *DataStore) UserViewList(q *Query) ([]*UserView, error) {
	records := make([]*UserView, 0)
	err := ds._select(&records, q)
	return records, err
}

// Private

// Find finds an object by ID
func (ds *DataStore) find(obj Model, id int64) error {
	orm := ds.ctx.DB.Model(obj)
	err := orm.Where("id = ?", id).Select()
	if err != nil {
		return err
	}
	return obj.Authorize(ds.actingUser, constants.ActionRead)
}

func (ds *DataStore) _select(models interface{}, q *Query) error {
	orm := ds.ctx.DB.Model(models)
	if !common.ListIsEmpty(q.GetColumns()) {
		orm.Column(q.GetColumns()...)
	}
	// Empty where clause causes orm to generate empty parens -> ()
	// which causes a SQL error. Include where only if non-empty.
	if q.WhereClause() != "" {
		orm.Where(q.WhereClause(), q.Params()...)
	}
	if q.GetOrderBy() != "" {
		orm.Order(q.GetOrderBy())
	}
	if q.GetLimit() > 0 {
		orm.Limit(q.GetLimit())
	}
	if q.GetOffset() >= 0 {
		orm.Offset(q.GetOffset())
	}
	return orm.Select()
}

func (ds *DataStore) save(model Model) error {
	if model.IsReadOnly() {
		return common.ErrNotSupported
	}
	if model.GetID() > 0 {
		return ds.update(model)
	}
	return ds.insert(model)
}

func (ds *DataStore) insert(model Model) error {
	err := model.Authorize(ds.actingUser, constants.ActionCreate)
	if err != nil {
		return err
	}
	model.SetTimestamps()

	err = model.BeforeSave()
	if err != nil {
		return err
	}
	db := ds.ctx.DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).Insert()
		ds.ctx.Log.Error().Msgf("Transaction on ID %d: %v", model.GetID(), err)
		return err
	})
}

func (ds *DataStore) update(model Model) error {
	if model.IsReadOnly() || model.UpdateIsForbidden() {
		return common.ErrNotSupported
	}
	err := model.Authorize(ds.actingUser, constants.ActionUpdate)
	if err != nil {
		return err
	}
	model.SetTimestamps()

	err = model.BeforeSave()
	if err != nil {
		return err
	}
	db := ds.ctx.DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Update()
		return err
	})
}

func (ds *DataStore) delete(model Model) error {
	if model.SupportsSoftDelete() {
		return ds.softDelete(model)
	}
	return ds.hardDelete(model)
}

func (ds *DataStore) checkDeleteAllowed(model Model) error {
	if model.IsReadOnly() || model.DeleteIsForbidden() {
		return common.ErrNotSupported
	}
	err := model.Authorize(ds.actingUser, constants.ActionDelete)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) softDelete(model Model) error {
	if err := ds.checkDeleteAllowed(model); err != nil {
		return err
	}
	model.SetSoftDeleteAttributes(ds.actingUser)
	model.SetTimestamps()
	db := ds.ctx.DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Update()
		return err
	})
}

func (ds *DataStore) hardDelete(model Model) error {
	if err := ds.checkDeleteAllowed(model); err != nil {
		return err
	}
	db := ds.ctx.DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		_, err := db.Model(model).WherePK().Delete()
		return err
	})
}

func (ds *DataStore) undelete(model Model) error {
	if model.IsReadOnly() || model.UpdateIsForbidden() {
		return common.ErrNotSupported
	}
	err := model.Authorize(ds.actingUser, constants.ActionUpdate)
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

// func Int64Value(obj interface{}, fieldName string) int64 {
// 	value := reflect.ValueOf(obj)
// 	if value.Type().Kind() != reflect.Ptr {
// 		value = reflect.New(reflect.TypeOf(obj))
// 	}
// 	field := value.Elem().FieldByName(fieldName)
// 	if field.IsValid() {
// 		return field.Int()
// 	}
// 	return int64(0)
// }

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
