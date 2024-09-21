package usecase

import "errors"

var (
	ErrAlreadyExists      = errors.New("entity already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("entity not found")
	ErrUnprocessable      = errors.New("unprocessable")
	ErrForbidden          = errors.New("forbidden")

	ErrIncorrectImageRef = errors.New("incorrect image reference provided")
	ErrIncorrectUserRef  = errors.New("incorrect user reference provided")
)
