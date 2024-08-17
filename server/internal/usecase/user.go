package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
)

type UserRepository interface {
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, id int, user *domain.User) (*domain.User, error)
}

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Update(ctx context.Context, id int, user *domain.User) (*domain.User, error) {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	existUser, err := uc.repo.GetUnique(ctx, user)
	if existUser != nil && existUser.ID != id {
		return nil, ErrAlreadyExists
	}

	if err := user.PrepareMutation(); err != nil {
		return nil, err
	}

	u, err := uc.repo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}
	u.HidePassword()

	return u, nil
}
