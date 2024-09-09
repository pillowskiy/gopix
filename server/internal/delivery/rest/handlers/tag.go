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

type TagUseCase interface {
	Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	UpsertImageTag(ctx context.Context, tag *domain.Tag, imageID int, executor *domain.User) error
	Search(ctx context.Context, query string) ([]domain.Tag, error)
	Delete(ctx context.Context, tagID int) error
}

type TagHandlers struct {
	uc     TagUseCase
	logger logger.Logger
}

func NewTagHandlers(uc TagUseCase, logger logger.Logger) *TagHandlers {
	return &TagHandlers{uc: uc, logger: logger}
}

func (h *TagHandlers) Create() echo.HandlerFunc {
	type createDTO struct {
		Name string `json:"name" validate:"required,gte=1,lte=32,lowercase"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		if _, err := GetContextUser(c); err != nil {
			h.logger.Errorf("TagHandlers.Create: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		tagInput := new(createDTO)
		if err := rest.DecodeEchoBody(c, tagInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Body").Response())
		}

		if err := validator.ValidateStruct(ctx, tagInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Body").Response())
		}

		tag := &domain.Tag{Name: tagInput.Name}
		createdTag, err := h.uc.Create(ctx, tag)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Create")
		}

		return c.JSON(http.StatusCreated, createdTag)
	}
}

func (h *TagHandlers) UpsertImageTag() echo.HandlerFunc {
	type upsertTagDTO struct {
		Name string `json:"name" validate:"required,gte=1,lte=32,lowercase"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.IntParam(c, "image_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid image ID").Response())
		}

		tagInput := new(upsertTagDTO)
		if err := rest.DecodeEchoBody(c, tagInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Body").Response())
		}

		if err := validator.ValidateStruct(ctx, tagInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Body").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("TagHandlers.UpsertImageTag: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		tag := &domain.Tag{Name: tagInput.Name}
		if err := h.uc.UpsertImageTag(ctx, tag, imageID, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "UpsertImageTag")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *TagHandlers) Search() echo.HandlerFunc {
	type searchDTO struct {
		Query string `query:"query" validate:"required,gte=1,lte=32,lowercase"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		queryInput := new(searchDTO)
		if err := rest.DecodeEchoBody(c, queryInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Query").Response())
		}

		if err := validator.ValidateStruct(ctx, queryInput); err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid Request Query").Response())
		}

		tags, err := h.uc.Search(ctx, queryInput.Query)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Search")
		}

		return c.JSON(http.StatusOK, tags)
	}
}

func (h *TagHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		tagID, err := rest.IntParam(c, "tag_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Invalid tag ID").Response())
		}

		if _, err := GetContextUser(c); err != nil {
			h.logger.Errorf("TagHandlers.Create: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.Delete(ctx, tagID); err != nil {
			return h.responseWithUseCaseErr(c, err, "Delete")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *TagHandlers) responseWithUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrIncorrectImageRef):
		restErr = rest.NewBadRequestError("Incorrect image reference provided")
		break
	case errors.Is(err, usecase.ErrForbidden):
		restErr = rest.NewForbiddenError("You don't have permission to perform this action")
		break
	case errors.Is(err, usecase.ErrAlreadyExists):
		restErr = rest.NewConflictError("Tag with this name already exists")
		break
	case errors.Is(err, usecase.ErrNotFound):
		restErr = rest.NewNotFoundError("Tag not found")
		break
	default:
		h.logger.Errorf("CommentUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
		break
	}

	return c.JSON(restErr.Response())
}
