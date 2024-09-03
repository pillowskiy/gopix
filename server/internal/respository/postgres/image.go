package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/batch"
	"github.com/pkg/errors"
)

var imagesSortQuery = NewSortQueryBuilder().
	AddField(string(domain.ImageNewestSort), SortField{Field: "uploaded_at", Order: sortOrderDESC}).
	AddField(string(domain.ImageOldestSort), SortField{Field: "uploaded_at", Order: sortOrderASC}).
	AddField(string(domain.ImagePopularSort), SortField{Field: "a.likes_count", Order: sortOrderDESC}).
	AddField(string(domain.ImageMostViewedSort), SortField{Field: "a.views_count", Order: sortOrderDESC})

type imageRepository struct {
	db             *sqlx.DB
	viewBatcher    batch.Batcher[viewBatchItem]
	likeBatcher    batch.Batcher[likeBatchItem]
	dislikeBatcher batch.Batcher[likeBatchItem]
}

func NewImageRepository(db *sqlx.DB) *imageRepository {
	repo := &imageRepository{db: db}

	repo.viewBatcher = batch.NewWithConfig(imageViewsBatchAgg, repo.processViewsBatch, &imageBatchConfig)
	go repo.viewBatcher.Ticker(imageBatchTickDuration)

	repo.likeBatcher = batch.NewWithConfig(imageLikesBatchAgg, repo.processLikesBatch, &imageBatchConfig)
	go repo.likeBatcher.Ticker(imageBatchTickDuration)

	return repo
}

func (r *imageRepository) Create(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(
		ctx,
		createImageQuery,
		image.AuthorID,
		image.Path,
		image.Title,
		image.Description,
		image.PHash,
		image.AccessLevel,
		image.ExpiresAt,
	)
	if err := rowx.StructScan(img); err != nil {
		return nil, errors.Wrap(err, "ImageRepository.Create.StructScan")
	}

	return img, nil
}

func (r *imageRepository) GetById(ctx context.Context, id int) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(ctx, getByIdImageQuery, id)

	if err := rowx.StructScan(img); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "ImageRepository.GetById.StructScan")
	}

	return img, nil
}

func (r *imageRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.db.ExecContext(ctx, deleteImageQuery, id); err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error) {
	var detailedImage domain.DetailedImage
	var tagsJSON []byte

	err := r.db.QueryRowxContext(ctx, getDetailedImageQuery, id).Scan(
		&detailedImage.ID,
		&detailedImage.AuthorID,
		&detailedImage.Path,
		&detailedImage.Title,
		&detailedImage.Description,
		&detailedImage.AccessLevel,
		&detailedImage.ExpiresAt,
		&detailedImage.CreatedAt,
		&detailedImage.UpdatedAt,
		&detailedImage.Author.ID,
		&detailedImage.Author.Username,
		&detailedImage.Author.AvatarURL,
		&detailedImage.Likes,
		&detailedImage.Views,
		&tagsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "imageRepository.GetDetailed.Scan")
	}

	if err := json.Unmarshal(tagsJSON, &detailedImage.Tags); err != nil {
		return nil, errors.Wrap(err, "imageRepository.GetDetailed.Unmarshal")
	}

	detailedImage.Views += r.viewBatcher.CountByGroup(imageGroupKey(id))
	detailedImage.Likes += r.likeBatcher.CountByGroup(imageGroupKey(id))

	return &detailedImage, nil
}

func (r *imageRepository) Update(ctx context.Context, id int, image *domain.Image) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(
		ctx,
		updateImageQuery,
		image.Title,
		image.Description,
		image.AccessLevel,
		image.ExpiresAt,
		id,
	)

	if err := rowx.StructScan(img); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "imageRepository.Update.Scan")
	}

	return img, nil
}

func (r *imageRepository) Discover(
	ctx context.Context,
	pagInput *domain.PaginationInput,
	sort domain.ImageSortMethod,
) (*domain.Pagination[domain.Image], error) {
	sortQuery, ok := imagesSortQuery.SortQuery(string(sort))
	if !ok {
		return nil, repository.ErrIncorrectInput
	}

	q := fmt.Sprintf(`
  SELECT
    i.id,
    i.author_id,
    i.path,
    i.title,
    i.description,
    i.access_level,
    i.expires_at,
    i.uploaded_at,
    i.updated_at
  FROM images i
    LEFT JOIN images_analytics a ON a.image_id = i.id
  WHERE access_level = 'public'::access_level
  ORDER BY %s LIMIT $1 OFFSET $2
  `, sortQuery)

	limit := pagInput.PerPage
	rowx, err := r.db.QueryxContext(ctx, q, limit, (pagInput.Page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.Discover.Queryx")
	}
	defer rowx.Close()

	images, err := scanToStructSliceOf[domain.Image](rowx)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.Discover.Scan")
	}

	pagination := &domain.Pagination[domain.Image]{
		PaginationInput: *pagInput,
		Items:           images,
	}

	countQuery := `SELECT COUNT(1) FROM images`
	_ = r.db.QueryRowxContext(ctx, countQuery).Scan(&pagination.Total)

	return pagination, nil
}

func (r *imageRepository) States(ctx context.Context, imageID int, userID int) (*domain.ImageStates, error) {
	states := new(domain.ImageStates)
	rowx := r.db.QueryRowxContext(ctx, statesImageQuery, imageID, userID)
	if err := rowx.StructScan(states); err != nil {
		return nil, errors.Wrap(err, "imageRepository.States.StructScan")
	}

	batchView := r.viewBatcher.Search(imageWithUserKey(imageID, userID), nil)
	if batchView != nil {
		states.Viewed = true
	}

	return states, nil
}

func (r *imageRepository) HasLike(ctx context.Context, imageID int, userID int) (bool, error) {
	var hasLike bool
	rowx := r.db.QueryRowxContext(ctx, hasLikeImageQuery, imageID, userID)
	if err := rowx.Scan(&hasLike); err != nil {
		return false, errors.Wrap(err, "imageRepository.HasLike.Scan")
	}

	return hasLike, nil
}

func (r *imageRepository) AddLike(ctx context.Context, imageID int, userID int) error {
	r.likeBatcher.Add(likeBatchItem{
		ImageID: imageID,
		UserID:  userID,
		Liked:   true,
	})
	return nil
}

func (r *imageRepository) RemoveLike(ctx context.Context, imageID int, userID int) error {
	r.likeBatcher.Add(likeBatchItem{
		ImageID: imageID,
		UserID:  userID,
		Liked:   false,
	})
	return nil
}

func (r *imageRepository) AddView(ctx context.Context, view *domain.ImageView) error {
	r.viewBatcher.Add(viewBatchItem{
		ImageID: view.ImageID,
		UserID:  view.UserID,
	})
	return nil
}
