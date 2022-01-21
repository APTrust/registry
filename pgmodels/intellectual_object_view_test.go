package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntellectualObjectView(t *testing.T) {
	db.LoadFixtures()
	objView, err := pgmodels.IntellectualObjectViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, objView)

	objView, err = pgmodels.IntellectualObjectViewByIdentifier(objView.Identifier)
	require.Nil(t, err)
	require.NotNil(t, objView)

	assert.Equal(t, int64(1), objView.GetID())

	query := pgmodels.NewQuery().
		Where("state", "=", constants.StateActive).
		OrderBy("created_at", "asc").
		Limit(10)
	objViews, err := pgmodels.IntellectualObjectViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 10, len(objViews))
}
