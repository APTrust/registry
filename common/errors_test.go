package common_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {
	valErr := common.NewValidationError()
	require.NotNil(t, valErr)
	require.NotNil(t, valErr.Errors)

	valErr.Errors["One"] = "First Error"
	valErr.Errors["Two"] = "Second Error"

	// Key order is not guaranteed, but string should
	// match one of these two.
	expected1 := "One: First Error\nTwo: Second Error"
	expected2 := "Two: Second Error\nOne: First Error"

	errStr := valErr.Error()

	assert.True(t, (errStr == expected1 || errStr == expected2))
}
