package models

import (
	"time"
)

type Role struct {
	ID        int64     `json:"id" form:"id" pg:"id"`
	Name      string    `json:"name" form:"name" pg:"name"`
	CreatedAt time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`
}

// GetID returns this record's ID.
func (role *Role) GetID() int64 {
	return role.ID
}

// Authorize - Anyone can retrieve a role record.
func (role *Role) Authorize(*Role, string) error {
	return nil
}

func (role *Role) DeleteIsForbidden() bool {
	return true
}

// IsReadOnly is true. We don't want anyone editing roles.
func (role *Role) IsReadOnly() bool {
	return true
}

// SupportsSoftDelete - false.
func (role *Role) SupportsSoftDelete() bool {
	return false
}

func (role *Role) SetSoftDeleteAttributes(*Role) {
	// No-op
}

func (role *Role) ClearSoftDeleteAttributes() {
	// No-op
}

func (role *Role) SetTimestamps() {
	now := time.Now().UTC()
	if role.CreatedAt.IsZero() {
		role.CreatedAt = now
	}
	role.UpdatedAt = now
}

func (role *Role) BeforeSave() error {
	// TODO: Validate
	return nil
}
