package domain

import "context"

type Permission string

const (
	PermEditRoles    Permission = "editRoles"
	PermEditUserRole Permission = "editUserRole"
)

type roleCtxKey struct{}

type Role struct {
	Name        string
	Permissions map[Permission]struct{}
}

func (u *Role) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

func RoleFromContext(ctx context.Context) (*Role, bool) {
	role, ok := ctx.Value(roleCtxKey{}).(*Role)
	if !ok {
		return nil, false
	}
	return role, true
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
