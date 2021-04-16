package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemRequeueForm struct {
	Form
	stageOptions []ListOption
}

func NewWorkItemRequeueForm(workItem *pgmodels.WorkItem) (*WorkItemRequeueForm, error) {
	if workItem.HasCompleted() {
		common.Context().Log.Error().Msgf("Invalid request for requeue form. WorkItem %d (%s) has completed %s and cannot be requeued.", workItem.ID, workItem.Name, workItem.Action)
		return nil, common.ErrNotSupported
	}

	itemForm := &WorkItemRequeueForm{
		Form: NewForm(workItem, "work_items/form.html", "/work_items"),
	}
	itemForm.init()
	return itemForm, nil
}

func (f *WorkItemRequeueForm) init() {
	options := []ListOption{
		{constants.StageRequested, constants.StageRequested},
	}
	f.Fields["Stage"] = &Field{
		Name:        "Stage",
		Label:       "Requeue to Stage",
		Placeholder: "Stage",
		ErrMsg:      pgmodels.ErrItemStage,
		Options:     options,
		Attrs: map[string]string{
			"required": "",
		},
	}
	if f.Model.(*pgmodels.WorkItem).Action == constants.ActionIngest {
		f.setIngestStages()
	}
}

// setIngestStages sets the stages we can requeue to. We can requeue
// an item to its current stage of ingest, or to any prior stage.
// We cannot requeue to future stages. E.g. We cannot requeue to
// the storage stage if the item hasn't even been validated yet.
func (f *WorkItemRequeueForm) setIngestStages() {
	item := f.Model.(*pgmodels.WorkItem)
	stages := make([]ListOption, 0)
	for _, stage := range constants.IngestStagesInOrder {
		if item.Stage != stage {
			stages = append(stages, ListOption{stage, stage})
		} else {
			stages = append(stages, ListOption{stage, stage})
			break
		}
	}

}

// setValues sets the form values to match the WorkItem values.
func (f *WorkItemRequeueForm) SetValues() {

}
