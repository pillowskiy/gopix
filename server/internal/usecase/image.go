package usecase

import (
	"context"
	"errors"

	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/logger"
)

const imageTTL = 3600

type ImageFileStorage interface {
	Put(ctx context.Context, file *domain.FileNode) error
	Delete(ctx context.Context, path string) error
}

type ImageCache interface {
	Get(ctx context.Context, id int) (*domain.Image, error)
	Set(ctx context.Context, id int, image *domain.Image, ttl int) error
	Del(ctx context.Context, id int) error
}

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (*domain.Image, error)
	GetById(ctx context.Context, id int) (*domain.Image, error)
	Delete(ctx context.Context, id int) error
	GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error)
	Update(ctx context.Context, id int, image *domain.Image) (*domain.Image, error)
	AddView(ctx context.Context, view *domain.ImageView) error
	States(ctx context.Context, imageID int, userID int) (*domain.ImageStates, error)
	Discover(
		ctx context.Context, page int, limit int, sort domain.ImageSortMethod,
	) (*domain.Pagination[domain.Image], error)
	HasLike(ctx context.Context, imageID int, userID int) (bool, error)
	AddLike(ctx context.Context, imageID int, userID int) error
	RemoveLike(ctx context.Context, imageID int, userID int) error
}

type imageUseCase struct {
	storage ImageFileStorage
	cache   ImageCache
	repo    ImageRepository
	logger  logger.Logger
}

func NewImageUseCase(
	storage ImageFileStorage,
	cache ImageCache,
	repo ImageRepository,
	logger logger.Logger,
) *imageUseCase {
	return &imageUseCase{storage: storage, repo: repo, cache: cache, logger: logger}
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

	uc.deleteCachedImage(ctx, id)
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

func (uc *imageUseCase) AddLike(ctx context.Context, imageID int, userID int) error {
	return uc.repo.AddLike(ctx, imageID, userID)
}

func (uc *imageUseCase) RemoveLike(ctx context.Context, imageID int, userID int) error {
	// We should check for like existence to make sure that we don't cause UX conflicts
	// For example, when the number of likes on an image is negative
	if hasLike := uc.HasLike(ctx, imageID, userID); hasLike {
		return ErrUnprocessable
	}

	return uc.repo.RemoveLike(ctx, imageID, userID)
}

func (uc *imageUseCase) States(ctx context.Context, imageID int, userID int) (*domain.ImageStates, error) {
	return uc.repo.States(ctx, imageID, userID)
}

func (uc *imageUseCase) Discover(
	ctx context.Context,
	page int,
	limit int,
	sort domain.ImageSortMethod,
) (*domain.Pagination[domain.Image], error) {
	pag, err := uc.repo.Discover(ctx, page, limit, sort)
	if err != nil && errors.Is(err, repository.ErrIncorrectInput) {
		return nil, ErrUnprocessable
	}

	return pag, err
}

func (uc *imageUseCase) HasLike(ctx context.Context, imageID int, userID int) bool {
	hasLike, err := uc.repo.HasLike(ctx, imageID, userID)
	if err != nil {
		uc.logger.Errorf("ImageUseCase.HasLike: %v", err)
	}

	return hasLike
}

func (uc *imageUseCase) GetById(ctx context.Context, id int) (*domain.Image, error) {
	cachedImg, err := uc.cache.Get(ctx, id)
	if cachedImg != nil {
		return cachedImg, err
	}

	img, err := uc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := uc.cache.Set(ctx, id, img, imageTTL); err != nil {
		uc.logger.Errorf("ImageUseCase.GetById.Set: %v", err)
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

	uc.deleteCachedImage(ctx, id)
	return uc.repo.Update(ctx, id, image)
}

func (uc *imageUseCase) deleteCachedImage(ctx context.Context, id int) {
	if err := uc.cache.Del(ctx, id); err != nil {
		uc.logger.Errorf("ImageUseCase.deleteCached: %v", err)
	}
}
