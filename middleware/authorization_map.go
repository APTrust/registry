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
	"AlertDelete":                        {"Alert", constants.AlertDelete},
	"AlertIndex":                         {"Alert", constants.AlertRead},
	"AlertNew":                           {"Alert", constants.AlertCreate},
	"AlertShow":                          {"Alert", constants.AlertRead},
	"AlertUpdate":                        {"Alert", constants.AlertUpdate},
	"ChecksumCreate":                     {"Checksum", constants.ChecksumCreate},
	"ChecksumDelete":                     {"Checksum", constants.ChecksumDelete},
	"ChecksumIndex":                      {"Checksum", constants.ChecksumRead},
	"ChecksumNew":                        {"Checksum", constants.ChecksumCreate},
	"ChecksumShow":                       {"Checksum", constants.ChecksumRead},
	"ChecksumUpdate":                     {"Checksum", constants.ChecksumUpdate},
	"DashboardShow":                      {"Dashboard", constants.DashboardShow},
	"DeletionRequestIndex":               {"DeletionRequest", constants.DeletionRequestList},
	"DeletionRequestShow":                {"DeletionRequest", constants.DeletionRequestShow},
	"EventCreate":                        {"Event", constants.EventCreate},
	"EventDelete":                        {"Event", constants.EventDelete},
	"EventIndex":                         {"Event", constants.EventRead},
	"EventNew":                           {"Event", constants.EventCreate},
	"EventShow":                          {"Event", constants.EventRead},
	"EventUpdate":                        {"Event", constants.EventUpdate},
	"GenericFileApproveDelete":           {"GenericFile", constants.GenericFileApproveDelete},
	"GenericFileCancelDelete":            {"GenericFile", constants.GenericFileApproveDelete},
	"GenericFileCreate":                  {"GenericFile", constants.GenericFileCreate},
	"GenericFileDelete":                  {"GenericFile", constants.GenericFileDelete},
	"GenericFileFinishBulkDelete":        {"GenericFile", constants.GenericFileFinishBulkDelete},
	"GenericFileIndex":                   {"GenericFile", constants.GenericFileRead},
	"GenericFileInitDelete":              {"GenericFile", constants.GenericFileRequestDelete},
	"GenericFileInitRestore":             {"GenericFile", constants.GenericFileRestore},
	"GenericFileNew":                     {"GenericFile", constants.GenericFileCreate},
	"GenericFileRequestDelete":           {"GenericFile", constants.GenericFileRequestDelete},
	"GenericFileRequestRestore":          {"GenericFile", constants.GenericFileRestore},
	"GenericFileReviewDelete":            {"GenericFile", constants.GenericFileApproveDelete},
	"GenericFileShow":                    {"GenericFile", constants.GenericFileRead},
	"GenericFileUpdate":                  {"GenericFile", constants.GenericFileUpdate},
	"InstitutionCreate":                  {"Institution", constants.InstitutionCreate},
	"InstitutionDelete":                  {"Institution", constants.InstitutionDelete},
	"InstitutionEdit":                    {"Institution", constants.InstitutionUpdate},
	"InstitutionIndex":                   {"Institution", constants.InstitutionRead},
	"InstitutionNew":                     {"Institution", constants.InstitutionCreate},
	"InstitutionShow":                    {"Institution", constants.InstitutionRead},
	"InstitutionUndelete":                {"Institution", constants.InstitutionUpdate},
	"InstitutionUpdate":                  {"Institution", constants.InstitutionUpdate},
	"IntellectualObjectApproveDelete":    {"IntellectualObject", constants.IntellectualObjectApproveDelete},
	"IntellectualObjectCreate":           {"IntellectualObject", constants.IntellectualObjectCreate},
	"IntellectualObjectDelete":           {"IntellectualObject", constants.IntellectualObjectDelete},
	"IntellectualObjectEvents":           {"PremisEvent", constants.EventRead},
	"IntellectualObjectFiles":            {"PremisEvent", constants.GenericFileRead},
	"IntellectualObjectFinishBulkDelete": {"IntellectualObject", constants.IntellectualObjectFinishBulkDelete},
	"IntellectualObjectIndex":            {"IntellectualObject", constants.IntellectualObjectRead},
	"IntellectualObjectInitDelete":       {"IntellectualObject", constants.IntellectualObjectRequestDelete},
	"IntellectualObjectInitRestore":      {"IntellectualObject", constants.IntellectualObjectRestore},
	"IntellectualObjectNew":              {"IntellectualObject", constants.IntellectualObjectCreate},
	"IntellectualObjectRequestDelete":    {"IntellectualObject", constants.IntellectualObjectRequestDelete},
	"IntellectualObjectRequestRestore":   {"IntellectualObject", constants.IntellectualObjectRestore},
	"IntellectualObjectShow":             {"IntellectualObject", constants.IntellectualObjectRead},
	"IntellectualObjectUpdate":           {"IntellectualObject", constants.IntellectualObjectUpdate},
	"PremisEventShowXHR":                 {"PremisEvent", constants.EventRead},
	"StorageRecordCreate":                {"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordDelete":                {"StorageRecord", constants.StorageRecordDelete},
	"StorageRecordIndex":                 {"StorageRecord", constants.StorageRecordRead},
	"StorageRecordNew":                   {"StorageRecord", constants.StorageRecordCreate},
	"StorageRecordShow":                  {"StorageRecord", constants.StorageRecordRead},
	"StorageRecordUpdate":                {"StorageRecord", constants.StorageRecordUpdate},
	"UserCreate":                         {"User", constants.UserCreate},
	"UserDelete":                         {"User", constants.UserDelete},
	"UserDeleteSelf":                     {"User", constants.UserDeleteSelf},
	"UserEdit":                           {"User", constants.UserUpdate},
	"UserIndex":                          {"User", constants.UserRead},
	"UserNew":                            {"User", constants.UserCreate},
	"UserReadSelf":                       {"User", constants.UserReadSelf},
	"UserShow":                           {"User", constants.UserRead},
	"UserUndelete":                       {"User", constants.UserUpdate},
	"UserUpdate":                         {"User", constants.UserUpdate},
	"UserUpdateSelf":                     {"User", constants.UserUpdateSelf},
	"WorkItemCreate":                     {"WorkItem", constants.WorkItemCreate},
	"WorkItemDelete":                     {"WorkItem", constants.WorkItemDelete},
	"WorkItemEdit":                       {"WorkItem", constants.WorkItemUpdate},
	"WorkItemIndex":                      {"WorkItem", constants.WorkItemRead},
	"WorkItemNew":                        {"WorkItem", constants.WorkItemCreate},
	"WorkItemRequeue":                    {"WorkItem", constants.WorkItemRequeue},
	"WorkItemShow":                       {"WorkItem", constants.WorkItemRead},
	"WorkItemShowRequeue":                {"WorkItem", constants.WorkItemRequeue},
	"WorkItemUpdate":                     {"WorkItem", constants.WorkItemUpdate},
}
