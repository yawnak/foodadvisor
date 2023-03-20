package domain

import "errors"

var (
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrPasswordTooLong = errors.New("password too long")
)
