package server

import "github.com/yawnak/foodadvisor/internal/domain"

type responseSuccess struct {
	SuccessMessage string `json:"success"`
	HTTPStatusCode int    `json:"-"`
}

func (s responseSuccess) Success() string {
	return s.SuccessMessage
}

func (s responseSuccess) Status() int {
	return s.HTTPStatusCode
}

type responseSignup struct {
	responseSuccess
	UserId int32 `json:"id"`
}

type responseLogin struct {
	responseSuccess
	UserId int32 `json:"id"`
}

type responseCreateMeal struct {
	responseSuccess
	MealId int32 `json:"id"`
}

type responseBasicAdvice struct {
	responseSuccess
	Meals []domain.Food `json:"meals"`
}
