package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type userUseCase interface {
	Update(ctx context.Context, id int, user *domain.User) (*domain.User, error)
}

type UserHandlers struct {
	uc     userUseCase
	logger logger.Logger
}

func NewUserHandlers(uc userUseCase, logger logger.Logger) *UserHandlers {
	return &UserHandlers{uc: uc, logger: logger}
}

func (h *UserHandlers) Update() echo.HandlerFunc {
	type updateDTO struct {
		Username  string `json:"username" validate:"gte=2,lte=60"`
		AvatarURL string `json:"avatarURL" validate:"lte=256"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Errorf("Update.ParseID: %v", err)
			return c.JSON(rest.NewBadRequestError("ID has incorrect type").Response())
		}

		upd := new(updateDTO)
		if err := rest.DecodeEchoBody(c, upd); err != nil {
			h.logger.Errorf("Update.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, upd); err != nil {
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		user := &domain.User{
			Username:  upd.Username,
			AvatarURL: upd.AvatarURL,
		}

		updatedUser, err := h.uc.Update(ctx, id, user)
		if err != nil {
			switch err {
			case usecase.ErrNotFound:
				return c.JSON(rest.NewNotFoundError("User not found").Response())
			case usecase.ErrAlreadyExists:
				return c.JSON(rest.NewConflictError("That username is already taken").Response())
			default:
				h.logger.Errorf("User.Update: %v", err)
				return c.JSON(rest.NewInternalServerError().Response())
			}
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *UserHandlers) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*domain.User)
		if !ok || user == nil {
			h.logger.Errorf("Cannot get user from context, make sure to use OnlyAuth middleware first")
			return c.JSON(rest.NewInternalServerError().Response())
		}
		return c.JSON(http.StatusOK, user)
	}
}
