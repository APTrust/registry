package constants

// These constants represent permissions. They are indexes into the
// boolean arrays defined below.
const (
	AlertCreate = iota
	AlertRead
	AlertUpdate
	AlertDelete
	ChecksumCreate
	ChecksumRead
	ChecksumUpdate
	ChecksumDelete
	EventCreate
	EventRead
	EventUpdate
	EventDelete
	FileCreate
	FileRead
	FileUpdate
	FileDelete
	FileRequestDelete
	FileApproveDelete    // For institutional admin approval
	FileFinishBulkDelete // Sys admin approval of bulk delete
	FileRestore
	InstitutionCreate
	InstitutionRead
	InstitutionUpdate
	InstitutionDelete
	ObjectCreate
	ObjectRead
	ObjectUpdate
	ObjectDelete
	ObjectRequestDelete
	ObjectApproveDelete    // For institutional admin approval
	ObjectFinishBulkDelete // Sys admin approval of bulk delete
	ObjectRestore
	StorageRecordCreate
	StorageRecordRead
	StorageRecordUpdate
	StorageRecordDelete
	UserCreate
	UserRead
	UserUpdate
	UserDelete
	UserReadSelf
	UserUpdateSelf
	WorkItemCreate
	WorkItemRead
	WorkItemUpdate
	WorkItemDelete

	// Future permissions for reports will go here.
)

var permissionsInitialized = false
var permissionCount = 46

// Permission lists for different roles. Bools default to false in Go,
// so roles will have only the permissions we explicitly grant below.
var instUser = make([]bool, permissionCount)
var instAdmin = make([]bool, permissionCount)
var sysAdmin = make([]bool, permissionCount)
var emptyList = make([]bool, permissionCount)

func initPermissions() {
	instUser[AlertRead] = true
	instUser[AlertUpdate] = true
	instUser[ChecksumRead] = true
	instUser[EventRead] = true
	instUser[FileRead] = true
	instUser[ObjectRead] = true
	instUser[StorageRecordRead] = true
	instUser[UserReadSelf] = true
	instUser[UserUpdateSelf] = true
	instUser[WorkItemRead] = true

	instAdmin[AlertRead] = true
	instAdmin[AlertUpdate] = true
	instAdmin[ChecksumRead] = true
	instAdmin[EventRead] = true
	instAdmin[FileRead] = true
	instAdmin[FileDelete] = true
	instAdmin[FileRequestDelete] = true
	instAdmin[FileApproveDelete] = true
	instAdmin[FileRestore] = true
	instAdmin[ObjectRead] = true
	instAdmin[ObjectDelete] = true
	instAdmin[ObjectRequestDelete] = true
	instAdmin[ObjectApproveDelete] = true
	instAdmin[ObjectRestore] = true
	instAdmin[StorageRecordRead] = true
	instAdmin[UserCreate] = true
	instAdmin[UserRead] = true
	instAdmin[UserUpdate] = true
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
	sysAdmin[UserUpdate] = true
	sysAdmin[UserDelete] = true
	sysAdmin[UserReadSelf] = true
	sysAdmin[UserUpdateSelf] = true
	sysAdmin[WorkItemCreate] = true
	sysAdmin[WorkItemRead] = true
	sysAdmin[WorkItemUpdate] = true
	sysAdmin[WorkItemDelete] = true

	permissionsInitialized = true
}

func CheckPermission(role string, permission int) bool {
	if !permissionsInitialized {
		initPermissions()
	}
	var permissionList []bool
	switch role {
	case RoleSysAdmin:
		permissionList = sysAdmin
	case RoleInstAdmin:
		permissionList = instAdmin
	case RoleInstUser:
		permissionList = instUser
	default:
		permissionList = emptyList
	}
	return permissionList[permission]
}
