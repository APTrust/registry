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
	"AlertCreate":                        {"Alert", constants.AlertCreate},
	"AlertNew":                           {"Alert", constants.AlertCreate},
	"AlertIndex":                         {"Alert", constants.AlertRead},
	"AlertShow":                          {"Alert", constants.AlertRead},
	"AlertUpdate":                        {"Alert", constants.AlertUpdate},
	"AlertDelete":                        {"Alert", constants.AlertDelete},
	"ChecksumNew":                        {"Checksum", constants.ChecksumCreate},
	"ChecksumCreate":                     {"Checksum", constants.ChecksumCreate},
	"ChecksumShow":                       {"Checksum", constants.ChecksumRead},
	"ChecksumIndex":                      {"Checksum", constants.ChecksumRead},
	"ChecksumUpdate":                     {"Checksum", constants.ChecksumUpdate},
	"ChecksumDelete":                     {"Checksum", constants.ChecksumDelete},
	"DashboardShow":                      {"Dashboard", constants.DashboardShow},
	"EventCreate":                        {"Event", constants.EventCreate},
	"EventNew":                           {"Event", constants.EventCreate},
	"EventShow":                          {"Event", constants.EventRead},
	"EventIndex":                         {"Event", constants.EventRead},
	"EventUpdate":                        {"Event", constants.EventUpdate},
	"EventDelete":                        {"Event", constants.EventDelete},
	"GenericFileNew":                     {"GenericFile", constants.GenericFileCreate},
	"GenericFileCreate":                  {"GenericFile", constants.GenericFileCreate},
	"GenericFileShow":                    {"GenericFile", constants.GenericFileRead},
	"GenericFileIndex":                   {"GenericFile", constants.GenericFileRead},
	"GenericFileUpdate":                  {"GenericFile", constants.GenericFileUpdate},
	"GenericFileDelete":                  {"GenericFile", constants.GenericFileDelete},
	"GenericFileRequestDelete":           {"GenericFile", constants.GenericFileRequestDelete},
	"GenericFileInitDelete":              {"GenericFile", constants.GenericFileRequestDelete},
	"GenericFileApproveDelete":           {"GenericFile", constants.GenericFileApproveDelete},
	"GenericFileFinishBulkDelete":        {"GenericFile", constants.GenericFileFinishBulkDelete},
	"GenericFileRequestRestore":          {"IntellectualObject", constants.GenericFileRestore},
	"GenericFileInitRestore":             {"IntellectualObject", constants.GenericFileRestore},
	"InstitutionNew":                     {"Institution", constants.InstitutionCreate},
	"InstitutionCreate":                  {"Institution", constants.InstitutionCreate},
	"InstitutionEdit":                    {"Institution", constants.InstitutionUpdate},
	"InstitutionIndex":                   {"Institution", constants.InstitutionRead},
	"InstitutionShow":                    {"Institution", constants.InstitutionRead},
	"InstitutionUpdate":                  {"Institution", constants.InstitutionUpdate},
	"InstitutionDelete":                  {"Institution", constants.InstitutionDelete},
	"InstitutionUndelete":                {"Institution", constants.InstitutionUpdate},
	"IntellectualObjectNew":              {"IntellectualObject", constants.IntellectualObjectCreate},
	"IntellectualObjectCreate":           {"IntellectualObject", constants.IntellectualObjectCreate},
	"IntellectualObjectIndex":            {"IntellectualObject", constants.IntellectualObjectRead},
	"IntellectualObjectShow":             {"IntellectualObject", constants.IntellectualObjectRead},
	"IntellectualObjectUpdate":           {"IntellectualObject", constants.IntellectualObjectUpdate},
	"IntellectualObjectDelete":           {"IntellectualObject", constants.IntellectualObjectDelete},
	"IntellectualObjectEvents":           {"PremisEvent", constants.EventRead},
	"IntellectualObjectFiles":            {"PremisEvent", constants.GenericFileRead},
	"IntellectualObjectRequestDelete":    {"IntellectualObject", constants.IntellectualObjectRequestDelete},
	"IntellectualObjectInitDelete":       {"IntellectualObject", constants.IntellectualObjectRequestDelete},
	"IntellectualObjectApproveDelete":    {"IntellectualObject", constants.IntellectualObjectApproveDelete},
	"IntellectualObjectFinishBulkDelete": {"IntellectualObject", constants.IntellectualObjectFinishBulkDelete},
	"IntellectualObjectRequestRestore":   {"IntellectualObject", constants.IntellectualObjectRestore},
	"IntellectualObjectInitRestore":      {"IntellectualObject", constants.IntellectualObjectRestore},
	"StorageRecordNew":                   {"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordCreate":                {"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordIndex":                 {"StorageRecord", constants.StorageRecordRead},
	"StorageRecordShow":                  {"StorageRecord", constants.StorageRecordRead},
	"StorageRecordUpdate":                {"StorageRecord", constants.StorageRecordUpdate},
	"StorageRecordDelete":                {"StorageRecord", constants.StorageRecordDelete},
	"UserNew":                            {"User", constants.UserCreate},
	"UserCreate":                         {"User", constants.UserCreate},
	"UserEdit":                           {"User", constants.UserUpdate},
	"UserIndex":                          {"User", constants.UserRead},
	"UserShow":                           {"User", constants.UserRead},
	"UserUpdate":                         {"User", constants.UserUpdate},
	"UserUndelete":                       {"User", constants.UserUpdate},
	"UserDelete":                         {"User", constants.UserDelete},
	"UserReadSelf":                       {"User", constants.UserReadSelf},
	"UserUpdateSelf":                     {"User", constants.UserUpdateSelf},
	"UserDeleteSelf":                     {"User", constants.UserDeleteSelf},
	"WorkItemNew":                        {"WorkItem", constants.WorkItemCreate},
	"WorkItemCreate":                     {"WorkItem", constants.WorkItemCreate},
	"WorkItemIndex":                      {"WorkItem", constants.WorkItemRead},
	"WorkItemRequeue":                    {"WorkItem", constants.WorkItemRequeue},
	"WorkItemShow":                       {"WorkItem", constants.WorkItemRead},
	"WorkItemShowRequeue":                {"WorkItem", constants.WorkItemRequeue},
	"WorkItemEdit":                       {"WorkItem", constants.WorkItemUpdate},
	"WorkItemUpdate":                     {"WorkItem", constants.WorkItemUpdate},
	"WorkItemDelete":                     {"WorkItem", constants.WorkItemDelete},
}
