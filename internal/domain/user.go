package domain

import "context"

type userCtxKey struct{}

type User struct {
	Id             int32   `json:"id"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	ExpirationDays int32   `json:"expiration"` //in days
	Role           *string `json:"role"`
}

func (u *User) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userCtxKey{}).(*User)
	if !ok {
		return nil, false
	}
	return user, true
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUserById(ctx context.Context, id int32) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int32) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserRole(ctx context.Context, id int32) (*Role, error)
	UpdateUserRole(ctx context.Context, id int32, role string) error
}