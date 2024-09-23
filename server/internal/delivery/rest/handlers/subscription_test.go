package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	handlersMock "github.com/pillowskiy/gopix/internal/delivery/rest/handlers/mock"
	"github.com/pillowskiy/gopix/internal/usecase"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionHandlers_Follow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubsUC := handlersMock.NewMockSubscriptionUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewSubscriptionHandlers(mockSubsUC, mockLog)

	e := echo.New()

	prepareFollowQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions/:user_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("user_id")
		c.SetParamValues(id)
		return c, rec
	}

	followingID := handlersMock.DomainID()
	itoaFollowingID := followingID.String()

	t.Run("SuccessFollow", func(t *testing.T) {
		c, rec := prepareFollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Follow(ctx, followingID, mockUser).Return(nil)

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SuccessFollow_AlreadyExists", func(t *testing.T) {
		c, rec := prepareFollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Follow(ctx, followingID, mockUser).Return(usecase.ErrAlreadyExists)

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareFollowQuery(itoaFollowingID)

		mockSubsUC.EXPECT().Follow(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectUserID", func(t *testing.T) {
		c, rec := prepareFollowQuery("")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Follow(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserRef", func(t *testing.T) {
		c, rec := prepareFollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Follow(ctx, gomock.Any(), gomock.Any()).Return(usecase.ErrIncorrectUserRef)

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareFollowQuery(itoaFollowingID)
		mockCtxUser(c)

		mockSubsUC.EXPECT().Follow(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Follow()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSubscriptionHandlers_Unfollow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubsUC := handlersMock.NewMockSubscriptionUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewSubscriptionHandlers(mockSubsUC, mockLog)

	e := echo.New()

	prepareUnfollowQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/subscriptions/:user_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("user_id")
		c.SetParamValues(id)
		return c, rec
	}

	followingID := handlersMock.DomainID()
	itoaFollowingID := followingID.String()

	t.Run("SuccessUnfollow", func(t *testing.T) {
		c, rec := prepareUnfollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Unfollow(ctx, followingID, mockUser).Return(nil)

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SuccessUnfollow_NotFound", func(t *testing.T) {
		c, rec := prepareUnfollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Unfollow(ctx, followingID, mockUser).Return(usecase.ErrNotFound)

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareUnfollowQuery(itoaFollowingID)

		mockSubsUC.EXPECT().Unfollow(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectUserID", func(t *testing.T) {
		c, rec := prepareUnfollowQuery("")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Unfollow(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserRef", func(t *testing.T) {
		c, rec := prepareUnfollowQuery(itoaFollowingID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockSubsUC.EXPECT().Unfollow(ctx, gomock.Any(), gomock.Any()).Return(usecase.ErrIncorrectUserRef)

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareUnfollowQuery(itoaFollowingID)
		mockCtxUser(c)

		mockSubsUC.EXPECT().Unfollow(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Unfollow()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
