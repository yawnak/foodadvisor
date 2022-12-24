package ui

import (
	"context"
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
	_, mealtype, err := prompts.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for meal type")
	}
	if mealtype != "indifferent" {
		questionary.MealType = &mealtype
	}

	prompts = promptui.Select{
		Label: "Select dish type",
		Items: append(domain.DishTypes, "indifferent"),
	}
	_, dishtype, err := prompts.Run()
	if err != nil {
		return nil, fmt.Errorf("error prompting for meal type")
	}
	if dishtype != "indifferent" {
		questionary.DishType = &dishtype
	}

	return &questionary, nil
}

func (cli *UICli) AdvisePrompt() error {
	questionary, err := cli.QuestionaryPrompt()
	if err != nil {
		return fmt.Errorf("error prompting questionary: %w", err)
	}
	log.Println(questionary)

	food, err := cli.adv.GetFoodByQuestionary(context.Background(), questionary)
	if err != nil {
		return fmt.Errorf("error getting food: %w", err)
	}
	log.Println(food)

	template := &promptui.SelectTemplates{
		Active:   "{{ .Name | cyan }}",
		Inactive: "{{ .Name | white }}",
	}

	prompts := promptui.Select{
		Label:     "Select food to show info",
		Templates: template,
		Items:     food,
	}

	_, res, err := prompts.Run()
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}
