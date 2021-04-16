package forms_test

import (
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkItemRequeueFormCompleted(t *testing.T) {
	item := &pgmodels.WorkItem{
		Status: constants.StatusSuccess,
	}
	_, err := forms.NewWorkItemRequeueForm(item)
	require.NotNil(t, err)
	assert.Equal(t, common.ErrNotSupported, err)
}

func TestWorkItemRequeueFormIngest(t *testing.T) {
	query := pgmodels.NewQuery().Where("action", "=", constants.ActionIngest).Limit(1)
	items, err := pgmodels.WorkItemSelect(query)
	require.Nil(t, err)
	require.Equal(t, 1, len(items))
	item := items[0]

	// Save as new item with desired stage.
	// Receive is first stage, so we should only be able to requeue
	// to that.
	item.ID = 0
	item.Stage = constants.StageReceive
	item.Status = constants.StatusPending

	form, err := forms.NewWorkItemRequeueForm(item)
	require.Nil(t, err)
	require.NotNil(t, form)
	require.NotNil(t, form.Fields["Stage"])
	require.Equal(t, 1, len(form.Fields["Stage"].Options))
	assert.Equal(t, constants.StageReceive, form.Fields["Stage"].Options[0].Value)

	// Store is sixth stage, so we should have six stage
	// options in the requeue list.
	item.Stage = constants.StageStore

	form, err = forms.NewWorkItemRequeueForm(item)
	require.Nil(t, err)
	require.NotNil(t, form)
	require.NotNil(t, form.Fields["Stage"])
	require.Equal(t, 6, len(form.Fields["Stage"].Options))
	opts := form.Fields["Stage"].Options
	assert.Equal(t, constants.StageReceive, opts[0].Value)
	assert.Equal(t, constants.StageValidate, opts[1].Value)
	assert.Equal(t, constants.StageReingestCheck, opts[2].Value)
	assert.Equal(t, constants.StageCopyToStaging, opts[3].Value)
	assert.Equal(t, constants.StageFormatIdentification, opts[4].Value)
	assert.Equal(t, constants.StageStore, opts[5].Value)

	// Cleanup is final stage, so all stage options
	// should appear.
	item.Stage = constants.StageCleanup

	form, err = forms.NewWorkItemRequeueForm(item)
	require.Nil(t, err)
	require.NotNil(t, form)
	require.NotNil(t, form.Fields["Stage"])
	require.Equal(t, 9, len(form.Fields["Stage"].Options))
	opts = form.Fields["Stage"].Options
	assert.Equal(t, constants.StageReceive, opts[0].Value)
	assert.Equal(t, constants.StageValidate, opts[1].Value)
	assert.Equal(t, constants.StageReingestCheck, opts[2].Value)
	assert.Equal(t, constants.StageCopyToStaging, opts[3].Value)
	assert.Equal(t, constants.StageFormatIdentification, opts[4].Value)
	assert.Equal(t, constants.StageStore, opts[5].Value)
	assert.Equal(t, constants.StageStorageValidation, opts[6].Value)
	assert.Equal(t, constants.StageRecord, opts[7].Value)
	assert.Equal(t, constants.StageCleanup, opts[8].Value)
}
