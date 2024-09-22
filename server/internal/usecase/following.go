package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
)

type FollowingRepository interface {
	Follow(ctx context.Context, userID domain.ID, executorID domain.ID) error
	Unfollow(ctx context.Context, userID domain.ID, executorID domain.ID) error
	IsFollowing(ctx context.Context, followerID, folowingID domain.ID) (bool, error)
	Stats(
		ctx context.Context, userID domain.ID, executorID *domain.ID,
	) (*domain.FollowingStats, error)
}

type FollowingUseCase struct {
	repo FollowingRepository
}

func NewFollowingUseCase(repo FollowingRepository) *FollowingUseCase {
	return &FollowingUseCase{repo: repo}
}

func (uc *FollowingUseCase) Follow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	isFollowing, err := uc.IsFollowing(ctx, executor.ID, userID)
	if err != nil {
		return err
	}

	if isFollowing {
		return ErrAlreadyExists
	}

	return uc.repo.Follow(ctx, userID, executor.ID)
}

func (uc *FollowingUseCase) Unfollow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	isFollowing, err := uc.IsFollowing(ctx, executor.ID, userID)
	if err != nil {
		return err
	}

	if !isFollowing {
		return ErrNotFound
	}

	return uc.repo.Unfollow(ctx, userID, executor.ID)
}

func (uc *FollowingUseCase) IsFollowing(ctx context.Context, followerID, followingID domain.ID) (bool, error) {
	isFollowing, err := uc.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return isFollowing, errors.Wrap(err, "FollowingUseCase.IsFollowing")
	}

	return isFollowing, nil
}

func (uc *FollowingUseCase) Stats(
	ctx context.Context, userID domain.ID, executorID *domain.ID,
) (*domain.FollowingStats, error) {
	return uc.repo.Stats(ctx, userID, executorID)
}
