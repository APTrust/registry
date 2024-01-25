package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserMin(t *testing.T) {
	db.LoadFixtures()
	user, err := pgmodels.UserByID(5)
	require.Nil(t, err)
	require.NotNil(t, user)

	userMin := user.ToMin()
	require.NotNil(t, userMin)
	assert.Equal(t, user.ID, userMin.ID)
	assert.Equal(t, user.Name, userMin.Name)
	assert.Equal(t, user.Email, userMin.Email)
}

func TestDeletionRequestMin(t *testing.T) {
	user, err := pgmodels.UserByID(5)
	require.Nil(t, err)
	require.NotNil(t, user)

	now := time.Now().UTC()
	req := &pgmodels.DeletionRequest{
		BaseModel: pgmodels.BaseModel{
			ID: 1234,
		},
		InstitutionID: 98,
		RequestedAt:   now,
		ConfirmedAt:   now,
		CancelledAt:   now,
		RequestedBy:   user,
		ConfirmedBy:   user,
		CancelledBy:   user,
		GenericFiles: []*pgmodels.GenericFile{
			pgmodels.RandomGenericFile(222, "object/id"),
			pgmodels.RandomGenericFile(222, "object/id"),
		},
		IntellectualObjects: []*pgmodels.IntellectualObject{
			pgmodels.RandomObject(),
			pgmodels.RandomObject(),
		},
		WorkItems: []*pgmodels.WorkItem{
			pgmodels.RandomWorkItem("object/id1", constants.ActionDelete, 888, 111),
			pgmodels.RandomWorkItem("object/id2", constants.ActionDelete, 999, 222),
		},
	}

	reqMin := req.ToMin()
	require.NotNil(t, reqMin)
}
