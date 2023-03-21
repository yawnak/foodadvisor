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
