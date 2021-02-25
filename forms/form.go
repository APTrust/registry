package forms

import (
	"github.com/APTrust/registry/models"
)

type Form struct {
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
	Name    string
	Label   string
	Value   interface{}
	Error   string
	Options []ListOption
	Attrs   map[string]string
}

type ListOption struct {
	Value string
	Text  string
}
