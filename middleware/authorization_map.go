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
	// PageTitle is the title of the page. This goes into the
	// HTML title element. This is a late addition to address
	// accessibility issues in the web UI.
	PageTitle string
}

// AuthMap maps HTTP handler names to the permissions
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
	"AlertCreate":                       {"Alert", constants.AlertCreate, "Create Alert"},
	"AlertDelete":                       {"Alert", constants.AlertDelete, "Delete Alert"},
	"AlertIndex":                        {"Alert", constants.AlertRead, "Alerts"},
	"AlertNew":                          {"Alert", constants.AlertCreate, "New Alert"},
	"AlertShow":                         {"Alert", constants.AlertRead, "Alert"},
	"AlertUpdate":                       {"Alert", constants.AlertUpdate, "Update Alert"},
	"AlertMarkAsReadXHR":                {"Alert", constants.AlertUpdate, "Mark Alert as Read"},
	"AlertMarkAllAsRead":                {"Alert", constants.AlertUpdate, "Mark All Alerts as Read"},
	"AlertMarkAsUnreadXHR":              {"Alert", constants.AlertUpdate, "Mark Alert as Unread"},
	"BillingReportShow":                 {"DepositStats", constants.BillingReportShow, "Billing Report"},
	"ChecksumCreate":                    {"Checksum", constants.ChecksumCreate, "Create Checksum"},
	"ChecksumDelete":                    {"Checksum", constants.ChecksumDelete, "Delete Checksum"},
	"ChecksumIndex":                     {"Checksum", constants.ChecksumRead, "File Checksums"},
	"ChecksumNew":                       {"Checksum", constants.ChecksumCreate, "New Checksum"},
	"ChecksumShow":                      {"Checksum", constants.ChecksumRead, "Checksum Detail"},
	"ChecksumUpdate":                    {"Checksum", constants.ChecksumUpdate, "Update Checksum"},
	"DashboardShow":                     {"Dashboard", constants.DashboardShow, "Dashboard"},
	"DeletionRequestApprove":            {"DeletionRequest", constants.DeletionRequestApprove, "Approve Deletion Request"},
	"DeletionRequestCancel":             {"DeletionRequest", constants.DeletionRequestApprove, "Cancel Deletion Request"},
	"DeletionRequestIndex":              {"DeletionRequest", constants.DeletionRequestList, "Deletion Requests"},
	"DeletionRequestReview":             {"DeletionRequest", constants.DeletionRequestApprove, "Review Deletion Request"},
	"DeletionRequestShow":               {"DeletionRequest", constants.DeletionRequestShow, "Deletion Request"},
	"DepositReportShow":                 {"DepositStats", constants.DepositReportShow, "Deposit Report"},
	"GenerateFailedFixityAlerts":        {"Alert", constants.GenerateFailedFixityAlert, "Generate Failed Fixity Check"},
	"GenericFileCreate":                 {"GenericFile", constants.FileCreate, "Create Generic File"},
	"GenericFileCreateBatch":            {"GenericFile", constants.FileCreate, "Create Generic File Batch"},
	"GenericFileDelete":                 {"GenericFile", constants.FileDelete, "Delete Generic File"},
	"GenericFileFinishBulkDelete":       {"GenericFile", constants.FileFinishBulkDelete, "Generic File Bulk Deletion Complete"},
	"GenericFileIndex":                  {"GenericFile", constants.FileRead, "Generic Files"},
	"GenericFileInitDelete":             {"GenericFile", constants.FileRequestDelete, "Generic File - Begin Deletion"},
	"GenericFileInitRestore":            {"GenericFile", constants.FileRestore, "Generic File - Begin Restoration"},
	"GenericFileNew":                    {"GenericFile", constants.FileCreate, "New Generic File"},
	"GenericFileRequestDelete":          {"GenericFile", constants.FileRequestDelete, "Generic File - Request Deletion"},
	"GenericFileRequestRestore":         {"GenericFile", constants.FileRestore, "Generic File - Request Restoration"},
	"GenericFileShow":                   {"GenericFile", constants.FileRead, "Generic File Detail"},
	"GenericFileUpdate":                 {"GenericFile", constants.FileUpdate, "Update Generic File"},
	"InstitutionCreate":                 {"Institution", constants.InstitutionCreate, "Create Institution"},
	"InstitutionDelete":                 {"Institution", constants.InstitutionDelete, "Deactivate Institution"},
	"InstitutionEdit":                   {"Institution", constants.InstitutionUpdate, "Edit Institution"},
	"InstitutionEditPrefs":              {"Institution", constants.InstitutionUpdatePrefs, "Edit Institution Preferences"},
	"InstitutionIndex":                  {"Institution", constants.InstitutionList, "Institutions"},
	"InstitutionNew":                    {"Institution", constants.InstitutionCreate, "New Institution"},
	"InstitutionShow":                   {"Institution", constants.InstitutionRead, "Institution Detail"},
	"InstitutionUndelete":               {"Institution", constants.InstitutionUpdate, "Reactivate Institution"},
	"InstitutionUpdate":                 {"Institution", constants.InstitutionUpdate, "Update Institution"},
	"InstitutionUpdatePrefs":            {"Institution", constants.InstitutionUpdatePrefs, "Update Institution Preferences"},
	"IntellectualObjectInitBatchDelete": {"IntellectualObject", constants.IntellectualObjectBatchDelete, "Intellectual Object Batch Delete"},
	"IntellectualObjectCreate":          {"IntellectualObject", constants.IntellectualObjectCreate, "Create Intellectual Object"},
	"IntellectualObjectDelete":          {"IntellectualObject", constants.IntellectualObjectDelete, "Delete Intellectual Object"},
	"IntellectualObjectEvents":          {"PremisEvent", constants.EventRead, "PREMIS Events"},
	// IntellectualObjectFiles gets an object ID and will look up that object to check
	// it's institution. The permission, however, is FileReade, because this endpoint
	// returns files. https://trello.com/c/n5asx3bj
	"IntellectualObjectFiles":            {"IntellectualObject", constants.FileRead, "Object Files"},
	"IntellectualObjectFinishBulkDelete": {"IntellectualObject", constants.IntellectualObjectFinishBulkDelete, "Intellectual Object - Finish Bulk Delete"},
	"IntellectualObjectIndex":            {"IntellectualObject", constants.IntellectualObjectRead, "Intellectual Objects"},
	"IntellectualObjectInitDelete":       {"IntellectualObject", constants.IntellectualObjectRequestDelete, "Initialize Object Deletion"},
	"IntellectualObjectInitRestore":      {"IntellectualObject", constants.IntellectualObjectRestore, "Initialize Object Restoration"},
	"IntellectualObjectNew":              {"IntellectualObject", constants.IntellectualObjectCreate, "New Intellectual Object"},
	"IntellectualObjectRequestDelete":    {"IntellectualObject", constants.IntellectualObjectRequestDelete, "Request Object Deletion"},
	"IntellectualObjectRequestRestore":   {"IntellectualObject", constants.IntellectualObjectRestore, "Request Object Restoration"},
	"IntellectualObjectShow":             {"IntellectualObject", constants.IntellectualObjectRead, "Intellectual Object Detail"},
	"IntellectualObjectUpdate":           {"IntellectualObject", constants.IntellectualObjectUpdate, "Update Intellectual Object"},
	"InternalMetadataIndex":              {"InternalMetadata", constants.InternalMetadataRead, "Internal Metadata"},
	"NsqShow":                            {"NSQ", constants.NsqAdmin, "NSQ Dashboard"},
	"NsqAdmin":                           {"NSQ", constants.NsqAdmin, "NSQ Admin"},
	"NsqInit":                            {"NSQ", constants.NsqAdmin, "NSQ"},
	"PremisEventCreate":                  {"PremisEvent", constants.EventCreate, "Create PREMIS Event"},
	"PremisEventIndex":                   {"PremisEvent", constants.EventRead, "PREMIS Events"},
	"PremisEventShow":                    {"PremisEvent", constants.EventRead, "PREMIS Event Detail"},
	"PremisEventShowXHR":                 {"PremisEvent", constants.EventRead, "PREMIS Event Detail"},
	"PrepareFileDelete":                  {"GenericFile", constants.PrepareFileDelete, "Prepare File Deletion"},
	"PrepareObjectDelete":                {"IntellectualObject", constants.PrepareObjectDelete, "Prepare Object Deletion"},
	"StorageRecordCreate":                {"StorageRecord", constants.StorageRecordCreate, "Create Storage Record"},
	"StorageRecordDelete":                {"StorageRecord", constants.StorageRecordDelete, "Delete Storage Record"},
	"StorageRecordIndex":                 {"StorageRecord", constants.StorageRecordRead, "Storage Records"},
	"StorageRecordNew":                   {"StorageRecord", constants.StorageRecordCreate, "New Storage Record"},
	"StorageRecordShow":                  {"StorageRecord", constants.StorageRecordRead, "Storage Record Detail"},
	"StorageRecordUpdate":                {"StorageRecord", constants.StorageRecordUpdate, "Update Storage Record"},
	"UserBeginLoginWithPasskey":          {"User", constants.UserBeginLoginWithPasskey, "Login With Passkey"},
	"UserBeginPasskeyRegistration":       {"User", constants.UserBeginPasskeyRegistration, "Set up a Passkey"},
	"UserChangePassword":                 {"User", constants.UserUpdateSelf, "Change Password"},
	"UserComplete2FASetup":               {"User", constants.UserComplete2FASetup, "Setup Two-Factor Authentication"},
	"UserConfirmPhone":                   {"User", constants.UserConfirmPhone, "Confirm Phone Number"},
	"UserCreate":                         {"User", constants.UserCreate, "Create User"},
	"UserDelete":                         {"User", constants.UserDelete, "Deactivate User"},
	"UserDeleteSelf":                     {"User", constants.UserDeleteSelf, "Deactivate Your Account"},
	"UserEdit":                           {"User", constants.UserUpdate, "Edit User"},
	"UserFinishLoginWithPasskey":         {"User", constants.UserFinishLoginWithPasskey, "Complete Passkey Login"},
	"UserFinishPasskeyRegistration":      {"User", constants.UserFinishPasskeyRegistration, "Complete Passkey Setup"},
	"UserGenerateBackupCodes":            {"User", constants.UserGenerateBackupCodes, "Create Two-Factor Backup Codes"},
	"UserGetAPIKey":                      {"User", constants.UserUpdateSelf, "Generate API Key"},
	"UserIndex":                          {"User", constants.UserRead, "Users"},
	"UserInit2FASetup":                   {"User", constants.UserInit2FASetup, "Start Two-Factor Setup"},
	"UserInitPasswordReset":              {"User", constants.UserUpdate, "Reset Password"},
	"UserMyAccount":                      {"User", constants.UserUpdateSelf, "My Account"},
	"UserNew":                            {"User", constants.UserCreate, "New User"},
	"UserReadSelf":                       {"User", constants.UserReadSelf, "User Detail"},
	"UserShow":                           {"User", constants.UserRead, "User Detail"},
	"UserShowChangePassword":             {"User", constants.UserUpdateSelf, "Change Password"},
	"UserTwoFactorBackup":                {"User", constants.UserTwoFactorBackup, "Generate Backup Codes"},
	"UserTwoFactorChoose":                {"User", constants.UserTwoFactorChoose, "Choose Two-Factor Method"},
	"UserTwoFactorGenerateSMS":           {"User", constants.UserTwoFactorGenerateSMS, "Generate SMS Message"},
	"UserTwoFactorPush":                  {"User", constants.UserTwoFactorPush, "Send Push Message"},
	"UserTwoFactorResend":                {"User", constants.UserTwoFactorResend, "Resent Two-Factor Token"},
	"UserTwoFactorVerify":                {"User", constants.UserTwoFactorVerify, "Verify Two-Factor Authentication Method"},
	"UserUndelete":                       {"User", constants.UserUpdate, "Reactivate User"},
	"UserUpdate":                         {"User", constants.UserUpdate, "Update User"},
	"UserUpdateXHR":                      {"User", constants.UserUpdate, "Update User"},
	"UserUpdateSelf":                     {"User", constants.UserUpdateSelf, "Update User"},
	"WorkItemCreate":                     {"WorkItem", constants.WorkItemCreate, "Create Work Item"},
	"WorkItemDelete":                     {"WorkItem", constants.WorkItemDelete, "Delete Work Item"},
	"WorkItemEdit":                       {"WorkItem", constants.WorkItemUpdate, "Edit Work Item"},
	"WorkItemIndex":                      {"WorkItem", constants.WorkItemRead, "Work Items"},
	"WorkItemNew":                        {"WorkItem", constants.WorkItemCreate, "New Work Item"},
	"WorkItemRedisDelete":                {"WorkItem", constants.WorkItemRedisDelete, "Delete Redis Data"},
	"WorkItemRedisIndex":                 {"WorkItem", constants.RedisList, "Redis Data"},
	"WorkItemRequeue":                    {"WorkItem", constants.WorkItemRequeue, "Requeue Work Item"},
	"WorkItemShow":                       {"WorkItem", constants.WorkItemRead, "Work Item Detail"},
	"WorkItemShowRequeue":                {"WorkItem", constants.WorkItemRequeue, "Requeue Work Item"},
	"WorkItemUpdate":                     {"WorkItem", constants.WorkItemUpdate, "Update Work Item"},
}
