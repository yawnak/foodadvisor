package db

import (
	"context"
	"fmt"
	"log"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yawnak/foodadvisor/internal/domain"
)

var (
	foodStruct = sqlbuilder.NewStruct(new(food))
	foodTable  = "food"
)

type food struct {
	Id       pgtype.Int4     `db:"id"`
	Name     pgtype.Text     `db:"name" fieldtag:"details"`
	CookTime pgtype.Interval `db:"cooktime" fieldtag:"details"`
}

func foodToFoodRepo(f *domain.Food) *food {
	res := food{}
	res.Id.Scan(int64(f.Id))
	res.Name.Scan(f.Name)
	res.CookTime.Valid = true
	res.CookTime.Microseconds = int64(f.CookTime) * 60 * 1_000_000
	return &res
}

func foodRepoToFood(f *food) *domain.Food {
	return &domain.Food{
		Id:       f.Id.Int32,
		Name:     f.Name.String,
		CookTime: int32(f.CookTime.Microseconds / 1000000 / 60),
	}
}

func (db *FoodDB) GetFoodById(ctx context.Context, id int32) (*domain.Food, error) {
	sb := foodStruct.SelectFrom(foodTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	var food food
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(foodStruct.Addr(&food)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return foodRepoToFood(&food), nil
}

func (db *FoodDB) CreateFood(ctx context.Context, food *domain.Food) (int32, error) {
	f := foodToFoodRepo(food)
	// sb := foodStruct.InsertIntoForTag(foodTable, "details", f)
	sb := foodStruct.WithTag("details").InsertInto(foodTable, f)
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	sql += " RETURNING id"
	log.Println(sql)
	log.Println(args)
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
	sb := foodStruct.WithTag("details").Update(foodTable, foodRepo)
	sb.Where(sb.Equal("id", foodRepo.Id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executring query: %w", err)
	}
	return err
}

func (db *FoodDB) GetFoodByQuestionary(ctx context.Context, questionary *domain.Questionary) ([]domain.Food, error) {
	sb := foodStruct.SelectFrom(foodTable)
	if questionary.MaxCookTime != nil {
		var ct pgtype.Interval
		ct.Microseconds = int64(int64(*questionary.MaxCookTime) * 1000000 * 60)
		ct.Valid = true
		sb.Where(sb.LessEqualThan("cooktime", ct))
	}
	if questionary.MaxPrice != nil {
		sb.Where(sb.LessEqualThan("price", *questionary.MaxPrice))
	}
	if questionary.MealType != nil {
		sb.Where(sb.Equal("mealtype", *questionary.MealType))
	}
	if questionary.DishType != nil {
		sb.Where(sb.Equal("dishtype", *questionary.DishType))
	}

	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	log.Println(sql)
	log.Println(args)

	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", err)
	}
	defer rows.Close()

	foods := make([]domain.Food, 0, 5)

	for rows.Next() {
		cur := food{}
		err = rows.Scan(foodStruct.Addr(&cur)...)
		if err != nil {
			return nil, fmt.Errorf("error scanning: %w", err)
		}
		foods = append(foods, *foodRepoToFood(&cur))
	}
	return foods, nil
}
