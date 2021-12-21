package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileIsGlacierOnly(t *testing.T) {
	gf := &pgmodels.GenericFile{}
	for _, option := range constants.GlacierOnlyOptions {
		gf.StorageOption = option
		assert.True(t, gf.IsGlacierOnly())
	}
	gf.StorageOption = constants.StorageOptionStandard
	assert.False(t, gf.IsGlacierOnly())
}

func TestIdForFileIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForFileIdentifier("institution2.edu/toads/toad19")
	require.Nil(t, err)
	assert.Equal(t, int64(42), id)

	id, err = pgmodels.IdForFileIdentifier("institution1.edu/gl-or/file1.epub")
	require.Nil(t, err)
	assert.Equal(t, int64(53), id)

	id, err = pgmodels.IdForFileIdentifier("bad identifier")
	require.NotNil(t, err)
}

func TestGenericFileByID(t *testing.T) {
	gf, err := pgmodels.GenericFileByID(1)
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, int64(1), gf.ID)

	gf, err = pgmodels.GenericFileByID(-999)
	assert.NotNil(t, err)
	assert.Nil(t, gf)
}

func TestGenericFileByIdentifier(t *testing.T) {
	gf, err := pgmodels.GenericFileByIdentifier("institution2.edu/toads/toad16")
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, "institution2.edu/toads/toad16", gf.Identifier)

	gf, err = pgmodels.GenericFileByIdentifier("-- does not exist --")
	assert.NotNil(t, err)
	assert.Nil(t, gf)
}

func TestGenericFileSelect(t *testing.T) {

}

func TestFileSaveGetUpdate(t *testing.T) {

}

func TestFileSelect(t *testing.T) {

}

func TestFileValidate(t *testing.T) {
	gf := &pgmodels.GenericFile{
		Size: -1,
	}
	err := gf.Validate()
	require.NotNil(t, err)

	assert.Equal(t, "FileFormat is required", err.Errors["FileFormat"])
	assert.Equal(t, "Identifier is required", err.Errors["Identifier"])
	assert.Equal(t, pgmodels.ErrInstState, err.Errors["State"])
	assert.Equal(t, "Size cannot be negative", err.Errors["Size"])
	assert.Equal(t, "Invalid institution id", err.Errors["InstitutionID"])
	assert.Equal(t, "Intellectual object ID is required", err.Errors["IntellectualObjectID"])
	assert.Equal(t, "Invalid storage option", err.Errors["StorageOption"])
	assert.Equal(t, "Valid UUID required", err.Errors["UUID"])

	gf.Size = 20
	gf.FileFormat = "text/html"
	gf.Identifier = "test.edu/some-html-file"
	gf.State = constants.StateActive
	gf.InstitutionID = InstTest
	gf.IntellectualObjectID = 20
	gf.StorageOption = constants.StorageOptionGlacierVA
	gf.UUID = "c464d6dd-9fa6-41d9-8cb5-cdc7b986d07d"

	err = gf.Validate()
	assert.Nil(t, err)
}

func TestObjectFileCount(t *testing.T) {

}

func TestObjectFiles(t *testing.T) {

}

func TestFileActiveDeletionWorkItem(t *testing.T) {

}

func TestFileDeletionRequest(t *testing.T) {

}

func TestFileAssertDeletionPreconditions(t *testing.T) {

}

func TestNewDeletionEvent(t *testing.T) {

}

func TestGenericFileDelete(t *testing.T) {

}
