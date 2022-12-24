package app

import (
	"context"

	"github.com/asstronom/foodadvisor/pkg/domain"
)

type FoodAdvisor struct {
	db domain.AdvisorRepo
}

func NewFoodAdvisor(repo domain.AdvisorRepo) (*FoodAdvisor, error) {
	var adv *FoodAdvisor
	adv.db = repo
	return adv, nil
}

func (adv *FoodAdvisor) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	id, err := adv.db.CreateUser(ctx, user)
	return id, err
}
