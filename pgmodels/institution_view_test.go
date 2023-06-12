package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstitutionView(t *testing.T) {
	db.ForceFixtureReload()
	instView, err := pgmodels.InstitutionViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, instView)

	instView, err = pgmodels.InstitutionViewByIdentifier(instView.Identifier)
	require.Nil(t, err)
	require.NotNil(t, instView)

	query := pgmodels.NewQuery().
		Where("state", "=", constants.StateActive).
		OrderBy("created_at", "asc")
	instViews, err := pgmodels.InstitutionViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 5, len(instViews))
}

func TestInstViewDisplayType(t *testing.T) {
	inst := &pgmodels.InstitutionView{
		Type: constants.InstTypeMember,
	}
	assert.Equal(t, "Member", inst.DisplayType())

	inst.Type = constants.InstTypeSubscriber
	assert.Equal(t, "Associate", inst.DisplayType())
}
