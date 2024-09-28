package handlers

import (
	"context"
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
}

type OAuthHandlers struct {
	uc     OAuthUseCase
	cfg    *config.Cookie
	logger logger.Logger
}

func NewOAuthHandlers(uc OAuthUseCase, cfg *config.Cookie, logger logger.Logger) *OAuthHandlers {
	return &OAuthHandlers{uc: uc, cfg: cfg, logger: logger}
}

func (h *OAuthHandlers) CallbackGoogle() echo.HandlerFunc {
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

		authUser, err := h.uc.Authenticate(ctx, dto.Code, domain.OAuthServiceGoogle)
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
		return c.JSON(http.StatusOK, authUser.User)
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
		SameSite: http.SameSiteNoneMode,
	})
}
