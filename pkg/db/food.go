package db

import (
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
