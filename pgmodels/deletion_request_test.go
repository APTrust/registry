package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ConfTokenPlain     = "ConfirmationToken"
	ConfTokenEncrypted = "$2a$10$TK8s1XnmWulSUdze8GN5uOgGmDDsnndQKF5/Rz1j0xaHT7AwXRVma"
)

func TestNewDeletionRequest(t *testing.T) {
	request, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, request)
	assert.Equal(t, 32, len(request.ConfirmationToken))
	assert.False(t, common.LooksEncrypted(request.ConfirmationToken))

	assert.True(t, len(request.EncryptedConfirmationToken) > 30)

	assert.True(t, common.LooksEncrypted(request.EncryptedConfirmationToken))
}

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

func TestDeletionRequestSelect(t *testing.T) {
	query := pgmodels.NewQuery().
		IsNull("confirmed_at").
		OrderBy("id", "asc").
		Limit(2)
	requests, err := pgmodels.DeletionRequestSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, requests)
	assert.Equal(t, 2, len(requests))
	assert.EqualValues(t, 1, requests[0].ID)
	assert.EqualValues(t, 3, requests[1].ID)
}

func TestDeletionRequestValidate(t *testing.T) {
	err := db.LoadFixtures()
	require.Nil(t, err)

	// These IDS are fixed in the fixture data under db/fixtures/users.sql
	//inst1AdminID := int64(2)
	inst1UserID := int64(3)
	inst2AdminID := int64(5)
	inst1ID := int64(2)

	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)
	valErr := req.Validate()
	require.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionInstitutionID, valErr.Errors["InstitutionID"])
	// assert.Equal(t, pgmodels.ErrDeletionRequesterID, valErr.Errors["RequestedByID"])

	// START HERE

	// ----------------------------------------------------------
	// TODO: Fix this. We're getting the wrong error here
	// ----------------------------------------------------------
	// Invalid requester, approver, canceller
	// These users don't exist
	// req.InstitutionID = inst1UserID
	// req.RequestedByID = 559988
	// req.ConfirmedByID = 331149
	// req.CancelledByID = 808080
	// valErr = req.Validate()
	// require.NotNil(t, valErr)
	// assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["RequestedByID"])
	// assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["ConfirmedByID"])
	// assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["CancelledByID"])

	// User role cannot approve or cancel deletions
	req.InstitutionID = inst1ID
	req.RequestedByID = inst1UserID
	req.ConfirmedByID = inst1UserID
	req.CancelledByID = inst1UserID
	valErr = req.Validate()
	require.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionWrongRole, valErr.Errors["ConfirmedByID"])
	assert.Equal(t, pgmodels.ErrDeletionWrongRole, valErr.Errors["CancelledByID"])

	// Admin can approve & cancel deletions, but only at their
	// own institution.
	req.InstitutionID = inst1ID
	req.RequestedByID = inst2AdminID
	req.ConfirmedByID = inst2AdminID
	req.CancelledByID = inst2AdminID
	valErr = req.Validate()
	require.NotNil(t, valErr)
	//assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["RequestedByID"])
	assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["ConfirmedByID"])
	assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["CancelledByID"])

	// ----------------------------------------------------------
	// TODO: Fix this. User is coming back as nil in Validate,
	// and they shouldn't because the user is there.
	// ----------------------------------------------------------
	// This is OK, though technically, a deletion would not
	// be both confirmed and cancelled.
	// req.InstitutionID = inst1ID
	// req.RequestedByID = inst1AdminID
	// req.ConfirmedByID = inst1AdminID
	// req.CancelledByID = inst1AdminID
	// valErr = req.Validate()
	// assert.Nil(t, valErr)
}

func TestDeletionRequestAddFile(t *testing.T) {

}

func TestDeletionRequestAddObject(t *testing.T) {

}

func TestDeletionRequestFirstFile(t *testing.T) {

}

func TestDeletionRequestFirstObject(t *testing.T) {

}

func TestDeletionRequestConfirm(t *testing.T) {

}

func TestDeletionRequestCancel(t *testing.T) {

}
