package simulator

import "errors"

var (
	ErrAlreadyExists      = errors.New("simulator already exists")
	ErrMissingField       = errors.New("missing field")
	ErrInvalidPermissions = errors.New("invalid permissions")
	ErrSimulatorNotFound  = errors.New("simulator not found")
	ErrNegativeWeight     = errors.New("weight cannot be negative")
	ErrWrongRange         = errors.New("weight range is invalid")
	ErrZeroIncrement      = errors.New("weight increment cannot be zero")
)
