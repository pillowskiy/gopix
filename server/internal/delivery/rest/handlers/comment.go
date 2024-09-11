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

type CommentUseCase interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetByImageID(
		ctx context.Context,
		imageID int,
		pagInput *domain.PaginationInput,
		sort domain.CommentSortMethod,
	) (*domain.Pagination[domain.DetailedComment], error)
	Update(ctx context.Context, commentID int, comment *domain.Comment, executor *domain.User) (*domain.Comment, error)
	Delete(ctx context.Context, commentID int, executor *domain.User) error
}

type CommentHandlers struct {
	uc     CommentUseCase
	logger logger.Logger
}

func NewCommentHandlers(uc CommentUseCase, logger logger.Logger) *CommentHandlers {
	return &CommentHandlers{uc: uc, logger: logger}
}

func (h *CommentHandlers) Create() echo.HandlerFunc {
	type createDTO struct {
		Text string `json:"text" validate:"required,gte=1,lte=512"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.IntParam(c, "image_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Image ID has incorrect type").Response())
		}

		cmt := new(createDTO)
		if err := rest.DecodeEchoBody(c, cmt); err != nil {
			return c.JSON(rest.NewBadRequestError("Comment body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, cmt); err != nil {
			return c.JSON(rest.NewBadRequestError("Comment body has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		comment := &domain.Comment{
			ImageID:  imageID,
			Text:     cmt.Text,
			AuthorID: user.ID,
		}

		createdCmt, err := h.uc.Create(ctx, comment)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Create")
		}

		return c.JSON(http.StatusCreated, createdCmt)
	}
}

func (h *CommentHandlers) GetByImageID() echo.HandlerFunc {
	type imageCommentsQuery struct {
		Limit int    `query:"limit" validate:"required,gte=1,lte=100"`
		Page  int    `query:"page" validate:"required,gte=1"`
		Sort  string `query:"sort" validate:"oneof=popular newest oldest mostViewed"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		imageID, err := rest.IntParam(c, "image_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Image ID has incorrect type").Response())
		}

		q := new(imageCommentsQuery)
		if err := rest.DecodeEchoBody(c, q); err != nil {
			return c.JSON(rest.NewBadRequestError("Query has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, q); err != nil {
			return c.JSON(rest.NewBadRequestError("Query has incorrect type").Response())
		}

		pagInput := &domain.PaginationInput{Page: q.Page, PerPage: q.Limit}
		comments, err := h.uc.GetByImageID(ctx, imageID, pagInput, domain.CommentSortMethod(q.Sort))
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "GetByImageID")
		}

		return c.JSON(http.StatusOK, comments)
	}
}

func (h *CommentHandlers) Update() echo.HandlerFunc {
	type updateDTO struct {
		Text string `json:"text" validate:"gte=1,lte=512"`
	}

	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		commentID, err := rest.IntParam(c, "comment_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Comment ID has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		cmt := new(updateDTO)
		if err := rest.DecodeEchoBody(c, cmt); err != nil {
			return c.JSON(rest.NewBadRequestError("Comment body has incorrect type").Response())
		}

		if err := validator.ValidateStruct(ctx, cmt); err != nil {
			return c.JSON(rest.NewBadRequestError("Comment body has incorrect type").Response())
		}

		comment := &domain.Comment{
			ID:       commentID,
			Text:     cmt.Text,
			AuthorID: user.ID,
		}

		updatedCmt, err := h.uc.Update(ctx, commentID, comment, user)
		if err != nil {
			return h.responseWithUseCaseErr(c, err, "Update")
		}

		return c.JSON(http.StatusOK, updatedCmt)
	}
}

func (h *CommentHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := rest.GetEchoRequestCtx(c)

		commentID, err := rest.IntParam(c, "comment_id")
		if err != nil {
			return c.JSON(rest.NewBadRequestError("Comment ID has incorrect type").Response())
		}

		user, err := GetContextUser(c)
		if err != nil {
			h.logger.Errorf("GetContextUser: %v", err)
			return c.JSON(rest.NewUnauthorizedError("Unauthorized").Response())
		}

		if err := h.uc.Delete(ctx, commentID, user); err != nil {
			return h.responseWithUseCaseErr(c, err, "Delete")
		}

		return c.JSON(http.StatusOK, true)
	}
}

func (h *CommentHandlers) responseWithUseCaseErr(c echo.Context, err error, trace string) error {
	var restErr *rest.Error
	switch {
	case errors.Is(err, usecase.ErrIncorrectImageRef):
		restErr = rest.NewBadRequestError("Incorrect image reference provided")
		break
	case errors.Is(err, usecase.ErrForbidden):
		restErr = rest.NewForbiddenError("You don't have permissions to perform this action")
		break
	case errors.Is(err, usecase.ErrUnprocessable):
		restErr = rest.NewBadRequestError("Incorrect data provided")
		break
	case errors.Is(err, usecase.ErrAlreadyExists):
		restErr = rest.NewConflictError("You've already commented this image")
		break
	case errors.Is(err, usecase.ErrNotFound):
		restErr = rest.NewNotFoundError("Comment not found")
		break
	default:
		h.logger.Errorf("CommentUseCase.%s: %v", trace, err)
		restErr = rest.NewInternalServerError()
		break
	}

	return c.JSON(restErr.Response())
}
