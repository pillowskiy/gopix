package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pkg/errors"
)

type UserCache interface {
	Del(ctx context.Context, id string) error
}

type UserRepository interface {
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id domain.ID) (*domain.User, error)
	Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error)
	SetPermissions(ctx context.Context, id domain.ID, permissions int) error
}

type UserFollowingUseCase interface {
	Stats(ctx context.Context, userID domain.ID, executorID *domain.ID) (*domain.FollowingStats, error)
}

type UserUseCase struct {
	repo        UserRepository
	cache       UserCache
	followingUC UserFollowingUseCase
	logger      logger.Logger
}

func NewUserUseCase(
	repo UserRepository, cache UserCache, followingUC UserFollowingUseCase, logger logger.Logger,
) *UserUseCase {
	return &UserUseCase{repo: repo, cache: cache, followingUC: followingUC, logger: logger}
}

func (uc *UserUseCase) GetDetailed(
	ctx context.Context, username string, executorID *domain.ID,
) (*domain.DetailedUser, error) {
	user, err := uc.repo.GetUnique(ctx, &domain.User{Username: username})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	user.HidePassword()

	stats, err := uc.followingUC.Stats(ctx, user.ID, executorID)
	if err != nil {
		return nil, err
	}

	return &domain.DetailedUser{User: *user, Subscription: *stats}, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error) {
	if _, err := uc.GetByID(ctx, id); err != nil {
		return nil, err
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
	ctx context.Context, id domain.ID, deny domain.Permission, allow domain.Permission,
) error {
	user, err := uc.GetByID(ctx, id)
	if err != nil {
		return err
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

func (uc *UserUseCase) GetByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "UserUseCase.GetByID")
	}

	return user, nil
}

func (uc *UserUseCase) deleteCachedUser(ctx context.Context, id domain.ID) {
	if err := uc.cache.Del(ctx, id.String()); err != nil {
		uc.logger.Errorf("UserUseCase.deleteCached: %v", err)
	}
}
