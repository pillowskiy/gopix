package usecase

import (
	"context"
	goErrors "errors"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
)

type AlbumRepository interface {
	Create(ctx context.Context, album *domain.Album) (*domain.Album, error)
	GetByID(ctx context.Context, albumID int) (*domain.Album, error)
	GetByAuthorID(ctx context.Context, authorID int) ([]domain.Album, error)
	GetAlbumImages(
		ctx context.Context, albumID int, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.Image], error)
	Delete(ctx context.Context, albumID int) error
	Update(ctx context.Context, albumID int, album *domain.Album) (*domain.Album, error)

	PutImage(ctx context.Context, albumID int, imageID int) error
	DeleteImage(ctx context.Context, albumID int, imageID int) error
}

type AlbumAccessPolicy interface {
	CanModify(user *domain.User, album *domain.Album) bool
}

type AlbumImageUseCase interface {
	GetByID(ctx context.Context, imageID int) (*domain.Image, error)
}

type albumUseCase struct {
	repo    AlbumRepository
	acl     AlbumAccessPolicy
	imageUC AlbumImageUseCase
}

func NewAlbumUseCase(
	repo AlbumRepository,
	acl AlbumAccessPolicy,
	imageUC AlbumImageUseCase,
) *albumUseCase {
	return &albumUseCase{repo: repo, acl: acl, imageUC: imageUC}
}

func (uc *albumUseCase) Create(ctx context.Context, album *domain.Album) (*domain.Album, error) {
	return uc.repo.Create(ctx, album)
}

func (uc *albumUseCase) GetByAuthorID(ctx context.Context, authorID int) ([]domain.Album, error) {
	album, err := uc.repo.GetByAuthorID(ctx, authorID)

	if err != nil {
		if goErrors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "AlbumUseCase.GetByAuthorID")
	}

	return album, nil
}

func (uc *albumUseCase) GetByID(ctx context.Context, albumID int) (*domain.Album, error) {
	album, err := uc.repo.GetByID(ctx, albumID)

	if err != nil {
		if goErrors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "AlbumUseCase.GetByID")
	}

	return album, nil
}

func (uc *albumUseCase) GetAlbumImages(
	ctx context.Context, albumID int, pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.Image], error) {
	if _, err := uc.GetByID(ctx, albumID); err != nil {
		return nil, err
	}

	return uc.repo.GetAlbumImages(ctx, albumID, pagInput)
}

func (uc *albumUseCase) Delete(ctx context.Context, albumID int, executor *domain.User) error {
	if err := uc.ExistsAndModifiable(ctx, executor, albumID); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, albumID)
}

func (uc *albumUseCase) Update(
	ctx context.Context,
	albumID int,
	album *domain.Album,
	executor *domain.User,
) (*domain.Album, error) {
	if err := uc.ExistsAndModifiable(ctx, executor, albumID); err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, albumID, album)
}

func (uc *albumUseCase) PutImage(
	ctx context.Context, albumID int, imageID int, executor *domain.User,
) error {
	if err := uc.ExistsAndModifiable(ctx, executor, albumID); err != nil {
		return err
	}

	if err := uc.correctImageRef(ctx, imageID); err != nil {
		return err
	}

	return uc.repo.PutImage(ctx, albumID, imageID)
}

func (uc *albumUseCase) DeleteImage(
	ctx context.Context, albumID int, imageID int, executor *domain.User,
) error {
	if err := uc.ExistsAndModifiable(ctx, executor, albumID); err != nil {
		return err
	}

	if err := uc.correctImageRef(ctx, imageID); err != nil {
		return err
	}

	return uc.repo.DeleteImage(ctx, albumID, imageID)

}

func (uc *albumUseCase) ExistsAndModifiable(
	ctx context.Context, user *domain.User, albumID int,
) error {
	album, err := uc.GetByID(ctx, albumID)
	if err != nil {
		return err
	}

	if canModify := uc.acl.CanModify(user, album); !canModify {
		return ErrForbidden
	}

	return nil
}

func (uc *albumUseCase) correctImageRef(ctx context.Context, imageID int) error {
	img, err := uc.imageUC.GetByID(ctx, imageID)

	isValidImage := img != nil && img.AccessLevel == domain.ImageAccessPublic
	if err != nil || !isValidImage {
		return ErrIncorrectImageRef
	}

	return nil
}
