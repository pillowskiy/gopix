package usecase

import (
	"context"
	"errors"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
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
	Update(ctx context.Context, id int, image *domain.Image) (*domain.Image, error)
	AddView(ctx context.Context, view *domain.ImageView) error
}

type imageUseCase struct {
	storage ImageFileStorage
	repo    ImageRepository
	logger  logger.Logger
}

func NewImageUseCase(storage ImageFileStorage, repo ImageRepository, logger logger.Logger) *imageUseCase {
	return &imageUseCase{storage: storage, repo: repo, logger: logger}
}

func (uc *imageUseCase) Create(
	ctx context.Context,
	image *domain.Image,
	file *domain.FileNode,
) (*domain.Image, error) {
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

func (uc *imageUseCase) Delete(
	ctx context.Context,
	id int,
) error {
	img, err := uc.GetById(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := uc.storage.Delete(ctx, img.Path); err != nil {
		uc.repo.Create(ctx, img)
		return err
	}

	return nil
}

func (uc *imageUseCase) GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error) {
	img, err := uc.repo.GetDetailed(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return img, nil
}

func (uc *imageUseCase) AddView(ctx context.Context, view *domain.ImageView) error {
	return uc.repo.AddView(ctx, view)
}

func (uc *imageUseCase) GetById(ctx context.Context, id int) (*domain.Image, error) {
	img, err := uc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return img, nil
}

func (uc *imageUseCase) Update(
	ctx context.Context,
	id int,
	image *domain.Image,
) (*domain.Image, error) {
	_, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, id, image)
}
