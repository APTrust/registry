package forms

import (
	"github.com/APTrust/registry/models"
)

type Form struct {
	Action string
	ds     *models.DataStore
	Fields map[string]*Field
}

func NewForm(ds *models.DataStore) Form {
	return Form{
		ds:     ds,
		Fields: make(map[string]*Field),
	}
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
