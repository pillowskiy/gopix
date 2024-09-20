package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
)

type TagRepository interface {
	Upsert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	UpsertImageTags(ctx context.Context, tag *domain.Tag, imageID domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Tag, error)
	GetByName(ctx context.Context, name string) (*domain.Tag, error)
	Search(ctx context.Context, name string) ([]domain.Tag, error)
	Delete(ctx context.Context, id domain.ID) error
}

type TagImageUseCase interface {
	GetByID(ctx context.Context, id domain.ID) (*domain.Image, error)
}

type TagAccessPolicy interface {
	CanModifyImageTags(user *domain.User, image *domain.Image) bool
}

type tagUseCase struct {
	repo    TagRepository
	imageUC TagImageUseCase
	acl     TagAccessPolicy
}

func NewTagUseCase(repo TagRepository, acl TagAccessPolicy, imageUC TagImageUseCase) *tagUseCase {
	return &tagUseCase{repo: repo, acl: acl, imageUC: imageUC}
}

func (uc *tagUseCase) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	existingTag, err := uc.repo.GetByName(ctx, tag.Name)
	if existingTag != nil || err == nil {
		return nil, ErrAlreadyExists
	}

	return uc.repo.Upsert(ctx, tag)
}

func (uc *tagUseCase) UpsertImageTag(
	ctx context.Context, tag *domain.Tag, imageID domain.ID, executor *domain.User,
) error {
	image, err := uc.imageUC.GetByID(ctx, imageID)
	if err != nil {
		return ErrIncorrectImageRef
	}

	if !uc.acl.CanModifyImageTags(executor, image) {
		return ErrForbidden
	}

	return uc.repo.UpsertImageTags(ctx, tag, imageID)
}

func (uc *tagUseCase) Search(ctx context.Context, query string) ([]domain.Tag, error) {
	return uc.repo.Search(ctx, query)
}

func (uc *tagUseCase) Delete(ctx context.Context, tagID domain.ID) error {
	tag, err := uc.GetByID(ctx, tagID)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, tag.ID)
}

func (uc *tagUseCase) GetByID(ctx context.Context, tagID domain.ID) (*domain.Tag, error) {
	tag, err := uc.repo.GetByID(ctx, tagID)
	if tag == nil || err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "tagUseCase.Delete.GetByID")
	}

	return tag, nil
}
