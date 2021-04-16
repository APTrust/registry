package web

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemForm struct {
	Form
	WorkItem *pgmodels.WorkItem
}

func NewWorkItemForm(workItem *pgmodels.WorkItem) (*WorkItemForm, error) {
	itemForm := &WorkItemForm{
		Form:     NewForm(),
		WorkItem: workItem,
	}
	itemForm.init()
	itemForm.SetValues()
	return itemForm, nil
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

func (f *WorkItemForm) SetValues() {
	f.Fields["Stage"].Value = f.WorkItem.Stage
	f.Fields["Status"].Value = f.WorkItem.Status
	f.Fields["Retry"].Value = f.WorkItem.Retry
	f.Fields["NeedsAdminReview"].Value = f.WorkItem.NeedsAdminReview
	f.Fields["Note"].Value = f.WorkItem.Note
	f.Fields["Node"].Value = f.WorkItem.Node
	f.Fields["PID"].Value = f.WorkItem.PID
}
