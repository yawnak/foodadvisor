package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yawnak/foodadvisor/pkg/domain"
	"golang.org/x/crypto/bcrypt"
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
	//retrieving user from db
	user, err := c.db.GetUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoUsername):
			return "", domain.ErrWrongCredentials
		default:
			return "", fmt.Errorf("error generating token: %w", err)
		}
	}
	//comparing using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return "", domain.ErrWrongCredentials
		default:
			return "", fmt.Errorf("error generating token: %w", err)
		}
	}
	//creating jwt token
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
	//hashing password using bcrypt
	bpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	switch {
	case errors.Is(err, bcrypt.ErrPasswordTooLong):
		return -1, fmt.Errorf("error creating user: %w", domain.ErrPasswordTooLong)
	case err != nil:
		return -1, fmt.Errorf("unknown error creating user: %w", err)
	}
	temp := *user
	temp.Password = string(bpass)

	//saving to db
	id, err := adv.db.CreateUser(ctx, &temp)
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
