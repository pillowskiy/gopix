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

func TestUserHanlers_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := handlersMock.NewMockuserUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewUserHandlers(mockUserUC, mockLog)

	e := echo.New()

	userID := 1
	itoaUserID := strconv.Itoa(userID)

	updatePath := func(id string) string {
		return fmt.Sprintf("/api/v1/users/%s", id)
	}

	prepareUpdateQuery := func(id string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPut, updatePath(id), body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("users/:id")
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	type UpdateInput struct {
		Username  string `json:"username"`
		AvatarURL string `json:"avatarUrl"`
	}

	validUpdateInput := UpdateInput{
		Username:  "test",
		AvatarURL: "https://test.com/test.png",
	}

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		updatedUser := &domain.User{
			Username:  validUpdateInput.Username,
			AvatarURL: validUpdateInput.AvatarURL,
		}

		mockUserUC.EXPECT().Update(ctx, userID, gomock.Any()).Return(updatedUser, nil)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.User)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, updatedUser, actual)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)
		c, rec := prepareUpdateQuery("abc", bytes.NewBuffer(body))

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareUpdateQuery(itoaUserID, nil)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidAvatarURL", func(t *testing.T) {
		invalidInput := UpdateInput{
			Username:  "username",
			AvatarURL: "a",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidUsername", func(t *testing.T) {
		invalidInput := UpdateInput{
			Username:  "a",
			AvatarURL: "https://test.com/test.png",
		}

		body, _ := json.Marshal(invalidInput)
		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockUserUC.EXPECT().Update(ctx, userID, gomock.Any()).Return(nil, errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockUserUC.EXPECT().Update(ctx, userID, gomock.Any()).Return(nil, usecase.ErrNotFound)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		body, _ := json.Marshal(validUpdateInput)

		c, rec := prepareUpdateQuery(itoaUserID, bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockUserUC.EXPECT().Update(ctx, userID, gomock.Any()).Return(nil, usecase.ErrAlreadyExists)

		assert.NoError(t, h.Update()(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestUserHanlers_Me(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	mockUserUC := handlersMock.NewMockuserUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewUserHandlers(mockUserUC, mockLog)

	e := echo.New()

	mePath := "/api/v1/users/@me"

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	prepareMeQuery := func() (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, mePath, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("users/@me")

		return c, rec
	}

	t.Run("SuccessMe", func(t *testing.T) {
		c, rec := prepareMeQuery()
		mockCtxUser(c)

		assert.NoError(t, h.Me()(c))
		assert.Equal(t, rec.Code, http.StatusOK)

		actual := new(domain.User)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, ctxUser, actual)
	})

	t.Run("InvalidUserCtx", func(t *testing.T) {
		c, rec := prepareMeQuery()

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.Me()(c))
		assert.Equal(t, rec.Code, http.StatusUnauthorized)
	})
}

func TestUserHanlers_OverwritePermissions(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type OPInput struct {
		Deny  domain.Permission `json:"deny"`
		Allow domain.Permission `json:"allow"`
	}

	validOPInput := OPInput{
		Deny:  1,
		Allow: 1024,
	}

	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}
	targetUserID := 2
	targetItoaUserID := strconv.Itoa(targetUserID)

	mockUserUC := handlersMock.NewMockuserUseCase(ctrl)
	mockLog := loggerMock.NewMockLogger(ctrl)
	h := handlers.NewUserHandlers(mockUserUC, mockLog)

	e := echo.New()

	overwritePermissionsPath := func(id string) string {
		return fmt.Sprintf("/api/v1/users/%s/permissions", id)
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	prepareOverwritePermissionQuery := func(
		id string, body io.Reader,
	) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPut, overwritePermissionsPath(id), body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("users/:id/permissions")
		c.SetParamNames("id")
		c.SetParamValues(id)
		return c, rec
	}

	t.Run("SuccessOverwritePermissions", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery(targetItoaUserID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockUserUC.EXPECT().OverwritePermissions(
			ctx, targetUserID, validOPInput.Deny, validOPInput.Allow,
		).Return(nil)

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, rec.Code, http.StatusOK)
	})

	t.Run("InvalidUserCtx", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery(targetItoaUserID, bytes.NewBuffer(body))

		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareOverwritePermissionQuery(targetItoaUserID, nil)
		mockCtxUser(c)

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("IncorrectUserID", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery("abc", bytes.NewBuffer(body))

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("DenySelfAction", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery(strconv.Itoa(ctxUser.ID), bytes.NewBuffer(body))
		mockCtxUser(c)

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery(targetItoaUserID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockUserUC.EXPECT().OverwritePermissions(
			ctx, targetUserID, validOPInput.Deny, validOPInput.Allow,
		).Return(errors.New("internal error"))
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any())

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		body, _ := json.Marshal(validOPInput)
		c, rec := prepareOverwritePermissionQuery(targetItoaUserID, bytes.NewBuffer(body))
		mockCtxUser(c)

		ctx := rest.GetEchoRequestCtx(c)
		mockUserUC.EXPECT().OverwritePermissions(
			ctx, targetUserID, validOPInput.Deny, validOPInput.Allow,
		).Return(usecase.ErrNotFound)

		assert.NoError(t, h.OverwritePermissions()(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
