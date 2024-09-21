package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/internal/usecase"
	usecaseMock "github.com/pillowskiy/gopix/internal/usecase/mock"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestImageUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	authorID := domain.ID(1)
	fakePath := "fake.png"

	mockImage := &domain.Image{
		AuthorID: authorID,
		Path:     fakePath,
	}

	mockFile := &domain.FileNode{
		Name: "test.png",
		Size: 1024,
		Data: []byte{1, 2, 3},
	}

	expectedTxCall := func(ctx context.Context) {
		mockRepo.EXPECT().
			DoInTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
	}

	t.Run("SuccessCreate", func(t *testing.T) {
		ctx := context.Background()
		expectedTxCall(ctx)
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(mockImage, nil)
		mockStorage.EXPECT().Put(ctx, mockFile).Return(nil)

		createdImage, err := imageUC.Create(context.Background(), mockImage, mockFile)
		if assert.NoError(t, err) {
			assert.NotEqual(t, fakePath, createdImage.Path)
			assert.Equal(t, authorID, createdImage.AuthorID)
			assert.Equal(t, mockFile.Name, createdImage.Path)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		expectedTxCall(context.Background())
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo error"))
		mockStorage.EXPECT().Put(gomock.Any(), mockFile).Times(0)
		mockLog.EXPECT().Error(gomock.Any())

		createdImage, err := imageUC.Create(context.Background(), mockImage, mockFile)
		assert.Error(t, err)
		assert.Nil(t, createdImage)
	})

	t.Run("StorageError", func(t *testing.T) {
		expectedTxCall(context.Background())
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockImage, nil)
		mockStorage.EXPECT().Put(gomock.Any(), mockFile).Return(errors.New("storage error"))
		mockLog.EXPECT().Error(gomock.Any())

		createdImage, err := imageUC.Create(context.Background(), mockImage, mockFile)
		assert.Error(t, err)
		assert.Nil(t, createdImage)
	})
}

func TestImageUseCase_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	authorID := domain.ID(1)

	mockImage := &domain.Image{ID: 1, AuthorID: authorID}
	mockUser := &domain.User{ID: authorID}

	expectGetByIDCall_Repo := func() {
		mockCache.EXPECT().Get(gomock.Any(), mockImage.ID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), mockImage.ID).Return(mockImage, nil)
		mockCache.EXPECT().Set(gomock.Any(), mockImage.ID.String(), mockImage, gomock.Any()).Return(nil)
	}

	expectGetByIDCall_Cached := func() {
		mockCache.EXPECT().Get(gomock.Any(), mockImage.ID.String()).Return(mockImage, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), mockImage.ID).Times(0)
		mockCache.EXPECT().Set(gomock.Any(), mockImage.ID.String(), mockImage, gomock.Any()).Times(0)
	}

	expectedTxCall := func(ctx context.Context) {
		mockRepo.EXPECT().
			DoInTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
	}

	t.Run("SuccessDelete_Repo", func(t *testing.T) {
		expectGetByIDCall_Repo()
		ctx := context.Background()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)

		expectedTxCall(ctx)
		mockRepo.EXPECT().Delete(ctx, mockImage.ID).Return(nil)
		mockStorage.EXPECT().Delete(ctx, mockImage.Path).Return(nil)

		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Return(nil)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)

		assert.NoError(t, err)
	})

	t.Run("SuccessDelete_Cached", func(t *testing.T) {
		expectGetByIDCall_Cached()
		ctx := context.Background()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)

		expectedTxCall(ctx)
		mockRepo.EXPECT().Delete(ctx, mockImage.ID).Return(nil)
		mockStorage.EXPECT().Delete(ctx, mockImage.Path).Return(nil)

		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Return(nil)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.NoError(t, err)
	})

	t.Run("ExistenceError", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), mockImage.ID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), mockImage.ID).Return(nil, repository.ErrNotFound)
		mockCache.EXPECT().Set(gomock.Any(), mockImage.ID.String(), mockImage, gomock.Any()).Times(0)

		mockACL.EXPECT().CanModify(mockUser, mockImage).Times(0)
		mockRepo.EXPECT().Delete(gomock.Any(), mockImage.ID.String()).Times(0)
		mockStorage.EXPECT().Delete(gomock.Any(), mockImage.Path).Times(0)
		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Times(0)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.Error(t, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(false)
		mockRepo.EXPECT().Delete(gomock.Any(), mockImage.ID).Times(0)
		mockStorage.EXPECT().Delete(gomock.Any(), mockImage.Path).Times(0)
		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Times(0)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrForbidden)
	})

	t.Run("RepoError", func(t *testing.T) {
		expectGetByIDCall_Cached()
		repoError := errors.New("repo error")

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)

		expectedTxCall(context.Background())
		mockRepo.EXPECT().Delete(gomock.Any(), mockImage.ID).Return(repoError)
		mockStorage.EXPECT().Delete(gomock.Any(), mockImage.Path).Times(0)
		mockLog.EXPECT().Error(gomock.Any())

		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Times(0)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.Error(t, err)
	})

	t.Run("CacheError_Del", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)

		expectedTxCall(context.Background())
		mockRepo.EXPECT().Delete(gomock.Any(), mockImage.ID).Return(nil)
		mockStorage.EXPECT().Delete(gomock.Any(), mockImage.Path).Return(nil)

		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Return(errors.New("cache error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.NoError(t, err, "Should ignore cache error")
	})

	t.Run("StorageError", func(t *testing.T) {
		expectGetByIDCall_Cached()
		storageError := errors.New("storage error")

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)

		expectedTxCall(context.Background())
		mockRepo.EXPECT().Delete(gomock.Any(), mockImage.ID).Return(nil)
		mockStorage.EXPECT().Delete(gomock.Any(), mockImage.Path).Return(storageError)
		mockLog.EXPECT().Error(gomock.Any())

		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Times(0)

		err := imageUC.Delete(context.Background(), mockImage.ID, mockUser)
		assert.Error(t, err)
	})
}

func TestImageUseCase_GetDetailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	mockDetailedImage := &domain.DetailedImage{
		ImageWithAuthor: domain.ImageWithAuthor{
			Image: domain.Image{
				ID: 1,
			},
		},
	}

	t.Run("SuccessGet", func(t *testing.T) {
		mockRepo.EXPECT().GetDetailed(gomock.Any(), gomock.Any()).Return(mockDetailedImage, nil)

		detailedImage, err := imageUC.GetDetailed(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, mockDetailedImage, detailedImage)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetDetailed(gomock.Any(), gomock.Any()).Return(nil, repository.ErrNotFound)

		detailedImage, err := imageUC.GetDetailed(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrNotFound, err)
		assert.Nil(t, detailedImage)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetDetailed(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo error"))

		detailedImage, err := imageUC.GetDetailed(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, detailedImage)
	})
}

func TestImageUseCase_AddView(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	userID := domain.ID(2)

	t.Run("SuccessAdd", func(t *testing.T) {
		mockRepo.EXPECT().AddView(gomock.Any(), imageID, &userID).Return(nil)

		err := imageUC.AddView(context.Background(), imageID, &userID)
		assert.NoError(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().AddView(gomock.Any(), imageID, &userID).Return(errors.New("repo error"))

		err := imageUC.AddView(context.Background(), imageID, &userID)
		assert.Error(t, err)
	})
}

func TestImageUseCase_AddLike(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	userID := domain.ID(2)

	t.Run("SuccessAdd", func(t *testing.T) {
		mockRepo.EXPECT().AddLike(gomock.Any(), imageID, userID).Return(nil)

		err := imageUC.AddLike(context.Background(), imageID, userID)
		assert.NoError(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().AddLike(gomock.Any(), imageID, userID).Return(errors.New("repo error"))

		err := imageUC.AddLike(context.Background(), imageID, userID)
		assert.Error(t, err)
	})
}

func TestImageUseCase_RemoveLike(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	userID := domain.ID(2)

	t.Run("SuccessRemove", func(t *testing.T) {
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(true, nil)
		mockRepo.EXPECT().RemoveLike(gomock.Any(), imageID, userID).Return(nil)

		err := imageUC.RemoveLike(context.Background(), imageID, userID)
		assert.NoError(t, err)
	})

	t.Run("AlreadyLiked", func(t *testing.T) {
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(false, nil)
		mockRepo.EXPECT().RemoveLike(gomock.Any(), imageID, userID).Times(0)

		err := imageUC.RemoveLike(context.Background(), imageID, userID)
		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrUnprocessable)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(true, nil)
		mockRepo.EXPECT().RemoveLike(gomock.Any(), imageID, userID).Return(errors.New("repo error"))

		err := imageUC.RemoveLike(context.Background(), imageID, userID)
		assert.Error(t, err)
	})
}

func TestImageUseCase_States(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	userID := domain.ID(2)

	mockStates := &domain.ImageStates{
		Viewed: true,
		Liked:  false,
	}

	t.Run("SuccessStates", func(t *testing.T) {
		mockRepo.EXPECT().States(gomock.Any(), imageID, userID).Return(mockStates, nil)

		states, err := imageUC.States(context.Background(), imageID, userID)
		assert.NoError(t, err)
		assert.NotNil(t, states)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().States(gomock.Any(), imageID, userID).Return(nil, errors.New("repo error"))

		states, err := imageUC.States(context.Background(), imageID, userID)
		assert.Error(t, err)
		assert.Nil(t, states)
	})
}

func TestImageUseCase_Discover(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	sort := domain.ImagePopularSort

	mockImage := &domain.Image{
		ID:   1,
		Path: "test.png",
	}

	pagInput := &domain.PaginationInput{
		Page:    1,
		PerPage: 10,
	}

	pag := &domain.Pagination[domain.ImageWithAuthor]{
		PaginationInput: *pagInput,
		Items: []domain.ImageWithAuthor{
			{
				Image: *mockImage,
			},
		},
		Total: 1,
	}

	t.Run("SuccessDiscover", func(t *testing.T) {
		mockRepo.EXPECT().Discover(gomock.Any(), pagInput, sort).Return(pag, nil)

		pag, err := imageUC.Discover(context.Background(), pagInput, sort)
		assert.NoError(t, err)
		assert.NotNil(t, pag)
	})

	t.Run("IncorrectInput", func(t *testing.T) {
		mockRepo.EXPECT().Discover(gomock.Any(), pagInput, sort).Return(nil, repository.ErrIncorrectInput)

		pag, err := imageUC.Discover(context.Background(), pagInput, sort)

		assert.Error(t, err)
		assert.Equal(t, usecase.ErrUnprocessable, err)
		assert.Nil(t, pag)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().Discover(gomock.Any(), pagInput, sort).Return(nil, errors.New("repo error"))

		pag, err := imageUC.Discover(context.Background(), pagInput, sort)
		assert.Error(t, err)
		assert.Nil(t, pag)
	})
}

func TestImageUseCase_HasLike(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	userID := domain.ID(2)

	t.Run("SuccessHasLike_True", func(t *testing.T) {
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(true, nil)

		hasLike := imageUC.HasLike(context.Background(), imageID, userID)
		assert.True(t, hasLike)
	})

	t.Run("SuccessHasLike_False", func(t *testing.T) {
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(false, nil)

		hasLike := imageUC.HasLike(context.Background(), imageID, userID)
		assert.False(t, hasLike)
	})

	t.Run("RepoError", func(t *testing.T) {
		liked := false
		mockRepo.EXPECT().HasLike(gomock.Any(), imageID, userID).Return(liked, errors.New("repo error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		hasLike := imageUC.HasLike(context.Background(), imageID, userID)
		assert.Equal(t, liked, hasLike)
	})
}

func TestImageUseCase_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	imageID := domain.ID(1)
	mockImage := &domain.Image{
		ID:   imageID,
		Path: "test.png",
	}

	t.Run("SuccessGetByID_Cached", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(mockImage, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Times(0)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Times(0)

		image, err := imageUC.GetByID(context.Background(), imageID)
		assert.NoError(t, err)
		assert.NotNil(t, image)
	})

	t.Run("SuccessGetByID_Repo", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Return(nil)

		image, err := imageUC.GetByID(context.Background(), imageID)
		assert.NoError(t, err)
		assert.NotNil(t, image)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, repository.ErrNotFound)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Times(0)

		image, err := imageUC.GetByID(context.Background(), imageID)

		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrNotFound)
		assert.Nil(t, image)
	})

	t.Run("CacheError_Read", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, errors.New("cache error"))
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Return(nil)

		image, err := imageUC.GetByID(context.Background(), imageID)
		assert.NoError(t, err, "Should ignore cache error and call repo")
		assert.NotNil(t, image)
	})

	t.Run("CacheError_Write", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Return(errors.New("cache error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		image, err := imageUC.GetByID(context.Background(), imageID)
		assert.NoError(t, err, "Should ignore cache error")
		assert.NotNil(t, image)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, errors.New("repo error"))
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Times(0)

		image, err := imageUC.GetByID(context.Background(), imageID)
		assert.Error(t, err)
		assert.Nil(t, image)
	})
}

func TestImageUseCase_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockImageRepository(ctrl)
	mockCache := usecaseMock.NewMockImageCache(ctrl)
	mockStorage := usecaseMock.NewMockImageFileStorage(ctrl)
	mockACL := usecaseMock.NewMockImageAccessPolicy(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	imageUC := usecase.NewImageUseCase(mockStorage, mockCache, mockRepo, mockACL, mockLog)

	authorID := domain.ID(1)
	imageID := domain.ID(2)

	mockUser := &domain.User{
		ID: authorID,
	}

	updateInput := &domain.Image{
		Path: "test.png",
	}

	mockImage := &domain.Image{
		ID:       imageID,
		AuthorID: authorID,
	}

	expectGetByIDCall_Repo := func() {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(mockImage, nil)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Return(nil)
	}

	expectGetByIDCall_Cached := func() {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(mockImage, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Times(0)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Times(0)
	}

	t.Run("SuccessUpdate_Repo", func(t *testing.T) {
		expectGetByIDCall_Repo()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), imageID, updateInput).Return(mockImage, nil)
		mockCache.EXPECT().Del(gomock.Any(), imageID.String()).Return(nil)

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockImage, updated)
	})

	t.Run("SuccessUpdate_Cached", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), mockImage.ID, updateInput).Return(mockImage, nil)
		mockCache.EXPECT().Del(gomock.Any(), mockImage.ID.String()).Return(nil)

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockImage, updated)
	})

	t.Run("ExistenceError", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), imageID.String()).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), imageID).Return(nil, repository.ErrNotFound)
		mockCache.EXPECT().Set(gomock.Any(), imageID.String(), mockImage, gomock.Any()).Times(0)

		mockACL.EXPECT().CanModify(mockUser, mockImage).Times(0)
		mockRepo.EXPECT().Update(gomock.Any(), imageID, updateInput).Times(0)
		mockCache.EXPECT().Del(gomock.Any(), imageID.String()).Times(0)

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("CacheError_Del", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), imageID, updateInput).Return(mockImage, nil)
		mockCache.EXPECT().Del(gomock.Any(), imageID.String()).Return(errors.New("cache error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.NoError(t, err, "Should ignore cache error")
		assert.Equal(t, mockImage, updated)
	})

	t.Run("Forbidden", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(false)
		mockRepo.EXPECT().Update(gomock.Any(), imageID, updateInput).Times(0)
		mockCache.EXPECT().Del(gomock.Any(), imageID.String()).Times(0)

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Equal(t, err, usecase.ErrForbidden)
		assert.Nil(t, updated)
	})

	t.Run("RepoError", func(t *testing.T) {
		expectGetByIDCall_Cached()

		mockACL.EXPECT().CanModify(mockUser, mockImage).Return(true)
		mockRepo.EXPECT().Update(gomock.Any(), imageID, updateInput).Return(nil, errors.New("repo error"))
		mockCache.EXPECT().Del(gomock.Any(), imageID.String()).Times(0)

		updated, err := imageUC.Update(context.Background(), imageID, updateInput, mockUser)

		assert.Error(t, err)
		assert.Nil(t, updated)
	})
}
