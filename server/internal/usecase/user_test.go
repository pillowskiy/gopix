package usecase_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/internal/usecase"
	usecaseMock "github.com/pillowskiy/gopix/internal/usecase/mock"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/stretchr/testify/assert"
)

func TestUseUseCase_GetDetailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := usecaseMock.NewMockUserRepository(ctrl)
	mockUserCache := usecaseMock.NewMockUserCache(ctrl)
	mockFollowingUC := usecaseMock.NewMockUserFollowingUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	userUC := usecase.NewUserUseCase(mockUserRepo, mockUserCache, mockFollowingUC, mockLog)

	username := "test"
	uniqueInput := &domain.User{Username: username}
	stats := &domain.FollowingStats{}

	executorID := new(domain.ID)

	mockUser := &domain.User{
		ID:          1,
		Username:    username,
		AvatarURL:   "https://test.com/test.png",
		Permissions: 1,
	}

	mockDetailedUser := &domain.DetailedUser{
		User:         *mockUser,
		Subscription: *stats,
	}

	t.Run("SuccessGetDetailed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), uniqueInput).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Stats(gomock.Any(), mockUser.ID, executorID).Return(stats, nil)

		detailedUser, err := userUC.GetDetailed(context.Background(), username, executorID)
		assert.NoError(t, err)
		assert.Equal(t, mockDetailedUser, detailedUser)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), uniqueInput).Return(nil, repository.ErrNotFound)
		mockFollowingUC.EXPECT().Stats(gomock.Any(), mockUser.ID, executorID).Times(0)

		detailedUser, err := userUC.GetDetailed(context.Background(), username, executorID)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, detailedUser)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), uniqueInput).Return(mockUser, nil)
		mockFollowingUC.EXPECT().Stats(gomock.Any(), mockUser.ID, executorID).Return(nil, errors.New("repo error"))

		detailedUser, err := userUC.GetDetailed(context.Background(), username, executorID)
		assert.Error(t, err)
		assert.Nil(t, detailedUser)
	})
}

func TestUserUseCase_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := usecaseMock.NewMockUserRepository(ctrl)
	mockUserCache := usecaseMock.NewMockUserCache(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	userUC := usecase.NewUserUseCase(mockUserRepo, mockUserCache, nil, mockLog)

	userID := domain.ID(1)
	validUserInput := &domain.User{
		Username:  "test",
		AvatarURL: "https://test.com/test.png",
	}

	mockUser := &domain.User{
		ID:          1,
		Username:    "test",
		AvatarURL:   "https://test.com/test.png",
		Permissions: 1,
	}

	t.Run("SuccessUpdate", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, nil)
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), validUserInput).Return(nil, repository.ErrNotFound)
		mockUserRepo.EXPECT().Update(gomock.Any(), userID, validUserInput).Return(mockUser, nil)
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Return(nil)

		updUser, err := userUC.Update(context.Background(), userID, validUserInput)
		if assert.NoError(t, err) {
			assert.Equal(t, mockUser, updUser)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, repository.ErrNotFound)
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), validUserInput).Times(0)
		mockUserRepo.EXPECT().Update(gomock.Any(), userID, validUserInput).Times(0)
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Times(0)

		updUser, err := userUC.Update(context.Background(), userID, validUserInput)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, updUser)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		anyExceptMockUser := &domain.User{ID: 2}
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, nil)
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), validUserInput).Return(anyExceptMockUser, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), userID, validUserInput).Times(0)
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Times(0)

		updUser, err := userUC.Update(context.Background(), userID, validUserInput)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrAlreadyExists, err)
		assert.Nil(t, updUser)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, nil)
		mockUserRepo.EXPECT().GetUnique(gomock.Any(), validUserInput).Return(nil, repository.ErrNotFound)
		mockUserRepo.EXPECT().Update(
			gomock.Any(), userID, validUserInput,
		).Return(nil, errors.New("repo error"))
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Times(0)

		updUser, err := userUC.Update(context.Background(), userID, validUserInput)

		assert.Error(t, err)
		assert.Nil(t, updUser)
	})
}

func TestUserUseCase_OverwritePermissions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := usecaseMock.NewMockUserRepository(ctrl)
	mockUserCache := usecaseMock.NewMockUserCache(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	userUC := usecase.NewUserUseCase(mockUserRepo, mockUserCache, nil, mockLog)

	userID := domain.ID(1)
	mockUser := &domain.User{
		ID:          userID,
		Username:    "test",
		AvatarURL:   "https://test.com/test.png",
		Permissions: 1,
	}

	initialPermissions := domain.PermissionsAdmin
	denyPerms := domain.PermissionsAdmin
	allowPerms := domain.PermissionsUploadImage

	// TEMP: magic operations
	expectPermissions := (initialPermissions | allowPerms) &^ denyPerms

	t.Run("SuccessPermissionsOverwrite", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, nil)
		mockUserRepo.EXPECT().SetPermissions(gomock.Any(), userID, int(expectPermissions)).Return(nil)
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Return(nil)

		err := userUC.OverwritePermissions(context.Background(), userID, denyPerms, allowPerms)

		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, repository.ErrNotFound)
		mockUserRepo.EXPECT().SetPermissions(gomock.Any(), userID, gomock.Any()).Times(0)
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Times(0)

		err := userUC.OverwritePermissions(context.Background(), userID, denyPerms, allowPerms)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(mockUser, nil)
		mockUserRepo.EXPECT().SetPermissions(
			gomock.Any(), userID, gomock.Any(),
		).Return(errors.New("repo error"))
		mockUserCache.EXPECT().Del(gomock.Any(), userID.String()).Times(0)

		err := userUC.OverwritePermissions(context.Background(), userID, denyPerms, allowPerms)
		assert.Error(t, err)
	})
}
