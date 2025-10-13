package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAlreadyExists      = errors.New("user already exists")
	ErrMissingField       = errors.New("missing field")
	ErrInvalidPermissions = errors.New("invalid permissions")
)
