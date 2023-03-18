package app

import (
	"context"
	"fmt"
	"time"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/golang-jwt/jwt"
)

const (
	signingKey = "pudgebooster"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int32 `json:"user_id"`
}

type FoodAdvisor struct {
	db domain.AdvisorRepo
}

func NewFoodAdvisor(repo domain.AdvisorRepo) (*FoodAdvisor, error) {
	var adv FoodAdvisor
	adv.db = repo
	return &adv, nil
}

func (c *FoodAdvisor) GenerateToken(ctx context.Context, username string, password string) (string, error) {
	user, err := c.db.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}
	if user.Username != username || user.Password != password {
		return "", domain.ErrWrongCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
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
		return nil, domain.ErrWrongCredentials
	}
	return user, nil
}

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
