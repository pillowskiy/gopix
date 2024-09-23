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
	_, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	prepareCreateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/:image_id/commnets", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("image_id")
		c.SetParamValues(id)
		return c, rec
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

		mockCommentUC.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)
		c, rec := prepareCreateQuery("abc", bytes.NewBuffer(body))

		mockCommentUC.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareCreateQuery(itoaImageID, nil)

		mockCommentUC.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyTextField", func(t *testing.T) {
		invalidInput := CreateInput{
			Text: "",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))

		mockCommentUC.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

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

	t.Run("AlreadyExists", func(t *testing.T) {
		body, _ := json.Marshal(validCreateInput)

		c, rec := prepareCreateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Create(ctx, gomock.Any()).Return(nil, usecase.ErrAlreadyExists)

		assert.NoError(t, h.Create()(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
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

	imageID := handlersMock.DomainID()
	itoaImageID := imageID.String()

	type ImageCommentsQuery struct {
		Limit int    `query:"limit"`
		Page  int    `query:"page"`
		Sort  string `query:"sort"`
	}

	validImageCommentsQuery := &ImageCommentsQuery{
		Limit: 10,
		Page:  1,
		Sort:  "newest",
	}

	prepareGetByImageIDQuery := func(
		id string, query *ImageCommentsQuery,
	) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images/:image_id/comments", nil)
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
		c.SetParamNames("image_id")
		c.SetParamValues(id)

		return c, rec
	}

	t.Run("SuccessByImageID", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, validImageCommentsQuery)

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

		mockCommentUC.EXPECT().GetByImageID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery("abs", validImageCommentsQuery)

		mockCommentUC.EXPECT().GetByImageID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, &ImageCommentsQuery{
			Limit: 0,
			Page:  0,
			Sort:  "",
		})

		mockCommentUC.EXPECT().GetByImageID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, validImageCommentsQuery)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("server error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("IncorrectImageRef", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, validImageCommentsQuery)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(nil, usecase.ErrIncorrectImageRef)

		assert.NoError(t, h.GetByImageID()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Unprocessable", func(t *testing.T) {
		c, rec := prepareGetByImageIDQuery(itoaImageID, validImageCommentsQuery)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetByImageID(
			ctx, imageID, gomock.Any(), gomock.Any(),
		).Return(nil, usecase.ErrUnprocessable)

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
	_, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)
	e := echo.New()

	commentID := handlersMock.DomainID()
	itoaCommentID := commentID.String()

	prepareUpdateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/comments/:comment_id", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
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

		mockCommentUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery("abc", bytes.NewBuffer(body))
		mockCtxUser(c)

		mockCommentUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareUpdateQuery(itoaCommentID, nil)
		mockCtxUser(c)

		mockCommentUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EmptyTextField", func(t *testing.T) {
		invalidInput := UpdateInput{
			Text: "",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)

		mockCommentUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

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

	t.Run("NotFound", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().Update(ctx, commentID, gomock.Any(), gomock.Any()).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaCommentID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().Update(ctx, commentID, gomock.Any(), gomock.Any()).Return(nil, usecase.ErrForbidden)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestCommentHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	_, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := handlersMock.DomainID()
	itoaCommentID := commentID.String()

	prepareDeleteQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/:comment_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
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

		mockCommentUC.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidCommentID", func(t *testing.T) {
		c, rec := prepareDeleteQuery("abc")
		mockCtxUser(c)

		mockCommentUC.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

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

	t.Run("NotFound", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Delete(
			ctx, gomock.Any(), gomock.Any(),
		).Return(usecase.ErrNotFound)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaCommentID)
		mockCtxUser(c)
		ctx := rest.GetEchoRequestCtx(c)

		mockCommentUC.EXPECT().Delete(
			ctx, gomock.Any(), gomock.Any(),
		).Return(usecase.ErrForbidden)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestCommentHandlers_GetReplies(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := handlersMock.DomainID()
	itoaCommentID := commentID.String()

	prepareGetRepliesQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/comments/:comment_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
	}

	mockReplies := []domain.DetailedComment{
		{
			Comment: domain.Comment{ID: 1},
		},
	}

	t.Run("SuccessGetReplies", func(t *testing.T) {
		c, rec := prepareGetRepliesQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetReplies(
			ctx, commentID, &mockUser.ID,
		).Return(mockReplies, nil)

		assert.NoError(t, h.GetReplies()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new([]domain.DetailedComment)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, mockReplies, *actual)
	})

	t.Run("SuccessGetReplies_Guest", func(t *testing.T) {
		c, rec := prepareGetRepliesQuery(itoaCommentID)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetReplies(
			ctx, commentID, nil,
		).Return(mockReplies, nil)

		assert.NoError(t, h.GetReplies()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new([]domain.DetailedComment)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, mockReplies, *actual)
	})

	t.Run("IncorrectCommentID", func(t *testing.T) {
		c, rec := prepareGetRepliesQuery("")

		mockCommentUC.EXPECT().GetReplies(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Times(0)

		assert.NoError(t, h.GetReplies()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		c, rec := prepareGetRepliesQuery(itoaCommentID)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetReplies(
			ctx, commentID, nil,
		).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.GetReplies()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetRepliesQuery(itoaCommentID)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().GetReplies(
			ctx, commentID, nil,
		).Return(nil, errors.New("internal server error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetReplies()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommentHandlers_LikeComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := handlersMock.DomainID()
	itoaCommentID := commentID.String()

	prepareLikeCommentQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments/:comment_id/like", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessLikeComment", func(t *testing.T) {
		c, rec := prepareLikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().LikeComment(ctx, commentID, mockUser).Return(nil)

		assert.NoError(t, h.LikeComment()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// When like already exists, delivery should emulate it as success
	t.Run("SuccessLikeComment_AlreadyExists", func(t *testing.T) {
		c, rec := prepareLikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().LikeComment(ctx, commentID, mockUser).Return(usecase.ErrAlreadyExists)

		assert.NoError(t, h.LikeComment()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorectUserContext", func(t *testing.T) {
		c, rec := prepareLikeCommentQuery(itoaCommentID)

		mockCommentUC.EXPECT().LikeComment(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.LikeComment()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectCommentID", func(t *testing.T) {
		c, rec := prepareLikeCommentQuery("")

		mockCommentUC.EXPECT().LikeComment(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.LikeComment()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareLikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().LikeComment(
			ctx, commentID, mockUser,
		).Return(errors.New("internal server error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.LikeComment()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCommentHandlers_UnlikeComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentUC := handlersMock.NewMockCommentUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	mockUser, mockCtxUser := handlersMock.NewMockCtxUser()
	h := handlers.NewCommentHandlers(mockCommentUC, mockLog)

	e := echo.New()

	commentID := handlersMock.DomainID()
	itoaCommentID := commentID.String()

	prepareUnlikeCommentQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/:comment_id/like", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("comment_id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessUnlikeComment", func(t *testing.T) {
		c, rec := prepareUnlikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().UnlikeComment(ctx, commentID, mockUser).Return(nil)

		assert.NoError(t, h.UnlikeComment()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// When like not found, delivery should emulate it as success
	t.Run("SuccessUnlikeComment_NotFound", func(t *testing.T) {
		c, rec := prepareUnlikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().UnlikeComment(ctx, commentID, mockUser).Return(usecase.ErrNotFound)

		assert.NoError(t, h.UnlikeComment()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorectUserContext", func(t *testing.T) {
		c, rec := prepareUnlikeCommentQuery(itoaCommentID)

		mockCommentUC.EXPECT().UnlikeComment(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.UnlikeComment()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectCommentID", func(t *testing.T) {
		c, rec := prepareUnlikeCommentQuery("")

		mockCommentUC.EXPECT().UnlikeComment(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.UnlikeComment()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareUnlikeCommentQuery(itoaCommentID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockCommentUC.EXPECT().UnlikeComment(
			ctx, commentID, mockUser,
		).Return(errors.New("internal server error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.UnlikeComment()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
