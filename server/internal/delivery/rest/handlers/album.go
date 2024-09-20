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

type albumUseCase interface {
	Create(ctx context.Context, album *domain.Album) (*domain.Album, error)
	GetByAuthorID(ctx context.Context, authorID domain.ID) ([]domain.Album, error)
	GetAlbumImages(
		ctx context.Context, albumID domain.ID, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.ImageWithAuthor], error)
	Delete(ctx context.Context, albumID domain.ID, executor *domain.User) error
	Update(
		ctx context.Context, albumID domain.ID, album *domain.Album, executor *domain.User,
	) (*domain.Album, error)

	PutImage(ctx context.Context, albumID domain.ID, imageID domain.ID, executor *domain.User) error
	DeleteImage(ctx context.Context, albumID domain.ID, imageID domain.ID, executor *domain.User) error
}

type AlbumHandlers struct {
	uc     albumUseCase
	logger logger.Logger
}

func NewAlbumHandlers(uc albumUseCase, logger logger.Logger) *AlbumHandlers {
	return &AlbumHandlers{uc: uc, logger: logger}
}

func (h *AlbumHandlers) Create() echo.HandlerFunc {
	type createDTO struct {
		Name        string `json:"name" validate:"required,gte=1,lte=128"`
		Description string `json:"description" validate:"gte=1,lte=512"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		cr := new(createDTO)
		if err := rest.DecodeEchoBody(c, cr); err != nil {
			h.logger.Errorf("Create.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("Create body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, cr); err != nil {
			return c.JSON(rest.NewBadRequestError("Create body has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		album := &domain.Album{
			Name:        cr.Name,
			Description: cr.Description,
			AuthorID:    user.ID,
		}

		createdAlbum, err := h.uc.Create(ctx, album)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Create")
		}

		return c.JSON(http.StatusCreated, createdAlbum)
	}
}

func (h *AlbumHandlers) GetByAuthorID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		userID, err := rest.PipeDomainIdentifier(c, "user_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid user ID").Response())
		}

		albums, err := h.uc.GetByAuthorID(ctx, userID)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "GetByAuthorID")
		}

		return c.JSON(http.StatusOK, albums)
	}
}

func (h *AlbumHandlers) GetAlbumImages() echo.HandlerFunc {
	type imageCommentsQuery struct {
		Limit int `query:"limit" validate:"required,gte=1,lte=100"`
		Page  int `query:"page" validate:"required,gte=1"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		albumID, err := rest.PipeDomainIdentifier(c, "album_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid album ID").Response())
		}

		pag := new(imageCommentsQuery)
		if err := rest.DecodeEchoBody(c, pag); err != nil {
			h.logger.Errorf("GetAlbumImages.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("GetAlbumImages body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, pag); err != nil {
			return c.JSON(rest.NewBadRequestError("GetAlbumImages body has incorrect type").Response())
		}

		pagInput := &domain.PaginationInput{
			PerPage: pag.Limit,
			Page:    pag.Page,
		}
		images, err := h.uc.GetAlbumImages(ctx, albumID, pagInput)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "GetAlbumImages")
		}

		return c.JSON(http.StatusOK, images)
	}
}

func (h *AlbumHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		albumID, err := rest.PipeDomainIdentifier(c, "album_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid album ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("AlbumHandlers.Delete.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.Delete(ctx, albumID, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "Delete")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *AlbumHandlers) Update() echo.HandlerFunc {
	type updateDTO struct {
		Name        string `json:"name" validate:"lte=128"`
		Description string `json:"description" validate:"lte=512"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		albumID, err := rest.PipeDomainIdentifier(c, "album_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid album ID").Response())
		}

		up := new(updateDTO)
		if err := rest.DecodeEchoBody(c, up); err != nil {
			h.logger.Errorf("AlbumHandlers.Update.DecodeBody: %v", err)
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, up); err != nil {
			return c.JSON(rest.NewBadRequestError("Update body has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("AlbumHandlers.Update.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		album := &domain.Album{
			Name:        up.Name,
			Description: up.Description,
		}

		updatedAlbum, err := h.uc.Update(ctx, albumID, album, user)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Update")
		}

		return c.JSON(http.StatusOK, updatedAlbum)
	}
}

func (h *AlbumHandlers) PutImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.PipeDomainIdentifier(c, "image_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Image ID").Response())
		}

		albumID, err := rest.PipeDomainIdentifier(c, "album_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Album ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("AlbumHandlers.PutImage.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.PutImage(ctx, albumID, imageID, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "PutImage")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *AlbumHandlers) DeleteImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.PipeDomainIdentifier(c, "image_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Image ID").Response())
		}

		albumID, err := rest.PipeDomainIdentifier(c, "album_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Album ID").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("AlbumHandlers.DeleteImage.GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.DeleteImage(ctx, albumID, imageID, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "DeleteImage")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *AlbumHandlers) responseWithUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrIncorrectImageRef):
		restErr = rest.NewBadRequestError("Incorrect image reference provided")
		break
	case errors.Is(err, usecase.ErrForbidden):
		restErr = rest.NewForbiddenError("You don't have permissions to perform this action")
		break
	case errors.Is(err, usecase.ErrNotFound):
		restErr = rest.NewNotFoundError("Album not found")
		break
	default:
		h.logger.Errorf("AlbumUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
		break
	}

	return c.JSON(restErr.Response())
}
