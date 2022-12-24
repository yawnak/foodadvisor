package ui

import (
	"fmt"

	"github.com/asstronom/foodadvisor/pkg/domain"
)

type UICli struct {
	user domain.User
}

func (cli *UICli) Run() error {
	err := cli.AuthenticationMenu()
	if err != nil {
		return fmt.Errorf("error authenticating: %w", err)
	}
	return nil
}