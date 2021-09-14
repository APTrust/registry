package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjIsGlacierOnly(t *testing.T) {
	obj := &pgmodels.IntellectualObject{}
	for _, option := range constants.GlacierOnlyOptions {
		obj.StorageOption = option
		assert.True(t, obj.IsGlacierOnly())
	}
	obj.StorageOption = constants.StorageOptionStandard
	assert.False(t, obj.IsGlacierOnly())
}

func TestIdForObjIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForObjIdentifier("institution1.edu/photos")
	require.Nil(t, err)
	assert.Equal(t, int64(1), id)

	id, err = pgmodels.IdForObjIdentifier("institution2.edu/coal")
	require.Nil(t, err)
	assert.Equal(t, int64(5), id)

	id, err = pgmodels.IdForObjIdentifier("bad identifier")
	require.NotNil(t, err)
}
