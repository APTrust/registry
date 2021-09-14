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

func TestIdForEventIdentifier(t *testing.T) {
	db.LoadFixtures()
	id, err := pgmodels.IdForEventIdentifier("77e16041-8887-4739-af04-9d35e5cab4dc")
	require.Nil(t, err)
	assert.Equal(t, int64(49), id)

	id, err = pgmodels.IdForEventIdentifier("274e230a-dc6b-48a1-a96c-709e5728632b")
	require.Nil(t, err)
	assert.Equal(t, int64(35), id)

	id, err = pgmodels.IdForEventIdentifier("bad identifier")
	require.NotNil(t, err)
}
