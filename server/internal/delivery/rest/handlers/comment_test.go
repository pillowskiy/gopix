package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func TestCommentHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	createPath := func(id string) string {
		return fmt.Sprintf("/api/v1/images/%s/comments", id)
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	prepareCreateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, createPath(id), body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/:image_id/comments")
		c.SetParamNames("image_id")
		c.SetParamValues(id)
		return c, rec
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	type CreateInput struct {
		Text string `json:"text"`
	}

	validCreateInput := CreateInput{
		Text: "test",
	}

	t.Run("SuccessCreate", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		createdComment := &domain.Comment{
			Text: validCreateInput.Text,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().Create(
			ctx, gomock.Any(),
		).Return(createdComment, nil)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		actual := new(domain.Comment)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, createdComment, actual)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery("abc", bytes.NewBuffer(body))

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareCreateQuery(itoaImageID, nil)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyTextField", func(t *testing.T) {
		invalidInput := CreateInput{
			Text: "",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)

		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)

		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCommentHandlers_GetByImageID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	getByImageIDPath := func(id string) string {
		return fmt.Sprintf("/api/v1/images/%s/comments", id)
	}

	type ImageCommentsQuery struct {
		Limit int    `query:"limit"`
		Page  int    `query:"page"`
		Sort  string `query:"sort"`
	}

	validImageCommentsQuery := ImageCommentsQuery{
		Limit: 10,
		Page:  1,
		Sort:  "newest",
	}

	prepareGetByImageIDQuery := func(
		id string, query *ImageCommentsQuery,
	) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, getByImageIDPath(id), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		if query != nil {
			params := make(url.Values)
			params.Add("limit", strconv.Itoa(query.Limit))
			params.Add("page", strconv.Itoa(query.Page))
			params.Add("sort", query.Sort)
			req.URL.RawQuery = params.Encode()
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/:image_id/comments")
		c.SetParamNames("image_id")
		c.SetParamValues(id)

		return c, rec
	}

	t.Run("SuccessByImageID", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, &validImageCommentsQuery)

		pag := &domain.Pagination[domain.DetailedComment]{
			Items: []domain.DetailedComment{
				{
					Comment: domain.Comment{},
					Author:  domain.CommentAuthor{},
				},
			},
			PaginationInput: domain.PaginationInput{
				PerPage: 10,
				Page:    1,
			},
			Total: 10,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(pag, nil)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.Pagination[domain.DetailedComment])
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, pag, actual)
	})

	t.Run("InvalidQuery", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, nil)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, &ImageCommentsQuery{
			Limit: 0,
			Page:  0,
			Sort:  "",
		})

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, &validImageCommentsQuery)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("server error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, &validImageCommentsQuery)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCommentHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := 1
	itoaCommentID := strconv.Itoa(commentID)

	updatePath := func(id string) string {
		return fmt.Sprintf("/api/v1/images/comments/%s", id)
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	prepareUpdateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, updatePath(id), body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("images/comments/:comment_id")
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	type UpdateInput struct {
		Text string `json:"text"`
	}

	validUpdateInput := UpdateInput{
		Text: "test",
	}

	t.Run("SuccessUpdate", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)

		updatedComment := &domain.Comment{
			Text: validUpdateInput.Text,
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().Update(
			ctx, commentID, gomock.Any(), gomock.Any(),
		).Return(updatedComment, nil)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.Comment)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, updatedComment, actual)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery("abc", bytes.NewBuffer(body))
		mockCtxUser(c)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareUpdateQuery(itoaCommentID, nil)
		mockCtxUser(c)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyTextField", func(t *testing.T) {
		invalidInput := UpdateInput{
			Text: "",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Update(
			ctx, commentID, gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Update(
			ctx, commentID, gomock.Any(), gomock.Any(),
		).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCommentHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := 1
	itoaCommentID := strconv.Itoa(commentID)

	deletePath := func(id string) string {
		return fmt.Sprintf("/api/v1/images/comments/%s", id)
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	prepareDeleteQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, deletePath(id), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("comments/:comment_id")
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	t.Run("SuccessDelete", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().Delete(
			ctx, gomock.Any(), gomock.Any(),
		).Return(nil)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidCommentID", func(t *testing.T) {
		c, rec := prepareDeleteQuery("abc")
		mockCtxUser(c)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Delete(
			ctx, gomock.Any(), gomock.Any(),
		).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Delete(
			ctx, gomock.Any(), gomock.Any(),
		).Return(usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
