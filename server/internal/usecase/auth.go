package usecase

import (
	"context"
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/token"
)

const userTTL = 3600

type AuthRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type AuthCache interface {
	Set(ctx context.Context, id int, user *domain.User, ttl int) error
	Get(ctx context.Context, id int) (*domain.User, error)
}

type AuthUseCase struct {
	repo     AuthRepository
	cache    AuthCache
	tokenGen token.TokenGenerator
	logger   logger.Logger
}

func NewAuthUseCase(
	repo AuthRepository,
	cache AuthCache,
	logger logger.Logger,
	tokenGen token.TokenGenerator,
) *AuthUseCase {
	return &AuthUseCase{repo: repo, logger: logger, tokenGen: tokenGen, cache: cache}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	uniqueUser, err := uc.repo.GetUnique(ctx, user)
	if uniqueUser != nil || err == nil {
		return nil, ErrAlreadyExists
	}

	if err = user.PrepareMutation(); err != nil {
		return nil, fmt.Errorf("AuthUseCase.Register.PreCreate: %v", err)
	}

	newUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	newUser.HidePassword()

	t, err := uc.generateToken(newUser)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithToken{
		User:  newUser,
		Token: t,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	uniqueUser, err := uc.repo.GetUnique(ctx, user)
	if uniqueUser == nil || err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := uniqueUser.ComparePassword(user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}
	uniqueUser.HidePassword()

	t, err := uc.generateToken(uniqueUser)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithToken{
		User:  uniqueUser,
		Token: t,
	}, nil
}

func (uc *AuthUseCase) Verify(ctx context.Context, token string) (*domain.User, error) {
	payload := new(domain.UserPayload)
	if err := uc.tokenGen.VerifyAndScan(token, payload); err != nil {
		return nil, err
	}

	cachedUser, err := uc.cache.Get(ctx, payload.ID)
	if cachedUser != nil {
		return cachedUser, nil
	}

	if err != nil {
		uc.logger.Errorf("authUseCase.cache.GetById: %v", err)
	}

	user, err := uc.repo.GetByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.Set(ctx, payload.ID, user, userTTL); err != nil {
		uc.logger.Errorf("authUseCase.cache.SetUser: %v", err)
	}

	return user, nil
}

func (uc *AuthUseCase) generateToken(user *domain.User) (string, error) {
	t, err := uc.tokenGen.Generate(&domain.UserPayload{
		ID:       user.ID,
		Username: user.Username,
	})
	if err != nil {
		return "", err
	}

	return t, nil
}
