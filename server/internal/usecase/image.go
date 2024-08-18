package usecase

import (
	"context"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/logger"
)

type ImageFileStorage interface {
	Put(ctx context.Context, file *domain.FileNode) error
	Delete(ctx context.Context, path string) error
}

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (*domain.Image, error)
	GetById(ctx context.Context, id int) (*domain.Image, error)
	Delete(ctx context.Context, id int) error
	GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error)
}

type imageUseCase struct {
	storage ImageFileStorage
	repo    ImageRepository
	logger  logger.Logger
}

func NewImageUseCase(storage ImageFileStorage, repo ImageRepository) *imageUseCase {
	return &imageUseCase{storage: storage, repo: repo}
}

func (uc *imageUseCase) Create(
	ctx context.Context,
	image *domain.Image,
	file *domain.FileNode,
) (*domain.Image, error) {
	if err := file.Prepare(); err != nil {
		return nil, ErrUnprocessableEntity
	}

	image.Path = file.Name
	img, err := uc.repo.Create(ctx, image)
	if err != nil {
		return nil, err
	}

	if err := uc.storage.Put(ctx, file); err != nil {
		uc.logger.Errorf("ImageUseCase.Create.Put: %v", err)
		uc.repo.Delete(ctx, img.ID)
		return nil, err
	}

	return img, nil
}
