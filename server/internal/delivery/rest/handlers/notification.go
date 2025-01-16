package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type notificationUseCase interface {
	GetNotifications(
		ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.Notification], error)
	WaitForNotification(ctx context.Context, userID domain.ID) (*domain.Notification, error)
	Notify(ctx context.Context, userID domain.ID, notif *domain.Notification) error
	GetStats(ctx context.Context, userID domain.ID) (*domain.NotificationStats, error)
}

type NotificationHandlers struct {
	uc     notificationUseCase
	logger logger.Logger
}

func NewNotificationHandlers(uc notificationUseCase, logger logger.Logger) *NotificationHandlers {
	return &NotificationHandlers{uc: uc, logger: logger}
}

func (h *NotificationHandlers) NotificationsSSE() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			notif, err := h.uc.WaitForNotification(ctx, user.ID)
			if err != nil {
				return c.JSON(rest.NewError(500, err.Error()).Response())
			}

			buf, err := json.Marshal(notif)
			if err != nil {
				return c.JSON(rest.NewError(500, err.Error()).Response())
			}

			event := rest.SSE{
				Data: buf,
			}
			if err := event.MarshalTo(w); err != nil {
				return c.JSON(rest.NewError(500, err.Error()).Response())
			}
			w.Flush()
		}
	}
}

func (h *NotificationHandlers) GetNotifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		notifs, err := h.uc.GetNotifications(ctx, user.ID, nil)
		if err != nil {
			return c.JSON(rest.NewError(http.StatusInternalServerError, err.Error()).Response())
		}

		return c.JSON(http.StatusOK, notifs)
	}
}

func (h *NotificationHandlers) GetStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		stats, err := h.uc.GetStats(ctx, user.ID)
		if err != nil {
			return c.JSON(rest.NewError(http.StatusInternalServerError, err.Error()).Response())
		}

		return c.JSON(http.StatusOK, stats)
	}
}

func (h *NotificationHandlers) Notify() echo.HandlerFunc {
	type notifyDTO struct {
		Title   string `json:"title"`
		Message string `json:"message"`
		UserID  string `json:"userId"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		dto := new(notifyDTO)
		if err := rest.DecodeEchoBody(c, dto); err != nil {
			h.logger.Errorf("Create.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("Create body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, dto); err != nil {
			return c.JSON(rest.NewBadRequestError("Create body has incorrect type").Response())
		}

		userID, err := domain.ParseID(dto.UserID)
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid user id").Response())
		}

		notif := &domain.Notification{Title: dto.Title, Message: dto.Message}
		err = h.uc.Notify(ctx, userID, notif)
		if err != nil {
			return c.JSON(rest.NewError(http.StatusInternalServerError, err.Error()).Response())
		}

		return c.NoContent(http.StatusCreated)
	}
}
