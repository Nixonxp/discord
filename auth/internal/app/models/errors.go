package models

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUnimplemented = errors.New("unimplemented")
	ErrCredInvalid   = errors.New("credentials invalid")
	Unauthenticated  = errors.New("unauthenticated")
)
