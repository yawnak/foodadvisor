package ui

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/manifoldco/promptui"
)

func (cli *UICli) FoodCreationPrompt() error {
	var food domain.Food
	var err error
	prompt := promptui.Prompt{
		Label: "Enter food name",
	}
	food.Name, err = prompt.Run()
	if err != nil {
		return fmt.Errorf("error running food name prompt: %w", err)
	}

	prompt = promptui.Prompt{
		Label:    "Enter cooktime",
		Validate: validateInt,
	}
	s, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error running food cooktime prompt: %w", err)
	}
	n, _ := strconv.Atoi(s)
	food.CookTime = int32(n)

	prompt = promptui.Prompt{
		Label:    "Enter price",
		Validate: validateInt,
	}
	s, err = prompt.Run()
	if err != nil {
		return fmt.Errorf("error running food cooktime prompt: %w", err)
	}
	n, _ = strconv.Atoi(s)
	food.Price = int32(n)

	prompts := promptui.Select{
		Label: "Choose meal type",
		Items: domain.MealTypes,
	}
	_, food.MealType, err = prompts.Run()
	if err != nil {
		return fmt.Errorf("error running select for mealtype: %w", err)
	}

	prompts = promptui.Select{
		Label: "Choose dish type",
		Items: domain.DishTypes,
	}
	_, food.DishType, err = prompts.Run()
	if err != nil {
		return fmt.Errorf("error running select for dishtype: %w", err)
	}

	food.Id, err = cli.adv.CreateFood(context.Background(), &food)
	if err != nil {
		return fmt.Errorf("error creating food: %w", err)
	}

	log.Printf("food: %v", food)
	return nil
}
