package web

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemRequeueForm struct {
	Form
	stageOptions []ListOption
}

func NewWorkItemRequeueForm(request *Request) (*WorkItemRequeueForm, error) {
	var err error
	item := &pgmodels.WorkItem{}
	if request.ResourceID > 0 {
		item, err = pgmodels.WorkItemByID(request.ResourceID)
		if err != nil {
			return nil, err
		}
		// We cannot requeue items after processing has completed.
		if item.HasCompleted() {
			common.Context().Log.Error().Msgf("Invalid request for requeue form. WorkItem %d (%s) has completed %s and cannot be requeued.", item.ID, item.Name, item.Action)
			return nil, common.ErrNotSupported
		}
	}
	// Bind submitted form values in case we have to
	// re-display the form with an error message.
	request.GinContext.ShouldBind(item)

	itemForm := &WorkItemRequeueForm{
		Form: NewForm(request, item),
	}
	itemForm.init()
	return itemForm, err
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
func (f *WorkItemRequeueForm) setValues() {

}
