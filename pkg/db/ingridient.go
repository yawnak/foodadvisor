package db

import sqlbuilder "github.com/huandu/go-sqlbuilder"

var (
	ingridientStruct = sqlbuilder.NewStruct(new(ingridient))
	ingridientsTable = "ingridients"
)

type ingridient struct {
}
