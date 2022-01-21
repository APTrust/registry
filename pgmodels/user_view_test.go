package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserView(t *testing.T) {
	db.LoadFixtures()
	userView, err := pgmodels.UserViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, userView)

	userView, err = pgmodels.UserViewByEmail(userView.Email)
	require.Nil(t, err)
	require.NotNil(t, userView)

	query := pgmodels.NewQuery().
		Where("institution_id", "=", 3).
		OrderBy("created_at", "asc")
	userViews, err := pgmodels.UserViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 2, len(userViews))
}
