package web

import (
	//"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type WorkItemForm struct {
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
	// Editable fields:
	//
	// Stage
}

// setValues sets the form values to match the WorkItem values.
func (f *WorkItemRequeueForm) setValues() {

}
