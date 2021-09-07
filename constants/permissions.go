package constants

// Permission is a string that keys into permission maps for
// different roles. We use string instead of bitmask or array index
// for a few reasons:
//
// 1. We may wind up with more than 64 of these, which would be too
//    many for a bitmask.
// 2. Our models need to construct permission names from strings made
//    up of model names and actions. E.g. "User" + "Create" or
//    "Object" + "Read".
// 3. We will likely insert permissions as the application grows, and
//    adding bitmasks and array indices is order-dependent, while
//    adding string keys is value-dependent, which means we can insert
//    them anywhere in the list.
type Permission string

const (
	AlertCreate                        = "AlertCreate"
	AlertDelete                        = "AlertDelete"
	AlertRead                          = "AlertRead"
	AlertUpdate                        = "AlertUpdate"
	ChecksumCreate                     = "ChecksumCreate"
	ChecksumDelete                     = "ChecksumDelete"
	ChecksumRead                       = "ChecksumRead"
	ChecksumUpdate                     = "ChecksumUpdate"
	DashboardShow                      = "DashboardShow"
	DeletionRequestApprove             = "DeletionRequestApprove"
	DeletionRequestList                = "DeletionRequestList"
	DeletionRequestShow                = "DeletionRequestShow"
	DepositReportShow                  = "DepositReportShow"
	EventCreate                        = "EventCreate"
	EventDelete                        = "EventDelete"
	EventRead                          = "EventRead"
	EventUpdate                        = "EventUpdate"
	FileCreate                         = "FileCreate"
	FileDelete                         = "FileDelete"
	FileFinishBulkDelete               = "FileFinishBulkDelete"
	FileRead                           = "FileRead"
	FileRequestDelete                  = "FileRequestDelete"
	FileRestore                        = "FileRestore"
	FileUpdate                         = "FileUpdate"
	InstitutionCreate                  = "InstitutionCreate"
	InstitutionDelete                  = "InstitutionDelete"
	InstitutionList                    = "InstitutionList"
	InstitutionRead                    = "InstitutionRead"
	InstitutionUpdate                  = "InstitutionUpdate"
	IntellectualObjectCreate           = "IntellectualObjectCreate"
	IntellectualObjectDelete           = "IntellectualObjectDelete"
	IntellectualObjectFinishBulkDelete = "IntellectualObjectFinishBulkDelete"
	IntellectualObjectRead             = "IntellectualObjectRead"
	IntellectualObjectRequestDelete    = "IntellectualObjectRequestDelete"
	IntellectualObjectRestore          = "IntellectualObjectRestore"
	IntellectualObjectUpdate           = "IntellectualObjectUpdate"
	ReportRead                         = "ReportRead"
	StorageRecordCreate                = "StorageRecordCreate"
	StorageRecordDelete                = "StorageRecordDelete"
	StorageRecordRead                  = "StorageRecordRead"
	StorageRecordUpdate                = "StorageRecordUpdate"
	UserComplete2FASetup               = "UserComplete2FASetup"
	UserConfirmPhone                   = "UserConfirmPhone"
	UserCreate                         = "UserCreate"
	UserDelete                         = "UserDelete"
	UserDeleteSelf                     = "UserDeleteSelf"
	UserGenerateBackupCodes            = "UserGenerateBackupCodes"
	UserInit2FASetup                   = "UserInit2FASetup"
	UserRead                           = "UserRead"
	UserReadSelf                       = "UserReadSelf"
	UserSignIn                         = "UserSignIn"
	UserSignOut                        = "UserSignOut"
	UserTwoFactorChoose                = "UserTwoFactorChoose"
	UserTwoFactorGenerateSMS           = "UserTwoFactorGenerateSMS"
	UserTwoFactorPush                  = "UserTwoFactorPush"
	UserTwoFactorResend                = "UserTwoFactorResend"
	UserTwoFactorVerify                = "UserTwoFactorVerify"
	UserUpdate                         = "UserUpdate"
	UserUpdateSelf                     = "UserUpdateSelf"
	WorkItemCreate                     = "WorkItemCreate"
	WorkItemDelete                     = "WorkItemDelete"
	WorkItemRead                       = "WorkItemRead"
	WorkItemRequeue                    = "WorkItemRequeue"
	WorkItemUpdate                     = "WorkItemUpdate"
)

var Permissions = []Permission{
	AlertCreate,
	AlertDelete,
	AlertRead,
	AlertUpdate,
	ChecksumCreate,
	ChecksumDelete,
	ChecksumRead,
	ChecksumUpdate,
	DashboardShow,
	DeletionRequestApprove,
	DeletionRequestList,
	DeletionRequestShow,
	DepositReportShow,
	EventCreate,
	EventDelete,
	EventRead,
	EventUpdate,
	FileCreate,
	FileDelete,
	FileFinishBulkDelete,
	FileRead,
	FileRequestDelete,
	FileRestore,
	FileUpdate,
	InstitutionCreate,
	InstitutionDelete,
	InstitutionList,
	InstitutionRead,
	InstitutionUpdate,
	IntellectualObjectCreate,
	IntellectualObjectDelete,
	IntellectualObjectFinishBulkDelete,
	IntellectualObjectRead,
	IntellectualObjectRequestDelete,
	IntellectualObjectRestore,
	IntellectualObjectUpdate,
	ReportRead,
	StorageRecordCreate,
	StorageRecordDelete,
	StorageRecordRead,
	StorageRecordUpdate,
	UserComplete2FASetup,
	UserConfirmPhone,
	UserCreate,
	UserDelete,
	UserDeleteSelf,
	UserGenerateBackupCodes,
	UserInit2FASetup,
	UserRead,
	UserReadSelf,
	UserSignIn,
	UserSignOut,
	UserTwoFactorChoose,
	UserTwoFactorGenerateSMS,
	UserTwoFactorPush,
	UserTwoFactorResend,
	UserTwoFactorVerify,
	UserUpdate,
	UserUpdateSelf,
	WorkItemCreate,
	WorkItemDelete,
	WorkItemRead,
	WorkItemRequeue,
	WorkItemUpdate,
}

var ForbiddenToAll = []Permission{
	ChecksumUpdate,
	ChecksumDelete,
	EventUpdate,
	EventDelete,
}

// SelfAccountPermissions are those which permit users to modify their
// own account information (password, email, name, API key, etc.).
// These are treated specially in the auth middleware because unlike
// other permissions, which are based on the subject resource's insitution id,
// these are based on the subject resource's user id. See
// ResourceAuthorization.checkPermission to understand how this specific
// set of permissions is used.
var SelfAccountPermissions = []Permission{
	UserComplete2FASetup,
	UserConfirmPhone,
	UserDeleteSelf,
	UserGenerateBackupCodes,
	UserInit2FASetup,
	UserReadSelf,
	UserUpdateSelf,
}

var permissionsInitialized = false

// Permission lists for different roles. Bools default to false in Go,
// so roles will have only the permissions we explicitly grant below.
// Note that emptyList gets zero permissions. This is the list we'll check
// if we can't figure out a user's role, or if a user has been deactivated.
var instUser = make(map[Permission]bool)
var instAdmin = make(map[Permission]bool)
var sysAdmin = make(map[Permission]bool)
var emptyList = make(map[Permission]bool)

func initPermissions() {
	// Institutional User Role
	instUser[AlertRead] = true
	instUser[AlertUpdate] = true
	instUser[ChecksumRead] = true
	instUser[DashboardShow] = true
	instUser[DeletionRequestList] = true
	instUser[DeletionRequestShow] = true
	instUser[DepositReportShow] = true
	instUser[EventRead] = true
	instUser[FileRead] = true
	instUser[FileRequestDelete] = true
	instUser[FileRestore] = true
	instUser[InstitutionRead] = true
	instUser[IntellectualObjectRead] = true
	instUser[IntellectualObjectRequestDelete] = true
	instUser[IntellectualObjectRestore] = true
	instUser[ReportRead] = true
	instUser[StorageRecordRead] = true
	instUser[UserComplete2FASetup] = true
	instUser[UserConfirmPhone] = true
	instUser[UserGenerateBackupCodes] = true
	instUser[UserInit2FASetup] = true
	instUser[UserReadSelf] = true
	instUser[UserSignIn] = true
	instUser[UserSignOut] = true
	instUser[UserTwoFactorChoose] = true
	instUser[UserTwoFactorGenerateSMS] = true
	instUser[UserTwoFactorPush] = true
	instUser[UserTwoFactorResend] = true
	instUser[UserTwoFactorVerify] = true
	instUser[UserUpdateSelf] = true
	instUser[WorkItemRead] = true

	// Institutional Admin Role
	instAdmin[AlertRead] = true
	instAdmin[AlertUpdate] = true
	instAdmin[ChecksumRead] = true
	instAdmin[DashboardShow] = true
	instAdmin[DeletionRequestApprove] = true
	instAdmin[DeletionRequestList] = true
	instAdmin[DeletionRequestShow] = true
	instAdmin[DepositReportShow] = true
	instAdmin[EventRead] = true
	instAdmin[FileDelete] = true
	instAdmin[FileRead] = true
	instAdmin[FileRequestDelete] = true
	instAdmin[FileRestore] = true
	instAdmin[InstitutionRead] = true
	instAdmin[IntellectualObjectDelete] = true
	instAdmin[IntellectualObjectRead] = true
	instAdmin[IntellectualObjectRequestDelete] = true
	instAdmin[IntellectualObjectRestore] = true
	instAdmin[ReportRead] = true
	instAdmin[StorageRecordRead] = true
	instAdmin[UserComplete2FASetup] = true
	instAdmin[UserConfirmPhone] = true
	instAdmin[UserCreate] = true
	instAdmin[UserDelete] = true
	instAdmin[UserGenerateBackupCodes] = true
	instAdmin[UserInit2FASetup] = true
	instAdmin[UserReadSelf] = true
	instAdmin[UserRead] = true
	instAdmin[UserSignIn] = true
	instAdmin[UserSignOut] = true
	instAdmin[UserTwoFactorChoose] = true
	instAdmin[UserTwoFactorGenerateSMS] = true
	instAdmin[UserTwoFactorPush] = true
	instAdmin[UserTwoFactorResend] = true
	instAdmin[UserTwoFactorVerify] = true
	instAdmin[UserUpdateSelf] = true
	instAdmin[UserUpdate] = true
	instAdmin[WorkItemRead] = true

	// Sys Admin Role
	sysAdmin[AlertCreate] = true
	sysAdmin[AlertDelete] = true
	sysAdmin[AlertRead] = true
	sysAdmin[AlertUpdate] = true
	sysAdmin[ChecksumCreate] = true
	sysAdmin[ChecksumDelete] = false // no one can do this
	sysAdmin[ChecksumRead] = true
	sysAdmin[ChecksumUpdate] = false // no one can do this
	sysAdmin[DashboardShow] = true
	sysAdmin[DeletionRequestApprove] = true
	sysAdmin[DeletionRequestList] = true
	sysAdmin[DeletionRequestShow] = true
	sysAdmin[DepositReportShow] = true
	sysAdmin[EventCreate] = true
	sysAdmin[EventDelete] = false // no one can do this
	sysAdmin[EventRead] = true
	sysAdmin[EventUpdate] = false // no one can do this
	sysAdmin[FileCreate] = true
	sysAdmin[FileDelete] = true
	sysAdmin[FileFinishBulkDelete] = true
	sysAdmin[FileRead] = true
	sysAdmin[FileRequestDelete] = true
	sysAdmin[FileRestore] = true
	sysAdmin[FileUpdate] = true
	sysAdmin[InstitutionCreate] = true
	sysAdmin[InstitutionDelete] = true
	sysAdmin[InstitutionList] = true
	sysAdmin[InstitutionRead] = true
	sysAdmin[InstitutionUpdate] = true
	sysAdmin[IntellectualObjectCreate] = true
	sysAdmin[IntellectualObjectDelete] = true
	sysAdmin[IntellectualObjectFinishBulkDelete] = true
	sysAdmin[IntellectualObjectRead] = true
	sysAdmin[IntellectualObjectRequestDelete] = true
	sysAdmin[IntellectualObjectRestore] = true
	sysAdmin[IntellectualObjectUpdate] = true
	sysAdmin[ReportRead] = true
	sysAdmin[StorageRecordCreate] = true
	sysAdmin[StorageRecordDelete] = true
	sysAdmin[StorageRecordRead] = true
	sysAdmin[StorageRecordUpdate] = true
	sysAdmin[UserComplete2FASetup] = true
	sysAdmin[UserConfirmPhone] = true
	sysAdmin[UserCreate] = true
	sysAdmin[UserDeleteSelf] = true
	sysAdmin[UserDelete] = true
	sysAdmin[UserGenerateBackupCodes] = true
	sysAdmin[UserInit2FASetup] = true
	sysAdmin[UserReadSelf] = true
	sysAdmin[UserRead] = true
	sysAdmin[UserSignIn] = true
	sysAdmin[UserSignOut] = true
	sysAdmin[UserTwoFactorChoose] = true
	sysAdmin[UserTwoFactorGenerateSMS] = true
	sysAdmin[UserTwoFactorPush] = true
	sysAdmin[UserTwoFactorResend] = true
	sysAdmin[UserTwoFactorVerify] = true
	sysAdmin[UserUpdateSelf] = true
	sysAdmin[UserUpdate] = true
	sysAdmin[WorkItemCreate] = true
	sysAdmin[WorkItemDelete] = true
	sysAdmin[WorkItemRead] = true
	sysAdmin[WorkItemRequeue] = true
	sysAdmin[WorkItemUpdate] = true

	permissionsInitialized = true
}

func CheckPermission(role string, permission Permission) bool {
	if !permissionsInitialized {
		initPermissions()
	}
	var permissions map[Permission]bool
	switch role {
	case RoleSysAdmin:
		permissions = sysAdmin
	case RoleInstAdmin:
		permissions = instAdmin
	case RoleInstUser:
		permissions = instUser
	default:
		permissions = emptyList
	}
	return permissions[permission]
}
