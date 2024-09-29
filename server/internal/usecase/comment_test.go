package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/internal/usecase"
	usecaseMock "github.com/pillowskiy/gopix/internal/usecase/mock"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCommentUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	imageID := domain.ID(2)
	authorID := domain.ID(3)

	mockComment := &domain.Comment{
		ID:       commentID,
		ImageID:  imageID,
		AuthorID: authorID,
		Text:     "test",
	}

	t.Run("SucessCreate", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().HasUserCommented(gomock.Any(), imageID, authorID).Return(false, nil)
		mockRepo.EXPECT().Create(gomock.Any(), mockComment).Return(mockComment, nil)

		createdComment, err := commentUC.Create(context.Background(), mockComment)
		if assert.NoError(t, err) {
			assert.Equal(t, mockComment, createdComment)
		}
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().HasUserCommented(gomock.Any(), imageID, authorID).Times(0)
		mockRepo.EXPECT().Create(gomock.Any(), mockComment).Times(0)

		createdComment, err := commentUC.Create(context.Background(), mockComment)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
		assert.Nil(t, createdComment)
	})

	t.Run("AlreadyCommented", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().HasUserCommented(gomock.Any(), imageID, authorID).Return(true, nil)
		mockRepo.EXPECT().Create(gomock.Any(), mockComment).Times(0)

		createdComment, err := commentUC.Create(context.Background(), mockComment)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrAlreadyExists, err)
		assert.Nil(t, createdComment)
	})

	t.Run("RepoError_HasUserComment", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().HasUserCommented(gomock.Any(), imageID, authorID).Return(false, errors.New("repo error"))
		mockRepo.EXPECT().Create(gomock.Any(), mockComment).Times(0)

		createdComment, err := commentUC.Create(context.Background(), mockComment)
		assert.Error(t, err)
		assert.Nil(t, createdComment)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().HasUserCommented(gomock.Any(), imageID, authorID).Return(false, nil)
		mockRepo.EXPECT().Create(gomock.Any(), mockComment).Return(nil, errors.New("repo error"))

		createdComment, err := commentUC.Create(context.Background(), mockComment)
		assert.Error(t, err)
		assert.Nil(t, createdComment)
	})
}

func TestCommentUseCase_GetByImageID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	imageID := domain.ID(1)

	pagInput := &domain.PaginationInput{
		PerPage: 10,
		Page:    1,
	}

	sortMethod := domain.CommentNewestSort

	pag := &domain.Pagination[domain.DetailedComment]{
		PaginationInput: *pagInput,
		Items: []domain.DetailedComment{
			{
				Comment: domain.Comment{
					ID:       1,
					ImageID:  imageID,
					AuthorID: 1,
					Text:     "test",
				},
			},
		},
		Total: 1,
	}

	t.Run("SuccessGetByImageID", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().GetByImageID(gomock.Any(), imageID, pagInput, sortMethod).Return(pag, nil)

		pag, err := commentUC.GetByImageID(context.Background(), imageID, pagInput, sortMethod)
		if assert.NoError(t, err) {
			assert.NotNil(t, pag)
		}
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().GetByImageID(gomock.Any(), imageID, pagInput, sortMethod).Times(0)

		pag, err := commentUC.GetByImageID(context.Background(), imageID, pagInput, sortMethod)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
		assert.Nil(t, pag)
	})

	t.Run("Unprocessable", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().GetByImageID(gomock.Any(), imageID, pagInput, sortMethod).Return(nil, repository.ErrIncorrectInput)

		pag, err := commentUC.GetByImageID(context.Background(), imageID, pagInput, sortMethod)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrUnprocessable, err)
		assert.Nil(t, pag)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(&domain.Image{ID: imageID}, nil)
		mockRepo.EXPECT().GetByImageID(gomock.Any(), imageID, pagInput, sortMethod).Return(nil, errors.New("repo error"))

		pag, err := commentUC.GetByImageID(context.Background(), imageID, pagInput, sortMethod)
		assert.Error(t, err)
		assert.Nil(t, pag)
	})
}

func TestCommentUseCase_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	mockComment := &domain.Comment{
		ID:       commentID,
		ImageID:  1,
		AuthorID: 1,
		Text:     "test",
	}

	t.Run("SuccessGetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)

		comment, err := commentUC.GetByID(context.Background(), commentID)
		if assert.NoError(t, err) {
			assert.Equal(t, mockComment, comment)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, repository.ErrNotFound)

		comment, err := commentUC.GetByID(context.Background(), commentID)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, comment)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, errors.New("repo error"))

		comment, err := commentUC.GetByID(context.Background(), commentID)
		assert.Error(t, err)
		assert.Nil(t, comment)
	})
}

func TestCommentUseCase_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	mockComment := &domain.Comment{ID: commentID}
	mockExecutor := &domain.User{ID: 1, Permissions: int(domain.PermissionsAdmin)}

	t.Run("SuccessDelete", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(true)
		mockRepo.EXPECT().Delete(gomock.Any(), commentID).Return(nil)

		err := commentUC.Delete(context.Background(), commentID, mockExecutor)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Times(0)
		mockRepo.EXPECT().Delete(gomock.Any(), commentID).Times(0)

		err := commentUC.Delete(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(false)
		mockRepo.EXPECT().Delete(gomock.Any(), commentID).Times(0)

		err := commentUC.Delete(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(true)
		mockRepo.EXPECT().Delete(gomock.Any(), commentID).Return(errors.New("repo error"))

		err := commentUC.Delete(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestCommentUseCase_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	mockComment := &domain.Comment{ID: commentID}
	mockExecutor := &domain.User{ID: 1, Permissions: int(domain.PermissionsAdmin)}

	t.Run("SuccessUpdate", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), commentID, mockComment).Return(mockComment, nil)

		updComment, err := commentUC.Update(context.Background(), commentID, mockComment, mockExecutor)
		if assert.NoError(t, err) {
			assert.Equal(t, mockComment, updComment)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Times(0)
		mockRepo.EXPECT().Update(gomock.Any(), commentID, mockComment).Times(0)

		updComment, err := commentUC.Update(context.Background(), commentID, mockComment, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, updComment)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(false)
		mockRepo.EXPECT().Update(gomock.Any(), commentID, mockComment).Times(0)

		updComment, err := commentUC.Update(context.Background(), commentID, mockComment, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
		assert.Nil(t, updComment)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockACL.EXPECT().CanModify(mockExecutor, mockComment).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), commentID, mockComment).Return(nil, errors.New("repo error"))

		updComment, err := commentUC.Update(context.Background(), commentID, mockComment, mockExecutor)
		assert.Error(t, err)
		assert.Nil(t, updComment)
	})
}

func TestCommentUseCase_LikeComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	mockComment := &domain.Comment{ID: commentID}
	mockExecutor := &domain.User{ID: 1, Permissions: int(domain.PermissionsAdmin)}

	t.Run("SuccessLikeComment", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(false, nil)
		mockRepo.EXPECT().LikeComment(gomock.Any(), commentID, mockExecutor.ID).Return(nil)

		err := commentUC.LikeComment(context.Background(), commentID, mockExecutor)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Times(0)
		mockRepo.EXPECT().LikeComment(gomock.Any(), commentID, mockExecutor.ID).Times(0)

		err := commentUC.LikeComment(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(true, nil)
		mockRepo.EXPECT().LikeComment(gomock.Any(), commentID, mockExecutor.ID).Times(0)

		err := commentUC.LikeComment(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrAlreadyExists, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(false, nil)
		mockRepo.EXPECT().LikeComment(gomock.Any(), commentID, mockExecutor.ID).Return(errors.New("repo error"))

		err := commentUC.LikeComment(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestCommentUseCase_UnlikeComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	mockExecutor := &domain.User{ID: 1, Permissions: int(domain.PermissionsAdmin)}

	t.Run("SuccessUnlikeComment", func(t *testing.T) {
		// TEMP: We can omit comment existence call, because if pair of commentID and userID doesn't exist, it returns ErrNotFound
		// mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(true, nil)
		mockRepo.EXPECT().UnlikeComment(gomock.Any(), commentID, mockExecutor.ID).Return(nil)

		err := commentUC.UnlikeComment(context.Background(), commentID, mockExecutor)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(false, nil)
		mockRepo.EXPECT().UnlikeComment(gomock.Any(), commentID, mockExecutor).Times(0)

		err := commentUC.UnlikeComment(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrNotFound)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().HasUserLikedComment(gomock.Any(), commentID, mockExecutor.ID).Return(true, nil)
		mockRepo.EXPECT().UnlikeComment(gomock.Any(), commentID, mockExecutor.ID).Return(errors.New("repo error"))

		err := commentUC.UnlikeComment(context.Background(), commentID, mockExecutor)
		assert.Error(t, err)
	})
}

func TestCommentUseCase_GetReplies(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := usecaseMock.NewMockCommentImageUseCase(ctrl)
	mockRepo := usecaseMock.NewMockCommentRepository(ctrl)
	mockACL := usecaseMock.NewMockCommentAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	commentUC := usecase.NewCommentUseCase(mockRepo, mockACL, mockImageUC, mockLog)

	commentID := domain.ID(1)
	executorID := new(domain.ID)

	mockComment := &domain.Comment{ID: commentID}
	mockComments := []domain.DetailedComment{
		{
			Comment: *mockComment,
		},
	}

	t.Run("SuccessGetReplies", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().GetReplies(gomock.Any(), commentID, executorID).Return(mockComments, nil)

		cmts, err := commentUC.GetReplies(context.Background(), commentID, executorID)
		assert.NoError(t, err)
		assert.Equal(t, mockComments, cmts)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().GetReplies(gomock.Any(), commentID, executorID).Times(0)

		cmts, err := commentUC.GetReplies(context.Background(), commentID, executorID)
		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrNotFound)
		assert.Nil(t, cmts)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), commentID).Return(mockComment, nil)
		mockRepo.EXPECT().GetReplies(gomock.Any(), commentID, executorID).Return(nil, errors.New("repo error"))

		cmts, err := commentUC.GetReplies(context.Background(), commentID, executorID)
		assert.Error(t, err)
		assert.Nil(t, cmts)
	})
}
