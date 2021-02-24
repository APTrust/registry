package models

import (
	"github.com/go-playground/validator/v10"
)

type Model interface {
	GetID() int64
	Authorize(*User, string) error
	DeleteIsForbidden() bool
	UpdateIsForbidden() bool
	IsReadOnly() bool // views are read-only
	SupportsSoftDelete() bool
	SetSoftDeleteAttributes(*User)
	ClearSoftDeleteAttributes()
	SetTimestamps()
	BeforeSave() error
	GetValidationErrors(validator.ValidationErrors) map[string]string
}
