package constants

const (
	ActionApproveDelete    = "ApproveDelete"
	ActionCreate           = "Create"
	ActionDelete           = "Delete"
	ActionFinishBulkDelete = "FinishBulkDelete"
	ActionRead             = "Read"
	ActionRequestDelete    = "RequestDelete"
	ActionRestore          = "Restore"
	ActionUpdate           = "Update"
	AlgMd5                 = "md5"
	AlgSha1                = "sha1"
	AlgSha256              = "sha256"
	AlgSha512              = "sha512"
	RoleInstAdmin          = "institutional_admin"
	RoleInstUser           = "institutional_user"
	RoleNone               = "none"
	RoleSysAdmin           = "admin"
)

var Actions = []string{
	ActionCreate,
	ActionRead,
	ActionUpdate,
	ActionDelete,
	ActionRequestDelete,
	ActionApproveDelete,
	ActionFinishBulkDelete,
	ActionRestore,
}

var DigestAlgs = []string{
	AlgMd5,
	AlgSha1,
	AlgSha256,
	AlgSha512,
}
