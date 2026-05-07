package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPremisEventView(t *testing.T) {
	db.LoadFixtures()
	eventView, err := pgmodels.PremisEventViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, eventView)

	eventView, err = pgmodels.PremisEventViewByIdentifier(eventView.Identifier)
	require.Nil(t, err)
	require.NotNil(t, eventView)

	query := pgmodels.NewQuery().
		Where("intellectual_object_id", "=", 1).
		OrderBy("date_time", "asc")
	eventViews, err := pgmodels.PremisEventViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 13, len(eventViews))
}
