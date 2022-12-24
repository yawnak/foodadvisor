package ui

import (
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/manifoldco/promptui"
)

type UICli struct {
	user domain.User
	adv  domain.Advisor
}

func NewUICli(adv domain.Advisor) (*UICli, error) {
	var cli UICli
	cli.adv = adv
	return &cli, nil
}

func (cli *UICli) MainMenu() error {
	prompt := promptui.Select{
		Label: "Choose",
		Items: []string{
			"Advise",
			"Create Food",
		},
	}
	_, res, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error running select: %w", err)
	}

	switch res {
	case "Advise":
		err = cli.AdvisePrompt()
		if err != nil {
			return fmt.Errorf("error advising: %w", err)
		}
	case "Create Food":
		err = cli.FoodCreationPrompt()
		if err != nil {
			return fmt.Errorf("error creating food: %w", err)
		}
	}
	return nil
}

func (cli *UICli) Run() error {
	err := cli.AuthenticationMenu()
	if err != nil {
		return fmt.Errorf("error authenticating: %w", err)
	}
	fmt.Println("AUTHORIZED")
	fmt.Println(cli.user)
	for {
		err = cli.MainMenu()
		if err != nil {
			break
		}
	}
	return nil
}
