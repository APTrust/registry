package forms_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	//"github.com/APTrust/registry/web/testutil"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFormAndItem() (forms.Form, *pgmodels.WorkItem) {
	workItem := &pgmodels.WorkItem{}
	form := forms.NewForm(workItem, "work_items/form.html", "/work_items")
	return form, workItem
}

func TestFormAction(t *testing.T) {
	form, workItem := getFormAndItem()
	assert.Equal(t, "/work_items/new", form.Action())

	workItem.ID = 1
	assert.Equal(t, "/work_items/edit/1", form.Action())
}

func TestFormPostSaveURL(t *testing.T) {
	form, workItem := getFormAndItem()
	workItem.ID = 1
	assert.Equal(t, "/work_items/show/1", form.PostSaveURL())
}

func TestFormHandleError(t *testing.T) {
	form, _ := getFormAndItem()
	form.Fields["test_field"] = &forms.Field{}
	form.Fields["other_field"] = &forms.Field{}
	valErr := &common.ValidationError{
		Errors: map[string]string{
			"test_field": "Oops!",
		},
	}
	form.HandleError(valErr)
	assert.Equal(t, valErr, form.Error)
	assert.True(t, form.Fields["test_field"].DisplayError)
	assert.Equal(t, form.Status, http.StatusBadRequest)

	form, _ = getFormAndItem()
	err := fmt.Errorf("Unexpected error")
	form.HandleError(err)
	assert.Equal(t, err, form.Error)
	assert.Equal(t, http.StatusInternalServerError, form.Status)

	form, _ = getFormAndItem()
	pgErr := pg.ErrNoRows
	form.HandleError(pgErr)
	assert.Equal(t, pgErr, form.Error)
	assert.Equal(t, http.StatusInternalServerError, form.Status)
}

func TestFormGetFields(t *testing.T) {
	form, _ := getFormAndItem()
	form.Fields["test_field"] = &forms.Field{}
	form.Fields["other_field"] = &forms.Field{}
	fields := form.GetFields()
	require.NotEmpty(t, fields)
	assert.NotNil(t, fields["test_field"])
	assert.NotNil(t, fields["other_field"])
}

func TestFormSetValues(t *testing.T) {
	form, _ := getFormAndItem()
	// Make sure there's no error.
	assert.Empty(t, form.GetFields())
}

func TestFormSave(t *testing.T) {
	workItem, err := pgmodels.WorkItemByID(32)
	require.Nil(t, err)
	form := forms.NewWorkItemForm(workItem)

	oldUpdatedAt := workItem.UpdatedAt
	require.True(t, form.Save())
	assert.Nil(t, form.Error)

	workItem, err = pgmodels.WorkItemByID(12)
	require.Nil(t, err)
	assert.True(t, workItem.UpdatedAt.After(oldUpdatedAt))

	(form.Model.(*pgmodels.WorkItem)).Action = "Invalid action to force error"
	require.False(t, form.Save())
	assert.NotNil(t, form.Error)
}
