package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yawnak/foodadvisor/pkg/domain"
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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrNoUsername
		default:
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
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
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr):
			if pgErr.Code == "23505" {
				return -1, domain.ErrDuplicateResourse
			}
		}
		return -1, fmt.Errorf("error scanning: %w", err)
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

func (db *FoodDB) UpdateUserRole(ctx context.Context, id int32, role string) error {
	usersT := goqu.T(usersTable)
	sql, args, err := goqu.Dialect("postgres").
		Update(usersT).
		Set(goqu.I("role").Set(role)).
		Where(usersT.Col("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("error building sql: %w", err)
	}
	_, err = db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error updating: %w", err)
	}
	return nil
}

func (db *FoodDB) GetUserRole(ctx context.Context, id int32) (*domain.Role, error) {
	withds := goqu.Dialect("postgres").
		Select("role").From(goqu.T(usersTable)).
		Where(goqu.T(usersTable).Col("id").Eq(id)).
		Prepared(true)
	permT := goqu.T(permissionsToRolesTable)
	sql, args, err := goqu.Dialect("postgres").
		Select(permT.Col("role"), permT.Col("permission")).
		With("t", withds).
		From(goqu.T(permissionsToRolesTable)).
		InnerJoin(goqu.T("t"),
			goqu.On(permT.Col("role").Eq(goqu.T("t").Col("role")))).Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error building sql: %w", err)
	}
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying rows: %w", err)
	}
	defer rows.Close()
	var roleName string
	permissions := make([]domain.Permission, 0, 1)
	for rows.Next() {
		var temp domain.Permission
		err := rows.Scan(&roleName, &temp)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, temp)
	}
	if len(permissions) == 0 {
		sql, args, err = goqu.Dialect("postgres").
			Select("role").From(usersTable).
			Where(
				goqu.T(usersTable).Col("id").Eq(id)).
			Prepared(true).ToSQL()
		if err != nil {
			return nil, fmt.Errorf("error building sql for select when role has 0 permissions: %w", err)
		}
		row := db.pool.QueryRow(ctx, sql, args...)
		err = row.Scan(&roleName)
		if err != nil {
			return nil, fmt.Errorf("error scanning single row: %w", err)
		}
	}
	return domain.NewRole(roleName, permissions...), nil
}
