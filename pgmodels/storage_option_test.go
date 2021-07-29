package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageOptionValidation(t *testing.T) {
	option := &pgmodels.StorageOption{}
	err := option.Validate()
	require.NotNil(t, err)

	assert.Equal(t, pgmodels.ErrStorageOptionProvider, err.Errors["Provider"])
	assert.Equal(t, pgmodels.ErrStorageOptionService, err.Errors["Service"])
	assert.Equal(t, pgmodels.ErrStorageOptionRegion, err.Errors["Region"])
	assert.Equal(t, pgmodels.ErrStorageOptionName, err.Errors["Name"])
	assert.Equal(t, pgmodels.ErrStorageOptionCost, err.Errors["CostGBPerMonth"])
	assert.Equal(t, pgmodels.ErrStorageOptionComment, err.Errors["Comment"])

	option.Provider = "Wasabi"
	option.Service = "S3"
	option.Region = "Mars"
	option.Name = "Wasabi-Mars"
	option.CostGBPerMonth = 0.021
	option.Comment = "Upload and download are slow."

	err = option.Validate()
	require.Nil(t, err)
}

func TestStorageOptionByID(t *testing.T) {
	db.LoadFixtures()
	option, err := pgmodels.StorageOptionByID(int64(1))
	require.Nil(t, err)
	require.NotNil(t, option)
	assert.Equal(t, int64(1), option.ID)
}

func TestStorageOptionByName(t *testing.T) {
	db.LoadFixtures()
	option, err := pgmodels.StorageOptionByName("Glacier-Deep-OH")
	require.Nil(t, err)
	require.NotNil(t, option)
	assert.Equal(t, "Glacier-Deep-OH", option.Name)
}

func TestStorageOptionGet(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where("name", "=", "Wasabi-VA")
	option, err := pgmodels.StorageOptionGet(query)
	require.Nil(t, err)
	require.NotNil(t, option)
	assert.Equal(t, "Wasabi-VA", option.Name)
}

func TestStorageOptionGetAll(t *testing.T) {
	db.LoadFixtures()
	options, err := pgmodels.StorageOptionGetAll()
	require.Nil(t, err)
	require.NotEmpty(t, options)
	require.Equal(t, 12, len(options))

	// Should be ordered by name
	assert.Equal(t, "Glacier-Deep-OH", options[0].Name)
	assert.Equal(t, "Wasabi-VA", options[11].Name)
}

func TestStorageOptionSelect(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery()
	query.Where("provider", "=", "AWS")
	query.OrderBy("name asc")
	options, err := pgmodels.StorageOptionSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, options)
	require.Equal(t, 10, len(options))
	assert.Equal(t, "Glacier-Deep-OH", options[0].Name)
	assert.Equal(t, "Standard", options[9].Name)
}

func TestStorageOptionSave(t *testing.T) {
	db.LoadFixtures()
	option, err := pgmodels.StorageOptionByName(constants.StorageOptionStandard)
	require.Nil(t, err)
	require.NotNil(t, option)
	option.Comment = option.Comment + " ** "
	err = option.Save()
	require.Nil(t, err)
}
