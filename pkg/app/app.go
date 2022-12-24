package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
)

var (
	ErrWrongCredentials = errors.New("wrong credentials")
)

type FoodAdvisor struct {
	db domain.AdvisorRepo
}

func NewFoodAdvisor(repo domain.AdvisorRepo) (*FoodAdvisor, error) {
	var adv FoodAdvisor
	adv.db = repo
	return &adv, nil
}

func (adv *FoodAdvisor) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	id, err := adv.db.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return id, nil
}

func (adv *FoodAdvisor) GetUserByCredentials(ctx context.Context, username string, password string) (*domain.User, error) {
	user, err := adv.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}
	if user.Username != username || user.Password != password {
		return nil, ErrWrongCredentials
	}
	return user, nil
}
