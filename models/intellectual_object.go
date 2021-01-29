package models

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
)

type IntellectualObject struct {
	ID                        int64     `json:"id" form:"id" pg:"id"`
	Title                     string    `json:"title" form:"title" pg:"title"`
	Description               string    `json:"description" form:"description" pg:"description"`
	Identifier                string    `json:"identifier" form:"identifier" pg:"identifier"`
	AltIdentifier             string    `json:"alt_identifier" form:"alt_identifier" pg:"alt_identifier"`
	Access                    string    `json:"access" form:"access" pg:"access"`
	BagName                   string    `json:"bag_name" form:"bag_name" pg:"bag_name"`
	InstitutionID             int64     `json:"institution_id" form:"institution_id" pg:"institution_id"`
	CreatedAt                 time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
	State                     string    `json:"state" form:"state" pg:"state"`
	ETag                      string    `json:"etag" form:"etag" pg:"etag"`
	BagGroupIdentifier        string    `json:"bag_group_identifier" form:"bag_group_identifier" pg:"bag_group_identifier"`
	StorageOption             string    `json:"storage_option" form:"storage_option" pg:"storage_option"`
	BagItProfileIdentifier    string    `json:"bagit_profile_identifier" form:"bag_it_profile_identifier" pg:"bagit_profile_identifier"`
	SourceOrganization        string    `json:"source_organization" form:"source_organization" pg:"source_organization"`
	InternalSenderIdentifier  string    `json:"internal_sender_identifier" form:"internal_sender_identifier" pg:"internal_sender_identifier"`
	InternalSenderDescription string    `json:"internal_sender_description" form:"internal_sender_description" pg:"internal_sender_description"`

	Institution  *Institution   `json:"institution" pg:"rel:has-one"`
	GenericFiles []*GenericFile `json:"generic_files" pg:"rel:has-many"`
	PremisEvents []*PremisEvent `json:"premis_events" pg:"rel:has-many"`
}

func (obj *IntellectualObject) GetID() int64 {
	return obj.ID
}

func (obj *IntellectualObject) Authorize(actingUser *User, action string) error {
	perm := "Object" + action
	if !actingUser.HasPermission(constants.Permission(perm), obj.InstitutionID) {
		ctx := common.Context()
		ctx.Log.Error().Msgf("Permission denied: acting user %d at inst %d can't %s object %d belonging to inst %d", actingUser.ID, actingUser.InstitutionID, perm, obj.ID, obj.InstitutionID)
		return common.ErrPermissionDenied
	}
	return nil
}

func (obj *IntellectualObject) DeleteIsForbidden() bool {
	return false
}

func (obj *IntellectualObject) UpdateIsForbidden() bool {
	return false
}

func (obj *IntellectualObject) IsReadOnly() bool {
	return false
}

func (obj *IntellectualObject) SupportsSoftDelete() bool {
	return true
}

func (obj *IntellectualObject) SetSoftDeleteAttributes(actingUser *User) {
	obj.State = "D"
}

func (obj *IntellectualObject) ClearSoftDeleteAttributes() {
	obj.State = "A"
}

func (obj *IntellectualObject) SetTimestamps() {
	now := time.Now().UTC()
	if obj.CreatedAt.IsZero() {
		obj.CreatedAt = now
	}
	obj.UpdatedAt = now
}

func (obj *IntellectualObject) BeforeSave() error {
	// TODO: Validate
	return nil
}

func IntellectualObjectFind(id int64) (*IntellectualObject, error) {
	ctx := common.Context()
	obj := &IntellectualObject{ID: id}
	err := ctx.DB.Model(obj).WherePK().Select()
	return obj, err
}

func IntellectualObjectFindByIdentifier(identifier string) (*IntellectualObject, error) {
	ctx := common.Context()
	obj := &IntellectualObject{}
	err := ctx.DB.Model(obj).Where(`"identifier" = ?`, identifier).First()
	return obj, err
}
