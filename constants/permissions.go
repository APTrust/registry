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
	AlertRead                          = "AlertRead"
	AlertUpdate                        = "AlertUpdate"
	AlertDelete                        = "AlertDelete"
	ChecksumCreate                     = "ChecksumCreate"
	ChecksumRead                       = "ChecksumRead"
	ChecksumUpdate                     = "ChecksumUpdate"
	ChecksumDelete                     = "ChecksumDelete"
	DashboardShow                      = "DashboardShow"
	DeletionRequestApprove             = "DeletionRequestApprove"
	DeletionRequestList                = "DeletionRequestList"
	DeletionRequestShow                = "DeletionRequestShow"
	DepositReportShow                  = "DepositReportShow"
	EventCreate                        = "EventCreate"
	EventRead                          = "EventRead"
	EventUpdate                        = "EventUpdate"
	EventDelete                        = "EventDelete"
	FileCreate                         = "FileCreate"
	FileRead                           = "FileRead"
	FileUpdate                         = "FileUpdate"
	FileDelete                         = "FileDelete"
	FileRequestDelete                  = "FileRequestDelete"
	FileFinishBulkDelete               = "FileFinishBulkDelete"
	FileRestore                        = "FileRestore"
	InstitutionCreate                  = "InstitutionCreate"
	InstitutionList                    = "InstitutionList"
	InstitutionRead                    = "InstitutionRead"
	InstitutionUpdate                  = "InstitutionUpdate"
	InstitutionDelete                  = "InstitutionDelete"
	IntellectualObjectCreate           = "IntellectualObjectCreate"
	IntellectualObjectRead             = "IntellectualObjectRead"
	IntellectualObjectUpdate           = "IntellectualObjectUpdate"
	IntellectualObjectDelete           = "IntellectualObjectDelete"
	IntellectualObjectRequestDelete    = "IntellectualObjectRequestDelete"
	IntellectualObjectFinishBulkDelete = "IntellectualObjectFinishBulkDelete"
	IntellectualObjectRestore          = "IntellectualObjectRestore"
	ReportRead                         = "ReportRead"
	StorageRecordCreate                = "StorageRecordCreate"
	StorageRecordRead                  = "StorageRecordRead"
	StorageRecordUpdate                = "StorageRecordUpdate"
	StorageRecordDelete                = "StorageRecordDelete"
	UserCreate                         = "UserCreate"
	UserRead                           = "UserRead"
	UserSignIn                         = "UserSignIn"
	UserSignOut                        = "UserSignOut"
	UserUpdate                         = "UserUpdate"
	UserDelete                         = "UserDelete"
	UserReadSelf                       = "UserReadSelf"
	UserUpdateSelf                     = "UserUpdateSelf"
	UserDeleteSelf                     = "UserDeleteSelf"
	UserTwoFactorChoose                = "UserTwoFactorChoose"
	UserTwoFactorGenerateSMS           = "UserTwoFactorGenerateSMS"
	UserTwoFactorPush                  = "UserTwoFactorPush"
	UserTwoFactorResend                = "UserTwoFactorResend"
	UserTwoFactorVerify                = "UserTwoFactorVerify"
	WorkItemCreate                     = "WorkItemCreate"
	WorkItemRead                       = "WorkItemRead"
	WorkItemRequeue                    = "WorkItemRequeue"
	WorkItemUpdate                     = "WorkItemUpdate"
	WorkItemDelete                     = "WorkItemDelete"
)

var Permissions = []Permission{
	AlertCreate,
	AlertRead,
	AlertUpdate,
	AlertDelete,
	ChecksumCreate,
	ChecksumRead,
	ChecksumUpdate,
	ChecksumDelete,
	DashboardShow,
	DeletionRequestApprove,
	DeletionRequestList,
	DeletionRequestShow,
	DepositReportShow,
	EventCreate,
	EventRead,
	EventUpdate,
	EventDelete,
	FileCreate,
	FileRead,
	FileUpdate,
	FileDelete,
	FileRequestDelete,
	FileFinishBulkDelete,
	FileRestore,
	InstitutionCreate,
	InstitutionList,
	InstitutionRead,
	InstitutionUpdate,
	InstitutionDelete,
	IntellectualObjectCreate,
	IntellectualObjectRead,
	IntellectualObjectUpdate,
	IntellectualObjectDelete,
	IntellectualObjectRequestDelete,
	IntellectualObjectFinishBulkDelete,
	IntellectualObjectRestore,
	ReportRead,
	StorageRecordCreate,
	StorageRecordRead,
	StorageRecordUpdate,
	StorageRecordDelete,
	UserCreate,
	UserRead,
	UserSignIn,
	UserSignOut,
	UserTwoFactorChoose,
	UserTwoFactorGenerateSMS,
	UserTwoFactorPush,
	UserTwoFactorResend,
	UserTwoFactorVerify,
	UserUpdate,
	UserDelete,
	UserReadSelf,
	UserUpdateSelf,
	UserDeleteSelf,
	WorkItemCreate,
	WorkItemRead,
	WorkItemRequeue,
	WorkItemUpdate,
	WorkItemDelete,
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
	UserReadSelf,
	UserUpdateSelf,
	UserDeleteSelf,
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
	instUser[DeletionRequestShow] = true
	instUser[DeletionRequestList] = true
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
	instUser[UserSignIn] = true
	instUser[UserSignOut] = true
	instUser[UserReadSelf] = true
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
	instAdmin[DeletionRequestShow] = true
	instAdmin[DeletionRequestList] = true
	instAdmin[DepositReportShow] = true
	instAdmin[EventRead] = true
	instAdmin[FileRead] = true
	instAdmin[FileDelete] = true
	instAdmin[FileRequestDelete] = true
	instAdmin[FileRestore] = true
	instAdmin[InstitutionRead] = true
	instAdmin[IntellectualObjectRead] = true
	instAdmin[IntellectualObjectDelete] = true
	instAdmin[IntellectualObjectRequestDelete] = true
	instAdmin[IntellectualObjectRestore] = true
	instAdmin[ReportRead] = true
	instAdmin[StorageRecordRead] = true
	instAdmin[UserCreate] = true
	instAdmin[UserRead] = true
	instAdmin[UserSignIn] = true
	instAdmin[UserSignOut] = true
	instAdmin[UserUpdate] = true
	instAdmin[UserDelete] = true
	instAdmin[UserReadSelf] = true
	instAdmin[UserTwoFactorChoose] = true
	instAdmin[UserTwoFactorGenerateSMS] = true
	instAdmin[UserTwoFactorPush] = true
	instAdmin[UserTwoFactorResend] = true
	instAdmin[UserTwoFactorVerify] = true
	instAdmin[UserUpdateSelf] = true
	instAdmin[WorkItemRead] = true

	// Sys Admin Role
	sysAdmin[AlertCreate] = true
	sysAdmin[AlertRead] = true
	sysAdmin[AlertUpdate] = true
	sysAdmin[AlertDelete] = true
	sysAdmin[ChecksumCreate] = true
	sysAdmin[ChecksumRead] = true
	sysAdmin[ChecksumUpdate] = false // no one can do this
	sysAdmin[ChecksumDelete] = false // no one can do this
	sysAdmin[DashboardShow] = true
	sysAdmin[DeletionRequestApprove] = true
	sysAdmin[DeletionRequestShow] = true
	sysAdmin[DeletionRequestList] = true
	sysAdmin[DepositReportShow] = true
	sysAdmin[EventCreate] = true
	sysAdmin[EventRead] = true
	sysAdmin[EventUpdate] = false // no one can do this
	sysAdmin[EventDelete] = false // no one can do this
	sysAdmin[FileCreate] = true
	sysAdmin[FileRead] = true
	sysAdmin[FileUpdate] = true
	sysAdmin[FileDelete] = true
	sysAdmin[FileRequestDelete] = true
	sysAdmin[FileFinishBulkDelete] = true
	sysAdmin[FileRestore] = true
	sysAdmin[InstitutionCreate] = true
	sysAdmin[InstitutionList] = true
	sysAdmin[InstitutionRead] = true
	sysAdmin[InstitutionUpdate] = true
	sysAdmin[InstitutionDelete] = true
	sysAdmin[IntellectualObjectCreate] = true
	sysAdmin[IntellectualObjectRead] = true
	sysAdmin[IntellectualObjectUpdate] = true
	sysAdmin[IntellectualObjectDelete] = true
	sysAdmin[IntellectualObjectRequestDelete] = true
	sysAdmin[IntellectualObjectFinishBulkDelete] = true
	sysAdmin[IntellectualObjectRestore] = true
	sysAdmin[ReportRead] = true
	sysAdmin[StorageRecordCreate] = true
	sysAdmin[StorageRecordRead] = true
	sysAdmin[StorageRecordUpdate] = true
	sysAdmin[StorageRecordDelete] = true
	sysAdmin[UserCreate] = true
	sysAdmin[UserRead] = true
	sysAdmin[UserSignIn] = true
	sysAdmin[UserSignOut] = true
	sysAdmin[UserTwoFactorChoose] = true
	sysAdmin[UserTwoFactorGenerateSMS] = true
	sysAdmin[UserTwoFactorPush] = true
	sysAdmin[UserTwoFactorResend] = true
	sysAdmin[UserTwoFactorVerify] = true
	sysAdmin[UserUpdate] = true
	sysAdmin[UserDelete] = true
	sysAdmin[UserDeleteSelf] = true
	sysAdmin[UserReadSelf] = true
	sysAdmin[UserUpdateSelf] = true
	sysAdmin[WorkItemCreate] = true
	sysAdmin[WorkItemRead] = true
	sysAdmin[WorkItemRequeue] = true
	sysAdmin[WorkItemUpdate] = true
	sysAdmin[WorkItemDelete] = true

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
