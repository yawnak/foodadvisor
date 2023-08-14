package domain

import "errors"

var (
	ErrUnknownError         = errors.New("unknown error")
	ErrWrongCredentials     = errors.New("wrong credentials")
	ErrNoUsername           = errors.New("username does not exist")
	ErrPasswordTooLong      = errors.New("password too long")
	ErrBadToken             = errors.New("auth token is of wrong format")
	ErrInvalidSigningMethod = errors.New("token is signed with wrong method")
	ErrDuplicateResourse    = errors.New("resourse already exists")
	ErrResourseNotFound     = errors.New("resourse not found")
)
