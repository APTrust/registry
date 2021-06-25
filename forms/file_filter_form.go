package forms

import (
	"github.com/APTrust/registry/pgmodels"
)

// FileFilterForm is the form that displays filtering options for
// the generic file list page.
type FileFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewFileFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (*FileFilterForm, error) {
	f := &FileFilterForm{
		Form:             NewForm(nil, "files/_filters.html", "/files"),
		FilterCollection: fc,
	}
	var err error
	if actingUser.IsAdmin() {
		// SysAdmin can view files at any institutions.
		f.instOptions, err = ListInstitutions(false)
		if err != nil {
			return nil, err
		}
	}
	f.init()
	f.SetValues()
	return f, nil
}

func (f *FileFilterForm) init() {
	f.Fields["identifier"] = &Field{
		Name:        "identifier",
		Label:       "File Identifier",
		Placeholder: "File Identifier",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["state"] = &Field{
		Name:        "state",
		Label:       "State",
		Placeholder: "State",
		Options:     ObjectStateList,
	}
	f.Fields["storage_option"] = &Field{
		Name:        "storage_option",
		Label:       "Storage Option",
		Placeholder: "Storage Option",
		Options:     StorageOptionList,
	}
	f.Fields["size__gteq"] = &Field{
		Name:        "size__gteq",
		Label:       "Min Size",
		Placeholder: "Min Size",
	}
	f.Fields["size__lteq"] = &Field{
		Name:        "size__lteq",
		Label:       "Max Size",
		Placeholder: "Max Size",
	}
	f.Fields["created_at__lteq"] = &Field{
		Name:        "created_at__lteq",
		Label:       "Created On or Before",
		Placeholder: "Created On or Before",
	}
	f.Fields["created_at__gteq"] = &Field{
		Name:        "created_at__gteq",
		Label:       "Created On or After",
		Placeholder: "Created On or After",
	}
}

// setValues sets the form values to match the Institution values.
func (f *FileFilterForm) SetValues() {
	f.Fields["identifier"].Value = f.FilterCollection.ValueOf("identifier")
	f.Fields["institution_id"].Value = f.FilterCollection.ValueOf("institution_id")
	f.Fields["state"].Value = f.FilterCollection.ValueOf("state")
	f.Fields["storage_option"].Value = f.FilterCollection.ValueOf("storage_option")
	f.Fields["size__gteq"].Value = f.FilterCollection.ValueOf("size__gteq")
	f.Fields["size__lteq"].Value = f.FilterCollection.ValueOf("size__lteq")
	f.Fields["created_at__gteq"].Value = f.FilterCollection.ValueOf("created_at__gteq")
	f.Fields["created_at__lteq"].Value = f.FilterCollection.ValueOf("created_at__lteq")
}
