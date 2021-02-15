package models_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note that ds is defined in common_test.go and is created
// with SysAdmin user, so that ds has all privileges.

func TestChecksumFind(t *testing.T) {
	cs, err := ds.ChecksumFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.Equal(t, int64(1), cs.ID)
	assert.EqualValues(t, 1, cs.GenericFileID)
	assert.EqualValues(t, "md5", cs.Algorithm)
	assert.Equal(t, "12345678", cs.Digest)
}

func TestChecksumsList(t *testing.T) {
	query := models.NewQuery().Where("generic_file_id", "=", int64(21)).OrderBy("created_at desc").OrderBy("algorithm asc")
	checksums, err := ds.ChecksumList(query)
	require.Nil(t, err)
	require.NotEmpty(t, checksums)
	algs := []string{
		"md5",
		"sha1",
		"sha256",
		"sha512",
	}
	for i, cs := range checksums {
		assert.Equal(t, int64(21), cs.GenericFileID)
		assert.Equal(t, algs[i], checksums[i].Algorithm)
	}
}

func TestChecksumSave(t *testing.T) {
	cs := &models.Checksum{
		Algorithm:     constants.AlgMd5,
		DateTime:      TestDate,
		Digest:        "0987654321abcdef",
		GenericFileID: int64(20),
	}

	err := ds.ChecksumSave(cs)
	require.Nil(t, err)
	require.NotNil(t, cs)
	assert.True(t, cs.ID > int64(0))
	assert.EqualValues(t, 20, cs.GenericFileID)
	assert.EqualValues(t, constants.AlgMd5, cs.Algorithm)
	assert.Equal(t, "0987654321abcdef", cs.Digest)

	// We should get an error here because we're not allowed
	// to update existing checksums.
	cs.Digest = "----------------"
	err = ds.ChecksumSave(cs)
	require.NotNil(t, err)
	assert.Equal(t, common.ErrNotSupported, err)
}

func TestGenericFileSaveDeleteUndelete(t *testing.T) {
	gf := &models.GenericFile{
		FileFormat:           "text/plain",
		Size:                 int64(400),
		Identifier:           "institution2.edu/toads/test-file.txt",
		IntellectualObjectID: int64(6),
		State:                "A",
		LastFixityCheck:      TestDate,
		InstitutionID:        int64(3),
		StorageOption:        constants.StorageOptionStandard,
		UUID:                 "811b9a46-f91f-4379-a2f7-7b1bc8125a7c",
	}

	err := ds.GenericFileSave(gf)
	require.Nil(t, err)
	assert.True(t, gf.ID > int64(0))
	assert.Equal(t, "A", gf.State)
	assert.False(t, gf.CreatedAt.IsZero())
	assert.False(t, gf.UpdatedAt.IsZero())

	err = ds.GenericFileDelete(gf)
	require.Nil(t, err)
	assert.Equal(t, "D", gf.State)

	err = ds.GenericFileUndelete(gf)
	require.Nil(t, err)
	assert.Equal(t, "A", gf.State)
}

func TestGenericFileFind(t *testing.T) {
	gf, err := ds.GenericFileFind(int64(1))
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)
	assert.Equal(t, "institution1.edu/photos/picture1", gf.Identifier)
	assert.Equal(t, int64(48771), gf.Size)
}

func TestGenericFileFindByIdentifier(t *testing.T) {
	gf, err := ds.GenericFileFindByIdentifier("institution1.edu/photos/picture1")
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)
	assert.Equal(t, "institution1.edu/photos/picture1", gf.Identifier)
	assert.Equal(t, int64(48771), gf.Size)
}

func TestGenericFileList(t *testing.T) {

}
