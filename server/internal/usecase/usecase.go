package usecase

import "errors"

var (
	ErrAlreadyExists      = errors.New("entity already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("entity not found")
)
