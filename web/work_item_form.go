package web

import (
	//"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemForm struct {
	Form
	instOptions []ListOption
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
	return itemForm, err
}

func (f *WorkItemForm) init() {
	// Editable fields:
	//
	// Stage, Status, Retry, NeedsAdminReview, Retry, Node, Pid
}

// setValues sets the form values to match the WorkItem values.
func (f *WorkItemForm) setValues() {

}
