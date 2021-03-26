package middleware

import (
	"github.com/APTrust/registry/constants"
)

// AuthMetadata contains information about what type of resource is being
// requested, and what action the user wants to take on that resource.
type AuthMetadata struct {
	// ResourceType is the type of resource the user is requesting.
	// E.g. "IntellectualObject", "GenericFile", etc.
	ResourceType string
	// Permission is the name of the permission required to access
	// the requested resources. E.g. "PremisEventCreate".
	Permission constants.Permission
}

// PermissionForHandler maps HTTP handler names to the permissions
// required to access that handler. We want to be explicit about permissions
// because explicitness is much easier to debug than assumptions and
// magic methods.
//
// When permissions are explicitly defined, we can check them in middleware
// instead of polluting request handler code with lots of security logic.
// We know that middleware always runs and cannot be accidentally omitted
// by the developer.
//
// We also have a failsafe built in to our Authorize middleware that
// raises an error if permissions are not checked and explicitly granted.
// If we forget to map a controller action to permission metadata below,
// requests that hit an unguarded route will return  an internal server
// error.
var AuthMap = map[string]AuthMetadata{
	"AlertCreate":            AuthMetadata{"Alert", constants.AlertCreate},
	"AlertNew":               AuthMetadata{"Alert", constants.AlertCreate},
	"AlertIndex":             AuthMetadata{"Alert", constants.AlertRead},
	"AlertShow":              AuthMetadata{"Alert", constants.AlertRead},
	"AlertUpdate":            AuthMetadata{"Alert", constants.AlertUpdate},
	"AlertDelete":            AuthMetadata{"Alert", constants.AlertDelete},
	"ChecksumNew":            AuthMetadata{"Checksum", constants.ChecksumCreate},
	"ChecksumCreate":         AuthMetadata{"Checksum", constants.ChecksumCreate},
	"ChecksumShow":           AuthMetadata{"Checksum", constants.ChecksumRead},
	"ChecksumIndex":          AuthMetadata{"Checksum", constants.ChecksumRead},
	"ChecksumUpdate":         AuthMetadata{"Checksum", constants.ChecksumUpdate},
	"ChecksumDelete":         AuthMetadata{"Checksum", constants.ChecksumDelete},
	"DashboardShow":          AuthMetadata{"Dashboard", constants.DashboardShow},
	"EventCreate":            AuthMetadata{"Event", constants.EventCreate},
	"EventNew":               AuthMetadata{"Event", constants.EventCreate},
	"EventShow":              AuthMetadata{"Event", constants.EventRead},
	"EventIndex":             AuthMetadata{"Event", constants.EventRead},
	"EventUpdate":            AuthMetadata{"Event", constants.EventUpdate},
	"EventDelete":            AuthMetadata{"Event", constants.EventDelete},
	"FileNew":                AuthMetadata{"GenericFile", constants.FileCreate},
	"FileCreate":             AuthMetadata{"GenericFile", constants.FileCreate},
	"FileShow":               AuthMetadata{"GenericFile", constants.FileRead},
	"FileIndex":              AuthMetadata{"GenericFile", constants.FileRead},
	"FileUpdate":             AuthMetadata{"GenericFile", constants.FileUpdate},
	"FileDelete":             AuthMetadata{"GenericFile", constants.FileDelete},
	"FileRequestDelete":      AuthMetadata{"GenericFile", constants.FileRequestDelete},
	"FileApproveDelete":      AuthMetadata{"GenericFile", constants.FileApproveDelete},
	"FileFinishBulkDelete":   AuthMetadata{"GenericFile", constants.FileFinishBulkDelete},
	"FileRestore":            AuthMetadata{"GenericFile", constants.FileRestore},
	"InstitutionNew":         AuthMetadata{"Institution", constants.InstitutionCreate},
	"InstitutionCreate":      AuthMetadata{"Institution", constants.InstitutionCreate},
	"InstitutionEdit":        AuthMetadata{"Institution", constants.InstitutionUpdate},
	"InstitutionIndex":       AuthMetadata{"Institution", constants.InstitutionRead},
	"InstitutionShow":        AuthMetadata{"Institution", constants.InstitutionRead},
	"InstitutionUpdate":      AuthMetadata{"Institution", constants.InstitutionUpdate},
	"InstitutionDelete":      AuthMetadata{"Institution", constants.InstitutionDelete},
	"ObjectNew":              AuthMetadata{"IntellectualObject", constants.ObjectCreate},
	"ObjectCreate":           AuthMetadata{"IntellectualObject", constants.ObjectCreate},
	"ObjectIndex":            AuthMetadata{"IntellectualObject", constants.ObjectRead},
	"ObjectShow":             AuthMetadata{"IntellectualObject", constants.ObjectRead},
	"ObjectUpdate":           AuthMetadata{"IntellectualObject", constants.ObjectUpdate},
	"ObjectDelete":           AuthMetadata{"IntellectualObject", constants.ObjectDelete},
	"ObjectRequestDelete":    AuthMetadata{"IntellectualObject", constants.ObjectRequestDelete},
	"ObjectApproveDelete":    AuthMetadata{"IntellectualObject", constants.ObjectApproveDelete},
	"ObjectFinishBulkDelete": AuthMetadata{"IntellectualObject", constants.ObjectFinishBulkDelete},
	"ObjectRestore":          AuthMetadata{"IntellectualObject", constants.ObjectRestore},
	"StorageRecordNew":       AuthMetadata{"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordCreate":    AuthMetadata{"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordIndex":     AuthMetadata{"StorageRecord", constants.StorageRecordRead},
	"StorageRecordShow":      AuthMetadata{"StorageRecord", constants.StorageRecordRead},
	"StorageRecordUpdate":    AuthMetadata{"StorageRecord", constants.StorageRecordUpdate},
	"StorageRecordDelete":    AuthMetadata{"StorageRecord", constants.StorageRecordDelete},
	"UserNew":                AuthMetadata{"User", constants.UserCreate},
	"UserCreate":             AuthMetadata{"User", constants.UserCreate},
	"UserEdit":               AuthMetadata{"User", constants.UserUpdate},
	"UserIndex":              AuthMetadata{"User", constants.UserRead},
	"UserShow":               AuthMetadata{"User", constants.UserRead},
	"UserUpdate":             AuthMetadata{"User", constants.UserUpdate},
	"UserUndelete":           AuthMetadata{"User", constants.UserUpdate},
	"UserDelete":             AuthMetadata{"User", constants.UserDelete},
	"UserReadSelf":           AuthMetadata{"User", constants.UserReadSelf},
	"UserUpdateSelf":         AuthMetadata{"User", constants.UserUpdateSelf},
	"UserDeleteSelf":         AuthMetadata{"User", constants.UserDeleteSelf},
	"WorkItemNew":            AuthMetadata{"WorkItem", constants.WorkItemCreate},
	"WorkItemCreate":         AuthMetadata{"WorkItem", constants.WorkItemCreate},
	"WorkItemIndex":          AuthMetadata{"WorkItem", constants.WorkItemRead},
	"WorkItemShow":           AuthMetadata{"WorkItem", constants.WorkItemRead},
	"WorkItemUpdate":         AuthMetadata{"WorkItem", constants.WorkItemUpdate},
	"WorkItemDelete":         AuthMetadata{"WorkItem", constants.WorkItemDelete},
}
