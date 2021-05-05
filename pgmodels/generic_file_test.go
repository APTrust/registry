package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
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
