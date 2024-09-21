package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	handlersMock "github.com/pillowskiy/gopix/internal/delivery/rest/handlers/mock"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	loggerMock "github.com/pillowskiy/gopix/pkg/logger/mock"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func NewMockCookieConfig() *config.Cookie {
	return &config.Cookie{
		Name:     "session",
		Expire:   time.Duration(0),
		Secure:   true,
		HttpOnly: true,
	}
}

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAuthUC := handlersMock.NewMockauthUseCase(ctrl)
	mockCfg := NewMockCookieConfig()

	h := handlers.NewAuthHandlers(mockAuthUC, mockLog, mockCfg)

	prepareRegisterQuery := func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	type RegisterInput struct {
		Username string
		Email    string
		Password string
	}

	validRegisterInput := RegisterInput{
		Username: "username",
		Email:    "username@gmail.com",
		Password: "password",
	}

	mockUser := &domain.User{
		ID:          handlersMock.DomainID(),
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	token := "token"

	t.Run("SuccessRegister", func(t *testing.T) {
		body, _ := json.Marshal(validRegisterInput)

		c, rec := prepareRegisterQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		userWithToken := &domain.UserWithToken{
			User:  mockUser,
			Token: token,
		}

		cookieToken := http.Cookie{
			Name:     mockCfg.Name,
			Value:    token,
			HttpOnly: mockCfg.HttpOnly,
			Path:     "/",
			MaxAge:   int((mockCfg.Expire * time.Second).Seconds()),
			SameSite: http.SameSiteNoneMode,
		}

		mockAuthUC.EXPECT().Register(ctx, gomock.Any()).Return(userWithToken, nil)

		assert.NoError(t, h.Register()(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		actual := new(domain.User)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, userWithToken.User, actual)

		assert.Equal(t, rec.Header().Get("Set-Cookie"), cookieToken.String())
		assert.Empty(t, actual.PasswordHash)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareRegisterQuery(nil)
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Register(ctx, gomock.Any()).Times(0)

		assert.NoError(t, h.Register()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ConflictCredentials", func(t *testing.T) {
		body, _ := json.Marshal(validRegisterInput)

		c, rec := prepareRegisterQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Register(ctx, gomock.Any()).Return(nil, usecase.ErrAlreadyExists)

		assert.NoError(t, h.Register()(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validRegisterInput)

		c, rec := prepareRegisterQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Register(ctx, gomock.Any()).Return(nil, errors.New("internal error"))

		assert.NoError(t, h.Register()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAuthHandlers_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAuthUC := handlersMock.NewMockauthUseCase(ctrl)
	mockCfg := NewMockCookieConfig()
	h := handlers.NewAuthHandlers(mockAuthUC, mockLog, mockCfg)

	prepareLoginQuery := func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	type LoginInput struct {
		Initials string
		Password string
	}

	validLoginInput := LoginInput{
		Initials: "username",
		Password: "password",
	}

	mockUser := &domain.User{
		ID:          handlersMock.DomainID(),
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	token := "token"

	t.Run("SuccessLogin", func(t *testing.T) {
		body, _ := json.Marshal(validLoginInput)

		c, rec := prepareLoginQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		userWithToken := &domain.UserWithToken{
			User:  mockUser,
			Token: token,
		}

		cookieToken := http.Cookie{
			Name:     mockCfg.Name,
			Value:    token,
			HttpOnly: mockCfg.HttpOnly,
			Path:     "/",
			MaxAge:   int((mockCfg.Expire * time.Second).Seconds()),
			SameSite: http.SameSiteNoneMode,
		}

		mockAuthUC.EXPECT().Login(ctx, gomock.Any()).Return(userWithToken, nil)

		assert.NoError(t, h.Login()(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := new(domain.User)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), actual))
		assert.Equal(t, userWithToken.User, actual)

		assert.Equal(t, rec.Header().Get("Set-Cookie"), cookieToken.String())
		assert.Empty(t, actual.PasswordHash)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		c, rec := prepareLoginQuery(nil)
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Login(ctx, gomock.Any()).Times(0)

		assert.NoError(t, h.Login()(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		body, _ := json.Marshal(validLoginInput)

		c, rec := prepareLoginQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Login(ctx, gomock.Any()).Return(nil, usecase.ErrInvalidCredentials)

		assert.NoError(t, h.Login()(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		body, _ := json.Marshal(validLoginInput)

		c, rec := prepareLoginQuery(bytes.NewBuffer(body))
		ctx := rest.GetEchoRequestCtx(c)

		mockAuthUC.EXPECT().Login(ctx, gomock.Any()).Return(nil, errors.New("internal error"))

		assert.NoError(t, h.Login()(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAuthHandlers_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	mockLog := loggerMock.NewMockLogger(ctrl)
	mockAuthUC := handlersMock.NewMockauthUseCase(ctrl)
	mockCfg := NewMockCookieConfig()
	h := handlers.NewAuthHandlers(mockAuthUC, mockLog, mockCfg)

	prepareLogoutQuery := func() (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/logout", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return c, rec
	}

	t.Run("SuccessLogout", func(t *testing.T) {
		c, rec := prepareLogoutQuery()

		cookieToken := &http.Cookie{
			Name:     mockCfg.Name,
			Value:    "",
			Path:     "/",
			HttpOnly: mockCfg.HttpOnly,
			MaxAge:   -1,
			SameSite: http.SameSiteStrictMode,
		}

		assert.NoError(t, h.Logout()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Header().Get("Set-Cookie"), cookieToken.String())
	})
}
