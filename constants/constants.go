package constants

const (
	RoleSysAdmin  = "admin"
	RoleInstAdmin = "institutional_admin"
	RoleInstUser  = "institutional_user"
	RoleNone      = "none"
	ActionList    = "list"
	ActionCreate  = "create"
	ActionView    = "view"
	ActionEdit    = "edit"
	ActionDelete  = "delete"
)

var Actions = []string{
	ActionList,
	ActionCreate,
	ActionView,
	ActionEdit,
	ActionDelete,
}
