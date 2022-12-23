package db

import (
	"context"
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/pgtype"
)

var (
	foodStruct = sqlbuilder.NewStruct(new(food))
	foodTable  = "food"
)

type food struct {
	Id          pgtype.Int4     `db:"id"`
	Name        pgtype.Varchar  `db:"name" fieldtag:"details"`
	CookTime    pgtype.Interval `db:"cooktime" fieldtag:"details"`
	Price       pgtype.Int4     `db:"price" fieldtag:"details"`
	IsBreakfast pgtype.Bool     `db:"isbreakfast" fieldtag:"details"`
	IsDinner    pgtype.Bool     `db:"isdinner" fieldtag:"details"`
	IsSupper    pgtype.Bool     `db:"issupper" fieldtag:"details"`
}

func foodToFoodRepo(f *domain.Food) *food {
	res := food{}
	res.Id.Set(f.Id)
	res.Name.Set(f.Name)
	res.CookTime.Set(f.CookTime)
	res.Price.Set(f.Price)
	res.IsBreakfast.Set(f.IsBreakfast)
	res.IsDinner.Set(f.IsDinner)
	res.IsSupper.Set(f.IsSupper)
	return &res
}

func foodRepoToFood(f *food) *domain.Food {
	return &domain.Food{
		Id:          f.Id.Int,
		Name:        f.Name.String,
		CookTime:    int32(f.CookTime.Microseconds / 1000),
		Price:       f.Price.Int,
		IsBreakfast: f.IsBreakfast.Bool,
		IsDinner:    f.IsDinner.Bool,
		IsSupper:    f.IsSupper.Bool,
	}
}

func (db *FoodDB) GetFoodById(ctx context.Context, id int32) (*domain.Food, error) {
	sb := foodStruct.SelectFrom(foodTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var food food
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&food)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return foodRepoToFood(&food), nil
}

func (db *FoodDB) CreateFood(ctx context.Context, food *domain.Food) (int32, error) {
	f := foodToFoodRepo(food)
	sb := foodStruct.InsertIntoForTag(foodTable, "details", userStruct.ValuesForTag("details", f)...)
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

func (db *FoodDB) DeleteFood(ctx context.Context, id int32) error {
	sb := foodStruct.DeleteFrom(foodTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	return nil
}

func (db *FoodDB) UpdateFood(ctx context.Context, food *domain.Food) error {
	foodRepo := foodToFoodRepo(food)
	sb := userStruct.UpdateForTag(usersTable, "details", userStruct.ValuesForTag("details", &foodRepo))
	sb.Where(sb.Equal("id", foodRepo.Id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executring query: %w", err)
	}
	return err
}
