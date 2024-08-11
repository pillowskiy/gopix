package usecase

import (
	"context"
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
)

type AuthRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
}

type AuthUseCase struct {
	repo   AuthRepository
	logger logger.Logger
}

func NewAuthUseCase(repo AuthRepository, logger logger.Logger) *AuthUseCase {
	return &AuthUseCase{repo: repo, logger: logger}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	uniqueUser, err := uc.repo.GetUnique(ctx, user)
	if uniqueUser != nil || err == nil {
		return nil, ErrAlreadyExists
	}

	if err = user.PreCreate(); err != nil {
		return nil, fmt.Errorf("AuthUseCase.Register.PreCreate: %v", err)
	}

	newUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	newUser.HidePassword()

	return &domain.UserWithToken{
		User:  newUser,
		Token: "",
	}, nil
}
