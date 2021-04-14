package web

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemForm struct {
	Form
}

func NewWorkItemForm(request *Request) (*WorkItemForm, error) {
	var err error
	item := &pgmodels.WorkItem{}
	if request.ResourceID > 0 {
		item, err = pgmodels.WorkItemByID(request.ResourceID)
		if err != nil {
			return nil, err
		}
	}
	// Bind submitted form values in case we have to
	// re-display the form with an error message.
	request.GinContext.ShouldBind(item)

	itemForm := &WorkItemForm{
		Form: NewForm(request, item),
	}
	itemForm.init()
	itemForm.setValues()
	return itemForm, err
}

func (f *WorkItemForm) init() {
	f.Fields["Stage"] = &Field{
		Name:        "Stage",
		Label:       "Stage",
		Placeholder: "Stage",
		ErrMsg:      pgmodels.ErrItemStage,
		Options:     Options(constants.Stages),
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Status"] = &Field{
		Name:        "Status",
		Label:       "Status",
		Placeholder: "Status",
		ErrMsg:      pgmodels.ErrItemStatus,
		Options:     Options(constants.Statuses),
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Retry"] = &Field{
		Name:        "Retry",
		Label:       "Retry",
		Placeholder: "Retry",
		ErrMsg:      "Please choose yes or no.",
		Options:     YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["NeedsAdminReview"] = &Field{
		Name:        "NeedsAdminReview",
		Label:       "NeedsAdminReview",
		Placeholder: "NeedsAdminReview",
		ErrMsg:      "Please choose yes or no.",
		Options:     YesNoList,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Note"] = &Field{
		Name:        "Note",
		Label:       "Note",
		Placeholder: "Note",
		ErrMsg:      pgmodels.ErrItemNote,
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Node"] = &Field{
		Name:        "Node",
		Label:       "Node",
		Placeholder: "Node",
		Attrs:       map[string]string{},
	}
	f.Fields["PID"] = &Field{
		Name:        "PID",
		Label:       "PID",
		Placeholder: "PID",
		Attrs: map[string]string{
			"required": "",
		},
	}
}

func (f *WorkItemForm) setValues() {
	item := f.Model.(*pgmodels.WorkItem)
	f.Fields["Stage"].Value = item.Stage
	f.Fields["Status"].Value = item.Status
	f.Fields["Retry"].Value = item.Retry
	f.Fields["NeedsAdminReview"].Value = item.NeedsAdminReview
	f.Fields["Note"].Value = item.Note
	f.Fields["Node"].Value = item.Node
	f.Fields["PID"].Value = item.PID
}
