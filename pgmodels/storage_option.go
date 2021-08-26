package pgmodels

import (
	"context"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

const (
	ErrStorageOptionProvider = "Provider is required."
	ErrStorageOptionService  = "Service is required."
	ErrStorageOptionRegion   = "Region is required."
	ErrStorageOptionName     = "Name is required."
	ErrStorageOptionCost     = "Cost is required."
	ErrStorageOptionComment  = "Comment is required."
)

// StorageOption contains information about APTrust storage option
// costs. This is used mainly in monthly cost reporting.
type StorageOption struct {
	ID             int64     `json:"id"`
	Provider       string    `json:"provider"`
	Service        string    `json:"service"`
	Region         string    `json:"region"`
	Name           string    `json:"name"`
	CostGBPerMonth float64   `json:"cost_gb_per_month"`
	Comment        string    `json:"comment"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// StorageOptionByID returns the option with the specified id.
// Returns pg.ErrNoRows if there is no match.
func StorageOptionByID(id int64) (*StorageOption, error) {
	query := NewQuery().Where("id", "=", id)
	return StorageOptionGet(query)
}

func StorageOptionByName(name string) (*StorageOption, error) {
	query := NewQuery().Where("name", "=", name)
	return StorageOptionGet(query)
}

// StorageOptionGet returns the first option matching the query.
func StorageOptionGet(query *Query) (*StorageOption, error) {
	var option StorageOption
	err := query.Select(&option)
	return &option, err
}

func StorageOptionGetAll() ([]*StorageOption, error) {
	query := NewQuery().OrderBy("name")
	return StorageOptionSelect(query)
}

// StorageOptionSelect returns all options matching the query.
func StorageOptionSelect(query *Query) ([]*StorageOption, error) {
	var options []*StorageOption
	err := query.Select(&options)
	return options, err
}

func (option *StorageOption) GetID() int64 {
	return option.ID
}

// Save saves this StorageOption to the database. This will peform an insert
// if StorageOption.ID is zero. Otherwise, it updates.
func (option *StorageOption) Save() error {
	if option.ID == int64(0) {
		return insert(option)
	}
	return update(option)
}

// The following statements have no effect other than to force a compile-time
// check that ensures our model properly implements these hook interfaces.
var (
	_ pg.BeforeInsertHook = (*StorageOption)(nil)
	_ pg.BeforeUpdateHook = (*StorageOption)(nil)
)

// BeforeInsert sets timestamps and bucket names on creation.
func (option *StorageOption) BeforeInsert(c context.Context) (context.Context, error) {
	option.UpdatedAt = time.Now().UTC()
	err := option.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (option *StorageOption) BeforeUpdate(c context.Context) (context.Context, error) {
	option.UpdatedAt = time.Now().UTC()
	err := option.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// Validate validates the model. This is called automatically on insert
// and update.
func (option *StorageOption) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if strings.TrimSpace(option.Provider) == "" {
		errors["Provider"] = ErrStorageOptionProvider
	}
	if strings.TrimSpace(option.Service) == "" {
		errors["Service"] = ErrStorageOptionService
	}
	if strings.TrimSpace(option.Region) == "" {
		errors["Region"] = ErrStorageOptionRegion
	}
	if strings.TrimSpace(option.Name) == "" {
		errors["Name"] = ErrStorageOptionName
	}
	if option.CostGBPerMonth <= 0.0 {
		errors["CostGBPerMonth"] = ErrStorageOptionCost
	}
	if strings.TrimSpace(option.Comment) == "" {
		errors["Comment"] = ErrStorageOptionComment
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}