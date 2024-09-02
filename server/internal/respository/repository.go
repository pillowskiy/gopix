package repository

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrIncorrectInput = errors.New("incorrect input")
)
