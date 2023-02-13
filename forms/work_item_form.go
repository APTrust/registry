package forms

import (
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemForm struct {
	Form
}

func NewWorkItemForm(workItem *pgmodels.WorkItem) *WorkItemForm {
	itemForm := &WorkItemForm{
		Form: NewForm(workItem, "work_items/form.html", "/work_items"),
	}
	itemForm.init()
	itemForm.SetValues()
	return itemForm
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
	f.Fields["IntellectualObjectID"] = &Field{
		Name:        "IntellectualObjectID",
		Label:       "Intellectual Object",
		Placeholder: "Intellectual Object",
	}
}

func (f *WorkItemForm) SetValues() {
	item := f.Model.(*pgmodels.WorkItem)
	f.Fields["Stage"].Value = item.Stage
	f.Fields["Status"].Value = item.Status
	f.Fields["Retry"].Value = item.Retry
	f.Fields["NeedsAdminReview"].Value = item.NeedsAdminReview
	f.Fields["Note"].Value = item.Note
	f.Fields["Node"].Value = item.Node
	f.Fields["PID"].Value = item.PID
	if item.IntellectualObjectID > 0 {
		obj, err := pgmodels.IntellectualObjectByID(item.IntellectualObjectID)
		if err == nil {
			f.Fields["IntellectualObjectID"].Value = strconv.FormatInt(item.IntellectualObjectID, 10)
			f.Fields["IntellectualObjectID"].Options = []*ListOption{
				{strconv.FormatInt(item.IntellectualObjectID, 10), obj.Identifier, true},
			}
			f.Fields["IntellectualObjectID"].Attrs = map[string]string{
				"disabled": "",
			}
		} else {
			common.Context().Log.Error().Msgf("WorkItem form could not retrieve associated object for work item %d. Obj id is %d. Error: %v", item.ID, item.IntellectualObjectID, err)
		}
	}
}
