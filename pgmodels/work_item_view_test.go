package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Validate is a no-op since these items are read-only.
func TestWorkItemViewValidation(t *testing.T) {
	item := &pgmodels.WorkItemView{}
	err := item.Validate()
	require.Nil(t, err)
}

func TestWorkItemViewGetID(t *testing.T) {
	item := &pgmodels.WorkItemView{
		ID: 199,
	}
	assert.Equal(t, int64(199), item.GetID())
}

func TestWorkItemViewByID(t *testing.T) {
	db.LoadFixtures()
	item, err := pgmodels.WorkItemViewByID(int64(23))
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, int64(23), item.ID)

	// The view has some extra fields that the regular
	// WorkItem doesn't have.
	assert.Equal(t, "Institution One", item.InstitutionName)
	assert.Equal(t, "institution1.edu", item.InstitutionIdentifier)
	assert.Equal(t, "institution1.edu/pdfs", item.ObjectIdentifier)
	assert.Equal(t, "pdf_docs_with_lots_of_words", item.AltIdentifier)
	assert.Equal(t, "", item.BagGroupIdentifier)
	assert.Equal(t, constants.StorageOptionStandard, item.StorageOption)
	assert.Equal(t, "https://example.com/profile.json", item.BagItProfileIdentifier)
	assert.Equal(t, "Institution One", item.SourceOrganization)
	assert.Equal(t, "Second internal identifier", item.InternalSenderIdentifier)
}

func TestWorkItemViewGet(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery().Where("name", "=", "pdfs.tar")
	item, err := pgmodels.WorkItemViewGet(query)
	require.Nil(t, err)
	require.NotNil(t, item)
	assert.Equal(t, "pdfs.tar", item.Name)
}

func TestWorkItemViewSelect(t *testing.T) {
	db.LoadFixtures()
	query := pgmodels.NewQuery()
	query.Where("name", "!=", "pdfs.tar")
	query.Where("name", "!=", "coal.tar")
	query.OrderBy("name", "asc")
	items, err := pgmodels.WorkItemViewSelect(query)
	require.Nil(t, err)
	require.NotEmpty(t, items)
	assert.True(t, (len(items) > 20 && len(items) < 40))
	for _, item := range items {
		assert.NotEqual(t, "pdfs.tar", item)
		assert.NotEqual(t, "coal.tar", item)
	}
}

func TestWorkItemHasViewCompleted(t *testing.T) {
	item := &pgmodels.WorkItemView{}
	for _, status := range constants.IncompleteStatusValues {
		item.Status = status
		assert.False(t, item.HasCompleted())
	}
	for _, status := range constants.CompletedStatusValues {
		item.Status = status
		assert.True(t, item.HasCompleted())
	}
}
