package auth

import "errors"

var (
	ErrInvalidCreds  = errors.New("invalid username or password")
	ErrAlreadyExists = errors.New("user already exists")
	ErrMissingField  = errors.New("missing field")
)
