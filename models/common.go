package models

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
	GetValidationErrors(map[string]interface{}) map[string]string
}
