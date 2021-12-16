package pgmodels

import (
	"time"

	"github.com/APTrust/registry/common"
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
	BaseModel
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
	query := NewQuery().OrderBy("name", "asc")
	return StorageOptionSelect(query)
}

// StorageOptionSelect returns all options matching the query.
func StorageOptionSelect(query *Query) ([]*StorageOption, error) {
	var options []*StorageOption
	err := query.Select(&options)
	return options, err
}

// Save saves this StorageOption to the database. This will peform an insert
// if StorageOption.ID is zero. Otherwise, it updates.
func (option *StorageOption) Save() error {
	option.UpdatedAt = time.Now().UTC()
	err := option.Validate()
	if err != nil {
		return err
	}
	if option.ID == int64(0) {
		return insert(option)
	}
	return update(option)
}

// Validate validates the model. This is called automatically on insert
// and update.
func (option *StorageOption) Validate() *common.ValidationError {
	errors := make(map[string]string)
	if common.IsEmptyString(option.Provider) {
		errors["Provider"] = ErrStorageOptionProvider
	}
	if common.IsEmptyString(option.Service) {
		errors["Service"] = ErrStorageOptionService
	}
	if common.IsEmptyString(option.Region) {
		errors["Region"] = ErrStorageOptionRegion
	}
	if common.IsEmptyString(option.Name) {
		errors["Name"] = ErrStorageOptionName
	}
	if option.CostGBPerMonth <= 0.0 {
		errors["CostGBPerMonth"] = ErrStorageOptionCost
	}
	if common.IsEmptyString(option.Comment) {
		errors["Comment"] = ErrStorageOptionComment
	}
	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
