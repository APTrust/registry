package forms_test

import (
	"testing"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkItemForm(t *testing.T) {
	item, err := pgmodels.WorkItemByID(32)
	require.Nil(t, err)
	require.NotNil(t, item)
	form := forms.NewWorkItemForm(item)
	require.NotNil(t, form)

	assert.Equal(t, item.Stage, form.Fields["Stage"].Value)
	assert.Equal(t, item.Status, form.Fields["Status"].Value)
	assert.Equal(t, item.Retry, form.Fields["Retry"].Value)
	assert.Equal(t, item.NeedsAdminReview, form.Fields["NeedsAdminReview"].Value)
	assert.Equal(t, item.Note, form.Fields["Note"].Value)
	assert.Equal(t, item.Node, form.Fields["Node"].Value)
	assert.Equal(t, item.PID, form.Fields["PID"].Value)
}
