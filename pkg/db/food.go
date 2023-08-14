package db

import (
	"context"
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yawnak/foodadvisor/internal/domain"
)

var (
	foodStruct = sqlbuilder.NewStruct(new(food))
	foodTable  = "meals"
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
		return nil, fmt.Errorf("error scanning food: %w", err)
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

func (db *FoodDB) GetFoodWithoutLastEaten(ctx context.Context, userid int32, limit uint, offset uint) ([]domain.Food, error) {
	uT := goqu.T(usersTable)
	uTId := uT.Col("id")
	uTExpiration := uT.Col("expiration")
	mT := goqu.T(foodTable)
	mTId := mT.Col("id")
	mTName := mT.Col("name")
	mTCookTime := mT.Col("cooktime")
	mTouT := goqu.T("meals_to_users")
	mTouTMealId := mTouT.Col("mealid")
	mTouTUserId := mTouT.Col("userid")
	mTouTLastEaten := mTouT.Col("lasteaten")

	//example query for userid 33, limit 10, offset 10
	// 	SELECT m.id, m.name, m.cooktime
	// FROM meals m
	// LEFT JOIN meals_to_users mu ON m.id = mu.mealid AND mu.userid = 33
	// JOIN users u ON u.id = 33
	// WHERE mu.lasteaten IS NULL OR mu.lasteaten < (CURRENT_DATE - u.expiration);
	// LIMIT 10
	// OFFSET 10
	sql, arg, err := pggoqu.Select(mTId, mTName, mTCookTime).From(mT).
		LeftJoin(mTouT, goqu.On(goqu.And(
			mTId.Eq(mTouTMealId), mTouTUserId.Eq(userid),
		))).Join(uT, goqu.On(uTId.Eq(userid))).
		Where(goqu.Or(
			mTouTLastEaten.IsNull(),
			mTouTLastEaten.Lt(goqu.L("CURRENT_DATE - ?", uTExpiration)),
		)).Limit(limit).Offset(offset).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error constructing sql: %w", err)
	}
	rows, err := db.pool.Query(ctx, sql, arg...)
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", err)
	}
	meals := make([]domain.Food, 0, limit)
	for rows.Next() {
		temp := new(food)
		err = rows.Scan(&temp.Id, &temp.Name, &temp.CookTime)
		if err != nil {
			return nil, fmt.Errorf("error scanning: %w", err)
		}
		meals = append(meals, *foodRepoToFood(temp))
	}
	return meals, nil
}

func (db *FoodDB) GetMeals(ctx context.Context, offset uint, limit uint) ([]domain.Food, error) {
	sb := foodStruct.SelectFrom(foodTable)
	sql, args := sb.Limit(int(limit)).Offset(int(offset)).BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", err)
	}
	meals := []domain.Food{}
	for rows.Next() {
		var temp food
		err = rows.Scan(foodStruct.Addr(&temp)...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		meals = append(meals, *foodRepoToFood(&temp))
	}
	return meals, nil
}
