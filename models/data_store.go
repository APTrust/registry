package models

import (
	"reflect"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
)

// DataStore provides methods for storing and retrieving data from our
// Postgres database. In addition to basic storage and retrieval, this class
// enforces some business rules, including the following:
//
// 1. Restricts access to data that the acting user is not allowed to see.
//
// 2. Restricts operations that the acting user is not allowed to peform,
//    such as updating or deleting records that belong to someone else.
//
// 3. Restricts operations that no one (including admin) is allowed to
//    perform. For example, records such as Checksums and PremisEvents
//    can never be altered or deleted by anyone.
type DataStore struct {
	ctx        *common.APTContext
	actingUser *User
}

// NewDataStore creates a new DataStore object for the actingUser. Note
// that this is a lightweight operation that uses existing DB connections
// from the connection pool. This class will check the acting user's
// permissions for all operations, and will not create, update, delete, or
// return data unless the acting user is authorized to access it.
func NewDataStore(actingUser *User) *DataStore {
	return &DataStore{
		ctx:        common.Context(),
		actingUser: actingUser,
	}
}

// ChecksumFind returns the Checksum with the specified ID.
func (ds *DataStore) ChecksumFind(id int64) (*Checksum, error) {
	checksum := &Checksum{}
	err := ds.find(checksum, id)
	return checksum, err
}

// ChecksumList returns a list of checksums matching the query.
func (ds *DataStore) ChecksumList(q *Query) ([]*Checksum, error) {
	checksums := make([]*Checksum, 0)
	err := ds._select(&checksums, q)
	return checksums, err
}

// ChecksumSave saves a new checksum. If you try to update an existing
// checksum, you'll get an error, because business rules say these
// records are immutable.
func (ds *DataStore) ChecksumSave(cs *Checksum) error {
	return ds.save(cs)
}

// GenericFileDelete sets the State attribute of the GenericFile to "D".
// This is a soft delete.
func (ds *DataStore) GenericFileDelete(gf *GenericFile) error {
	return ds.delete(gf)
}

// GenericFileFind returns the GenericFile with the specified ID.
func (ds *DataStore) GenericFileFind(id int64) (*GenericFile, error) {
	gf := &GenericFile{}
	err := ds.find(gf, id)
	return gf, err
}

// GenericFileFindByIdentifier returns the GenericFile with the specified
// identifier.
func (ds *DataStore) GenericFileFindByIdentifier(identifier string) (*GenericFile, error) {
	gf := &GenericFile{}
	query := NewQuery().Where("identifier", "=", identifier).Limit(1)
	err := ds._select(gf, query)
	return gf, err
}

// GenericFileList returns the GenericFiles matching the specified query.
func (ds *DataStore) GenericFileList(q *Query) ([]*GenericFile, error) {
	gfs := make([]*GenericFile, 0)
	err := ds._select(&gfs, q)
	return gfs, err
}

// GenericFileSave saves a GenericFile.
func (ds *DataStore) GenericFileSave(gf *GenericFile) error {
	return ds.save(gf)
}

// GenericFileUndelete sets the State attribute of the GenericFile to "A".
// This undoes the soft delete.
func (ds *DataStore) GenericFileUnelete(gf *GenericFile) error {
	return ds.undelete(gf)
}

// InstitutionDelete sets the DeactivatedAt timestamp on an Institution,
// marking it as no longer active. This is a soft delete.
func (ds *DataStore) InstitutionDelete(inst *Institution) error {
	return ds.delete(inst)
}

// InstitutionFind returns the Institution with the specified ID.
func (ds *DataStore) InstitutionFind(id int64) (*Institution, error) {
	inst := &Institution{}
	err := ds.find(inst, id)
	return inst, err
}

// InstitutionFindByIdentifier returns the Institution with the specified
// identifier.
func (ds *DataStore) InstitutionFindByIdentifier(identifier string) (*Institution, error) {
	inst := &Institution{}
	query := NewQuery().Where("identifier", "=", identifier).Limit(1)
	err := ds._select(inst, query)
	return inst, err
}

// InstitutionList returns the Institutions matching the query.
func (ds *DataStore) InstitutionList(q *Query) ([]*Institution, error) {
	institutions := make([]*Institution, 0)
	err := ds._select(&institutions, q)
	return institutions, err
}

// InstitutionSave saves an Institution record.
func (ds *DataStore) InstitutionSave(inst *Institution) error {
	return ds.save(inst)
}

// InstitutionUnelete clears the DeactivatedAt timestamp on an Institution,
// marking it as once again active. This undoes soft delete.
func (ds *DataStore) InstitutionUndelete(inst *Institution) error {
	return ds.undelete(inst)
}

// IntellectualObjectDelete marks an object as deleted by setting its
// State to "D". This is a soft delete.
func (ds *DataStore) IntellectualObjectDelete(obj *IntellectualObject) error {
	return ds.delete(obj)
}

// IntellectualObjectFind returns the IntellectualObject with the specified
// ID.
func (ds *DataStore) IntellectualObjectFind(id int64) (*IntellectualObject, error) {
	obj := &IntellectualObject{}
	err := ds.find(obj, id)
	return obj, err
}

// IntellectualObjectFindByIdentifier returns the IntellectualObject
// with the specified identifier.
func (ds *DataStore) IntellectualObjectFindByIdentifier(identifier string) (*IntellectualObject, error) {
	obj := &IntellectualObject{}
	query := NewQuery().Where("identifier", "=", identifier).Limit(1)
	err := ds._select(obj, query)
	return obj, err
}

// IntellectualObjectList returns IntellectualObjects that match the query.
func (ds *DataStore) IntellectualObjectList(q *Query) ([]*IntellectualObject, error) {
	objitutions := make([]*IntellectualObject, 0)
	err := ds._select(&objitutions, q)
	return objitutions, err
}

// IntellectualObjectSave inserts or updates an IntellectualObject.
func (ds *DataStore) IntellectualObjectSave(obj *IntellectualObject) error {
	return ds.save(obj)
}

// IntellectualObjectUndelete marks an object as active by setting its
// State to "A". This undoes soft delete.
func (ds *DataStore) IntellectualObjectUndelete(obj *IntellectualObject) error {
	return ds.undelete(obj)
}

// PremisEventFind returns the PremisEvent with the specified ID.
func (ds *DataStore) PremisEventFind(id int64) (*PremisEvent, error) {
	event := &PremisEvent{}
	err := ds.find(event, id)
	return event, err
}

// PremisEventFindByIdentifier returns the PremisEvent with the specified
// identifier.
func (ds *DataStore) PremisEventFindByIdentifier(identifier string) (*PremisEvent, error) {
	event := &PremisEvent{}
	query := NewQuery().Where("identifier", "=", identifier).Limit(1)
	err := ds._select(event, query)
	return event, err
}

// PremisEventList returns the PremisEvents that match the query.
func (ds *DataStore) PremisEventList(q *Query) ([]*PremisEvent, error) {
	events := make([]*PremisEvent, 0)
	err := ds._select(&events, q)
	return events, err
}

// PremisEventSave saves a new PremisEvent. Attempting to update an existing
// PremisEvent will cause an error as events are read-only.
func (ds *DataStore) PremisEventSave(event *PremisEvent) error {
	return ds.save(event)
}

// StorageRecordFind returns the StorageRecord with the specified ID.
func (ds *DataStore) StorageRecordFind(id int64) (*StorageRecord, error) {
	sr := &StorageRecord{}
	err := ds.find(sr, id)
	return sr, err
}

// StorageRecordList returns a list of StorageRecords matching the query.
func (ds *DataStore) StorageRecordList(q *Query) ([]*StorageRecord, error) {
	records := make([]*StorageRecord, 0)
	err := ds._select(&records, q)
	return records, err
}

// StorageRecordSave saves a new StorageRecord or updates an existing one.
func (ds *DataStore) StorageRecordSave(sr *StorageRecord) error {
	return ds.save(sr)
}

// StorageRecordDelete deletes a StorageRecord. Note that this is a hard
// delete and cannot be undone.
func (ds *DataStore) StorageRecordDelete(sr *StorageRecord) error {
	return ds.delete(sr)
}

// UserDelete sets the DeactivatedAt timestamp on a User to indicate their
// account is no longer active. This is a soft delete and can be undone later.
func (ds *DataStore) UserDelete(user *User) error {
	return ds.delete(user)
}

// UserFind returns the User with the specified ID. The User record will
// include the related Instution record.
func (ds *DataStore) UserFind(id int64) (*User, error) {
	user := &User{}
	err := ds.ctx.DB.Model(user).Relation("Institution").Where(`"user"."id" = ?`, id).Select()
	if err != nil {
		return nil, err
	}
	err = user.Authorize(ds.actingUser, constants.ActionRead)
	return user, err
}

// UserFindByEmail returns the User with the specified email address.
// The User record will include the related Instution record.
func (ds *DataStore) UserFindByEmail(email string) (*User, error) {
	user := &User{}
	err := ds.ctx.DB.Model(user).Relation("Institution").Where(`"user"."email" = ?`, email).Select()
	if err != nil {
		return nil, err
	}
	err = user.Authorize(ds.actingUser, constants.ActionRead)
	return user, err
}

// UserList returns a list of Users matching the specified query.
func (ds *DataStore) UserList(q *Query) ([]*User, error) {
	users := make([]*User, 0)
	err := ds._select(&users, q)
	return users, err
}

// UserSave inserts or updates a User record.
func (ds *DataStore) UserSave(user *User) error {
	return ds.save(user)
}

// UserSignIn signs a user in. If successful, it returns the User
// record with User.Institution properly set. If it fails, check
// the error.
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

// UserSignOut signs a user out.
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

// UserUndelete cleans the DeactivatedAt timestamp on a user record to
// indicate that their account is active.
func (ds *DataStore) UserUndelete(user *User) error {
	return ds.undelete(user)
}

// UserViewList returns a list of UserView objects.
func (ds *DataStore) UserViewList(q *Query) ([]*UserView, error) {
	records := make([]*UserView, 0)
	err := ds._select(&records, q)
	return records, err
}

// WorkItemFind returns the WorkItem with the specified ID.
func (ds *DataStore) WorkItemFind(id int64) (*WorkItem, error) {
	item := &WorkItem{}
	err := ds.find(item, id)
	return item, err
}

// WorkItemList returns a list of WorkItems matching the query.
func (ds *DataStore) WorkItemList(q *Query) ([]*WorkItem, error) {
	items := make([]*WorkItem, 0)
	err := ds._select(&items, q)
	return items, err
}

// WorkItemSave inserts or updates an WorkItem.
func (ds *DataStore) WorkItemSave(item *WorkItem) error {
	return ds.save(item)
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
	for _, rel := range q.GetRelations() {
		orm.Relation(rel)
	}
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
		if err != nil {
			ds.ctx.Log.Error().Msgf("Transaction on ID %d: %v", model.GetID(), err)
		}
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
