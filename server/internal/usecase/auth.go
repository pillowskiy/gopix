package usecase

import (
	"context"
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/token"
)

const userTTL = 3600

type AuthRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id domain.ID) (*domain.User, error)

	repository.Transactional
}

type AuthCache interface {
	Set(ctx context.Context, id string, user *domain.User, ttl int) error
	Get(ctx context.Context, id string) (*domain.User, error)
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

	t, err := uc.GenerateToken(newUser)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithToken{
		User:  newUser,
		Token: t,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	uniqueUser, err := uc.GetUnique(ctx, user)
	if err != nil {
		return nil, err
	}

	if uniqueUser.External {
		return nil, ErrInvalidCredentials
	}

	if err := uniqueUser.ComparePassword(user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}
	uniqueUser.HidePassword()

	t, err := uc.GenerateToken(uniqueUser)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithToken{
		User:  uniqueUser,
		Token: t,
	}, nil
}

func (uc *AuthUseCase) GetUnique(ctx context.Context, user *domain.User) (*domain.User, error) {
	uniqueUser, err := uc.repo.GetUnique(ctx, user)
	if uniqueUser == nil || err != nil {
		return nil, ErrInvalidCredentials
	}

	return uniqueUser, nil
}

func (uc *AuthUseCase) Verify(ctx context.Context, token string) (*domain.User, error) {
	payload := new(domain.UserPayload)
	if err := uc.tokenGen.VerifyAndScan(token, payload); err != nil {
		return nil, err
	}

	cachedUser, err := uc.cache.Get(ctx, payload.ID.String())
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

	if err := uc.cache.Set(ctx, payload.ID.String(), user, userTTL); err != nil {
		uc.logger.Errorf("authUseCase.cache.SetUser: %v", err)
	}

	return user, nil
}

func (uc *AuthUseCase) GenerateToken(user *domain.User) (string, error) {
	t, err := uc.tokenGen.Generate(&domain.UserPayload{
		ID:       user.ID,
		Username: user.Username,
	})
	if err != nil {
		return "", err
	}

	return t, nil
}
