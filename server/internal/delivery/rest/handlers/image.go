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
	"github.com/pillowskiy/gopix/pkg/validator"
)

type imageUseCase interface {
	Create(ctx context.Context, image *domain.Image, file *domain.File) (*domain.Image, error)
	Delete(ctx context.Context, id domain.ID, executor *domain.User) error
	Similar(ctx context.Context, id domain.ID) ([]domain.ImageWithMeta, error)
	GetDetailed(ctx context.Context, id domain.ID) (*domain.DetailedImage, error)
	Update(ctx context.Context, id domain.ID, image *domain.Image, executor *domain.User) (*domain.Image, error)
	AddView(ctx context.Context, imageID domain.ID, userID *domain.ID) error
	Discover(
		ctx context.Context, pagInput *domain.PaginationInput, sort domain.ImageSortMethod,
	) (*domain.Pagination[domain.ImageWithMeta], error)
	Favorites(
		ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.ImageWithMeta], error)

	States(ctx context.Context, imageID domain.ID, userID domain.ID) (*domain.ImageStates, error)
	AddLike(ctx context.Context, imageID domain.ID, userID domain.ID) error
	RemoveLike(ctx context.Context, imageID domain.ID, userID domain.ID) error
}

type ImageHandlers struct {
	uc     imageUseCase
	logger logger.Logger
}

func NewImageHandlers(uc imageUseCase, logger logger.Logger) *ImageHandlers {
	return &ImageHandlers{uc: uc, logger: logger}
}

// TEMP: This handler seems too complex for delivery layer, maybe it should be separated
func (h *ImageHandlers) Upload() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("Upload.GetContextUser: %v", err)
			return err
		}

		fileHeader, err := rest.ReadEchoImage(c, "file")
		if err != nil {
			if restErr, ok := err.(*rest.Error); ok {
				return c.JSON(restErr.Response())
			}

			h.logger.Errorf("Upload.ReadEchoImage: %v", err)
			return c.JSON(rest.NewInternalServerError().Response())
		}

		file, err := fileHeader.Open()
		if err != nil {
			h.logger.Errorf("Upload.Open: %v", err)
			return c.JSON(rest.NewInternalServerError().Response())
		}
		defer file.Close()

		dFile := &domain.File{
			Reader: file,
			Size:   fileHeader.Size,
		}

		img := &domain.Image{AuthorID: user.ID}
		createdImg, err := h.uc.Create(ctx, img, dFile)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Create")
		}

		return c.JSON(http.StatusCreated, createdImg)
	}
}

func (h *ImageHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("Delete.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.Delete(ctx, id, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "Delete")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *ImageHandlers) Similar() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		images, err := h.uc.Similar(ctx, id)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Similar")
		}

		return c.JSON(http.StatusOK, images)
	}
}

func (h *ImageHandlers) GetDetailed() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		img, err := h.uc.GetDetailed(ctx, imageID)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "GetDetailed")
		}

		var userID *domain.ID
		if user, ok := c.Get("user").(*domain.User); ok {
			userID = &user.ID
		}

		if err := h.uc.AddView(ctx, imageID, userID); err != nil {
			h.logger.Errorf("imagesHandlers.UseCase.AddView: %v", err)
		}

		return c.JSON(http.StatusOK, img)
	}
}

func (h *ImageHandlers) Update() echo.HandlerFunc {
	type updateDTO struct {
		Title       string `json:"title" validate:"lte=256"`
		Description string `json:"description" validate:"lte=1024"`
		AccessLevel string `json:"accessLevel" validate:"oneof=link private public"`
	}
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		dto := new(updateDTO)
		if err := rest.DecodeEchoBody(c, dto); err != nil {
			h.logger.Errorf("Update.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, dto); err != nil {
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("Update.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		image := &domain.Image{
			Title:       dto.Title,
			Description: dto.Description,
			AccessLevel: domain.ImageAccessLevel(dto.AccessLevel),
		}

		img, err := h.uc.Update(ctx, id, image, user)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Update")
		}

		return c.JSON(http.StatusOK, img)
	}
}

func (h *ImageHandlers) GetStates() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetStates.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		states, err := h.uc.States(ctx, id, user.ID)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "GetStates")
		}

		return c.JSON(http.StatusOK, states)
	}
}

func (h *ImageHandlers) GetDiscover() echo.HandlerFunc {
	type discoverQuery struct {
		Limit int    `query:"limit" validate:"required,gte=1,lte=100"`
		Page  int    `query:"page" validate:"required,gte=1"`
		Sort  string `query:"sort" validate:"oneof=newest oldest popular mostViewed"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		query := new(discoverQuery)
		if err := rest.DecodeEchoBody(c, query); err != nil {
			h.logger.Errorf("GetDiscover.DecodeQuery: %v", err)
			return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, query); err != nil {
			return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
		}

		pagInput := &domain.PaginationInput{Page: query.Page, PerPage: query.Limit}
		images, err := h.uc.Discover(ctx, pagInput, domain.ImageSortMethod(query.Sort))
		if err != nil {
			if errors.Is(err, usecase.ErrUnprocessable) {
				return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
			}
			return h.responseWithUseCaseErr(c, err, "GetDiscover")
		}

		return c.JSON(http.StatusOK, images)
	}
}

func (h *ImageHandlers) Favorites() echo.HandlerFunc {
	type discoverQuery struct {
		Limit int    `query:"limit" validate:"required,gte=1,lte=100"`
		Page  int    `query:"page" validate:"required,gte=1"`
		Sort  string `query:"sort" validate:"oneof=newest oldest popular mostViewed"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		userId, err := rest.PipeDomainIdentifier(c, "user_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid user ID").Response())
		}

		query := new(discoverQuery)
		if err := rest.DecodeEchoBody(c, query); err != nil {
			h.logger.Errorf("GetDiscover.DecodeQuery: %v", err)
			return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, query); err != nil {
			return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
		}

		pagInput := &domain.PaginationInput{Page: query.Page, PerPage: query.Limit}
		images, err := h.uc.Favorites(ctx, userId, pagInput)
		if err != nil {
			if errors.Is(err, usecase.ErrUnprocessable) {
				return c.JSON(rest.NewBadRequestError("Discover query has incorrect type").Response())
			}
			return h.responseWithUseCaseErr(c, err, "GetFavorites")
		}

		return c.JSON(http.StatusOK, images)
	}
}

func (h *ImageHandlers) AddLike() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("AddLike.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.AddLike(ctx, id, user.ID); err != nil {
			return h.responseWithUseCaseErr(c, err, "Like")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *ImageHandlers) RemoveLike() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		id, err := rest.PipeDomainIdentifier(c, "id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("RemoveLike.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.RemoveLike(ctx, id, user.ID); err != nil {
			return h.responseWithUseCaseErr(c, err, "Dislike")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *ImageHandlers) responseWithUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrForbidden):
		restErr = rest.NewForbiddenError("You don't have permissions to perform this action")
	case errors.Is(err, usecase.ErrUnprocessable):
		restErr = rest.NewBadRequestError("Image cannot be processed because it may conflict")
	case errors.Is(err, usecase.ErrNotFound):
		restErr = rest.NewNotFoundError("Image not found")
	default:
		h.logger.Errorf("ImageUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
	}

	return c.JSON(restErr.Response())
}
