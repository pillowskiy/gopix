package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type userUseCase interface {
	Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error)
	OverwritePermissions(
		ctx context.Context, id domain.ID, deny domain.Permission, allow domain.Permission,
	) error
	GetDetailed(ctx context.Context, username string, executorID *domain.ID) (*domain.DetailedUser, error)
}

type UserHandlers struct {
	uc     userUseCase
	logger logger.Logger
}

func NewUserHandlers(uc userUseCase, logger logger.Logger) *UserHandlers {
	return &UserHandlers{uc: uc, logger: logger}
}

func (h *UserHandlers) GetDetailed() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		username := c.Param("username")
		if len(username) < 2 {
			return c.JSON(rest.NewBadRequestError("Username is invalid").Response())
		}

		var executorID *domain.ID
		executor, err := GetContextUser(c)
		if err == nil && executor != nil {
			executorID = &executor.ID
		}

		user, err := h.uc.GetDetailed(ctx, username, executorID)
		if err != nil {
			switch err {
			case usecase.ErrNotFound:
				return c.JSON(rest.NewNotFoundError("User not found").Response())
			default:
				h.logger.Errorf("User.Update: %v", err)
				return c.JSON(rest.NewInternalServerError().Response())
			}
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *UserHandlers) Update() echo.HandlerFunc {
	// TODO: To form data & upload
	type updateDTO struct {
		Username  string `json:"username" validate:"gte=2,lte=60"`
		AvatarURL string `json:"avatarURL" validate:"lte=256,http_url"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
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

func (h *UserHandlers) OverwritePermissions() echo.HandlerFunc {
	type overwriteDTO struct {
		Deny  domain.Permission `json:"deny" validate:"min=1,number"`
		Allow domain.Permission `json:"allow" validate:"min=1,number"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("ID has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("OverwritePermissions.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if user.ID == id {
			return c.JSON(rest.NewBadRequestError("You can't overwrite your own permissions").Response())
		}

		ovr := new(overwriteDTO)
		if err := rest.DecodeEchoBody(c, ovr); err != nil {
			h.logger.Errorf("OverwritePermissions.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("OverwritePermissions body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, ovr); err != nil {
			return c.JSON(rest.NewBadRequestError("OverwritePermissions body has incorrect type").Response())
		}

		if ovr.Deny <= 0 && ovr.Allow <= 0 {
			return c.JSON(rest.NewBadRequestError("OverwritePermissions body has incorrect type").Response())
		}

		err = h.uc.OverwritePermissions(ctx, id, ovr.Deny, ovr.Allow)
		if err != nil {
			switch err {
			case usecase.ErrNotFound:
				return c.JSON(rest.NewNotFoundError("User not found").Response())
			default:
				h.logger.Errorf("User.OverwritePermissions: %v", err)
				return c.JSON(rest.NewInternalServerError().Response())
			}
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *UserHandlers) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("UserHandlers.Me: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		return c.JSON(http.StatusOK, user)
	}
}
