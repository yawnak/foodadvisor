package db

import "github.com/huandu/go-sqlbuilder"

var (
	foodStruct = sqlbuilder.NewStruct(new(food))
	foodTable  = "food"
)

type food struct {
}
