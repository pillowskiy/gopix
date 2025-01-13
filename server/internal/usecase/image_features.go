package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/worker"
)

type ImageVecRepository interface {
	Similar(ctx context.Context, imageID domain.ID) ([]domain.ID, error)
	Features(ctx context.Context, imageID domain.ID, file *domain.FileNode) error
	DeleteFeatures(ctx context.Context, imageID domain.ID) error
}

type ImagePropsRepository interface {
	Create(ctx context.Context, imageID domain.ID, props *domain.ImageProperties) error
	Properties(ctx context.Context, imageID domain.ID) (*domain.ImageProperties, error)
	Delete(ctx context.Context, imageID domain.ID) error

	repository.Transactional
}

type FeaturesExtractor interface {
	MakeFileNode(ctx context.Context, file *domain.File) (*domain.FileNode, error)
	Features(ctx context.Context, fileNode *domain.FileNode) (*domain.ImageProperties, error)
}

type featureExtractionTask struct {
	imageID  domain.ID
	fileNode *domain.FileNode
}

type imageFeaturesUseCase struct {
	vecRepo       ImageVecRepository
	imgPropsRepo  ImagePropsRepository
	featExtractor FeaturesExtractor
	logger        logger.Logger
	wrk           *worker.Worker[featureExtractionTask]
}

func NewImageFeaturesUseCase(
	vecRepo ImageVecRepository,
	imgPropsRepo ImagePropsRepository,
	featExtractor FeaturesExtractor,
	logger logger.Logger,
) *imageFeaturesUseCase {
	return &imageFeaturesUseCase{
		vecRepo:       vecRepo,
		imgPropsRepo:  imgPropsRepo,
		featExtractor: featExtractor,
		logger:        logger,
	}
}

func (uc *imageFeaturesUseCase) CreateFileNode(ctx context.Context, file *domain.File) (*domain.FileNode, error) {
	defer file.Restore()
	return uc.featExtractor.MakeFileNode(ctx, file)
}

func (uc *imageFeaturesUseCase) ExtractFeatures(ctx context.Context, imageID domain.ID, fileNode *domain.FileNode) error {
	extProps, err := uc.imgPropsRepo.Properties(ctx, imageID)
	if extProps != nil || err == nil {
		return ErrAlreadyExists
	}
	defer fileNode.Restore()

	return uc.imgPropsRepo.DoInTransaction(ctx, func(ctx context.Context) error {
		imgProps, err := uc.featExtractor.Features(ctx, fileNode)
		if err != nil {
			return fmt.Errorf("failed to extract features: %w", err)
		}

		if err := uc.imgPropsRepo.Create(ctx, imageID, imgProps); err != nil {
			return fmt.Errorf("failed to store image properties: %w", err)
		}

		if err := uc.vecRepo.Features(ctx, imageID, fileNode); err != nil {
			uc.logger.Errorf("Failed to extract features vector: %v", err)
		}

		return nil
	})
}

func (uc *imageFeaturesUseCase) DeleteFeatures(ctx context.Context, imageID domain.ID) error {
	// NOTE: We just delete the potential data, we don't care if it exists or not
	return uc.imgPropsRepo.DoInTransaction(ctx, func(ctx context.Context) (err error) {
		if err := uc.vecRepo.DeleteFeatures(ctx, imageID); err != nil {
			err = fmt.Errorf("failed to delete features vector: %w", err)
		}

		if err := uc.imgPropsRepo.Delete(ctx, imageID); err != nil {
			err = fmt.Errorf("failed to delete image properties: %w", err)
		}

		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}

		return err
	})
}

func (uc *imageFeaturesUseCase) Similar(ctx context.Context, imageID domain.ID) ([]domain.ID, error) {
	// TODO: Add search by image properties if the vector repo found nothing
	return uc.vecRepo.Similar(ctx, imageID)
}
