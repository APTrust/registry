package db_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/stretchr/testify/require"
)

// Make sure LoadFixtures returns no error.
func TestLoadFixtures(t *testing.T) {
	require.Nil(t, db.LoadFixtures())
}
