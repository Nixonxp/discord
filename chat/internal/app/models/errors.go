package models

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUnimplemented = errors.New("unimplemented")
	ErrCreate        = errors.New("error create")
)
