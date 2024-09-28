package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	_, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	prepareCreateQuery := func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/tags", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
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
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	prepareUpsertImageTagQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/images/:image_id/tags", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("image_id")
		c.SetParamValues(id)
		return c, rec
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

func TestTagHandlers_DeleteImageTag(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagUC := handlersMock.NewMockTagUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewTagHandlers(mockTagUC, mockLog)

	e := echo.New()

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	tagID := handlersMock.DomainID()
	itoaTagID := tagID.String()

	prepareDeleteImageTagQuery := func(
		imageID, tagID string,
	) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/images/:image_id/tags/:tag_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("image_id", "tag_id")
		c.SetParamValues(imageID, tagID)

		return c, rec
	}

	t.Run("SuccessDeleteTag", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(ctx, tagID, imageID, ctxUser).Return(nil)

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidUserContext", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, itoaTagID)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery("abs", itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidTagID", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, "abs")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(
			ctx, gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(usecase.ErrForbidden)

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(
			ctx, gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteImageTagQuery(itoaImageID, itoaTagID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockTagUC.EXPECT().DeleteImageTag(
			ctx, gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.DeleteImageTag()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
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
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images/tags", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		if query != nil {
			req.URL.RawQuery = url.Values{
				"query": []string{query.Query},
			}.Encode()
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
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
	_, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewTagHandlers(mockTagUC, mockLog)
	e := echo.New()

	tagID := handlersMock.DomainID()
	itoaTagID := tagID.String()

	prepareDeleteQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/images/tags/:tag_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tag_id")
		c.SetParamValues(id)
		return c, rec
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
