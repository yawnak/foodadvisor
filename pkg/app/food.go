package app

import (
	"context"
	"fmt"

	"github.com/yawnak/foodadvisor/internal/domain"
)

func (adv *FoodAdvisor) CreateFood(ctx context.Context, food *domain.Food) (int32, error) {
	id, err := adv.db.CreateFood(ctx, food)
	if err != nil {
		return 0, fmt.Errorf("error creating food: %w", err)
	}
	return id, err
}

func (adv *FoodAdvisor) GetFoodByQuestionary(ctx context.Context, questionary *domain.Questionary) ([]domain.Food, error) {
	return adv.db.GetFoodByQuestionary(ctx, questionary)
}

func (adv *FoodAdvisor) GetMeals(ctx context.Context, offset uint, limit uint) ([]domain.Food, error) {
	meals, err := adv.db.GetMeals(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return meals, nil
}

func (adv *FoodAdvisor) BasicAdvise(ctx context.Context, userid int32, limit uint, offset uint) ([]domain.Food, error) {
	meals, err := adv.db.GetFoodWithoutLastEaten(ctx, userid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting food without last eaten: %w", err)
	}
	return meals, nil
}

func (adv *FoodAdvisor) GetMealById(ctx context.Context, mealid int32) (*domain.Food, error) {
	meal, err := adv.db.GetFoodById(ctx, mealid)
	if err != nil {
		return nil, fmt.Errorf("error querying db: %w", err)
	}
	return meal, nil
}
