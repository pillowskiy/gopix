package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
)

type SubscriptionFollowingUseCase interface {
	Follow(ctx context.Context, userID domain.ID, executor *domain.User) error
	Unfollow(ctx context.Context, userID domain.ID, executor *domain.User) error
}

type SubscriptionUserUseCase interface {
	GetByID(ctx context.Context, userID domain.ID) (*domain.User, error)
}

type subscriptionUseCase struct {
	followingUC SubscriptionFollowingUseCase
	userUC      SubscriptionUserUseCase
}

func NewSubscriptionUseCase(
	followingUC SubscriptionFollowingUseCase,
	userUC SubscriptionUserUseCase,
) *subscriptionUseCase {
	return &subscriptionUseCase{followingUC: followingUC, userUC: userUC}
}

func (uc *subscriptionUseCase) Follow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	if err := uc.correctUserRef(ctx, userID); err != nil {
		return err
	}

	return uc.followingUC.Follow(ctx, userID, executor)
}

func (uc *subscriptionUseCase) Unfollow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	if err := uc.correctUserRef(ctx, userID); err != nil {
		return err
	}

	return uc.followingUC.Unfollow(ctx, userID, executor)
}

func (uc *subscriptionUseCase) correctUserRef(ctx context.Context, userID domain.ID) error {
	if _, err := uc.userUC.GetByID(ctx, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrIncorrectUserRef
		}

		return errors.Wrap(err, "SubscriptionUseCase.correctUserRef")
	}

	return nil
}
