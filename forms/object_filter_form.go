package forms

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// ObjectFilterForm is the form that displays filtering options for
// the generic file list page.
type ObjectFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewObjectFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (*ObjectFilterForm, error) {
	f := &ObjectFilterForm{
		Form:             NewForm(nil, "objects/_filters.html", "/objects"),
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

func (f *ObjectFilterForm) init() {
	f.Fields["access"] = &Field{
		Name:        "access",
		Label:       "Access",
		Placeholder: "Access",
		Options:     Options(constants.AccessSettings),
	}
	f.Fields["alt_identifier"] = &Field{
		Name:        "alt_identifier",
		Label:       "Alternate Identifier",
		Placeholder: "Alternate Identifier",
	}
	f.Fields["bag_group_identifier"] = &Field{
		Name:        "bag_group_identifier",
		Label:       "Bag Group Identifier",
		Placeholder: "Bag Group Identifier",
	}
	f.Fields["bag_name"] = &Field{
		Name:        "bag_name",
		Label:       "Bag Name",
		Placeholder: "Bag Name",
	}
	f.Fields["bagit_profile_identifier"] = &Field{
		Name:        "bagit_profile_identifier",
		Label:       "BagIt Profile",
		Placeholder: "BagIt Profile",
		Options:     BagItProfileIdentifiers,
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
	f.Fields["etag"] = &Field{
		Name:        "etag",
		Label:       "ETag",
		Placeholder: "ETag",
	}
	f.Fields["file_count__gteq"] = &Field{
		Name:        "file_count__gteq",
		Label:       "Min File Count",
		Placeholder: "Min File Count",
	}
	f.Fields["file_count__lteq"] = &Field{
		Name:        "file_count__lteq",
		Label:       "Max File Count",
		Placeholder: "Max File Count",
	}
	f.Fields["identifier"] = &Field{
		Name:        "identifier",
		Label:       "Object Identifier",
		Placeholder: "Object Identifier",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["institution_parent_id"] = &Field{
		Name:        "institution_parent_id",
		Label:       "Parent Institution",
		Placeholder: "Parent Institution",
		Options:     f.instOptions,
	}
	f.Fields["internal_sender_identifier"] = &Field{
		Name:        "internal_sender_identifier",
		Label:       "Internal Sender Identifier",
		Placeholder: "Internal Sender Identifier",
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
	f.Fields["source_organization"] = &Field{
		Name:        "source_organization",
		Label:       "Source Organization",
		Placeholder: "Source Organization",
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
	f.Fields["updated_at__lteq"] = &Field{
		Name:        "updated_at__lteq",
		Label:       "Updated On or Before",
		Placeholder: "Updated On or Before",
	}
	f.Fields["updated_at__gteq"] = &Field{
		Name:        "updated_at__gteq",
		Label:       "Updated On or After",
		Placeholder: "Updated On or After",
	}
}

// setValues sets the form values to match the Institution values.
func (f *ObjectFilterForm) SetValues() {
	f.Fields["access"].Value = f.FilterCollection.ValueOf("access")
	f.Fields["alt_identifier"].Value = f.FilterCollection.ValueOf("alt_identifier")
	f.Fields["bag_group_identifier"].Value = f.FilterCollection.ValueOf("bag_group_identifier")
	f.Fields["bagit_profile_identifier"].Value = f.FilterCollection.ValueOf("bagit_profile_identifier")
	f.Fields["etag"].Value = f.FilterCollection.ValueOf("etag")
	f.Fields["file_count__gteq"].Value = f.FilterCollection.ValueOf("file_count__gteq")
	f.Fields["file_count__lteq"].Value = f.FilterCollection.ValueOf("file_count__lteq")
	f.Fields["identifier"].Value = f.FilterCollection.ValueOf("identifier")
	f.Fields["institution_id"].Value = f.FilterCollection.ValueOf("institution_id")
	f.Fields["institution_parent_id"].Value = f.FilterCollection.ValueOf("institution_parent_id")
	f.Fields["internal_sender_identifier"].Value = f.FilterCollection.ValueOf("internal_sender_identifier")

	f.Fields["size__gteq"].Value = f.FilterCollection.ValueOf("size__gteq")
	f.Fields["size__lteq"].Value = f.FilterCollection.ValueOf("size__lteq")
	f.Fields["source_organization"].Value = f.FilterCollection.ValueOf("source_organization")
	f.Fields["state"].Value = f.FilterCollection.ValueOf("state")
	f.Fields["storage_option"].Value = f.FilterCollection.ValueOf("storage_option")
	f.Fields["created_at__gteq"].Value = f.FilterCollection.ValueOf("created_at__gteq")
	f.Fields["created_at__lteq"].Value = f.FilterCollection.ValueOf("created_at__lteq")
	f.Fields["updated_at__gteq"].Value = f.FilterCollection.ValueOf("updated_at__gteq")
	f.Fields["updated_at__lteq"].Value = f.FilterCollection.ValueOf("updated_at__lteq")
}
