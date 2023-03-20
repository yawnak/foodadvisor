package db

import (
	"context"
	"fmt"

	"github.com/yawnak/foodadvisor/pkg/domain"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	userStruct = sqlbuilder.NewStruct(new(userRepo))
	usersTable = "users"
)

type userRepo struct {
	Id             pgtype.Int4     `db:"id"`
	Username       pgtype.Text     `db:"username" fieldtag:"details"`
	Password       pgtype.Text     `db:"passhash" fieldtag:"details"`
	ExpirationDays pgtype.Interval `db:"expiration" fieldtag:"details"` //in days
}

func userToUserRepo(user *domain.User) *userRepo {
	var res userRepo
	res.Id.Int32 = user.Id
	res.Id.Valid = true
	res.Username.Scan(user.Username)
	res.Password.Scan(user.Password)
	res.ExpirationDays.Days = user.ExpirationDays
	res.ExpirationDays.Valid = true
	return &res
}

func userRepotouser(u *userRepo) *domain.User {
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
	var user userRepo
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
	var user userRepo
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return userRepotouser(&user), nil
}

func (db *FoodDB) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	userRepo := userToUserRepo(user)
	fmt.Println(userRepo)
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
