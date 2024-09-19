package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
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

func TestImageHandlers_Upload(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)

	type mockImage struct {
		Name string
		Data []byte
	}

	mockImages := map[string]mockImage{
		"image/jpeg": {
			Name: "test.jpg",
			Data: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01},
		},
		"image/png": {
			Name: "test.png",
			Data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		},
		"image/gif": {
			Name: "test.gif",
			Data: []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61},
		},
		"image/webp": {
			Name: "test.webp",
			Data: []byte{0x52, 0x49, 0x46, 0x46, 0x20, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x20},
		},
		"image/x-icon": {
			Name: "test.ico",
			Data: []byte{0x00, 0x00, 0x01, 0x00, 0x01, 0x00},
		},
		"image/bmp": {
			Name: "test.bmp",
			Data: []byte{0x42, 0x4D},
		},
	}

	_, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	prepareUploadQuery := func(
		file mockImage, mime string, field string,
	) (echo.Context, *httptest.ResponseRecorder) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		imageHeader := make(textproto.MIMEHeader)

		dataHeader := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, file.Name)
		imageHeader.Set("Content-Disposition", dataHeader)
		imageHeader.Set("Content-Type", mime)

		part, _ := writer.CreatePart(imageHeader)
		part.Write(file.Data)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/images", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	for contentType, file := range mockImages {
		t.Run(fmt.Sprintf("SuccessUpload_%s", contentType), func(t *testing.T) {
			c, rec := prepareUploadQuery(file, contentType, "file")
			mockCtxUser(c)

			img := &domain.Image{Path: "anypath.png"}
			ctx := rest.GetEchoRequestCtx(c)
			mockImageUC.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(img, nil)

			assert.NoError(t, h.Upload()(c))
			assert.Equal(t, http.StatusCreated, rec.Code)

			actual := new(domain.Image)
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
			assert.Equal(t, img, actual)
		})
	}

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareUploadQuery(mockImages["image/jpeg"], "image/jpeg", "file")

		mockImageUC.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Upload()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("IncorrectInput", func(t *testing.T) {
		invalidInput := mockImage{
			Name: "test.txt",
			Data: []byte{0x00, 0x00, 0x01, 0x00, 0x01, 0x00},
		}
		c, rec := prepareUploadQuery(invalidInput, "text/plain", "file")
		mockCtxUser(c)

		mockImageUC.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Upload()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectFormField", func(t *testing.T) {
		ct := "image/jpeg"
		c, rec := prepareUploadQuery(mockImages[ct], ct, "wrong")
		mockCtxUser(c)

		mockImageUC.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Upload()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectContentType", func(t *testing.T) {
		invalidInput := mockImage{
			Name: "test.svg",
			Data: []byte{
				0x3C, 0x3F, 0x78, 0x6D, 0x6C, 0x20, 0x76, 0x65, 0x72,
				0x73, 0x69, 0x6F, 0x6E, 0x3D, 0x22, 0x31, 0x2E, 0x30,
				0x22, 0x3F, 0x3E, 0x3C, 0x73, 0x76, 0x67, 0x20,
			},
		}
		ct := "image/svg+xml"
		c, rec := prepareUploadQuery(invalidInput, ct, "file")
		mockCtxUser(c)

		mockImageUC.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Upload()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareUploadQuery(mockImages["image/jpeg"], "image/jpeg", "file")
		mockCtxUser(c)

		mockImageUC.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Upload()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestImageHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareDeleteQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/images/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessDelete", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Delete(ctx, imageID, ctxUser).Return(nil)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidUserContext", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaImageID)

		mockImageUC.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		c, rec := prepareDeleteQuery("abs")
		mockCtxUser(c)

		mockImageUC.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Delete(ctx, imageID, ctxUser).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Delete(ctx, imageID, ctxUser).Return(usecase.ErrNotFound)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		c, rec := prepareDeleteQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Delete(ctx, imageID, ctxUser).Return(usecase.ErrForbidden)

		assert.NoError(t, h.Delete()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestImageHandlers_GetDetailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareGetDetailedQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	img := &domain.DetailedImage{
		ImageWithAuthor: domain.ImageWithAuthor{
			Image: domain.Image{ID: 1, Path: "path.png"},
		},
	}

	t.Run("SuccessGetDetailed", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery(itoaImageID)
		ctx := rest.GetEchoRequestCtx(c)

		mockImageUC.EXPECT().GetDetailed(ctx, imageID).Return(img, nil)
		mockImageUC.EXPECT().AddView(ctx, gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := &domain.DetailedImage{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, img, actual)
	})

	t.Run("AddImageView", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().GetDetailed(ctx, imageID).Return(img, nil)

		mockImageUC.EXPECT().AddView(ctx, imageID, &ctxUser.ID)

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("AddImageGuestView", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().GetDetailed(ctx, imageID).Return(img, nil)

		mockImageUC.EXPECT().AddView(ctx, imageID, nil)

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("AddViewErr", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().GetDetailed(ctx, imageID).Return(img, nil)

		mockImageUC.EXPECT().AddView(ctx, gomock.Any(), gomock.Any()).Return(errors.New("any error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery("abs")

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		c, rec := prepareGetDetailedQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().GetDetailed(ctx, imageID).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.GetDetailed()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestImageHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareUpdateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/images/:id", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	type UpdateInput struct {
		Title       string
		Description string
		AccessLevel string
	}

	validUpdateInput := &UpdateInput{
		Title:       "title",
		Description: "description",
		AccessLevel: "public",
	}

	t.Run("SuccessUpdate", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		img := &domain.Image{
			Title:       validUpdateInput.Title,
			Description: validUpdateInput.Description,
			AccessLevel: domain.ImageAccessLevel(validUpdateInput.AccessLevel),
		}

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Update(ctx, imageID, img, ctxUser).Return(img, nil)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := &domain.Image{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, img, actual)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery("abs", bytes.NewBuffer(body))
		mockCtxUser(c)

		mockImageUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareUpdateQuery(itoaImageID, nil)
		mockCtxUser(c)

		mockImageUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaImageID, bytes.NewBuffer(body))

		mockImageUC.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Update(ctx, imageID, gomock.Any(), ctxUser).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Update(ctx, imageID, gomock.Any(), ctxUser).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery(itoaImageID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Update(ctx, imageID, gomock.Any(), ctxUser).Return(nil, usecase.ErrForbidden)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestImageHandlers_GetStates(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareGetStatesQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images/:id/states", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	states := &domain.ImageStates{Viewed: true, Liked: true}

	t.Run("SuccessGetStates", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().States(ctx, imageID, ctxUser.ID).Return(states, nil)

		assert.NoError(t, h.GetStates()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.ImageStates)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, states, actual)
	})

	t.Run("IncorrectImageID", func(t *testing.T) {
		c, rec := prepareGetStatesQuery("abs")
		mockCtxUser(c)

		mockImageUC.EXPECT().States(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetStates()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().States(ctx, gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetStates()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().States(ctx, imageID, ctxUser.ID).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetStates()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestImageHandlers_GetDiscover(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewImageHandlers(mockImageUC, mockLog)

	e := echo.New()

	type DiscoverInput struct {
		Page    int
		PerPage int
		Sort    domain.ImageSortMethod
	}

	prepareGetStatesQuery := func(query *DiscoverInput) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		if query != nil {
			q := req.URL.Query()
			q.Add("page", strconv.Itoa(query.Page))
			q.Add("limit", strconv.Itoa(query.PerPage))
			q.Add("sort", string(query.Sort))
			req.URL.RawQuery = q.Encode()
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	pagInput := &domain.PaginationInput{
		Page:    1,
		PerPage: 10,
	}

	validDiscoverInput := &DiscoverInput{
		Page:    pagInput.Page,
		PerPage: pagInput.PerPage,
		Sort:    domain.ImagePopularSort,
	}

	pag := &domain.Pagination[domain.ImageWithAuthor]{
		PaginationInput: domain.PaginationInput{
			Page:    1,
			PerPage: 10,
		},
		Items: []domain.ImageWithAuthor{
			{
				Image: domain.Image{
					ID:   1,
					Path: "path.png",
				},
			},
		},
		Total: 100,
	}

	t.Run("SuccessGetDiscover", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(validDiscoverInput)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Discover(ctx, pagInput, validDiscoverInput.Sort).Return(pag, nil)

		assert.NoError(t, h.GetDiscover()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.Pagination[domain.ImageWithAuthor])
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, pag, actual)
	})

	t.Run("IncorrectInput", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(nil)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Discover(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.GetDiscover()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Unprocessable", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(validDiscoverInput)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Discover(ctx, gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUnprocessable)

		assert.NoError(t, h.GetDiscover()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareGetStatesQuery(validDiscoverInput)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().Discover(ctx, gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.GetDiscover()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestImageHandlers_AddLike(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareAddLikeQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/:id/like", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessAddLike", func(t *testing.T) {
		c, rec := prepareAddLikeQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().AddLike(ctx, imageID, ctxUser.ID).Return(nil)

		assert.NoError(t, h.AddLike()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		c, rec := prepareAddLikeQuery("abs")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().AddLike(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.AddLike()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareAddLikeQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().AddLike(ctx, gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.AddLike()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareAddLikeQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().AddLike(ctx, imageID, ctxUser.ID).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.AddLike()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestImageHandlers_RemoveLike(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageUC := handlersMock.NewMockimageUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	ctxUser, mockCtxUser := handlersMock.NewMockCtxUser()

	h := handlers.NewImageHandlers(mockImageUC, mockLog)
	e := echo.New()

	imageID := 1
	itoaImageID := strconv.Itoa(imageID)

	prepareRemoveLikeQuery := func(id string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/images/:id/like", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessRemoveLike", func(t *testing.T) {
		c, rec := prepareRemoveLikeQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().RemoveLike(ctx, imageID, ctxUser.ID).Return(nil)

		assert.NoError(t, h.RemoveLike()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LikeNotExists", func(t *testing.T) {
		c, rec := prepareRemoveLikeQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().RemoveLike(ctx, imageID, ctxUser.ID).Return(usecase.ErrUnprocessable)

		assert.NoError(t, h.RemoveLike()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidImageID", func(t *testing.T) {
		c, rec := prepareRemoveLikeQuery("abs")
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().RemoveLike(ctx, gomock.Any(), gomock.Any()).Times(0)

		assert.NoError(t, h.RemoveLike()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("IncorrectUserContext", func(t *testing.T) {
		c, rec := prepareRemoveLikeQuery(itoaImageID)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().RemoveLike(ctx, gomock.Any(), gomock.Any()).Times(0)
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.RemoveLike()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		c, rec := prepareRemoveLikeQuery(itoaImageID)
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockImageUC.EXPECT().RemoveLike(ctx, imageID, ctxUser.ID).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.RemoveLike()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
