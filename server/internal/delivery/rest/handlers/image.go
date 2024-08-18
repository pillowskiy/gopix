package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
)

type imageUseCase interface {
	Create(ctx context.Context, image *domain.Image, file *domain.FileNode) (*domain.Image, error)
}

type ImageHandlers struct {
	uc     imageUseCase
	logger logger.Logger
}

func NewImageHandlers(uc imageUseCase, logger logger.Logger) *ImageHandlers {
	return &ImageHandlers{uc: uc, logger: logger}
}

func (h *ImageHandlers) Create() echo.HandlerFunc {
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

		img := &domain.Image{
			AuthorID: user.ID,
		}
		createdImg, err := h.uc.Create(ctx, img, file)
		if err != nil {
			if err == usecase.ErrUnprocessableEntity {
				return c.JSON(rest.NewBadRequestError("Image body has incorrect type").Response())
			}

			h.logger.Errorf("ImageUseCase.Create: %v", err)
			return c.JSON(rest.NewInternalServerError().Response())
		}

		return c.JSON(http.StatusCreated, createdImg)
	}
}
