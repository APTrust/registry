package forms

import (
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InstitutionForm struct {
	Form
	Institution *models.Institution
	instOptions []ListOption
}

var InstitutionFormErrors = map[string]string{
	"Name":       "Name must contain at least two letters.",
	"Identifier": "Identifier must be a domain name.",
}

func NewInstitutionForm(ds *models.DataStore, institution *models.Institution) (*InstitutionForm, error) {
	var err error
	institutionForm := &InstitutionForm{
		Form:        NewForm(ds),
		Institution: institution,
	}
	institutionForm.instOptions, err = ListInstitutions(ds)
	if err != nil {
		return nil, err
	}
	institutionForm.init()
	institutionForm.setValues()
	return institutionForm, err
}

func (f *InstitutionForm) init() {
	f.Fields["Name"] = &Field{
		Name:        "Name",
		Placeholder: "Name",
		Attrs: map[string]string{
			"required": "",
		},
	}
	f.Fields["Identifier"] = &Field{
		Name:        "Identifier",
		Placeholder: "Identifier",
		Attrs: map[string]string{
			"required": "",
			"pattern":  "[A-Za-z0-9]{2,}\\.[A-Za-z0-9]{2,}",
		},
	}
}

func (f *InstitutionForm) Bind(c *gin.Context) error {
	err := c.ShouldBind(f.Institution)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				f.Fields[fieldErr.Field()].Error = InstitutionFormErrors[fieldErr.Field()]
			}
		}
	}
	f.setValues()
	return err
}

// setValues sets the form values to match the Institution values.
func (f *InstitutionForm) setValues() {
	f.Fields["Name"].Value = f.Institution.Name
	f.Fields["Identifier"].Value = f.Institution.Identifier
}
