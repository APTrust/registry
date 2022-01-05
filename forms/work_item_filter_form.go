package forms

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// WorkItemFilterForm is the form that displays filtering options for
// the work item list page.
type WorkItemFilterForm struct {
	Form
	FilterCollection  *pgmodels.FilterCollection
	actingUserIsAdmin bool
	instOptions       []ListOption
}

func NewWorkItemFilterForm(fc *pgmodels.FilterCollection, actingUser *pgmodels.User) (FilterForm, error) {
	f := &WorkItemFilterForm{
		Form:             NewForm(nil, "work_items/_filters.html", "/work_items"),
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

func (f *WorkItemFilterForm) init() {
	f.Fields["action"] = &Field{
		Name:        "action",
		Label:       "Action",
		Placeholder: "Action",
		Options:     Options(constants.WorkItemActions),
	}
	f.Fields["alt_identifier"] = &Field{
		Name:        "alt_identifier",
		Label:       "Alternate Identifier",
		Placeholder: "Alternate Identifier",
	}
	f.Fields["bag_date__gteq"] = &Field{
		Name:        "bag_date__gteq",
		Label:       "Bag Date On or After",
		Placeholder: "Bag Date On or After",
	}
	f.Fields["bag_date__lteq"] = &Field{
		Name:        "bag_date__lteq",
		Label:       "Bag Date On or Before",
		Placeholder: "Bag Date On or Before",
	}
	f.Fields["bag_group_identifier"] = &Field{
		Name:        "bag_group_identifier",
		Label:       "Bag Group Identifier",
		Placeholder: "Bag Group Identifier",
	}
	f.Fields["bagit_profile_identifier"] = &Field{
		Name:        "bagit_profile_identifier",
		Label:       "BagIt Profile",
		Placeholder: "BagIt Profile",
		Options:     BagItProfileIdentifiers,
	}
	f.Fields["bucket"] = &Field{
		Name:        "bucket",
		Label:       "Bucket",
		Placeholder: "Bucket",
	}
	f.Fields["date_processed__gteq"] = &Field{
		Name:        "date_processed__gteq",
		Label:       "Processed On or After",
		Placeholder: "Processed On or After",
	}
	f.Fields["date_processed__lteq"] = &Field{
		Name:        "date_processed__lteq",
		Label:       "Processed On or Before",
		Placeholder: "Processed On or Before",
	}
	f.Fields["etag"] = &Field{
		Name:        "etag",
		Label:       "ETag (from S3 upload)",
		Placeholder: "ETag",
	}
	f.Fields["generic_file_identifier"] = &Field{
		Name:        "generic_file_identifier",
		Label:       "File Identifier",
		Placeholder: "File Identifier",
	}
	f.Fields["institution_id"] = &Field{
		Name:        "institution_id",
		Label:       "Institution",
		Placeholder: "Institution",
		Options:     f.instOptions,
	}
	f.Fields["name"] = &Field{
		Name:        "name",
		Label:       "Name of tar file",
		Placeholder: "Name of tar file",
	}
	f.Fields["needs_admin_review"] = &Field{
		Name:        "needs_admin_review",
		Label:       "Needs Admin Review",
		Placeholder: "Needs Admin Review",
		Options:     YesNoList,
	}
	f.Fields["node__not_null"] = &Field{
		Name:        "node__not_null",
		Label:       "Has Worker",
		Placeholder: "Has Worker",
		Options:     YesNoList,
	}
	f.Fields["object_identifier"] = &Field{
		Name:        "object_identifier",
		Label:       "Object Identifier",
		Placeholder: "Object Identifier",
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
	f.Fields["stage"] = &Field{
		Name:        "stage",
		Label:       "Work Item Stage",
		Placeholder: "Work Item Stage",
		Options:     Options(constants.Stages),
	}
	f.Fields["status"] = &Field{
		Name:        "status",
		Label:       "Work Item Status",
		Placeholder: "Work Item Status",
		Options:     Options(constants.Statuses),
	}
	f.Fields["storage_option"] = &Field{
		Name:        "storage_option",
		Label:       "Storage Option",
		Placeholder: "Storage Option",
		Options:     StorageOptionList,
	}
	f.Fields["user"] = &Field{
		Name:        "user",
		Label:       "Initiated By",
		Placeholder: "User email address",
	}
}

// setValues sets the form values to match the Institution values.
func (f *WorkItemFilterForm) SetValues() {
	for _, fieldName := range pgmodels.WorkItemFilters {
		if f.Fields[fieldName] == nil {
			common.ConsoleDebug("No filter for", fieldName)
			continue
		}
		f.Fields[fieldName].Value = f.FilterCollection.ValueOf(fieldName)
	}
}
