package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
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

	for _, status := range constants.Statuses {
		req.CancelledByID = 0
		req.Status = status
		assert.Equal(t, status, req.DisplayStatus())

		req.CancelledByID = 1000
		assert.Equal(t, "Rejected", req.DisplayStatus())
	}
}
