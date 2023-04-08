package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/yawnak/foodadvisor/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (adv *FoodAdvisor) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	//hashing password using bcrypt
	bpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	switch {
	case errors.Is(err, bcrypt.ErrPasswordTooLong):
		return -1, domain.ErrPasswordTooLong
	case err != nil:
		return -1, fmt.Errorf("unknown error generating passhash: %w", err)
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

func (adv *FoodAdvisor) SetUserRole(ctx context.Context, id int32, role string) error {
	return adv.db.UpdateUserRole(ctx, id, role)
}
