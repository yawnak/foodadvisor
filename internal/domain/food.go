package domain

import (
	"context"
)

type Food struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	//in minutes
	CookTime int32 `json:"cooktime"`
}

type FoodRepo interface {
	GetFoodById(ctx context.Context, id int32) (*Food, error)
	CreateFood(ctx context.Context, food *Food) (int32, error)
	DeleteFood(ctx context.Context, id int32) error
	UpdateFood(ctx context.Context, food *Food) error
	GetFoodByQuestionary(ctx context.Context, questionary *Questionary) ([]Food, error)
	GetFoodWithoutLastEaten(ctx context.Context, userid int32, limit uint, offset uint) ([]Food, error)
}

type FoodUsecase interface {
	CreateFood(ctx context.Context, food *Food) (int32, error)
	BasicAdvise(ctx context.Context, userid int32, limit uint, offset uint) ([]Food, error)
}
