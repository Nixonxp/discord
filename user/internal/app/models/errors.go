package models

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUnimplemented = errors.New("unimplemented")
	ErrNotFound      = errors.New("not found")
	ErrCredInvalid   = errors.New("credentials invalid")
	Unauthenticated  = errors.New("unauthenticated")
	PermissionDenied = errors.New("permission denied")
)
