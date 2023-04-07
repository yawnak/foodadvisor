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

type AdvisorRepo interface {
	UserRepo
	FoodRepo
	RoleRepo
}

type Advisor interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	SetUserRole(ctx context.Context, id int32, role string) error
	GetUserByCredentials(ctx context.Context, username string, password string) (*User, error)
	CreateFood(ctx context.Context, food *Food) (int32, error)
	GetFoodByQuestionary(ctx context.Context, questionary *Questionary) ([]Food, error)
	GenerateToken(ctx context.Context, username string, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int32, error)
	ParseTokenWithRole(ctx context.Context, token string) (int32, *Role, error)
}

type RoleRepo interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, name string) error
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
