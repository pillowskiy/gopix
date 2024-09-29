package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrIncorrectInput = errors.New("incorrect input")
)

type InTransactionalCall func(ctx context.Context) error

type Transactional interface {
	DoInTransaction(context.Context, InTransactionalCall) error
}
