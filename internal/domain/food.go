package domain

import "context"

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

type FoodUsecase interface {
	CreateFood(ctx context.Context, food *Food) (int32, error)
}
