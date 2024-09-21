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
	tokenMock "github.com/pillowskiy/gopix/pkg/token/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAuthRepository(ctrl)
	mockCache := usecaseMock.NewMockAuthCache(ctrl)
	mockLogger := loggerMock.NewMockLogger(ctrl)
	mockTokenGen := tokenMock.NewMockTokenGenerator(ctrl)

	authUC := usecase.NewAuthUseCase(mockRepo, mockCache, mockLogger, mockTokenGen)

	userID := domain.ID(1)
	token := "test"
	pwd := "test"

	mockUser := &domain.User{ID: 1, Username: "test", PasswordHash: "test"}
	mockPayload := &domain.UserPayload{ID: userID, Username: mockUser.Username}
	mockUserWithToken := &domain.UserWithToken{User: mockUser, Token: token}

	registerInput := &domain.User{Username: "test", PasswordHash: pwd}

	t.Run("SuccessRegister", func(t *testing.T) {
		mockRepo.EXPECT().GetUnique(gomock.Any(), registerInput).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().Create(gomock.Any(), registerInput).Return(mockUser, nil)
		mockTokenGen.EXPECT().Generate(mockPayload).Return(token, nil)

		created, err := authUC.Register(context.Background(), registerInput)
		assert.NoError(t, err)
		assert.Equal(t, mockUserWithToken, created)

		assert.NotEqual(t, pwd, registerInput.PasswordHash, "Input password should be hashed")
		assert.Equal(t, "", created.User.PasswordHash, "Created user password should be empty")
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		mockRepo.EXPECT().GetUnique(gomock.Any(), registerInput).Return(mockUser, nil)

		mockRepo.EXPECT().Create(gomock.Any(), registerInput).Times(0)
		mockTokenGen.EXPECT().Generate(gomock.Any()).Times(0)

		created, err := authUC.Register(context.Background(), registerInput)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrAlreadyExists, err)
		assert.Nil(t, created)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetUnique(gomock.Any(), registerInput).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().Create(gomock.Any(), registerInput).Return(nil, errors.New("repo error"))
		mockTokenGen.EXPECT().Generate(gomock.Any()).Times(0)

		created, err := authUC.Register(context.Background(), registerInput)
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	t.Run("TokenGenError", func(t *testing.T) {
		mockRepo.EXPECT().GetUnique(gomock.Any(), registerInput).Return(nil, repository.ErrNotFound)
		mockRepo.EXPECT().Create(gomock.Any(), registerInput).Return(mockUser, nil)
		mockTokenGen.EXPECT().Generate(gomock.Any()).Return("", errors.New("tokenGen error"))

		created, err := authUC.Register(context.Background(), registerInput)
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAuthRepository(ctrl)
	mockCache := usecaseMock.NewMockAuthCache(ctrl)
	mockLogger := loggerMock.NewMockLogger(ctrl)
	mockTokenGen := tokenMock.NewMockTokenGenerator(ctrl)

	authUC := usecase.NewAuthUseCase(mockRepo, mockCache, mockLogger, mockTokenGen)

	userID := domain.ID(1)
	token := "test"
	pwd := "test"

	mockUser := func() *domain.User {
		u := &domain.User{ID: userID, Username: "test", PasswordHash: pwd}
		err := u.PrepareMutation()
		if err != nil {
			panic(err)
		}
		return u
	}

	t.Run("SuccessLogin", func(t *testing.T) {
		user := mockUser()

		mockPayload := &domain.UserPayload{ID: userID, Username: user.Username}
		mockUserWithToken := &domain.UserWithToken{User: user, Token: token}
		loginInput := &domain.User{Username: "test", PasswordHash: pwd}

		mockRepo.EXPECT().GetUnique(gomock.Any(), loginInput).Return(user, nil)
		mockTokenGen.EXPECT().Generate(mockPayload).Return(token, nil)

		authUser, err := authUC.Login(context.Background(), loginInput)
		assert.NoError(t, err)
		assert.Equal(t, mockUserWithToken, authUser)

		assert.Equal(t, "", authUser.User.PasswordHash, "User password should be empty")
	})

	t.Run("ErrInvalidLogin", func(t *testing.T) {
		loginInput := &domain.User{Username: "test", PasswordHash: pwd}

		mockRepo.EXPECT().GetUnique(gomock.Any(), loginInput).Return(nil, repository.ErrNotFound)
		mockTokenGen.EXPECT().Generate(gomock.Any()).Times(0)

		authUser, err := authUC.Login(context.Background(), loginInput)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrInvalidCredentials, err)
		assert.Nil(t, authUser)
	})

	t.Run("ErrInvalidPassword", func(t *testing.T) {
		user := mockUser()
		loginInput := &domain.User{Username: "test", PasswordHash: "incorrect"}

		mockRepo.EXPECT().GetUnique(gomock.Any(), loginInput).Return(user, nil)
		mockTokenGen.EXPECT().Generate(gomock.Any()).Times(0)

		authUser, err := authUC.Login(context.Background(), loginInput)
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrInvalidCredentials, err)
		assert.Nil(t, authUser)
	})

	t.Run("TokenGenError", func(t *testing.T) {
		user := mockUser()
		mockPayload := &domain.UserPayload{ID: userID, Username: user.Username}
		loginInput := &domain.User{Username: "test", PasswordHash: pwd}

		mockRepo.EXPECT().GetUnique(gomock.Any(), loginInput).Return(user, nil)
		mockTokenGen.EXPECT().Generate(mockPayload).Return("", errors.New("tokenGen error"))

		authUser, err := authUC.Login(context.Background(), loginInput)
		assert.Error(t, err)
		assert.Nil(t, authUser)
	})
}

func TestAuthUseCase_Verify(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecaseMock.NewMockAuthRepository(ctrl)
	mockCache := usecaseMock.NewMockAuthCache(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockTokenGen := tokenMock.NewMockTokenGenerator(ctrl)

	authUC := usecase.NewAuthUseCase(mockRepo, mockCache, mockLog, mockTokenGen)

	token := "test"
	payload := &domain.UserPayload{}
	cacheKey := payload.ID.String()
	mockUser := &domain.User{ID: payload.ID, Username: payload.Username}

	t.Run("SuccessVerify_Cached", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(nil)
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Return(mockUser, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Times(0)
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Times(0)

		user, err := authUC.Verify(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, payload.ID, user.ID)
		assert.Equal(t, payload.Username, user.Username)
	})

	t.Run("SuccessVerify_Repo", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(nil)
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Return(mockUser, nil)
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Return(nil)

		user, err := authUC.Verify(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, payload.ID, user.ID)
		assert.Equal(t, payload.Username, user.Username)
	})

	t.Run("ErrorVerify", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(errors.New("tokenGen error"))
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Times(0)
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Times(0)
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Times(0)

		user, err := authUC.Verify(context.Background(), token)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("CacheReadError", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(nil)
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, errors.New("cache error"))
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Return(mockUser, nil)
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Return(nil)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		user, err := authUC.Verify(context.Background(), token)
		assert.NoError(t, err, "cache error should not be propagated")
		assert.Equal(t, payload.ID, user.ID)
		assert.Equal(t, payload.Username, user.Username)
	})

	t.Run("CacheWriteError", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(nil)
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Return(mockUser, nil)
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Return(errors.New("cache error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		user, err := authUC.Verify(context.Background(), token)
		assert.NoError(t, err, "cache error should not be propagated")
		assert.Equal(t, payload.ID, user.ID)
		assert.Equal(t, payload.Username, user.Username)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockTokenGen.EXPECT().VerifyAndScan(token, payload).Return(nil)
		mockCache.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), payload.ID).Return(nil, errors.New("repo error"))
		mockCache.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).Times(0)

		user, err := authUC.Verify(context.Background(), token)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
