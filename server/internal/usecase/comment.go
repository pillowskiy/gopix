package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pkg/errors"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetByImageID(ctx context.Context, imageID int) ([]domain.DetailedComment, error)
	GetByID(ctx context.Context, imageID int) (*domain.Comment, error)
	Delete(ctx context.Context, commentID int) error
	Update(ctx context.Context, commentID int, comment *domain.Comment) (*domain.Comment, error)
	HasUserCommented(ctx context.Context, commentID int, userID int) (bool, error)
}

type CommentAccessPolicy interface {
	CanModify(user *domain.User, comment *domain.Comment) bool
}

type commentUseCase struct {
	repo   CommentRepository
	acl    CommentAccessPolicy
	logger logger.Logger
}

func NewCommentUseCase(repo CommentRepository, acl CommentAccessPolicy, logger logger.Logger) *commentUseCase {
	return &commentUseCase{repo: repo, acl: acl, logger: logger}
}

func (uc *commentUseCase) Create(
	ctx context.Context,
	comment *domain.Comment,
) (*domain.Comment, error) {
	commented, err := uc.repo.HasUserCommented(ctx, comment.ImageID, comment.AuthorID)
	if err != nil {
		return nil, errors.Wrap(err, "commentUseCase.Create.HasUserCommented")
	}

	if commented {
		return nil, ErrAlreadyExists
	}

	return uc.repo.Create(ctx, comment)
}

func (uc *commentUseCase) GetByImageID(
	ctx context.Context,
	imageID int,
) ([]domain.DetailedComment, error) {
	return uc.repo.GetByImageID(ctx, imageID)
}

func (uc *commentUseCase) GetByID(ctx context.Context, commentID int) (*domain.Comment, error) {
	cmt, err := uc.repo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return cmt, nil
}

func (uc *commentUseCase) Delete(ctx context.Context, commentID int, executor *domain.User) error {
	cmt, err := uc.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if canModify := uc.acl.CanModify(executor, cmt); !canModify {
		return ErrForbidden
	}

	return uc.repo.Delete(ctx, commentID)
}

func (uc *commentUseCase) Update(
	ctx context.Context,
	commentID int,
	comment *domain.Comment,
	executor *domain.User,
) (*domain.Comment, error) {
	cmt, err := uc.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if canModify := uc.acl.CanModify(executor, cmt); !canModify {
		return nil, ErrForbidden
	}

	return uc.repo.Update(ctx, commentID, comment)
}
