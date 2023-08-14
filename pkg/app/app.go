package app

import (
	"github.com/yawnak/foodadvisor/internal/domain"
)

type FoodAdvisor struct {
	db domain.AdvisorRepo
}

func NewFoodAdvisor(repo domain.AdvisorRepo) (*FoodAdvisor, error) {
	var adv FoodAdvisor
	adv.db = repo
	return &adv, nil
}
