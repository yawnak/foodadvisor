package domain

import "errors"

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrNoUsername       = errors.New("username does not exist")
	ErrPasswordTooLong  = errors.New("password too long")
)
