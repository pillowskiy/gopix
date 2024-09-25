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
	PostgresRepository
	viewBatcher batch.Batcher[viewBatchItem]
	likeBatcher batch.Batcher[likeBatchItem]
}

func NewImageRepository(db *sqlx.DB) *imageRepository {
	repo := &imageRepository{
		PostgresRepository: PostgresRepository{db},
	}

	repo.viewBatcher = batch.NewWithConfig(imageViewsBatchAgg, repo.processViewsBatch, &imageBatchConfig)
	go repo.viewBatcher.Ticker(imageBatchTickDuration)

	repo.likeBatcher = batch.NewWithConfig(imageLikesBatchAgg, repo.processLikesBatch, &imageBatchConfig)
	go repo.likeBatcher.Ticker(imageBatchTickDuration)

	return repo
}

func (r *imageRepository) Create(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.ext(ctx).QueryRowxContext(
		ctx,
		createImageQuery,
		image.AuthorID,
		image.Path,
		image.Title,
		image.Description,
		image.AccessLevel,
		image.ExpiresAt,
		image.Mime,
		image.Ext,
	)

	if err := rowx.StructScan(img); err != nil {
		return nil, errors.Wrap(err, "ImageRepository.Create.StructScan")
	}

	return img, nil
}

func (r *imageRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.ext(ctx).QueryRowxContext(ctx, getByIdImageQuery, id)

	if err := rowx.StructScan(img); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "ImageRepository.GetById.StructScan")
	}

	return img, nil
}

func (r *imageRepository) FindMany(
	ctx context.Context, ids []domain.ID,
) ([]domain.ImageWithAuthor, error) {
	query, args, err := sqlx.In(findManyImagesQuery, ids)
	if err != nil {
		return nil, errors.Wrap(err, "ImageRepository.FindMany.In")
	}
	query = r.db.Rebind(query)

	rows, err := r.ext(ctx).QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "ImageRepository.FindMany.QueryxContext")
	}

	images, err := scanToStructSliceOf[domain.ImageWithAuthor](rows)
	if err != nil {
		return nil, errors.Wrap(err, "ImageRepository.FindMany.scanToStructSliceOf")
	}

	return images, nil
}

func (r *imageRepository) Delete(ctx context.Context, id domain.ID) error {
	if _, err := r.ext(ctx).ExecContext(ctx, deleteImageQuery, id); err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) GetDetailed(ctx context.Context, id domain.ID) (*domain.DetailedImage, error) {
	var detailedImage domain.DetailedImage
	var tagsJSON []byte

	err := r.ext(ctx).QueryRowxContext(ctx, getDetailedImageQuery, id).Scan(
		&detailedImage.ID,
		&detailedImage.AuthorID,
		&detailedImage.Path,
		&detailedImage.Title,
		&detailedImage.Description,
		&detailedImage.AccessLevel,
		&detailedImage.Ext,
		&detailedImage.Mime,
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

	detailedImage.Tags = []domain.ImageTag{}
	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &detailedImage.Tags); err != nil {
			return nil, errors.Wrap(err, "imageRepository.GetDetailed.Unmarshal")
		}
	}

	detailedImage.Views += r.viewBatcher.CountByGroup(imageGroupKey(id))
	detailedImage.Likes += r.likeBatcher.CountByGroup(imageGroupKey(id))

	return &detailedImage, nil
}

func (r *imageRepository) Update(ctx context.Context, id domain.ID, image *domain.Image) (*domain.Image, error) {
	img := new(domain.Image)
	rowx := r.ext(ctx).QueryRowxContext(
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
) (*domain.Pagination[domain.ImageWithAuthor], error) {
	sortQuery, ok := imagesSortQuery.SortQuery(string(sort))
	if !ok {
		return nil, repository.ErrIncorrectInput
	}

	q := fmt.Sprintf(`
  SELECT
    i.*,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url"
  FROM images i
  LEFT JOIN users u ON i.author_id = u.id
  JOIN images_analytics a ON a.image_id = i.id
  WHERE access_level = 'public'::access_level
  ORDER BY %s LIMIT $1 OFFSET $2
  `, sortQuery)

	limit := pagInput.PerPage
	rowx, err := r.ext(ctx).QueryxContext(ctx, q, limit, (pagInput.Page-1)*limit)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.Discover.Queryx")
	}
	defer rowx.Close()

	images, err := scanToStructSliceOf[domain.ImageWithAuthor](rowx)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.Discover.Scan")
	}

	pagination := &domain.Pagination[domain.ImageWithAuthor]{
		PaginationInput: *pagInput,
		Items:           images,
	}

	countQuery := `SELECT COUNT(1) FROM images`
	_ = r.ext(ctx).QueryRowxContext(ctx, countQuery).Scan(&pagination.Total)

	return pagination, nil
}

func (r *imageRepository) States(ctx context.Context, imageID domain.ID, userID domain.ID) (*domain.ImageStates, error) {
	states := new(domain.ImageStates)
	rowx := r.ext(ctx).QueryRowxContext(ctx, statesImageQuery, imageID, userID)
	if err := rowx.StructScan(states); err != nil {
		return nil, errors.Wrap(err, "imageRepository.States.StructScan")
	}

	batchView := r.getBatchView(imageID, userID)
	if batchView != nil {
		states.Viewed = true
	}

	batchLike := r.getBatchLike(imageID, userID)
	if batchLike != nil {
		states.Liked = batchLike.Liked
	}

	return states, nil
}

func (r *imageRepository) HasLike(ctx context.Context, imageID domain.ID, userID domain.ID) (bool, error) {
	var hasLike bool

	batchLike := r.getBatchLike(imageID, userID)
	if batchLike != nil {
		return batchLike.Liked, nil
	}

	rowx := r.ext(ctx).QueryRowxContext(ctx, hasLikeImageQuery, imageID, userID)
	if err := rowx.Scan(&hasLike); err != nil {
		return false, errors.Wrap(err, "imageRepository.HasLike.Scan")
	}

	return hasLike, nil
}

func (r *imageRepository) getBatchView(imageID domain.ID, userID domain.ID) *viewBatchItem {
	return r.viewBatcher.Search(imageWithUserKey(imageID, userID), nil)
}

func (r *imageRepository) getBatchLike(imageID domain.ID, userID domain.ID) *likeBatchItem {
	return r.likeBatcher.Search(imageWithUserKey(imageID, userID), nil)
}

func (r *imageRepository) AddLike(ctx context.Context, imageID domain.ID, userID domain.ID) error {
	r.likeBatcher.Add(likeBatchItem{
		ImageID: imageID,
		UserID:  userID,
		Liked:   true,
	})
	return nil
}

func (r *imageRepository) RemoveLike(ctx context.Context, imageID domain.ID, userID domain.ID) error {
	r.likeBatcher.Add(likeBatchItem{
		ImageID: imageID,
		UserID:  userID,
		Liked:   false,
	})
	return nil
}

func (r *imageRepository) AddView(ctx context.Context, imageID domain.ID, userID *domain.ID) error {
	r.viewBatcher.Add(viewBatchItem{
		ImageID: imageID,
		UserID:  userID,
	})
	return nil
}
