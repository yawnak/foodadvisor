package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yawnak/foodadvisor/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// TODO: make signingKey not a global var
const (
	signingKey = "pudgebooster"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int32        `json:"user_id"`
	Role   *domain.Role `json:"role"`
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
	//retrieving user role
	role, err := c.db.GetUserRole(ctx, user.Id)
	if err != nil {
		return "", fmt.Errorf("error getting user role: %w", err)
	}

	//creating jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.Id,
		Role:   role,
	})
	return token.SignedString([]byte(signingKey))
}

func (adv *FoodAdvisor) parseClaims(ctx context.Context, token string) (*tokenClaims, error) {
	accessToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidSigningMethod
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		var validationError *jwt.ValidationError
		switch {
		case errors.As(err, &validationError):
			return nil, domain.ErrBadToken
		default:
			return nil, err
		}
	}
	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return nil, domain.ErrBadToken
	}
	return claims, nil
}

func (adv *FoodAdvisor) ParseToken(ctx context.Context, token string) (int32, error) {
	claims, err := adv.parseClaims(ctx, token)
	if err != nil {
		return -1, err
	}
	return claims.UserID, nil
}

func (adv *FoodAdvisor) ParseTokenWithRole(ctx context.Context, token string) (int32, *domain.Role, error) {
	claims, err := adv.parseClaims(ctx, token)
	if err != nil {
		return -1, nil, err
	}
	return claims.UserID, claims.Role, nil
}
