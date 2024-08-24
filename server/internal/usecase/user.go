package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
)

type UserCache interface {
	Del(ctx context.Context, id int) error
}

type UserRepository interface {
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, id int, user *domain.User) (*domain.User, error)
	SetPermissions(ctx context.Context, id int, permissions int) error
}

type UserUseCase struct {
	repo   UserRepository
	cache  UserCache
	logger logger.Logger
}

func NewUserUseCase(repo UserRepository, cache UserCache, logger logger.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, cache: cache, logger: logger}
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

	uc.deleteCachedUser(ctx, id)
	return u, nil
}

func (uc *UserUseCase) OverwritePermissions(
	ctx context.Context, id int, deny domain.Permission, allow domain.Permission,
) error {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	perms := user.Permissions

	if deny != 0 {
		perms = perms &^ int(deny)
	}

	if allow != 0 {
		perms |= int(allow)
	}

	if err := uc.repo.SetPermissions(ctx, id, perms); err != nil {
		return err
	}

	uc.deleteCachedUser(ctx, id)
	return nil
}

func (uc *UserUseCase) deleteCachedUser(ctx context.Context, id int) {
	if err := uc.cache.Del(ctx, id); err != nil {
		uc.logger.Errorf("UserUseCase.deleteCached: %v", err)
	}
}
