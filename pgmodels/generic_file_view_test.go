package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericFileView(t *testing.T) {
	db.LoadFixtures()
	gfView, err := pgmodels.GenericFileViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, gfView)

	gfView, err = pgmodels.GenericFileViewByIdentifier(gfView.Identifier)
	require.Nil(t, err)
	require.NotNil(t, gfView)

	query := pgmodels.NewQuery().
		Where("state", "=", constants.StateActive).
		OrderBy("created_at", "asc").
		Limit(40)
	gfViews, err := pgmodels.GenericFileViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 40, len(gfViews))
}
