package middlewares

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
)

type AuthUseCase interface {
	Verify(ctx context.Context, token string) (*domain.User, error)
}

type AuthMiddlewares struct {
	authUC AuthUseCase
	logger logger.Logger
	cfg    *config.Cookie
}

type UserCtxKey struct{}

func NewAuthMiddlewares(authUC AuthUseCase, logger logger.Logger, cfg *config.Cookie) *AuthMiddlewares {
	return &AuthMiddlewares{authUC: authUC, logger: logger, cfg: cfg}
}

func (mw *AuthMiddlewares) OnlyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := rest.GetEchoRequestID(c)
		cookie, err := c.Cookie(mw.cfg.Name)
		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				reqID,
				err.Error(),
			)
			return c.JSON(rest.NewUnauthorizedError("No token provided").Response())
		}

		user, err := mw.authUC.Verify(rest.GetEchoRequestCtx(c), cookie.Value)
		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, CookieValue: %s, Error: %s",
				reqID,
				cookie.Value,
				err.Error(),
			)
			return c.JSON(rest.NewUnauthorizedError("Invalid token").Response())
		}

		c.Set("user", user)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (mw *AuthMiddlewares) OwnerOrAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		restErr := rest.NewForbiddenError("Only owner or admin can access this resource")

		user, ok := c.Get("user").(*domain.User)
		if !ok || user == nil {
			mw.logger.Errorf("Cannot get user from context, make sure to use OnlyAuth middleware first")
			return c.JSON(restErr.Response())
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			mw.logger.Errorf("OwnerOrAdminMiddleware: Error: %s", err.Error())
			return c.JSON(restErr.Response())
		}

		isAdmin := user.HasPermission(domain.PermissionsAdmin)
		if user.ID != id && !isAdmin {
			return c.JSON(restErr.Response())
		}

		return next(c)
	}
}
