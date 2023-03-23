package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yawnak/foodadvisor/pkg/domain"
)

const (
	roleTable               = "roles"
	permissionsToRolesTable = "permissions_to_roles"
)

func (db *FoodDB) CreateRole(ctx context.Context, role domain.Role) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting tx: %w", err)
	}
	defer tx.Rollback(ctx)

	sql := fmt.Sprintf("INSERT INTO %s (name) VALUES($1)", roleTable)
	_, err = tx.Exec(ctx, sql, role.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr):
			if pgErr.Code == "23505" {
				return domain.ErrDuplicateResourse
			}
		default:
			return fmt.Errorf("error inserting role: %w", err)
		}
	}
	rows := make([][]interface{}, 0, len(role.Permissions))
	for permission, _ := range role.Permissions {
		rows = append(rows, []interface{}{role.Name, permission})
	}

	n, err := tx.CopyFrom(ctx, pgx.Identifier{permissionsToRolesTable}, []string{"role", "permission"}, pgx.CopyFromRows(rows))

	if err != nil {
		fmt.Println("number of insertions", n)
		return fmt.Errorf("error inserting permissions: %w", err)
	}

	tx.Commit(ctx)
	return nil
}

func (db *FoodDB) GetRole(ctx context.Context, name string) (*domain.Role, error) {
	t1 := goqu.T(permissionsToRolesTable)
	t2 := goqu.T(roleTable)
	t1role := t1.Col("role")
	t1perm := t1.Col("permission")
	t2name := t2.Col("name")
	bldr := goqu.Dialect("postgres")
	ds := bldr.Select(t1role, t1perm).From(t1).Prepared(true).
		InnerJoin(t2, goqu.On(t1role.Eq(t2name))).Where(t1role.Eq(name))
	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error building sql: %w", err)
	}
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting roles and permissions: %w", err)
	}
	defer rows.Close()
	var role domain.Role
	role.Permissions = make(map[domain.Permission]struct{})

	for rows.Next() {
		var temp domain.Permission
		err := rows.Scan(&role.Name, &temp)
		if err != nil {
			return nil, fmt.Errorf("error scanning permission row: %w", err)
		}
		role.Permissions[temp] = struct{}{}
	}

	if len(role.Permissions) == 0 {
		return nil, domain.ErrResourseNotFound
	}

	return &role, nil
}
