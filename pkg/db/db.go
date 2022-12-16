package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type FoodDB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, connStr string) (*FoodDB, error) {
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %s", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging db: %s", err)
	}
	return &FoodDB{pool: pool}, nil
}

func (db *FoodDB) Close() {
	db.pool.Close()
}
