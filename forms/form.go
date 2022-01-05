package forms

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/go-pg/pg/v10"
)

// FilterFormConstructor describes the function signature for the
// constructors of all our filter forms. Filter forms appear on
// index pages, allowing users to filter lists by name, date, etc.
type FilterFormConstructor func(*pgmodels.FilterCollection, *pgmodels.User) (FilterForm, error)

type FilterForm interface {
	SetValues()
	GetFields() map[string]*Field
}

type Form struct {
	BaseURL  string
	Error    error
	Fields   map[string]*Field
	Model    pgmodels.Model
	Status   int
	Template string
}

func NewForm(model pgmodels.Model, template, baseURL string) Form {
	return Form{
		BaseURL:  baseURL,
		Fields:   make(map[string]*Field),
		Model:    model,
		Template: template,
	}
}

// Save saves the underlying object and sets the error value
// and http status code as necessary.
func (f *Form) Save() bool {
	f.Status = http.StatusSeeOther
	err := f.Model.Save()
	if err != nil {
		f.HandleError(err)
	}
	return err == nil
}

// Action returns the html form.action attribute for this form.
func (f *Form) Action() string {
	if f.Model.GetID() > 0 {
		return fmt.Sprintf("%s/edit/%d", f.BaseURL, f.Model.GetID())
	}
	return fmt.Sprintf("%s/new", f.BaseURL)
}

// PostSaveURL is the url to redirect to after successful save.
func (f *Form) PostSaveURL() string {
	return fmt.Sprintf("%s/show/%d", f.BaseURL, f.Model.GetID())
}

// Handle error sets the error property and http status code on error.
func (f *Form) HandleError(err error) {
	f.Error = err
	f.Status = http.StatusBadRequest
	if valErr, ok := err.(*common.ValidationError); ok {
		for fieldName := range valErr.Errors {
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
			f.Status = http.StatusInternalServerError
		}
	}
}

// GetFields returns a slice of all of the form's fields.
// This is used primarily in testing.
func (f *Form) GetFields() map[string]*Field {
	return f.Fields
}

func (f *Form) SetValues() {
	// no-op
	common.ConsoleDebug("***** Called base.setValues() *****")
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
