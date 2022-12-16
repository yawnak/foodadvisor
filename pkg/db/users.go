package db

import (
	"context"
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/pgtype"
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
	res.Id.Set(u.Id)
	res.Username.Set(u.Username)
	res.Password.Set(u.Password)
	res.ExpirationDays.Days = u.ExpirationDays
	res.ExpirationDays.Status = pgtype.Present
	return &res
}

func userRepotouser(u *user) *domain.User {
	var res domain.User
	res.Id = u.Id.Int
	res.Username = u.Username.String
	res.Password = u.Username.String
	res.ExpirationDays = u.ExpirationDays.Days
	return &res
}

func (db *FoodDB) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
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

func (db *FoodDB) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	sb := userStruct.InsertIntoForTag(usersTable, "details", userStruct.ValuesForTag("details", user)...)
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	sql += " RETURNING id"
	row := db.pool.QueryRow(ctx, sql, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error scanning returning id: %w", err)
	}
	return id, nil
}

func (db *FoodDB) DeleteUser(ctx context.Context, id int64) error {
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
	sb := userStruct.UpdateForTag(usersTable, "details", userStruct.ValuesForTag("details", &userRepo))
	sb.Where(sb.Equal("id", userRepo.Id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executring query: %w", err)
	}
	return err
}
