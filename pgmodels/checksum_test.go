package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksumByID(t *testing.T) {
	cs, err := pgmodels.ChecksumByID(1)
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.Equal(t, int64(1), cs.ID)
}

func TestChecksumGet(t *testing.T) {
	query := pgmodels.NewQuery().
		Where("algorithm", "=", constants.AlgSha256).
		Offset(0).
		Limit(1)
	cs, err := pgmodels.ChecksumGet(query)
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.Equal(t, constants.AlgSha256, cs.Algorithm)
}

func TestChecksumSelect(t *testing.T) {
	query := pgmodels.NewQuery().
		Where("algorithm", "=", constants.AlgSha256).
		Offset(0).
		Limit(4)
	checksums, err := pgmodels.ChecksumSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, checksums)
	assert.Equal(t, 4, len(checksums))
	for _, cs := range checksums {
		assert.Equal(t, constants.AlgSha256, cs.Algorithm)
	}
}

func TestChecksumGetID(t *testing.T) {
	cs := &pgmodels.Checksum{}
	assert.EqualValues(t, 0, cs.GetID())
	cs.ID = 19
	assert.EqualValues(t, 19, cs.GetID())
}

func TestChecksumSave(t *testing.T) {
	cs := &pgmodels.Checksum{}
	err := cs.Save()
	require.NotNil(t, err)
	valErr, ok := err.(*common.ValidationError)
	require.True(t, ok)
	testChecksumValidationErrors(t, valErr)

	cs = pgmodels.RandomChecksum(constants.AlgMd5)
	cs.GenericFileID = 1
	err = cs.Save()
	require.Nil(t, err)
	assert.True(t, cs.ID > int64(0))
}

func TestChecksumValidate(t *testing.T) {
	cs := &pgmodels.Checksum{}
	valErr := cs.Validate()
	require.NotNil(t, valErr)
	testChecksumValidationErrors(t, valErr)
}

func testChecksumValidationErrors(t *testing.T, valErr *common.ValidationError) {
	assert.Equal(t, "Checksum Digest cannot be empty", valErr.Errors["Digest"])
	assert.Equal(t, "Checksum DateTime is required", valErr.Errors["DateTime"])
	assert.Equal(t, "Checksum requires a valid algorithm", valErr.Errors["Algorithm"])
	assert.Equal(t, "Checksum requires a valid file id", valErr.Errors["GenericFileID"])
}
