package domain

import (
	"context"
	"time"
)

const (
	TokenTTL = time.Hour * 12
)

type Questionary struct {
	MaxCookTime *int32
	MaxPrice    *int32
	MealType    *string
	DishType    *string
}

type User struct {
	Id             int32  `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ExpirationDays int32  `json:"expiration"` //in days
}

type AdvisorRepo interface {
	UserRepo
	FoodRepo
	RoleRepo
}

type Advisor interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUserByCredentials(ctx context.Context, username string, password string) (*User, error)
	CreateFood(ctx context.Context, food *Food) (int32, error)
	GetFoodByQuestionary(ctx context.Context, questionary *Questionary) ([]Food, error)
	GenerateToken(ctx context.Context, username string, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int32, error)
}

type RoleRepo interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, name string) error
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
	GetFoodByQuestionary(ctx context.Context, questionary *Questionary) ([]Food, error)
}
