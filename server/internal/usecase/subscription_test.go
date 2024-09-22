package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	usecaseMock "github.com/pillowskiy/gopix/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubscriptionUseCase_Follow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFollowingUC := usecaseMock.NewMockSubscriptionFollowingUseCase(ctrl)
	mockUserUC := usecaseMock.NewMockSubscriptionUserUseCase(ctrl)
	subscriptionUC := usecase.NewSubscriptionUseCase(mockFollowingUC, mockUserUC)

	mockUserID := domain.ID(1)
	mockUser := &domain.User{ID: mockUserID}
	mockExecutor := &domain.User{ID: 2}

	mockCorrectUserRefCall := func(userID domain.ID) *gomock.Call {
		return mockUserUC.EXPECT().GetByID(gomock.Any(), userID)
	}

	t.Run("SuccessFollow", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor).Return(nil)

		err := subscriptionUC.Follow(context.Background(), mockUserID, mockExecutor)

		assert.NoError(t, err)
	})

	t.Run("ErrUserRef", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(nil, usecase.ErrNotFound)
		mockFollowingUC.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor).Times(0)

		err := subscriptionUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectUserRef, err)
	})

	t.Run("UserUCError", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(nil, errors.New("repo error"))
		mockFollowingUC.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor).Times(0)

		err := subscriptionUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})

	t.Run("FollowingUCError", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor).Return(errors.New("repo error"))

		err := subscriptionUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestSubscriptionUseCase_Unfollow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFollowingUC := usecaseMock.NewMockSubscriptionFollowingUseCase(ctrl)
	mockUserUC := usecaseMock.NewMockSubscriptionUserUseCase(ctrl)
	subscriptionUC := usecase.NewSubscriptionUseCase(mockFollowingUC, mockUserUC)

	mockUserID := domain.ID(1)
	mockUser := &domain.User{ID: mockUserID}
	mockExecutor := &domain.User{ID: 2}

	mockCorrectUserRefCall := func(userID domain.ID) *gomock.Call {
		return mockUserUC.EXPECT().GetByID(gomock.Any(), userID)
	}

	t.Run("SuccessUnfollow", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor).Return(nil)

		err := subscriptionUC.Unfollow(context.Background(), mockUserID, mockExecutor)

		assert.NoError(t, err)
	})

	t.Run("ErrUserRef", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(nil, usecase.ErrNotFound)
		mockFollowingUC.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor).Times(0)

		err := subscriptionUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectUserRef, err)
	})

	t.Run("UserUCError", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(nil, errors.New("repo error"))
		mockFollowingUC.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor).Times(0)

		err := subscriptionUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})

	t.Run("FollowingUCError", func(t *testing.T) {
		mockCorrectUserRefCall(mockUserID).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor).Return(errors.New("repo error"))

		err := subscriptionUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})
}
