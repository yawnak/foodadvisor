package db

import (
	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/pgtype"
)

var (
	foodStruct = sqlbuilder.NewStruct(new(food))
	foodTable  = "food"
)

type food struct {
	Id          pgtype.Int4
	Name        pgtype.Varchar
	CookTime    pgtype.Interval
	Price       pgtype.Int4
	IsBreakfast pgtype.Bool
	IsDinner    pgtype.Bool
	IsSupper    pgtype.Bool
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

func foodRepoToFood(f food) *domain.Food {
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
