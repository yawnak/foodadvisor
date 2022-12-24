package db

import (
	"context"
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	userStruct = sqlbuilder.NewStruct(new(user))
	usersTable = "users"
)

type user struct {
	Id             pgtype.Int4     `db:"id"`
	Username       pgtype.Text     `db:"username" fieldtag:"details"`
	Password       pgtype.Text     `db:"passwrd" fieldtag:"details"`
	ExpirationDays pgtype.Interval `db:"expiration" fieldtag:"details"` //in days
}

func userToUserRepo(u *domain.User) *user {
	var res user
	res.Id.Int32 = u.Id
	res.Username.String = u.Username
	res.Password.String = u.Password
	res.ExpirationDays.Days = u.ExpirationDays
	res.Id.Valid = true
	res.Username.Valid = true
	res.Password.Valid = true
	res.ExpirationDays.Valid = true
	return &res
}

func userRepotouser(u *user) *domain.User {
	var res domain.User
	res.Id = u.Id.Int32
	res.Username = u.Username.String
	res.Password = u.Password.String
	res.ExpirationDays = u.ExpirationDays.Days
	return &res
}

func (db *FoodDB) GetUserById(ctx context.Context, id int32) (*domain.User, error) {
	sb := userStruct.SelectFrom(usersTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	var user user
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return userRepotouser(&user), nil
}

func (db *FoodDB) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	sb := userStruct.SelectFrom(usersTable)
	sb.Where(sb.Equal("username", username))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	var user user
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return userRepotouser(&user), nil
}

func (db *FoodDB) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	userRepo := userToUserRepo(user)
	sb := userStruct.InsertIntoForTag(usersTable, "details", userRepo)
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	sql += " RETURNING id"
	row := db.pool.QueryRow(ctx, sql, args...)
	var id int32
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error scanning returning id: %w", err)
	}
	return id, nil
}

func (db *FoodDB) DeleteUser(ctx context.Context, id int32) error {
	sb := userStruct.DeleteFrom(usersTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	return nil
}

func (db *FoodDB) UpdateUser(ctx context.Context, user *domain.User) error {
	userRepo := userToUserRepo(user)
	sb := userStruct.UpdateForTag(usersTable, "details", userRepo)
	sb.Where(sb.Equal("id", userRepo.Id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executring query: %w", err)
	}
	return err
}
