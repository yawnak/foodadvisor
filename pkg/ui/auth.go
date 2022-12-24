package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func (cli *UICli) AuthenticationMenu() error {
	prompt := promptui.Select{
		Label: "Choose:",
		Items: []string{
			"Login",
			"Register"},
	}
	_, res, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error selecting: %w", err)
	}

	switch res {
	case "Login":
		err = cli.LoginPrompt()
		if err != nil {
			return fmt.Errorf("error during logging in: %w", err)
		}
	case "Register":
		err = cli.RegisterPrompt()
		if err != nil {
			return fmt.Errorf("error during logging in: %w", err)
		}
	default:
		return fmt.Errorf("unknown result from select prompt: %s", res)
	}
	return nil
}

func (cli *UICli) LoginPrompt() error {
	prompt := promptui.Prompt{
		Label: "Enter username:",
	}
	username, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error getting username: %w", err)
	}
	prompt = promptui.Prompt{
		Label: "Enter password",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error getting password: %w", err)
	}
	fmt.Printf("username: %s, password: %s\n", username, password)
	return nil
}

func (cli *UICli) RegisterPrompt() error {
	prompt := promptui.Prompt{
		Label: "Enter username:",
	}
	username, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error getting username: %w", err)
	}
	prompt = promptui.Prompt{
		Label: "Enter password",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error getting password: %w", err)
	}

	prompt = promptui.Prompt{
		Label: "Confirm password",
		Mask:  '*',
	}
	confirmPassword, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error getting password: %w", err)
	}

	if confirmPassword != password {
		fmt.Println("passwords don't match, try again")
		err = cli.RegisterPrompt()
		if err != nil {
			return fmt.Errorf("error retrying: %w", err)
		}
	} else {
		fmt.Printf("username: %s, password: %s\n", username, password)
	}
	return nil
}
