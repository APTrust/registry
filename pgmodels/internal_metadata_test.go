package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalMetadata(t *testing.T) {

	// Test create and save
	im := pgmodels.NewInteralMetadata("Test Key 1", "Test Value 1")
	err := im.Save()
	require.Nil(t, err)

	// Test by key - not exists
	im, err = pgmodels.InternalMetadataByKey("does not exist")
	assert.NotNil(t, err)
	assert.Empty(t, im.ID)

	// Test by key - exists
	im, err = pgmodels.InternalMetadataByKey("Test Key 1")
	assert.Nil(t, err)
	assert.NotEmpty(t, im.ID)
	assert.NotEmpty(t, im.CreatedAt)
	assert.NotEmpty(t, im.UpdatedAt)
	assert.Equal(t, "Test Key 1", im.Key)
	assert.Equal(t, "Test Value 1", im.Value)
	originalID := im.ID

	// Test update
	im.Value = "Test Value 1 - Updated"
	assert.Nil(t, im.Save())

	im, err = pgmodels.InternalMetadataByKey("Test Key 1")
	assert.Nil(t, err)
	assert.Equal(t, "Test Value 1 - Updated", im.Value)
	assert.Equal(t, originalID, im.ID)

	// Test validation
	emptyRecord := pgmodels.NewInteralMetadata("", "")
	err = emptyRecord.Save()
	require.NotNil(t, err)

	valErr := err.(*common.ValidationError)
	assert.Equal(t, "InternalMetadata key cannot be empty", valErr.Errors["Key"])
	assert.Equal(t, "InternalMetadata value cannot be empty", valErr.Errors["Value"])

	// Delete should throw error, because it's not allowed
	assert.Error(t, im.Delete())

	// Test select
	query := pgmodels.NewQuery().IsNotNull("value")
	records, err := pgmodels.InternalMetadataSelect(query)
	require.Nil(t, err)
	assert.NotEmpty(t, records)
}
