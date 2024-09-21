package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/internal/usecase"
	usecaseMock "github.com/pillowskiy/gopix/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTagUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockTagImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockTagRepository(ctrl)
	mockACL := usecaseMock.NewMockTagAccessPolicy(ctrl)

	tagUC := usecase.NewTagUseCase(mockRepo, mockACL, mockImageUC)

	tagInput := &domain.Tag{Name: "test"}
	mockTag := &domain.Tag{ID: 1, Name: "test"}

	t.Run("SuccessCreate", func(t *testing.T) {
		mockRepo.EXPECT().GetByName(gomock.Any(), tagInput.Name).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().Create(gomock.Any(), tagInput).Return(mockTag, nil)

		createdTag, err := tagUC.Create(context.Background(), tagInput)
		assert.NoError(t, err)
		assert.Equal(t, mockTag, createdTag)
	})

	t.Run("AlreadyExist", func(t *testing.T) {
		mockRepo.EXPECT().GetByName(gomock.Any(), tagInput.Name).Return(mockTag, nil)

		createdTag, err := tagUC.Create(context.Background(), tagInput)
		assert.Error(t, err)
		assert.Nil(t, createdTag)
	})
}

func TestTagUseCase_UpsertImageTag(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockTagImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockTagRepository(ctrl)
	mockACL := usecaseMock.NewMockTagAccessPolicy(ctrl)

	tagUC := usecase.NewTagUseCase(mockRepo, mockACL, mockImageUC)

	tagInput := &domain.Tag{Name: "test"}

	imageID := domain.ID(1)
	userID := domain.ID(2)

	mockUser := &domain.User{ID: userID, Permissions: int(domain.PermissionsAdmin)}
	mockImage := &domain.Image{ID: imageID, AuthorID: userID}

	t.Run("SuccessUpsertImageTag", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockACL.EXPECT().CanModifyImageTags(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().UpsertImageTags(gomock.Any(), tagInput, imageID).Return(nil)

		err := tagUC.UpsertImageTag(context.Background(), tagInput, imageID, mockUser)

		assert.NoError(t, err)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, usecase.ErrNotFound)
		mockACL.EXPECT().CanModifyImageTags(mockUser, mockImage).Times(0)
		mockRepo.EXPECT().UpsertImageTags(gomock.Any(), tagInput, imageID).Times(0)

		err := tagUC.UpsertImageTag(context.Background(), tagInput, imageID, mockUser)
		if assert.Error(t, err) {
			assert.Equal(t, usecase.ErrIncorrectImageRef, err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockACL.EXPECT().CanModifyImageTags(mockUser, mockImage).Return(false)
		mockRepo.EXPECT().UpsertImageTags(gomock.Any(), tagInput, imageID).Times(0)

		err := tagUC.UpsertImageTag(context.Background(), tagInput, imageID, mockUser)
		if assert.Error(t, err) {
			assert.Equal(t, usecase.ErrForbidden, err)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockACL.EXPECT().CanModifyImageTags(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().UpsertImageTags(gomock.Any(), tagInput, imageID).Return(errors.New("repo error"))

		err := tagUC.UpsertImageTag(context.Background(), tagInput, imageID, mockUser)
		assert.Error(t, err)
	})
}

func TestTagUseCase_Search(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockTagImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockTagRepository(ctrl)
	mockACL := usecaseMock.NewMockTagAccessPolicy(ctrl)

	tagUC := usecase.NewTagUseCase(mockRepo, mockACL, mockImageUC)

	query := "test"
	tags := []domain.Tag{
		{ID: 1},
		{ID: 2},
	}

	t.Run("SuccessSearch", func(t *testing.T) {
		mockRepo.EXPECT().Search(gomock.Any(), query).Return(tags, nil)

		sTags, err := tagUC.Search(context.Background(), query)
		assert.NoError(t, err)
		assert.Equal(t, tags, sTags)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().Search(gomock.Any(), query).Return(nil, errors.New("repo error"))

		sTags, err := tagUC.Search(context.Background(), query)
		assert.Error(t, err)
		assert.Nil(t, sTags)
	})
}

func TestTagUseCase_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockTagImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockTagRepository(ctrl)
	mockACL := usecaseMock.NewMockTagAccessPolicy(ctrl)

	tagUC := usecase.NewTagUseCase(mockRepo, mockACL, mockImageUC)

	tagID := domain.ID(1)
	tag := &domain.Tag{ID: tagID}

	t.Run("SuccessDelete", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(tag, nil)
		mockRepo.EXPECT().Delete(gomock.Any(), tagID).Return(nil)

		err := tagUC.Delete(context.Background(), tagID)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(nil, usecase.ErrNotFound)
		mockRepo.EXPECT().Delete(gomock.Any(), tagID).Times(0)

		err := tagUC.Delete(context.Background(), tagID)
		assert.Error(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(tag, nil)
		mockRepo.EXPECT().Delete(gomock.Any(), tagID).Return(errors.New("repo error"))

		err := tagUC.Delete(context.Background(), tagID)
		assert.Error(t, err)
	})
}

func TestTagUseCase_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockTagImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockTagRepository(ctrl)
	mockACL := usecaseMock.NewMockTagAccessPolicy(ctrl)

	tagUC := usecase.NewTagUseCase(mockRepo, mockACL, mockImageUC)

	tagID := domain.ID(1)
	tag := &domain.Tag{ID: tagID}

	t.Run("SuccessGetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(tag, nil)

		tTag, err := tagUC.GetByID(context.Background(), tagID)
		assert.NoError(t, err)
		assert.Equal(t, tag, tTag)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(nil, repository.ErrNotFound)

		tTag, err := tagUC.GetByID(context.Background(), tagID)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, tTag)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), tagID).Return(nil, errors.New("repo error"))

		tTag, err := tagUC.GetByID(context.Background(), tagID)
		assert.Error(t, err)
		assert.Nil(t, tTag)
	})
}
