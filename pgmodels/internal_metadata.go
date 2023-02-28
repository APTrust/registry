package pgmodels

import (
	"github.com/APTrust/registry/common"
)

// InternalMetadata contains metadata about the state of the
// database. This is a holdover from Rails and ActiveRecord
// which turns out to be useful for housekeeping.
//
// We primarily use it as a kind of locking mechanism to ensure
// that the database's internal processes don't overlap.
// Internal processes set a key and value when they're running,
// then clear the value when they're done.
//
// For example, you'll see pairs like "update_counts is running"
// = true when the long-running update_counts function is doing
// its work. This prevents a second process from starting
// update_counts before the last process has completed, thus
// avoiding deadlock.
//
// Registry uses this table when running restoration spot tests,
// which should run once per day. While we may have two or more
// Registry instances running at once, we want to make sure
// restoration spot tests run **only once** per day. This table
// helps with that.
type InternalMetadata struct {
	tableName struct{} `pg:"ar_internal_metadata"`
	TimestampModel
	Key   string
	Value string
}

// NewInternalMetadata creates a new InternalMetadata record.
func NewInteralMetadata(key, value string) *InternalMetadata {
	return &InternalMetadata{
		Key:   key,
		Value: value,
	}
}

// InternalMetadataByKey returns the file with the specified key.
// The key column has a unique constraint, so this should return
// one record, max. Returns pg.ErrNoRows if there is no match.
func InternalMetadataByKey(key string) (*InternalMetadata, error) {
	query := NewQuery().Where(`"internal_metadata"."key"`, "=", key)
	return InternalMetadataGet(query)
}

// InternalMetadataGet returns the first file matching the query.
func InternalMetadataGet(query *Query) (*InternalMetadata, error) {
	var im InternalMetadata
	err := query.Select(&im)
	return &im, err
}

// InternalMetadataSelect returns all files matching the query.
func InternalMetadataSelect(query *Query) ([]*InternalMetadata, error) {
	var records []*InternalMetadata
	err := query.Select(&records)
	return records, err
}

// Save saves this record to the database. This will peform an insert
// if InternalMetadata.ID is zero or an update otherwise.
func (im *InternalMetadata) Save() error {
	err := im.Validate()
	if err != nil {
		return err
	}
	im.SetTimestamps()
	if im.ID == int64(0) {
		return insert(im)
	}
	return update(im)
}

// Delete is not supported because these records are required
// for internal housekeeping. This method will always return
// an error. We define to tell developers explicitly not to do
// this.
func (im *InternalMetadata) Delete() error {
	return common.ErrNotSupported
}

func (im *InternalMetadata) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(im.Key) {
		errors["Key"] = "InternalMetadata key cannot be empty"
	}
	if common.IsEmptyString(im.Value) {
		errors["Value"] = "InternalMetadata value cannot be empty"
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
