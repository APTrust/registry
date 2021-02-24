package models

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type StorageRecord struct {
	ID            int64  `json:"id" form:"id" pg:"id"`
	GenericFileID int64  `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	URL           string `json:"url" form:"url" pg:"url"`
}

func (sr *StorageRecord) GetID() int64 {
	return sr.ID
}

func (sr *StorageRecord) Authorize(actingUser *User, action string) error {
	ctx := common.Context()
	gf, err := NewDataStore(actingUser).GenericFileFind(sr.GenericFileID)
	if err != nil {
		ctx.Log.Error().Msgf("While checking permissions for StorageRecord %d, user %s, error on GenericFile %d: %v", sr.ID, actingUser.Email, sr.GenericFileID, err)
		return err
	}
	perm := "StorageRecord" + action
	if !actingUser.HasPermission(constants.Permission(perm), gf.InstitutionID) {
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s storage record %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, sr.ID, gf.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (sr *StorageRecord) DeleteIsForbidden() bool {
	return false
}

func (sr *StorageRecord) UpdateIsForbidden() bool {
	return false
}

func (sr *StorageRecord) IsReadOnly() bool {
	return false
}

func (sr *StorageRecord) SupportsSoftDelete() bool {
	return false
}

func (sr *StorageRecord) SetSoftDeleteAttributes(actingUser *User) {
	// No-op
}

func (sr *StorageRecord) ClearSoftDeleteAttributes() {
	// No-op
}

func (sr *StorageRecord) SetTimestamps() {
	// No-op
}

func (sr *StorageRecord) BeforeSave() error {
	// TODO: Validate
	return nil
}

func (sr *StorageRecord) GetValidationErrors(map[string]interface{}) map[string]string {
	return nil
}
