package helpers_test

import (
	"testing"

	"github.com/APTrust/registry/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCover(t *testing.T) {
	lastSource := ""
	coverChanged := false
	for i := 0; i < 12; i++ {
		photo := helpers.GetCover()
		require.NotNil(t, photo)
		if lastSource == "" {
			lastSource = photo.Source
			continue
		}
		if lastSource != photo.Source {
			coverChanged = true
			break
		}
	}
	assert.True(t, coverChanged)
}
