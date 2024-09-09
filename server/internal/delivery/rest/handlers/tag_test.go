package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/stretchr/testify/assert"

	handlersMock "github.com/pillowskiy/gopix/internal/delivery/rest/handlers/mock"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/pillowskiy/gopix/pkg/rest"

	"go.uber.org/mock/gomock"
)

func TestTagHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagUC := handlersMock.NewMockTagUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	prepareCreateQuery := func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/tags")
		return c, rec
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	type CreateInput struct {
		Name string
	}

	validCreateInput := CreateInput{
		Name: "tag",
	}

	t.Run("SuccessCreate", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewBuffer(body))
		mockCtxUser(c)

		createdComment := &domain.Tag{
			Name: validCreateInput.Name,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Create(
			ctx, gomock.Any(),
		).Return(createdComment, nil)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		actual := new(domain.Tag)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, createdComment, actual)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewBuffer(body))

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Create(ctx, gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareCreateQuery(nil)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Create(ctx, gomock.Any()).Times(0)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, usecase.ErrAlreadyExists)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestTagHandlers_UpsertImageTag(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagUC := handlersMock.NewMockTagUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareUpsertImageTagQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPut, "/", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/:image_id/tags")
		c.SetParamNames("image_id")
		c.SetParamValues(id)
		return c, rec
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	type UpsertInput struct {
		Name string
	}

	validUpsertInput := UpsertInput{
		Name: "tag",
	}

	t.Run("SuccessUpsert", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		tag := &domain.Tag{Name: validUpsertInput.Name}
		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().UpsertImageTag(ctx, tag, imageID, ctxUser).Return(nil)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery("abs", bytes.NewBuffer(body))
		mockCtxUser(c)

		mockTagUC.EXPECT().UpsertImageTag(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareUpsertImageTagQuery(itoaImageID, nil)
		mockCtxUser(c)

		mockTagUC.EXPECT().UpsertImageTag(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery(itoaImageID, bytes.NewBuffer(body))

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockTagUC.EXPECT().UpsertImageTag(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockTagUC.EXPECT().UpsertImageTag(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(errors.New("internal error"))

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ErrIncorrectImageRef", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().UpsertImageTag(ctx, gomock.Any(), imageID, ctxUser).Return(usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		body, _ := json.Marshal(validUpsertInput)
		c, rec := prepareUpsertImageTagQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().UpsertImageTag(ctx, gomock.Any(), imageID, ctxUser).Return(usecase.ErrForbidden)

		assert.NoError(t, h.UpsertImageTag()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestTagHandlers_Seacrh(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagUC := handlersMock.NewMockTagUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	type SearchInput struct {
		Query string
	}

	validSearchInput := SearchInput{
		Query: "tag",
	}

	prepareSearchQuery := func(query *SearchInput) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		if query != nil {
			req.URL.RawQuery = url.Values{
				"query": []string{query.Query},
			}.Encode()
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/tags")
		return c, rec
	}

	t.Run("SuccessSearch", func(t *testing.T) {
		c, rec := prepareSearchQuery(&validSearchInput)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Search(ctx, validSearchInput.Query).Return([]domain.Tag{}, nil)

		assert.NoError(t, h.Search()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new([]domain.Tag)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, []domain.Tag{}, *actual)
	})

	t.Run("InvalidQuery", func(t *testing.T) {
		c, rec := prepareSearchQuery(nil)

		mockTagUC.EXPECT().Search(gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Search()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareSearchQuery(&validSearchInput)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Search(ctx, validSearchInput.Query).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Search()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTagHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagUC := handlersMock.NewMockTagUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	tagID := 1
	itoaTagID := strconv.Itoa(tagID)

	prepareDeleteQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/tags/:tag_id")
		c.SetParamNames("tag_id")
		c.SetParamValues(id)
		return c, rec
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	t.Run("SuccessDelete", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Delete(ctx, tagID).Return(nil)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidTagID", func(t *testing.T) {
		c, rec := prepareDeleteQuery("abs")

		mockTagUC.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaTagID)

		mockTagUC.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Delete(ctx, tagID).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("TagNotFound", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().Delete(ctx, tagID).Return(usecase.ErrNotFound)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
