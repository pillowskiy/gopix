package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
)

type SubscriptionUseCase interface {
	Follow(ctx context.Context, userID domain.ID, executor *domain.User) error
	Unfollow(ctx context.Context, userID domain.ID, executor *domain.User) error
}

type SubscriptionHandlers struct {
	uc     SubscriptionUseCase
	logger logger.Logger
}

func NewSubscriptionHandlers(uc SubscriptionUseCase, logger logger.Logger) *SubscriptionHandlers {
	return &SubscriptionHandlers{uc: uc, logger: logger}
}

func (h *SubscriptionHandlers) Follow() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		id, err := rest.PipeDomainIdentifier(c, "user_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("ID has incorrect type").Response())
		}

		if err := h.uc.Follow(ctx, id, user); err != nil {
			switch {
			case errors.Is(err, usecase.ErrAlreadyExists):
				return c.NoContent(http.StatusOK)
			default:
				return h.responseWithUseCaseErr(c, err, "Follow")
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *SubscriptionHandlers) Unfollow() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		id, err := rest.PipeDomainIdentifier(c, "user_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("ID has incorrect type").Response())
		}

		if err := h.uc.Unfollow(ctx, id, user); err != nil {
			switch {
			case errors.Is(err, usecase.ErrNotFound):
				return c.NoContent(http.StatusOK)
			default:
				return h.responseWithUseCaseErr(c, err, "Unfollow")
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *SubscriptionHandlers) responseWithUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrIncorrectUserRef):
		restErr = rest.NewBadRequestError("Incorrect user reference provided")
	default:
		h.logger.Errorf("SubscriptionUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
	}

	return c.JSON(restErr.Response())
}
