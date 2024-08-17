package middlewares

import (
	"context"
	"errors"
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

		user, err := mw.getUserFromCtx(c)
		if err != nil {
			return c.JSON(restErr.Response())
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(restErr.Response())
		}

		isAdmin := user.HasPermission(domain.PermissionsAdmin)
		if user.ID != id && !isAdmin {
			return c.JSON(restErr.Response())
		}

		return next(c)
	}
}

func (mw *AuthMiddlewares) WithSomePermission(permissions ...domain.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			restErr := rest.NewForbiddenError("You don't have permission to access this resource")
			user, err := mw.getUserFromCtx(c)
			if err != nil {
				return c.JSON(restErr.Response())
			}

			if user.HasPermission(domain.PermissionsAdmin) {
				return next(c)
			}

			for _, perm := range permissions {
				if user.HasPermission(perm) {
					return next(c)
				}
			}

			return c.JSON(restErr.Response())
		}
	}
}

func (mw *AuthMiddlewares) WithEveryPermission(permissions ...domain.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			restErr := rest.NewForbiddenError("You don't have permission to access this resource")
			user, err := mw.getUserFromCtx(c)
			if err != nil {
				return c.JSON(restErr.Response())
			}

			if user.HasPermission(domain.PermissionsAdmin) {
				return next(c)
			}

			for _, perm := range permissions {
				if !user.HasPermission(perm) {
					return c.JSON(restErr.Response())
				}
			}

			return next(c)
		}
	}
}

func (mw *AuthMiddlewares) OnlyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return mw.WithSomePermission(domain.PermissionsAdmin)(next)
}

func (mw *AuthMiddlewares) getUserFromCtx(c echo.Context) (*domain.User, error) {
	user, ok := c.Get("user").(*domain.User)
	if !ok || user == nil {
		mw.logger.Errorf("Cannot get user from context, make sure to use OnlyAuth middleware first")
		return nil, errors.New("cannot get user from context")
	}
	return user, nil
}