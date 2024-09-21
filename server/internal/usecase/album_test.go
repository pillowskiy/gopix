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

func TestAlbumUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	authorID := domain.ID(1)
	albumID := domain.ID(2)
	albumInput := &domain.Album{Name: "test", AuthorID: authorID}
	mockAlbum := &domain.Album{ID: albumID, Name: albumInput.Name}

	t.Run("SuccessCreate", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), albumInput).Return(mockAlbum, nil)

		createdAlbum, err := albumUC.Create(context.Background(), albumInput)
		assert.NoError(t, err)
		assert.Equal(t, mockAlbum, createdAlbum)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), albumInput).Return(nil, errors.New("repo error"))

		createdAlbum, err := albumUC.Create(context.Background(), albumInput)

		assert.Error(t, err)
		assert.Nil(t, createdAlbum)
	})
}

func TestAlbumUseCase_GetByAuthorID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	authorID := domain.ID(1)
	albumID := domain.ID(2)
	mockAlbums := []domain.Album{
		{
			ID:       albumID,
			Name:     "test",
			AuthorID: authorID,
		},
	}

	t.Run("SuccessGetByAuthorID", func(t *testing.T) {
		mockRepo.EXPECT().GetByAuthorID(gomock.Any(), authorID).Return(mockAlbums, nil)

		albums, err := albumUC.GetByAuthorID(context.Background(), authorID)

		assert.NoError(t, err)
		assert.Equal(t, mockAlbums, albums)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByAuthorID(gomock.Any(), authorID).Return(nil, repository.ErrNotFound)

		albums, err := albumUC.GetByAuthorID(context.Background(), authorID)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, albums)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByAuthorID(gomock.Any(), authorID).Return(nil, errors.New("repo error"))

		albums, err := albumUC.GetByAuthorID(context.Background(), authorID)

		assert.Error(t, err)
		assert.Nil(t, albums)
	})
}

func TestAlbumUseCase_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	albumID := domain.ID(1)
	mockAlbum := &domain.Album{
		ID:   albumID,
		Name: "test",
	}

	t.Run("SuccessGetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)

		album, err := albumUC.GetByID(context.Background(), albumID)

		assert.NoError(t, err)
		assert.Equal(t, mockAlbum, album)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)

		albums, err := albumUC.GetByID(context.Background(), albumID)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, albums)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByAuthorID(gomock.Any(), albumID).Return(nil, errors.New("repo error"))

		albums, err := albumUC.GetByAuthorID(context.Background(), albumID)

		assert.Error(t, err)
		assert.Nil(t, albums)
	})
}

func TestAlbumUseCase_GetAlbumImages(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	albumID := domain.ID(1)
	mockAlbum := &domain.Album{
		ID:   albumID,
		Name: "test",
	}

	pagInput := &domain.PaginationInput{
		PerPage: 10,
		Page:    1,
	}

	mockPag := &domain.Pagination[domain.ImageWithAuthor]{
		PaginationInput: *pagInput,
		Total:           10,
		Items: []domain.ImageWithAuthor{
			{
				Image: domain.Image{
					ID:   1,
					Path: "test.png",
				},
			},
		},
	}

	t.Run("SuccessGetAlbumImages", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockRepo.EXPECT().GetAlbumImages(gomock.Any(), albumID, pagInput).Return(mockPag, nil)

		pag, err := albumUC.GetAlbumImages(context.Background(), albumID, pagInput)

		assert.NoError(t, err)
		assert.Equal(t, mockPag, pag)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().GetAlbumImages(gomock.Any(), albumID, pagInput).Times(0)

		pag, err := albumUC.GetAlbumImages(context.Background(), albumID, pagInput)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, pag)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockRepo.EXPECT().GetAlbumImages(gomock.Any(), albumID, pagInput).Return(nil, errors.New("repo error"))

		pag, err := albumUC.GetAlbumImages(context.Background(), albumID, pagInput)

		assert.Error(t, err)
		assert.Nil(t, pag)
	})
}

func TestAlbumUseCase_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	albumID := domain.ID(1)
	authorID := domain.ID(2)

	mockUser := &domain.User{ID: authorID}
	mockAlbum := &domain.Album{
		ID:       albumID,
		Name:     "test",
		AuthorID: authorID,
	}

	t.Run("SuccessDelete", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)
		mockRepo.EXPECT().Delete(gomock.Any(), albumID).Return(nil)

		err := albumUC.Delete(context.Background(), albumID, mockUser)

		assert.NoError(t, err)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Times(0)
		mockRepo.EXPECT().Delete(gomock.Any(), albumID).Times(0)

		err := albumUC.Delete(context.Background(), albumID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(false)
		mockRepo.EXPECT().Delete(gomock.Any(), albumID).Times(0)

		err := albumUC.Delete(context.Background(), albumID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)
		mockRepo.EXPECT().Delete(gomock.Any(), albumID).Return(errors.New("repo error"))

		err := albumUC.Delete(context.Background(), albumID, mockUser)

		assert.Error(t, err)
	})
}

func TestAlbumUseCase_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	albumID := domain.ID(1)
	authorID := domain.ID(2)

	mockUser := &domain.User{ID: authorID}
	mockAlbum := &domain.Album{
		ID:       albumID,
		Name:     "test",
		AuthorID: authorID,
	}
	updateInput := &domain.Album{Name: "test"}

	t.Run("SuccessUpdate", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), albumID, updateInput).Return(mockAlbum, nil)

		_, err := albumUC.Update(context.Background(), albumID, updateInput, mockUser)

		assert.NoError(t, err)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Times(0)
		mockRepo.EXPECT().Update(gomock.Any(), albumID, updateInput).Times(0)

		album, err := albumUC.Update(context.Background(), albumID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, album)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(false)
		mockRepo.EXPECT().Update(gomock.Any(), albumID, updateInput).Times(0)

		album, err := albumUC.Update(context.Background(), albumID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
		assert.Nil(t, album)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), albumID, updateInput).Return(nil, errors.New("repo error"))

		album, err := albumUC.Update(context.Background(), albumID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Nil(t, album)
	})
}

func TestAlbumUseCase_ExistsAndModifiable(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	albumID := domain.ID(1)
	authorID := domain.ID(2)

	mockUser := &domain.User{ID: authorID}
	mockAlbum := &domain.Album{
		ID:       albumID,
		Name:     "test",
		AuthorID: authorID,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		err := albumUC.ExistsAndModifiable(context.Background(), mockUser, albumID)

		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Times(0)

		err := albumUC.ExistsAndModifiable(context.Background(), mockUser, albumID)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(false)

		err := albumUC.ExistsAndModifiable(context.Background(), mockUser, albumID)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
	})
}

func TestAlbumUseCase_PutImage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	imageID := domain.ID(1)
	albumID := domain.ID(2)

	mockImage := &domain.Image{ID: imageID, AccessLevel: domain.ImageAccessPublic}
	mockAlbum := &domain.Album{ID: albumID}
	mockUser := &domain.User{ID: 1}

	t.Run("SuccessPutImage", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Return(nil)

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.NoError(t, err)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Times(0)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Times(0)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(false)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Times(0)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, usecase.ErrNotFound)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
	})

	t.Run("InvalidImage", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		invalidImage := &domain.Image{ID: imageID, AccessLevel: domain.ImageAccessPrivate}
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(invalidImage, nil)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)

		mockRepo.EXPECT().PutImage(gomock.Any(), albumID, imageID).Return(errors.New("repo error"))

		err := albumUC.PutImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
	})
}

func TestAlbumUseCase_DeleteImage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAlbumRepository(ctrl)
	mockACL := usecaseMock.NewMockAlbumAccessPolicy(ctrl)
	mockImageUC := usecaseMock.NewMockAlbumImageUseCase(ctrl)

	albumUC := usecase.NewAlbumUseCase(mockRepo, mockACL, mockImageUC)

	imageID := domain.ID(1)
	albumID := domain.ID(2)

	mockImage := &domain.Image{ID: imageID, AccessLevel: domain.ImageAccessPublic}
	mockAlbum := &domain.Album{ID: albumID}
	mockUser := &domain.User{ID: 1}

	t.Run("SuccessDeleteImage", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Return(nil)

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.NoError(t, err)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(nil, repository.ErrNotFound)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Times(0)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Times(0)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(false)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Times(0)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrForbidden, err)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, usecase.ErrNotFound)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
	})

	t.Run("InvalidImage", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		invalidImage := &domain.Image{ID: imageID, AccessLevel: domain.ImageAccessPrivate}
		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(invalidImage, nil)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Times(0)

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrIncorrectImageRef, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), albumID).Return(mockAlbum, nil)
		mockACL.EXPECT().CanModify(mockUser, mockAlbum).Return(true)

		mockImageUC.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)

		mockRepo.EXPECT().DeleteImage(gomock.Any(), albumID, imageID).Return(errors.New("repo error"))

		err := albumUC.DeleteImage(context.Background(), albumID, imageID, mockUser)

		assert.Error(t, err)
	})
}
