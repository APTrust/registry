package pgmodels_test

import (
	"testing"

	//"github.com/APTrust/registry/constants"
	//"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ConfTokenPlain       = "ConfirmationToken"
	CancelTokenPlain     = "CancelToken"
	ConfTokenEncrypted   = "$2a$10$TK8s1XnmWulSUdze8GN5uOgGmDDsnndQKF5/Rz1j0xaHT7AwXRVma"
	CancelTokenEncrypted = "$2a$10$xwxTFn.k1TbfbNSW3/udduwtjwo7nQSBlIlARHvTXADAhCfQtZt46"
)

func TestDeletionRequestByID(t *testing.T) {
	request, err := pgmodels.DeletionRequestByID(1)
	require.Nil(t, err)
	require.NotNil(t, request)
	assert.EqualValues(t, 2, request.InstitutionID)
	assert.NotEmpty(t, request.RequestedAt)
	require.NotNil(t, request.RequestedBy)
	assert.Equal(t, "user@inst1.edu", request.RequestedBy.Email)

	assert.Nil(t, request.ConfirmedBy)
	assert.Empty(t, request.ConfirmedAt)
	assert.Nil(t, request.CancelledBy)
	assert.Empty(t, request.CancelledAt)

	require.NotNil(t, request.GenericFiles)
	require.Equal(t, 3, len(request.GenericFiles))

	for _, gf := range request.GenericFiles {
		assert.True(t, (gf.ID == 7 || gf.ID == 8 || gf.ID == 9))
	}

	// Check a confirmed request with objects
	request, err = pgmodels.DeletionRequestByID(2)
	require.Nil(t, err)
	require.NotNil(t, request)
	assert.NotNil(t, request.ConfirmedBy)
	assert.Equal(t, 1, len(request.IntellectualObjects))

	// Check a cancelled request with objects
	request, err = pgmodels.DeletionRequestByID(3)
	require.Nil(t, err)
	require.NotNil(t, request)
	assert.NotNil(t, request.CancelledBy)
	assert.Equal(t, 1, len(request.IntellectualObjects))
}
