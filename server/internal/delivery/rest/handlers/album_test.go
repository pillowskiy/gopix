package handlers_test

import (
	"bytes"
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
	handlersMock "github.com/pillowskiy/gopix/internal/delivery/rest/handlers/mock"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAlbumHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAlbumUC := handlersMock.NewMockalbumUseCase(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	prepareCreateQuery := func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/albums", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	type CreateInput struct {
		Name        string
		Description string
	}

	validCreateInput := CreateInput{
		Name:        "test",
		Description: "test description",
	}

	t.Run("SuccessCreate", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewReader(body))
		mockCtxUser(c)

		withCorrectAuthor := &domain.Album{
			Name:        validCreateInput.Name,
			Description: validCreateInput.Description,
			AuthorID:    ctxUser.ID,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Create(ctx, withCorrectAuthor).Return(withCorrectAuthor, nil)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(bytes.NewBuffer(body))

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Create(ctx, gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareCreateQuery(nil)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Create(ctx, gomock.Any()).Times(0)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)

		c, rec := prepareCreateQuery(bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)

		c, rec := prepareCreateQuery(bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAlbumHandlers_GetByAuthorID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAlbumUC := handlersMock.NewMockalbumUseCase(ctrl)

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	authorID := handlersMock.DomainID()
	itoaAuthorID := authorID.String()

	prepareGetByAuthorIDQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/albums/users/:user_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("user_id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessGetByAuthorID", func(t *testing.T) {
		c, rec := prepareGetByAuthorIDQuery(itoaAuthorID)

		mockAlbumUC.EXPECT().GetByAuthorID(gomock.Any(), authorID).Return([]domain.Album{}, nil)

		assert.NoError(t, h.GetByAuthorID()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectAuthorID", func(t *testing.T) {
		c, rec := prepareGetByAuthorIDQuery("abs")

		mockAlbumUC.EXPECT().GetByAuthorID(gomock.Any(), authorID).Times(0)

		assert.NoError(t, h.GetByAuthorID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetByAuthorIDQuery(itoaAuthorID)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockAlbumUC.EXPECT().GetByAuthorID(gomock.Any(), authorID).Return(nil, errors.New("internal error"))

		assert.NoError(t, h.GetByAuthorID()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAlbumHandlers_GetAlbumImages(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAlbumUC := handlersMock.NewMockalbumUseCase(ctrl)

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	albumID := handlersMock.DomainID()
	itoaAlbumID := albumID.String()

	type AlbumImagesQuery struct {
		Limit int `query:"limit"`
		Page  int `query:"page"`
	}

	validAlbumImagesQuery := &AlbumImagesQuery{
		Limit: 10,
		Page:  1,
	}

	prepareGetAlbumImagesQuery := func(id string, query *AlbumImagesQuery) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/albums/:album_id/images", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		if query != nil {
			params := make(url.Values)
			params.Add("limit", strconv.Itoa(query.Limit))
			params.Add("page", strconv.Itoa(query.Page))
			req.URL.RawQuery = params.Encode()
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("album_id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessGetAlbumImages", func(t *testing.T) {
		c, rec := prepareGetAlbumImagesQuery(itoaAlbumID, validAlbumImagesQuery)

		pagInput := &domain.PaginationInput{
			PerPage: validAlbumImagesQuery.Limit,
			Page:    validAlbumImagesQuery.Page,
		}

		pag := &domain.Pagination[domain.ImageWithAuthor]{
			Items: []domain.ImageWithAuthor{
				{
					Image: domain.Image{
						ID: 1,
					},
				},
			},
			PaginationInput: *pagInput,
			Total:           10,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().GetAlbumImages(ctx, albumID, pagInput).Return(pag, nil)

		assert.NoError(t, h.GetAlbumImages()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.Pagination[domain.ImageWithAuthor])
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, pag, actual)
	})

	t.Run("IncorrectQuery", func(t *testing.T) {
		c, rec := prepareGetAlbumImagesQuery(itoaAlbumID, nil)

		assert.NoError(t, h.GetAlbumImages()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := prepareGetAlbumImagesQuery("abs", validAlbumImagesQuery)
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().GetAlbumImages(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetAlbumImages()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		c, rec := prepareGetAlbumImagesQuery(itoaAlbumID, validAlbumImagesQuery)
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().GetAlbumImages(ctx, albumID, gomock.Any()).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.GetAlbumImages()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetAlbumImagesQuery(itoaAlbumID, validAlbumImagesQuery)
		ctx := rest.GetEchoRequestCtx(c)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockAlbumUC.EXPECT().GetAlbumImages(ctx, albumID, gomock.Any()).Return(nil, errors.New("internal error"))

		assert.NoError(t, h.GetAlbumImages()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAlbumHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAlbumUC := handlersMock.NewMockalbumUseCase(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	albumID := handlersMock.DomainID()
	itoaAlbumID := albumID.String()

	prepareDeleteAlbumQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/albums/:album_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("album_id")
		c.SetParamValues(id)

		return c, rec
	}

	t.Run("SuccessDeleteAlbum", func(t *testing.T) {
		c, rec := prepareDeleteAlbumQuery(itoaAlbumID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Delete(ctx, albumID, ctxUser).Return(nil)
		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectAlbumID", func(t *testing.T) {
		c, rec := prepareDeleteAlbumQuery("abs")
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().Delete(ctx, gomock.Any(), ctxUser).Times(0)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareDeleteAlbumQuery(itoaAlbumID)
		ctx := rest.GetEchoRequestCtx(c)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockAlbumUC.EXPECT().Delete(ctx, albumID, gomock.Any()).Times(0)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteAlbumQuery(itoaAlbumID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockAlbumUC.EXPECT().Delete(ctx, albumID, ctxUser).Return(errors.New("internal error"))

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := prepareDeleteAlbumQuery(itoaAlbumID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockAlbumUC.EXPECT().Delete(ctx, albumID, ctxUser).Return(usecase.ErrForbidden)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestAlbumHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumUC := handlersMock.NewMockalbumUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	albumID := handlersMock.DomainID()
	itoaAlbumID := albumID.String()

	prepareUpdateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/albums/:album_id", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("album_id")
		c.SetParamValues(id)
		return c, rec
	}

	type UpdateInput struct {
		Name        string
		Description string
	}

	validUpdateInput := UpdateInput{
		Name:        "test",
		Description: "test",
	}

	t.Run("SuccessUpdate", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaAlbumID, bytes.NewBuffer(body))
		mockCtxUser(c)

		updatedAlbum := &domain.Album{
			Name:        validUpdateInput.Name,
			Description: validUpdateInput.Description,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(
			ctx, albumID, updatedAlbum, ctxUser,
		).Return(updatedAlbum, nil)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.Album)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, updatedAlbum, actual)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaAlbumID, bytes.NewBuffer(body))

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(ctx, albumID, gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidAlbumID", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery("abc", bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(ctx, gomock.Any(), gomock.Any(), ctxUser).Times(0)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaAlbumID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(ctx, albumID, gomock.Any(), ctxUser).Return(nil, usecase.ErrForbidden)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaAlbumID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(ctx, albumID, gomock.Any(), ctxUser).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaAlbumID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().Update(ctx, albumID, gomock.Any(), ctxUser).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAlbumHandlers_PutImage(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAlbumUC := handlersMock.NewMockalbumUseCase(mockCtrl)
	mockLog := loggerMock.NewMockLogger(mockCtrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	albumID := handlersMock.DomainID()
	itoaAlbumID := albumID.String()

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	preparePutImageQuery := func(albumID string, imageID string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/albums/:album_id/images/:image_id")
		c.SetParamNames("album_id", "image_id")
		c.SetParamValues(albumID, imageID)
		return c, rec
	}

	t.Run("SuccessPutImage", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, imageID, ctxUser).Return(nil)

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, imageID, gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, "abs")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, gomock.Any(), ctxUser).Times(0)

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectAlbumID", func(t *testing.T) {
		c, rec := preparePutImageQuery("abs", itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, gomock.Any(), albumID, ctxUser).Times(0)

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, imageID, ctxUser).Return(usecase.ErrForbidden)

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, imageID, ctxUser).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		c, rec := preparePutImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().PutImage(ctx, albumID, imageID, ctxUser).Return(usecase.ErrNotFound)

		assert.NoError(t, h.PutImage()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAlbumHandlers_DeleteImage(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAlbumUC := handlersMock.NewMockalbumUseCase(mockCtrl)
	mockLog := loggerMock.NewMockLogger(mockCtrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewAlbumHandlers(mockAlbumUC, mockLog)

	e := echo.New()

	albumID := handlersMock.DomainID()
	itoaAlbumID := albumID.String()

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	prepareDeleteImageQuery := func(albumID string, imageID string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/albums/:album_id/images/:image_id")
		c.SetParamNames("album_id", "image_id")
		c.SetParamValues(albumID, imageID)
		return c, rec
	}

	t.Run("SuccessDeleteImage", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, imageID, ctxUser).Return(nil)

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, imageID, gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, "abs")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, gomock.Any(), ctxUser).Times(0)

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectAlbumID", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery("abs", itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, gomock.Any(), albumID, ctxUser).Times(0)

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, imageID, ctxUser).Return(usecase.ErrForbidden)

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, imageID, ctxUser).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AlbumNotFound", func(t *testing.T) {
		c, rec := prepareDeleteImageQuery(itoaAlbumID, itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockAlbumUC.EXPECT().DeleteImage(ctx, albumID, imageID, ctxUser).Return(usecase.ErrNotFound)

		assert.NoError(t, h.DeleteImage()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
