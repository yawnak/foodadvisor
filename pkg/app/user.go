package app

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// if date is nil time.Now will be used
func (adv *FoodAdvisor) UpdateUserEaten(ctx context.Context, userid int32, foodid int32, date *time.Time) error {
	var err error
	if date == nil {
		err = adv.db.UpdateUserEatenFood(ctx, userid, foodid, time.Now())
	} else {
		err = adv.db.UpdateUserEatenFood(ctx, userid, foodid, *date)
	}
	if err != nil {
		return fmt.Errorf("error requesting db: %w", err)
	}
	return nil
}

func (adv *FoodAdvisor) GetUserById(ctx context.Context, userid int32) (*domain.User, error) {
	u, err := adv.db.GetUserById(ctx, userid)
	if err != nil {
		return nil, fmt.Errorf("error getting user by id from db: %w", err)
	}
	u.Password = ""
	return u, err
}
