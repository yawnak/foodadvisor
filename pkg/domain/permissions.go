package domain

type Permission string

const (
	PermEditRoles    Permission = "editRoles"
	PermEditUserRole Permission = "editUserRole"
)

type Role struct {
	Name        string
	Permissions map[Permission]struct{}
}

func NewRole(name string, permissions ...Permission) *Role {
	role := Role{
		Name:        name,
		Permissions: make(map[Permission]struct{}),
	}
	for _, p := range permissions {
		role.Permissions[p] = struct{}{}
	}
	return &role
}
