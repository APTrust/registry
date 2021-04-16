package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

type Form struct {
	Action string
	Fields map[string]*Field
	Error  error
}

func NewForm() Form {
	return Form{
		Fields: make(map[string]*Field),
	}
}

func (f *Form) HandleError(err error) int {
	status := http.StatusBadRequest
	if valErr, ok := err.(*common.ValidationError); ok {
		for fieldName, _ := range valErr.Errors {
			f.Fields[fieldName].DisplayError = true
		}
	} else {
		// Integrity violation is usually someone entering
		// an identifier or email address that is already
		// in use.
		isIntegrityViolation := false
		f.Error = err
		if pgErr, ok := err.(pg.Error); ok {
			isIntegrityViolation = pgErr.IntegrityViolation()
		}
		if !isIntegrityViolation {
			status = http.StatusInternalServerError
		}
	}
	return status
}

func (f *Form) SetValues() {
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
