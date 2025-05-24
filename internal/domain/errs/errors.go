package errs

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserCreation = errors.New("failed to create user")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
)
