package user

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidCreds  = errors.New("invalid username or password")
	ErrAlreadyExists = errors.New("user already exists")
	ErrMissingField  = errors.New("missing field")
)
