package ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/manifoldco/promptui"
)

func (cli *UICli) QuestionaryPrompt() (*domain.Questionary, error) {
	var questionary domain.Questionary
	prompt := promptui.Prompt{
		Label:    "Enter max cooking time (-1 if doesn't matter)",
		Validate: validateInt,
	}
	s, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for max cook time: %w", err)
	}
	n, _ := strconv.Atoi(s)
	if n != -1 {
		number := int32(n)
		questionary.MaxCookTime = &number
	}

	prompt = promptui.Prompt{
		Label:    "Enter max price (-1 if doesn't matter)",
		Validate: validateInt,
	}
	s, err = prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for max price: %w", err)
	}
	n, _ = strconv.Atoi(s)
	if n != -1 {
		number := int32(n)
		questionary.MaxPrice = &number
	}

	prompts := promptui.Select{
		Label: "Select meal type",
		Items: append(domain.MealTypes, "indifferent"),
	}
	_, s, err = prompts.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for meal type")
	}
	if s != "indifferent" {
		questionary.MealType = &s
	}

	prompts = promptui.Select{
		Label: "Select dish type",
		Items: append(domain.DishTypes, "indifferent"),
	}
	_, s, err = prompts.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for meal type")
	}
	if s != "indifferent" {
		questionary.DishType = &s
	}

	log.Println(questionary)

	return &questionary, nil
}
