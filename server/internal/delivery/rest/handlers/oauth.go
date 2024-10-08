package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type OAuthUseCase interface {
	Authenticate(ctx context.Context, code string, service domain.OAuthService) (*domain.UserWithToken, error)
	GetAuthURL(ctx context.Context, service domain.OAuthService) (string, error)
}

type OAuthHandlers struct {
	uc     OAuthUseCase
	cfg    *config.Cookie
	logger logger.Logger
}

func NewOAuthHandlers(uc OAuthUseCase, cfg *config.Cookie, logger logger.Logger) *OAuthHandlers {
	return &OAuthHandlers{uc: uc, cfg: cfg, logger: logger}
}

func (h *OAuthHandlers) Callback(service domain.OAuthService) echo.HandlerFunc {
	type callbackDTO struct {
		Code string `query:"code" validate:"required,gte=1"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		dto := new(callbackDTO)
		if err := rest.DecodeEchoBody(c, dto); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Query").Response())
		}

		if err := validator.ValidateStruct(ctx, dto); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Query").Response())
		}

		authUser, err := h.uc.Authenticate(ctx, dto.Code, service)
		if err != nil {
			switch err {
			case usecase.ErrInvalidCredentials:
				return c.JSON(rest.NewUnauthorizedError("Invalid credentials").Response())
			case usecase.ErrUnprocessable:
				return c.JSON(rest.NewBadRequestError("Invalid Request Query").Response())
			default:
				h.logger.Errorf("OAuthHandlers.CallbackGoogle: %v", err)
				return c.JSON(rest.NewInternalServerError().Response())
			}
		}

		h.storeToken(c, authUser.Token)
		// See https://stackoverflow.com/questions/42216700/how-can-i-redirect-after-oauth2-with-samesite-strict-and-still-get-my-cookies
		return c.HTML(
			http.StatusOK,
			fmt.Sprintf(`<html>
      <head>
      <meta http-equiv="refresh" content="0;URL='%s'"/>
      </head>
      </html>`, h.cfg.Origin),
		)
	}
}

func (h *OAuthHandlers) AuthRedirect(service domain.OAuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectURL, err := h.uc.GetAuthURL(context.Background(), service)

		if err != nil {
			h.logger.Errorf("OAuthHandlers.GetAuthURL: %v (service: %s)", err, service)
			return c.JSON(rest.NewInternalServerError().Response())
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

// TEMP: to cookie manager
func (h *OAuthHandlers) storeToken(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     h.cfg.Name,
		Value:    token,
		HttpOnly: h.cfg.HttpOnly,
		Path:     "/",
		MaxAge:   int((h.cfg.Expire * time.Second).Seconds()),
		SameSite: http.SameSiteStrictMode,
	})
}
