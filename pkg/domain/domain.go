package domain

import (
	"context"
)

type User struct {
	Id             int32
	Username       string
	Password       string
	ExpirationDays int32 //in days
}

type AdvisorRepo interface {
	UserRepo
	FoodRepo
}

type Advisor interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUserByCredentials(ctx context.Context, username string, password string) (*User, error)
	CreateFood(ctx context.Context, food *Food) (int32, error)
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUserById(ctx context.Context, id int32) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int32) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type Ingridient struct {
	Id   int32
	Name string
}

type Question struct {
	Id   int32
	Text string
}

var MealTypes = []string{
	"breakfast", "dinner", "supper",
}

var DishTypes = []string{
	"soup", "porridge", "puree", "desert", "cake", "cutlet", "dumpling",
}

type Food struct {
	Id       int32
	Name     string
	CookTime int32
	Price    int32
	MealType string
	DishType string
}

type FoodRepo interface {
	GetFoodById(ctx context.Context, id int32) (*Food, error)
	CreateFood(ctx context.Context, food *Food) (int32, error)
	DeleteFood(ctx context.Context, id int32) error
	UpdateFood(ctx context.Context, food *Food) error
}
