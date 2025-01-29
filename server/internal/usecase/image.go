package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/pkg/logger"
)

const imageTTL = 3600

type ImageFileStorage interface {
	Put(ctx context.Context, file *domain.FileNode) error
	Delete(ctx context.Context, path string) error
}

type ImageCache interface {
	Get(ctx context.Context, id string) (*domain.Image, error)
	Set(ctx context.Context, id string, image *domain.Image, ttl int) error
	Del(ctx context.Context, id string) error
}

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (*domain.Image, error)
	GetByID(ctx context.Context, id domain.ID) (*domain.Image, error)
	FindMany(ctx context.Context, ids []domain.ID) ([]domain.ImageWithMeta, error)
	Delete(ctx context.Context, id domain.ID) error
	GetDetailed(ctx context.Context, id domain.ID) (*domain.DetailedImage, error)
	Update(ctx context.Context, id domain.ID, image *domain.Image) (*domain.Image, error)
	AddView(ctx context.Context, imageID domain.ID, userID *domain.ID) error
	States(ctx context.Context, imageID domain.ID, userID domain.ID) (*domain.ImageStates, error)
	Discover(
		ctx context.Context, pagInput *domain.PaginationInput, sort domain.ImageSortMethod,
	) (*domain.Pagination[domain.ImageWithMeta], error)
	HasLike(ctx context.Context, imageID domain.ID, userID domain.ID) (bool, error)
	AddLike(ctx context.Context, imageID domain.ID, userID domain.ID) error
	RemoveLike(ctx context.Context, imageID domain.ID, userID domain.ID) error
	Favorites(
		ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.ImageWithMeta], error)

	repository.Transactional
}

type ImageFeaturesUseCase interface {
	CreateFileNode(ctx context.Context, file *domain.File) (*domain.FileNode, error)
	ExtractFeatures(ctx context.Context, imageID domain.ID, file *domain.FileNode) error
	Similar(ctx context.Context, imageID domain.ID) ([]domain.ID, error)
	DeleteFeatures(ctx context.Context, imageID domain.ID) error
}

type ImageAccessPolicy interface {
	CanModify(user *domain.User, image *domain.Image) bool
}

type NotificationManager interface {
	Notify(ctx context.Context, userID domain.ID, notif *domain.Notification) error
}

type imageUseCase struct {
	storage    ImageFileStorage
	cache      ImageCache
	repo       ImageRepository
	featuresUC ImageFeaturesUseCase
	acl        ImageAccessPolicy
	notifMng   NotificationManager
	logger     logger.Logger
}

func NewImageUseCase(
	storage ImageFileStorage,
	cache ImageCache,
	repo ImageRepository,
	featuresUC ImageFeaturesUseCase,
	acl ImageAccessPolicy,
	notifMng NotificationManager,
	logger logger.Logger,
) *imageUseCase {
	return &imageUseCase{
		storage:    storage,
		repo:       repo,
		featuresUC: featuresUC,
		cache:      cache,
		acl:        acl,
		notifMng:   notifMng,
		logger:     logger,
	}
}

func (uc *imageUseCase) Create(
	ctx context.Context,
	image *domain.Image,
	file *domain.File,
	executor *domain.User,
) (img *domain.Image, err error) {
	err = uc.repo.DoInTransaction(ctx, func(ctx context.Context) error {
		fileNode, err := uc.featuresUC.CreateFileNode(ctx, file)
		if err != nil {
			return fmt.Errorf("failed to create file node: %w", err)
		}
		image.Path = fileNode.Name

		createdImg, err := uc.repo.Create(ctx, image)
		if err != nil {
			return fmt.Errorf("failed to create image: %w", err)
		}

		if err := uc.featuresUC.ExtractFeatures(ctx, createdImg.ID, fileNode); err != nil {
			if featErr, ok := err.(*extractImgFeatErr); ok && !featErr.fatal {
				uc.notifMng.Notify(ctx, executor.ID, &domain.Notification{
					Title:   "Failed to extract features from image",
					Message: fmt.Sprintf("We're sorry, but we were unable to extract the features of your image %s, which may affect the indexing of your image. Please, try again manually", image.Path),
				})
			} else {
				return fmt.Errorf("failed to extract features: %w", err)
			}
		}

		if err := uc.storage.Put(ctx, fileNode); err != nil {
			return fmt.Errorf("failed to store raw image: %w", err)
		}

		img = createdImg
		return nil
	})
	if err != nil {
		uc.logger.Error(err)
	}

	return
}

func (uc *imageUseCase) Similar(
	ctx context.Context, id domain.ID,
) ([]domain.ImageWithMeta, error) {
	if _, err := uc.GetByID(ctx, id); err != nil {
		return nil, err
	}

	ids, err := uc.featuresUC.Similar(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.repo.FindMany(ctx, ids)
}

func (uc *imageUseCase) Delete(
	ctx context.Context,
	id domain.ID,
	executor *domain.User,
) error {
	img, err := uc.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if canEdit := uc.acl.CanModify(executor, img); !canEdit {
		return ErrForbidden
	}

	err = uc.repo.DoInTransaction(ctx, func(ctx context.Context) error {
		if err := uc.repo.Delete(ctx, id); err != nil {
			return err
		}

		if err := uc.storage.Delete(ctx, img.Path); err != nil {
			return err
		}

		if err := uc.featuresUC.DeleteFeatures(ctx, id); err != nil {
			uc.logger.Errorf("Failed to delete features: %v", err)
		}

		return nil
	})
	if err != nil {
		uc.logger.Error(err)
		return err
	}

	uc.deleteCachedImage(ctx, id)
	return nil
}

func (uc *imageUseCase) GetDetailed(ctx context.Context, id domain.ID) (*domain.DetailedImage, error) {
	img, err := uc.repo.GetDetailed(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return img, nil
}

func (uc *imageUseCase) AddView(ctx context.Context, imageID domain.ID, userID *domain.ID) error {
	return uc.repo.AddView(ctx, imageID, userID)
}

func (uc *imageUseCase) AddLike(ctx context.Context, imageID domain.ID, userID domain.ID) error {
	return uc.repo.AddLike(ctx, imageID, userID)
}

func (uc *imageUseCase) RemoveLike(ctx context.Context, imageID domain.ID, userID domain.ID) error {
	// We should check for like existence to make sure that we don't cause UX conflicts
	// For example, when the number of likes is negative
	if hasLike := uc.HasLike(ctx, imageID, userID); !hasLike {
		return ErrUnprocessable
	}

	return uc.repo.RemoveLike(ctx, imageID, userID)
}

// Since you don't need to provide the data that exists for the existence check,
// we don't check the existence of the related data, because it may affect performance
func (uc *imageUseCase) States(
	ctx context.Context, imageID domain.ID, userID domain.ID,
) (*domain.ImageStates, error) {
	return uc.repo.States(ctx, imageID, userID)
}

func (uc *imageUseCase) Discover(
	ctx context.Context,
	pagInput *domain.PaginationInput,
	sort domain.ImageSortMethod,
) (*domain.Pagination[domain.ImageWithMeta], error) {
	pag, err := uc.repo.Discover(ctx, pagInput, sort)
	if err != nil && errors.Is(err, repository.ErrIncorrectInput) {
		return nil, ErrUnprocessable
	}

	return pag, err
}

func (uc *imageUseCase) Favorites(
	ctx context.Context,
	userID domain.ID,
	pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.ImageWithMeta], error) {
	pag, err := uc.repo.Favorites(ctx, userID, pagInput)
	if err != nil && errors.Is(err, repository.ErrIncorrectInput) {
		return nil, ErrUnprocessable
	}

	return pag, err
}

func (uc *imageUseCase) HasLike(ctx context.Context, imageID domain.ID, userID domain.ID) bool {
	hasLike, err := uc.repo.HasLike(ctx, imageID, userID)
	if err != nil {
		uc.logger.Errorf("ImageUseCase.HasLike: %v", err)
	}

	return hasLike
}

func (uc *imageUseCase) GetByID(ctx context.Context, id domain.ID) (*domain.Image, error) {
	cachedImg, err := uc.cache.Get(ctx, id.String())
	if cachedImg != nil {
		return cachedImg, err
	}

	img, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := uc.cache.Set(ctx, id.String(), img, imageTTL); err != nil {
		uc.logger.Errorf("ImageUseCase.GetById.Set: %v", err)
	}

	return img, nil
}

func (uc *imageUseCase) Update(
	ctx context.Context,
	id domain.ID,
	image *domain.Image,
	executor *domain.User,
) (*domain.Image, error) {
	img, err := uc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if canEdit := uc.acl.CanModify(executor, img); !canEdit {
		return nil, ErrForbidden
	}

	updated, err := uc.repo.Update(ctx, id, image)
	if err != nil {
		return nil, err
	}

	uc.deleteCachedImage(ctx, id)
	return updated, nil
}

func (uc *imageUseCase) deleteCachedImage(ctx context.Context, id domain.ID) {
	if err := uc.cache.Del(ctx, id.String()); err != nil {
		uc.logger.Errorf("ImageUseCase.deleteCached: %v", err)
	}
}
