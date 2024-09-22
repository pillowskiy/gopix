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

func TestFollowingUseCase_Follow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockFollowingRepository(ctrl)
	followingUC := usecase.NewFollowingUseCase(mockRepo)

	mockUserID := domain.ID(1)
	mockExecutor := &domain.User{ID: 2}

	t.Run("SuccessFollow", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(false, nil)
		mockRepo.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor.ID).Return(nil)

		err := followingUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.NoError(t, err)
	})

	t.Run("AlreadyFollowing", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(true, nil)
		mockRepo.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor.ID).Times(0)

		err := followingUC.Follow(context.Background(), mockUserID, mockExecutor)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrAlreadyExists, err)
	})

	t.Run("RepoErrror_IsFollowing", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(false, errors.New("repo error"))
		mockRepo.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor.ID).Times(0)

		err := followingUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(false, nil)
		mockRepo.EXPECT().Follow(gomock.Any(), mockUserID, mockExecutor.ID).Return(errors.New("repo error"))

		err := followingUC.Follow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestFollowingUseCase_Unfollow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockFollowingRepository(ctrl)
	followingUC := usecase.NewFollowingUseCase(mockRepo)

	mockUserID := domain.ID(1)
	mockExecutor := &domain.User{ID: 2}

	t.Run("SuccessUnfolow", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(true, nil)
		mockRepo.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor.ID).Return(nil)

		err := followingUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.NoError(t, err)
	})

	t.Run("NotFollowing", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(false, nil)
		mockRepo.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor.ID).Times(0)

		err := followingUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("RepoErrror_IsFollowing", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(false, errors.New("repo error"))
		mockRepo.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor.ID).Times(0)

		err := followingUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockExecutor.ID, mockUserID).Return(true, nil)
		mockRepo.EXPECT().Unfollow(gomock.Any(), mockUserID, mockExecutor.ID).Return(errors.New("repo error"))

		err := followingUC.Unfollow(context.Background(), mockUserID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestFollowingUseCase_IsFollowing(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockFollowingRepository(ctrl)
	followingUC := usecase.NewFollowingUseCase(mockRepo)

	mockFollowerID := domain.ID(1)
	mockFollowingID := domain.ID(2)

	t.Run("SuccessIsFollowing", func(t *testing.T) {
		isFollowing := true
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockFollowerID, mockFollowingID).Return(isFollowing, nil)

		actual, err := followingUC.IsFollowing(context.Background(), mockFollowerID, mockFollowingID)
		assert.NoError(t, err)
		assert.Equal(t, isFollowing, actual)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().IsFollowing(gomock.Any(), mockFollowerID, mockFollowingID).Return(false, errors.New("repo error"))

		_, err := followingUC.IsFollowing(context.Background(), mockFollowerID, mockFollowingID)
		assert.Error(t, err)
	})
}

func TestFollowingUseCase_Stats(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockFollowingRepository(ctrl)
	followingUC := usecase.NewFollowingUseCase(mockRepo)

	mockUserID := domain.ID(1)
	mockExecutorID := new(domain.ID)
	mockStats := &domain.FollowingStats{}

	t.Run("SuccessStats", func(t *testing.T) {
		mockRepo.EXPECT().Stats(gomock.Any(), mockUserID, mockExecutorID).Return(mockStats, nil)

		actual, err := followingUC.Stats(context.Background(), mockUserID, mockExecutorID)
		assert.NoError(t, err)
		assert.Equal(t, mockStats, actual)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().Stats(gomock.Any(), mockUserID, mockExecutorID).Return(nil, errors.New("repo error"))

		_, err := followingUC.Stats(context.Background(), mockUserID, mockExecutorID)
		assert.Error(t, err)
	})
}
