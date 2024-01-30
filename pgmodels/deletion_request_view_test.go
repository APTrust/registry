package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletionRequestView(t *testing.T) {
	db.LoadFixtures()
	requestView, err := pgmodels.DeletionRequestViewByID(1)
	require.Nil(t, err)
	require.NotNil(t, requestView)

	query := pgmodels.NewQuery().
		Where("institution_id", "=", 2)
	requests, err := pgmodels.DeletionRequestViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 3, len(requests))
}

func TestDeletionRequestViewDisplayStatus(t *testing.T) {
	req := &pgmodels.DeletionRequestView{}

	req.CancelledByID = 0
	req.ConfirmedByID = 0
	assert.Equal(t, "Awaiting Approval", req.DisplayStatus())

	req.CancelledByID = 0
	req.ConfirmedByID = 100
	assert.Equal(t, "Approved", req.DisplayStatus())

	req.ConfirmedByID = 0
	req.CancelledByID = 1000
	assert.Equal(t, "Rejected", req.DisplayStatus())
}
