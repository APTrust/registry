package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageRecordByID(t *testing.T) {
	sr, err := pgmodels.StorageRecordByID(1)
	require.Nil(t, err)
	require.NotNil(t, sr)
	assert.EqualValues(t, 1, sr.ID)
}

func TestStorageRecordGet(t *testing.T) {
	query := pgmodels.NewQuery().
		Where("generic_file_id", "=", 1).
		Limit(1)
	sr, err := pgmodels.StorageRecordGet(query)
	require.Nil(t, err)
	require.NotNil(t, sr)
	assert.EqualValues(t, 1, sr.GenericFileID)
}

func TestStorageRecordSelect(t *testing.T) {
	query := pgmodels.NewQuery().
		Where("generic_file_id", "=", 1)
	records, err := pgmodels.StorageRecordSelect(query)
	require.Nil(t, err)
	require.Equal(t, 2, len(records))
	for _, sr := range records {
		assert.EqualValues(t, 1, sr.GenericFileID)
	}
}

func TestStorageRecordGetID(t *testing.T) {
	sr := &pgmodels.StorageRecord{}
	assert.EqualValues(t, 0, sr.GetID())
	sr.ID = 916
	assert.EqualValues(t, 916, sr.GetID())
}

func TestStorageRecordSave(t *testing.T) {
	sr := &pgmodels.StorageRecord{}
	err := sr.Save()
	require.NotNil(t, err)
	valErr, ok := err.(*common.ValidationError)
	require.True(t, ok)
	testStorageRecordValError(t, valErr)

	sr.GenericFileID = 1
	sr.URL = "https://example.com/s3/54321"
	err = sr.Save()
	require.Nil(t, err)
	assert.True(t, sr.ID > 1)
}

func TestStorageRecordValidate(t *testing.T) {
	sr := &pgmodels.StorageRecord{}
	valErr := sr.Validate()
	testStorageRecordValError(t, valErr)
}

func testStorageRecordValError(t *testing.T, valErr *common.ValidationError) {
	assert.Equal(t, "Storage record requires a valid file id", valErr.Errors["GenericFileID"])
	assert.Equal(t, "Storage record requires a valid URL", valErr.Errors["URL"])
}
