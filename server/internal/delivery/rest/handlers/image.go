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

type imageUseCase interface {
	Create(ctx context.Context, image *domain.Image, file *domain.FileNode) (*domain.Image, error)
	Delete(ctx context.Context, id int) error
	GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error)
}

type ImageHandlers struct {
	uc     imageUseCase
	logger logger.Logger
}

func NewImageHandlers(uc imageUseCase, logger logger.Logger) *ImageHandlers {
	return &ImageHandlers{uc: uc, logger: logger}
}

func (h *ImageHandlers) Upload() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, ok := c.Get("user").(*domain.User)
		if !ok || user == nil {
			h.logger.Errorf("Cannot get user from context, make sure to use OnlyAuth middleware first")
			return c.JSON(rest.NewInternalServerError().Response())
		}

		file, err := rest.ReadEchoImage(c, "file")
		if err != nil {
			if restErr, ok := err.(*rest.Error); ok {
				return c.JSON(restErr.Response())
			}

			h.logger.Errorf("Create.ReadEchoImage: %v", err)
			return c.JSON(rest.NewInternalServerError().Response())
		}

		img := &domain.Image{AuthorID: user.ID}
		createdImg, err := h.uc.Create(ctx, img, file)
		if err != nil {
			return h.parseAndWriteUseCaseErr(c, err, "Create")
		}

		return c.JSON(http.StatusCreated, createdImg)
	}
}

func (h *ImageHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, ok := c.Get("user").(*domain.User)
		if !ok || user == nil {
			h.logger.Errorf("Cannot get user from context, make sure to use OnlyAuth middleware first")
			return c.JSON(rest.NewInternalServerError().Response())
		}

		id, err := rest.IntParam(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		if err := h.uc.Delete(ctx, id); err != nil {
			return h.parseAndWriteUseCaseErr(c, err, "Delete")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *ImageHandlers) GetDetailed() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.IntParam(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		img, err := h.uc.GetDetailed(ctx, id)
		if err != nil {
			return h.parseAndWriteUseCaseErr(c, err, "GetDetailed")
		}

		return c.JSON(http.StatusOK, img)
	}
}

func (h *ImageHandlers) parseAndWriteUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrUnprocessableEntity):
		restErr = rest.NewBadRequestError("Image body has incorrect type")
		break
	case errors.Is(err, usecase.ErrNotFound):
		restErr = rest.NewNotFoundError("Image not found")
		break
	default:
		h.logger.Errorf("ImageUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
		break
	}

	return c.JSON(restErr.Response())
}
