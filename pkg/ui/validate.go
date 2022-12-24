package ui

import (
	"errors"
	"strconv"
)

func validateInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("not a number")
	}
	return nil
}
