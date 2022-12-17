package domain

import (
	"context"
	"time"
)

type User struct {
	Id             int32
	Username       string
	Password       string
	ExpirationDays int32 //in days
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (int64, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type Ingridient struct {
	Id   int32
	Name string
}

type Food struct {
	Id          int32
	Name        string
	CookTime    time.Time
	Price       int32
	IsBreakfast bool
	IsDinner    bool
	IsSupper    bool
}

type FoodRepo interface {
}
