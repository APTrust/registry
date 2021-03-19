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
	AlertCreate            Permission = "AlertCreate"
	AlertRead                         = "AlertRead"
	AlertUpdate                       = "AlertUpdate"
	AlertDelete                       = "AlertDelete"
	ChecksumCreate                    = "ChecksumCreate"
	ChecksumRead                      = "ChecksumRead"
	ChecksumUpdate                    = "ChecksumUpdate"
	ChecksumDelete                    = "ChecksumDelete"
	DashboardShow                     = "DashboardShow"
	EventCreate                       = "EventCreate"
	EventRead                         = "EventRead"
	EventUpdate                       = "EventUpdate"
	EventDelete                       = "EventDelete"
	FileCreate                        = "FileCreate"
	FileRead                          = "FileRead"
	FileUpdate                        = "FileUpdate"
	FileDelete                        = "FileDelete"
	FileRequestDelete                 = "FileRequestDelete"
	FileApproveDelete                 = "FileApproveDelete"
	FileFinishBulkDelete              = "FileFinishBulkDelete"
	FileRestore                       = "FileRestore"
	InstitutionCreate                 = "InstitutionCreate"
	InstitutionRead                   = "InstitutionRead"
	InstitutionUpdate                 = "InstitutionUpdate"
	InstitutionDelete                 = "InstitutionDelete"
	ObjectCreate                      = "ObjectCreate"
	ObjectRead                        = "ObjectRead"
	ObjectUpdate                      = "ObjectUpdate"
	ObjectDelete                      = "ObjectDelete"
	ObjectRequestDelete               = "ObjectRequestDelete"
	ObjectApproveDelete               = "ObjectApproveDelete"
	ObjectFinishBulkDelete            = "ObjectFinishBulkDelete"
	ObjectRestore                     = "ObjectRestore"
	StorageRecordCreate               = "StorageRecordCreate"
	StorageRecordRead                 = "StorageRecordRead"
	StorageRecordUpdate               = "StorageRecordUpdate"
	StorageRecordDelete               = "StorageRecordDelete"
	UserCreate                        = "UserCreate"
	UserRead                          = "UserRead"
	UserSignIn                        = "UserSignIn"
	UserSignOut                       = "UserSignOut"
	UserUpdate                        = "UserUpdate"
	UserDelete                        = "UserDelete"
	UserReadSelf                      = "UserReadSelf"
	UserUpdateSelf                    = "UserUpdateSelf"
	UserDeleteSelf                    = "UserDeleteSelf"
	WorkItemCreate                    = "WorkItemCreate"
	WorkItemRead                      = "WorkItemRead"
	WorkItemUpdate                    = "WorkItemUpdate"
	WorkItemDelete                    = "WorkItemDelete"
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
	EventCreate,
	EventRead,
	EventUpdate,
	EventDelete,
	FileCreate,
	FileRead,
	FileUpdate,
	FileDelete,
	FileRequestDelete,
	FileApproveDelete,
	FileFinishBulkDelete,
	FileRestore,
	InstitutionCreate,
	InstitutionRead,
	InstitutionUpdate,
	InstitutionDelete,
	ObjectCreate,
	ObjectRead,
	ObjectUpdate,
	ObjectDelete,
	ObjectRequestDelete,
	ObjectApproveDelete,
	ObjectFinishBulkDelete,
	ObjectRestore,
	StorageRecordCreate,
	StorageRecordRead,
	StorageRecordUpdate,
	StorageRecordDelete,
	UserCreate,
	UserRead,
	UserSignIn,
	UserSignOut,
	UserUpdate,
	UserDelete,
	UserReadSelf,
	UserUpdateSelf,
	UserDeleteSelf,
	WorkItemCreate,
	WorkItemRead,
	WorkItemUpdate,
	WorkItemDelete,
}

var ForbiddenToAll = []Permission{
	ChecksumUpdate,
	ChecksumDelete,
	EventUpdate,
	EventDelete,
}

var permissionsInitialized = false
var permissionCount = len(Permissions)

// Permission lists for different roles. Bools default to false in Go,
// so roles will have only the permissions we explicitly grant below.
// Note that emptyList gets zero permissions. This is the list we'll check
// if we can't figure out a user's role, or if a user has been deactivated.
var instUser = make(map[Permission]bool)
var instAdmin = make(map[Permission]bool)
var sysAdmin = make(map[Permission]bool)
var emptyList = make(map[Permission]bool)

func initPermissions() {
	instUser[AlertRead] = true
	instUser[AlertUpdate] = true
	instUser[ChecksumRead] = true
	instUser[DashboardShow] = true
	instUser[EventRead] = true
	instUser[FileRead] = true
	instUser[InstitutionRead] = true
	instUser[ObjectRead] = true
	instUser[StorageRecordRead] = true
	instUser[UserSignIn] = true
	instUser[UserSignOut] = true
	instUser[UserReadSelf] = true
	instUser[UserUpdateSelf] = true
	instUser[WorkItemRead] = true

	instAdmin[AlertRead] = true
	instAdmin[AlertUpdate] = true
	instAdmin[ChecksumRead] = true
	instAdmin[DashboardShow] = true
	instAdmin[EventRead] = true
	instAdmin[FileRead] = true
	instAdmin[FileDelete] = true
	instAdmin[FileRequestDelete] = true
	instAdmin[FileApproveDelete] = true
	instAdmin[FileRestore] = true
	instAdmin[InstitutionRead] = true
	instAdmin[ObjectRead] = true
	instAdmin[ObjectDelete] = true
	instAdmin[ObjectRequestDelete] = true
	instAdmin[ObjectApproveDelete] = true
	instAdmin[ObjectRestore] = true
	instAdmin[StorageRecordRead] = true
	instAdmin[UserCreate] = true
	instAdmin[UserRead] = true
	instAdmin[UserSignIn] = true
	instAdmin[UserSignOut] = true
	instAdmin[UserUpdate] = true
	instAdmin[UserDelete] = true
	instAdmin[UserReadSelf] = true
	instAdmin[UserUpdateSelf] = true
	instAdmin[WorkItemRead] = true

	sysAdmin[AlertCreate] = true
	sysAdmin[AlertRead] = true
	sysAdmin[AlertUpdate] = true
	sysAdmin[AlertDelete] = true
	sysAdmin[ChecksumCreate] = true
	sysAdmin[ChecksumRead] = true
	sysAdmin[ChecksumUpdate] = false // no one can do this
	sysAdmin[ChecksumDelete] = false // no one can do this
	sysAdmin[DashboardShow] = true
	sysAdmin[EventCreate] = true
	sysAdmin[EventRead] = true
	sysAdmin[EventUpdate] = false // no one can do this
	sysAdmin[EventDelete] = false // no one can do this
	sysAdmin[FileCreate] = true
	sysAdmin[FileRead] = true
	sysAdmin[FileUpdate] = true
	sysAdmin[FileDelete] = true
	sysAdmin[FileRequestDelete] = true
	sysAdmin[FileApproveDelete] = true
	sysAdmin[FileFinishBulkDelete] = true
	sysAdmin[FileRestore] = true
	sysAdmin[InstitutionCreate] = true
	sysAdmin[InstitutionRead] = true
	sysAdmin[InstitutionUpdate] = true
	sysAdmin[InstitutionDelete] = true
	sysAdmin[ObjectCreate] = true
	sysAdmin[ObjectRead] = true
	sysAdmin[ObjectUpdate] = true
	sysAdmin[ObjectDelete] = true
	sysAdmin[ObjectRequestDelete] = true
	sysAdmin[ObjectApproveDelete] = true
	sysAdmin[ObjectFinishBulkDelete] = true
	sysAdmin[ObjectRestore] = true
	sysAdmin[StorageRecordCreate] = true
	sysAdmin[StorageRecordRead] = true
	sysAdmin[StorageRecordUpdate] = true
	sysAdmin[StorageRecordDelete] = true
	sysAdmin[UserCreate] = true
	sysAdmin[UserRead] = true
	sysAdmin[UserSignIn] = true
	sysAdmin[UserSignOut] = true
	sysAdmin[UserUpdate] = true
	sysAdmin[UserDelete] = true
	sysAdmin[UserDeleteSelf] = true
	sysAdmin[UserReadSelf] = true
	sysAdmin[UserUpdateSelf] = true
	sysAdmin[WorkItemCreate] = true
	sysAdmin[WorkItemRead] = true
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

// PermissionForHandler maps HTTP handler names to the permissions
// required to access that handler.
var PermissionForHandler = map[string]Permission{
	"AlertCreate":            AlertCreate,
	"AlertNew":               AlertCreate,
	"AlertIndex":             AlertRead,
	"AlertShow":              AlertRead,
	"AlertUpdate":            AlertUpdate,
	"AlertDelete":            AlertDelete,
	"ChecksumNew":            ChecksumCreate,
	"ChecksumCreate":         ChecksumCreate,
	"ChecksumShow":           ChecksumRead,
	"ChecksumIndex":          ChecksumRead,
	"ChecksumUpdate":         ChecksumUpdate,
	"ChecksumDelete":         ChecksumDelete,
	"DashboardShow":          DashboardShow,
	"EventCreate":            EventCreate,
	"EventNew":               EventCreate,
	"EventShow":              EventRead,
	"EventIndex":             EventRead,
	"EventUpdate":            EventUpdate,
	"EventDelete":            EventDelete,
	"FileNew":                FileCreate,
	"FileCreate":             FileCreate,
	"FileShow":               FileRead,
	"FileIndex":              FileRead,
	"FileUpdate":             FileUpdate,
	"FileDelete":             FileDelete,
	"FileRequestDelete":      FileRequestDelete,
	"FileApproveDelete":      FileApproveDelete,
	"FileFinishBulkDelete":   FileFinishBulkDelete,
	"FileRestore":            FileRestore,
	"InstitutionNew":         InstitutionCreate,
	"InstitutionCreate":      InstitutionCreate,
	"InstitutionIndex":       InstitutionRead,
	"InstitutionShow":        InstitutionRead,
	"InstitutionUpdate":      InstitutionUpdate,
	"InstitutionDelete":      InstitutionDelete,
	"ObjectNew":              ObjectCreate,
	"ObjectCreate":           ObjectCreate,
	"ObjectIndex":            ObjectRead,
	"ObjectShow":             ObjectRead,
	"ObjectUpdate":           ObjectUpdate,
	"ObjectDelete":           ObjectDelete,
	"ObjectRequestDelete":    ObjectRequestDelete,
	"ObjectApproveDelete":    ObjectApproveDelete,
	"ObjectFinishBulkDelete": ObjectFinishBulkDelete,
	"ObjectRestore":          ObjectRestore,
	"StorageRecordNew":       StorageRecordCreate,
	"StorageRecordCreate":    StorageRecordCreate,
	"StorageRecordIndex":     StorageRecordRead,
	"StorageRecordShow":      StorageRecordRead,
	"StorageRecordUpdate":    StorageRecordUpdate,
	"StorageRecordDelete":    StorageRecordDelete,
	"UserNew":                UserCreate,
	"UserCreate":             UserCreate,
	"UserIndex":              UserRead,
	"UserShow":               UserRead,
	"UserUpdate":             UserUpdate,
	"UserDelete":             UserDelete,
	"UserReadSelf":           UserReadSelf,
	"UserUpdateSelf":         UserUpdateSelf,
	"UserDeleteSelf":         UserDeleteSelf,
	"WorkItemNew":            WorkItemCreate,
	"WorkItemCreate":         WorkItemCreate,
	"WorkItemIndex":          WorkItemRead,
	"WorkItemShow":           WorkItemRead,
	"WorkItemUpdate":         WorkItemUpdate,
	"WorkItemDelete":         WorkItemDelete,
}
