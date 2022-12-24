package app

import (
	"context"

	"github.com/asstronom/foodadvisor/pkg/domain"
)

type FoodAdvisor struct {
	db domain.AdvisorRepo
}

func (adv *FoodAdvisor) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	id, err := adv.db.CreateUser(ctx, user)
	return id, err
}
