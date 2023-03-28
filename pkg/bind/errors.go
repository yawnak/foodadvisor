package bind

import (
	"errors"
	"fmt"
)

type ErrUnknown struct {
	Err error
}

func (err *ErrUnknown) Error() string {
	return fmt.Sprintf("unknown error: %v", err.Err)
}

func (err *ErrUnknown) Unwrap() error {
	return err.Err
}

type ErrUnkownField struct {
	Field string
}

func (err *ErrUnkownField) Error() string {
	return fmt.Sprintf("unknown field: %s", err.Field)
}

type ErrSyntax struct {
	Offset int64
}

func (err *ErrSyntax) Error() string {
	return fmt.Sprintf("bad syntax at position %d", err.Offset)
}

type ErrUnmarshalType struct {
	Field  string
	Offset int64
	Type   string
}

func (err *ErrUnmarshalType) Error() string {
	return fmt.Sprintf("invalid value for field %s at position %d (expected type %v)", err.Field, err.Offset, err.Type)
}

var (
	ErrEmptyBody    = errors.New("empty request body")
	ErrBodyTooLarge = errors.New("request body too large")
	ErrBadFormat    = errors.New("bad format")
)
