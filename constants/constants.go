package constants

const (
	RoleSysAdmin  = "admin"
	RoleInstAdmin = "institutional_admin"
	RoleInstUser  = "institutional_user"
	RoleNone      = "none"
	ActionCreate  = "Create"
	ActionRead    = "Read"
	ActionUpdate  = "Update"
	ActionDelete  = "Delete"
)

var Actions = []string{
	ActionCreate,
	ActionRead,
	ActionUpdate,
	ActionDelete,
}
