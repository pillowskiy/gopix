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
	GetByImageID(
		ctx context.Context, imageID domain.ID, pagInput *domain.PaginationInput, sort domain.CommentSortMethod,
	) (*domain.Pagination[domain.DetailedComment], error)
	GetByID(ctx context.Context, imageID domain.ID) (*domain.Comment, error)
	Delete(ctx context.Context, commentID domain.ID) error
	Update(ctx context.Context, commentID domain.ID, comment *domain.Comment) (*domain.Comment, error)
	HasUserCommented(ctx context.Context, commentID domain.ID, userID domain.ID) (bool, error)
	GetReplies(ctx context.Context, commentID domain.ID, userID *domain.ID) ([]domain.DetailedComment, error)

	LikeComment(ctx context.Context, commentID domain.ID, userID domain.ID) error
	UnlikeComment(ctx context.Context, commentID domain.ID, userID domain.ID) error
	HasUserLikedComment(ctx context.Context, commentID domain.ID, userID domain.ID) (bool, error)
}

type CommentAccessPolicy interface {
	CanModify(user *domain.User, comment *domain.Comment) bool
}

type CommentImageUseCase interface {
	GetByID(ctx context.Context, imageID domain.ID) (*domain.Image, error)
}

type commentUseCase struct {
	repo    CommentRepository
	acl     CommentAccessPolicy
	imageUC CommentImageUseCase
	logger  logger.Logger
}

func NewCommentUseCase(
	repo CommentRepository,
	acl CommentAccessPolicy,
	imageUC CommentImageUseCase,
	logger logger.Logger,
) *commentUseCase {
	return &commentUseCase{repo: repo, acl: acl, imageUC: imageUC, logger: logger}
}

func (uc *commentUseCase) Create(
	ctx context.Context,
	comment *domain.Comment,
) (*domain.Comment, error) {
	if _, err := uc.imageUC.GetByID(ctx, comment.ImageID); err != nil {
		return nil, ErrIncorrectImageRef
	}

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
	imageID domain.ID,
	pagInput *domain.PaginationInput,
	sort domain.CommentSortMethod,
) (*domain.Pagination[domain.DetailedComment], error) {
	if _, err := uc.imageUC.GetByID(ctx, imageID); err != nil {
		return nil, ErrIncorrectImageRef
	}

	pag, err := uc.repo.GetByImageID(ctx, imageID, pagInput, sort)
	if err != nil && errors.Is(err, repository.ErrIncorrectInput) {
		return nil, ErrUnprocessable
	}

	return pag, err
}

func (uc *commentUseCase) GetByID(ctx context.Context, commentID domain.ID) (*domain.Comment, error) {
	cmt, err := uc.repo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return cmt, nil
}

func (uc *commentUseCase) Delete(ctx context.Context, commentID domain.ID, executor *domain.User) error {
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
	commentID domain.ID,
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

func (uc *commentUseCase) GetReplies(ctx context.Context, commentID domain.ID, executorID *domain.ID) ([]domain.DetailedComment, error) {
	// TODO: potentially omit comment existence check
	if _, err := uc.GetByID(ctx, commentID); err != nil {
		return nil, err
	}

	return uc.repo.GetReplies(ctx, commentID, executorID)
}

func (uc *commentUseCase) LikeComment(ctx context.Context, commentID domain.ID, executor *domain.User) error {
	_, err := uc.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	liked, err := uc.repo.HasUserLikedComment(ctx, commentID, executor.ID)
	if liked {
		return ErrAlreadyExists
	}

	if err := uc.repo.LikeComment(ctx, commentID, executor.ID); err != nil {
		return errors.Wrap(err, "commentUseCase.LikeComment")
	}

	return nil
}

func (uc *commentUseCase) UnlikeComment(ctx context.Context, commentID domain.ID, executor *domain.User) error {
	liked, err := uc.repo.HasUserLikedComment(ctx, commentID, executor.ID)
	if !liked || err != nil {
		return ErrNotFound
	}

	if err := uc.repo.UnlikeComment(ctx, commentID, executor.ID); err != nil {
		return errors.Wrap(err, "commentUseCase.UnlikeComment")
	}

	return nil
}
