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
	UserUsecase
	FoodUsecase
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
