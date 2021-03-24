package forms

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type Form struct {
	Action       string
	Fields       map[string]*Field
	Model        pgmodels.Model
	ginCtx       *gin.Context
	templateData gin.H
}

func NewForm(c *gin.Context, t gin.H, model pgmodels.Model) Form {
	return Form{
		Fields:       make(map[string]*Field),
		Model:        model,
		ginCtx:       c,
		templateData: t,
	}
}

func (f *Form) Save() (int, error) {
	status := http.StatusCreated
	if f.Model.GetID() > 0 {
		status = http.StatusOK
	}
	_ = f.ginCtx.ShouldBind(f.Model)
	err := f.Model.Save()
	if err != nil {
		status = f.handleError(err)
	}
	f.setValues()
	return status, err
}

func (f *Form) handleError(err error) int {
	status := http.StatusBadRequest
	if valErr, ok := err.(*common.ValidationError); ok {
		for fieldName, _ := range valErr.Errors {
			f.Fields[fieldName].DisplayError = true
		}
	} else {
		f.templateData["FormError"] = err.Error()
		status = http.StatusInternalServerError
	}
	return status
}

func (f *Form) setValues() {
	// no-op
	fmt.Println("***** Called base.setValues() *****")
}

type Field struct {
	Name         string
	Label        string
	Placeholder  string
	Value        interface{}
	ErrMsg       string
	DisplayError bool
	Options      []ListOption
	Attrs        map[string]string
}

type ListOption struct {
	Value string
	Text  string
}
