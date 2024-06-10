package models

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUnimplemented = errors.New("unimplemented")
	ErrCreate        = errors.New("error create")
	ErrEmpty         = errors.New("error empty")
	ErrNotFound      = errors.New("not found")
	Unauthenticated  = errors.New("unauthenticated")
	PermissionDenied = errors.New("permission denied")
)
