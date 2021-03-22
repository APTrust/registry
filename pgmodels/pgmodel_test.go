package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
)

func TestInstIDFor(t *testing.T) {
	db.LoadFixtures()

	// This is all known fixture data...
	id, err := pgmodels.InstIDFor("Checksum", 8)
	assert.Nil(t, err)
	assert.EqualValues(t, 4, id)

	id, err = pgmodels.InstIDFor("GenericFile", 5)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, id)

	id, err = pgmodels.InstIDFor("Institution", 3)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, id)

	id, err = pgmodels.InstIDFor("IntellectualObject", 6)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, id)

	id, err = pgmodels.InstIDFor("PremisEvent", 8)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, id)

	id, err = pgmodels.InstIDFor("StorageRecord", 8)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, id)

	id, err = pgmodels.InstIDFor("User", 4)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, id)

	id, err = pgmodels.InstIDFor("WorkItem", 4)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, id)

	// Expected failure cases
	id, err = pgmodels.InstIDFor("WorkItem", -100)
	assert.NotNil(t, err)
	assert.True(t, pgmodels.IsNoRowError(err))
	assert.EqualValues(t, 0, id)

	id, err = pgmodels.InstIDFor("UnknownType", 1)
	assert.NotNil(t, err)
	assert.Equal(t, common.ErrInvalidParam, err)
	assert.EqualValues(t, 0, id)

}
