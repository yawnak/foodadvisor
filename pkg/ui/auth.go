package ui

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/manifoldco/promptui"
)

func (cli *UICli) CredentialsPrompt() (string, string, error) {
	prompt := promptui.Prompt{
		Label: "Enter username:",
	}
	username, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("error getting username: %w", err)
	}
	prompt = promptui.Prompt{
		Label: "Enter password",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("error getting password: %w", err)
	}
	return username, password, nil
}

func (cli *UICli) ExpirationPrompt() (int32, error) {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("not a number")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter expiration days",
		Validate: validate,
	}
	s, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	n, _ := strconv.Atoi(s)
	return int32(n), nil
}

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
	username, password, err := cli.CredentialsPrompt()
	if err != nil {
		return fmt.Errorf("error getting credentials: %w", err)
	}
	fmt.Printf("username: %s, password: %s\n", username, password)
	return nil
}

func (cli *UICli) RegisterPrompt() error {
	var user domain.User
	var err error

	user.Username, user.Password, err = cli.CredentialsPrompt()

	if err != nil {
		return fmt.Errorf("error getting credentials: %w", err)
	}

	user.ExpirationDays, err = cli.ExpirationPrompt()

	if err != nil {
		return fmt.Errorf("error getting expiration: %w", err)
	}
	user.Id, err = cli.adv.CreateUser(context.Background(), &user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	cli.user = user
	return nil
}
