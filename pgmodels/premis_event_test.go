package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectEventCount(t *testing.T) {
	db.LoadFixtures()
	count, err := pgmodels.ObjectEventCount(3)
	require.Nil(t, err)
	assert.Equal(t, 1, count)

	// No events for non-existent object
	count, err = pgmodels.ObjectEventCount(0)
	require.Nil(t, err)
	assert.Equal(t, 0, count)

}
