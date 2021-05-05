package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
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
