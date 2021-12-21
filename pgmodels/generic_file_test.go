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

func TestFileSaveGetUpdate(t *testing.T) {

}

func TestFileSelect(t *testing.T) {

}

func TestFileValidate(t *testing.T) {

}

func TestObjectFileCount(t *testing.T) {

}

func TestObjectFiles(t *testing.T) {

}

func TestAssertDeletionPreconditions(t *testing.T) {

}

func TestNewDeletionEvent(t *testing.T) {

}

func TestGenericFileDelete(t *testing.T) {

}
