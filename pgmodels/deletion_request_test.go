package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/testutil"
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
	inst1AdminID := int64(2)
	inst1UserID := int64(3)
	inst2AdminID := int64(5)
	inst1ID := int64(2)

	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)
	valErr := req.Validate()
	require.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionInstitutionID, valErr.Errors["InstitutionID"])
	assert.Equal(t, pgmodels.ErrDeletionRequesterID, valErr.Errors["RequestedByID"])

	// Invalid requester, approver, canceller
	// These users don't exist
	req.InstitutionID = inst1UserID
	req.RequestedByID = 559988
	req.ConfirmedByID = 331149
	req.CancelledByID = 808080
	valErr = req.Validate()
	require.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["RequestedByID"])
	assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["ConfirmedByID"])
	assert.Equal(t, pgmodels.ErrDeletionUserNotFound, valErr.Errors["CancelledByID"])

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
	assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["RequestedByID"])
	assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["ConfirmedByID"])
	assert.Equal(t, pgmodels.ErrDeletionWrongInst, valErr.Errors["CancelledByID"])

	// This is OK, though technically, a deletion would not
	// be both confirmed and cancelled.
	req, err = pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)
	req.InstitutionID = inst1ID
	req.RequestedByID = inst1AdminID
	req.ConfirmedByID = inst1AdminID
	req.CancelledByID = inst1AdminID
	valErr = req.Validate()
	assert.Nil(t, valErr)

	// Make sure we catch illegal files and objects.
	// User can't delete file or object belonging
	// to another institution.
	req.GenericFiles = []*pgmodels.GenericFile{
		&pgmodels.GenericFile{
			InstitutionID: 9999,
		},
	}
	valErr = req.Validate()
	assert.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionIllegalObject, valErr.Errors["GenericFiles"])
	req.GenericFiles = []*pgmodels.GenericFile{}

	req.IntellectualObjects = []*pgmodels.IntellectualObject{
		&pgmodels.IntellectualObject{
			InstitutionID: 9999,
		},
	}
	valErr = req.Validate()
	assert.NotNil(t, valErr)
	assert.Equal(t, pgmodels.ErrDeletionIllegalObject, valErr.Errors["IntellectualObjects"])
}

func TestDeletionConfirmationOneAdmin(t *testing.T) {
	err := db.LoadFixtures()
	require.Nil(t, err)

	// These IDS are fixed in the fixture data under db/fixtures/users.sql
	inst1AdminID := int64(2)
	inst1ID := int64(2)

	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)

	// Admin can approve their own deletion, if they are the
	// only admin at the institution.
	req.InstitutionID = inst1ID
	req.RequestedByID = inst1AdminID
	req.ConfirmedByID = inst1AdminID
	valErr := req.Validate()
	require.Nil(t, valErr)

	// Now, if there is more than one admin at the institution,
	// the the admin cannot approve their own deletion request.
	// The other admin has to approve it.
	encPassword, _ := common.EncryptPassword("password")
	newAdmin := pgmodels.User{
		Name:              "Temp Inst 1 Admin",
		Email:             "admin_temp@inst1.edu",
		InstitutionID:     inst1ID,
		Role:              constants.RoleInstAdmin,
		EncryptedPassword: encPassword,
	}
	assert.Nil(t, newAdmin.Save())
	defer common.Context().DB.Exec("delete from users where email=?", newAdmin.Email)

	valErr = req.Validate()
	require.NotNil(t, valErr)
	assert.Equal(t, 1, len(valErr.Errors))
	assert.Equal(t, pgmodels.ErrDeletionBadAdmin, valErr.Errors["ConfirmedByID"])
}

func TestDeletionRequestAddFile(t *testing.T) {
	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)
	assert.Empty(t, req.GenericFiles)
	assert.Nil(t, req.FirstFile())

	for i := 0; i < 5; i++ {
		req.AddFile(pgmodels.RandomGenericFile(9999, "obj/identifier"))
	}
	assert.Equal(t, 5, len(req.GenericFiles))
	assert.Equal(t, req.GenericFiles[0], req.FirstFile())
}

func TestDeletionRequestAddObject(t *testing.T) {
	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)
	assert.Empty(t, req.IntellectualObjects)
	assert.Nil(t, req.FirstObject())

	for i := 0; i < 5; i++ {
		req.AddObject(pgmodels.RandomObject())
	}
	assert.Equal(t, 5, len(req.IntellectualObjects))
	assert.Equal(t, req.IntellectualObjects[0], req.FirstObject())
}

func TestDeletionRequestIncludesFilesAndObjects(t *testing.T) {
	defer db.ForceFixtureReload()

	// Create a deletion request
	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)

	// Add the requestor, institution id, etc.
	inst2Admin := testutil.InitUser(t, "admin@inst2.edu")
	req.RequestedBy = inst2Admin
	req.RequestedByID = inst2Admin.ID
	req.InstitutionID = inst2Admin.InstitutionID
	req.RequestedAt = time.Now().UTC()

	// Now add three objects and three files to the request
	for i := 0; i < 3; i++ {
		obj := pgmodels.RandomObject()
		obj.InstitutionID = inst2Admin.InstitutionID
		require.NoError(t, obj.Save())
		req.AddObject(obj)

		gf := pgmodels.RandomGenericFile(obj.ID, obj.Identifier)
		gf.InstitutionID = inst2Admin.InstitutionID
		require.NoError(t, gf.Save())
		req.AddFile(gf)
	}

	// Save the request, with associated files and objects
	require.NoError(t, req.Save())

	// Now test the IncludesFile and IncludesObject functions
	for _, obj := range req.IntellectualObjects {
		isIncluded, err := pgmodels.DeletionRequestIncludesObject(req.ID, obj.ID)
		require.NoError(t, err)
		assert.True(t, isIncluded)
	}
	// Random, bogus object should not be included
	isIncluded, err := pgmodels.DeletionRequestIncludesObject(req.ID, 99999999)
	require.NoError(t, err)
	assert.False(t, isIncluded)

	// Check the files
	for _, gf := range req.GenericFiles {
		isIncluded, err := pgmodels.DeletionRequestIncludesFile(req.ID, gf.ID)
		require.NoError(t, err)
		assert.True(t, isIncluded)
	}
	// Random, bogus file should not be included
	isIncluded, err = pgmodels.DeletionRequestIncludesFile(req.ID, 99999999)
	require.NoError(t, err)
	assert.False(t, isIncluded)
}

func TestDeletionRequestConfirm(t *testing.T) {
	user, err := pgmodels.UserByID(5)
	require.Nil(t, err)
	require.NotNil(t, user)

	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)

	assert.Empty(t, req.ConfirmedBy)
	assert.Empty(t, req.ConfirmedByID)
	assert.Empty(t, req.ConfirmedAt)

	req.Confirm(user)
	assert.Equal(t, user, req.ConfirmedBy)
	assert.Equal(t, user.ID, req.ConfirmedByID)
	assert.NotEmpty(t, req.ConfirmedAt)
}

func TestDeletionRequestCancel(t *testing.T) {
	user, err := pgmodels.UserByID(5)
	require.Nil(t, err)
	require.NotNil(t, user)

	req, err := pgmodels.NewDeletionRequest()
	require.Nil(t, err)
	require.NotNil(t, req)

	assert.Empty(t, req.CancelledBy)
	assert.Empty(t, req.CancelledByID)
	assert.Empty(t, req.CancelledAt)

	req.Cancel(user)
	assert.Equal(t, user, req.CancelledBy)
	assert.Equal(t, user.ID, req.CancelledByID)
	assert.NotEmpty(t, req.CancelledAt)
}
