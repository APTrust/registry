package pgmodels

import (
	"time"
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
