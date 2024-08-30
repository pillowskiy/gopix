package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/batch"
	"github.com/pkg/errors"
)

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
	q := `
  INSERT INTO images (author_id, path, title, description, access_level, expires_at)
  VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, '')::access_level, 'public'::access_level), $6) RETURNING *`

	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(
		ctx, q,
		image.AuthorID,
		image.Path,
		image.Title,
		image.Description,
		image.AccessLevel,
		image.ExpiresAt,
	)
	if err := rowx.StructScan(img); err != nil {
		return nil, errors.Wrap(err, "ImageRepository.Create.StructScan")
	}

	return img, nil
}

func (r *imageRepository) GetById(ctx context.Context, id int) (*domain.Image, error) {
	q := `SELECT * FROM images WHERE id = $1`

	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(ctx, q, id)

	if err := rowx.StructScan(img); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "ImageRepository.GetById.StructScan")
	}

	return img, nil
}

func (r *imageRepository) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM images WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error) {
	var detailedImage domain.DetailedImage
	var tagsJSON []byte

	query := `
  SELECT
    i.id,
    i.author_id,
    i.path,
    i.title,
    i.description,
    i.access_level,
    i.expires_at,
    i.uploaded_at,
    i.updated_at,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url",
    COALESCE(a.likes_count, 0) AS likes,
    COALESCE(a.views_count, 0) AS views,
    TO_JSON(COALESCE(
      ARRAY_AGG(
        json_build_object('id', t.id, 'name', t.name)
      ) FILTER (WHERE t.id IS NOT NULL),
      '{}'
    )) AS tags
    FROM
        images i
    JOIN
      users u ON i.author_id = u.id
    LEFT JOIN
      images_to_tags it ON i.id = it.image_id
    LEFT JOIN
      tags t ON it.tag_id = t.id
    LEFT JOIN
      images_analytics a ON a.image_id = i.id
    WHERE
      i.id = $1
    GROUP BY
      i.id, u.id, a.likes_count, a.views_count
  `

	err := r.db.QueryRowxContext(ctx, query, id).Scan(
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
	q := `UPDATE images SET
    title = COALESCE(NULLIF($1, ''), title),
    description = COALESCE(NULLIF($2, ''), description),
    access_level = COALESCE(NULLIF($3, '')::access_level, access_level)::access_level,
    expires_at = COALESCE($4, expires_at)
  WHERE id = $5 RETURNING *`

	img := new(domain.Image)
	rowx := r.db.QueryRowxContext(ctx, q,
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

func (r *imageRepository) States(ctx context.Context, imageID int, userID int) (*domain.ImageStates, error) {
	q := `
  WITH params AS (
    SELECT $1::int AS image_id, $2::int AS user_id
  )
  SELECT 
    EXISTS (
      SELECT 1
      FROM images_to_views v
      JOIN params p ON v.image_id = p.image_id AND v.user_id = p.user_id
    ) AS viewed,
    EXISTS (
      SELECT 1
      FROM images_to_likes l
      JOIN params p ON l.image_id = p.image_id AND l.user_id = p.user_id
    ) AS liked;
  `

	states := new(domain.ImageStates)
	rowx := r.db.QueryRowxContext(ctx, q, imageID, userID)
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
	q := `
  SELECT EXISTS (
    SELECT 1 FROM images_to_likes WHERE image_id = $1 AND user_id = $2
  )
  `

	var hasLike bool
	rowx := r.db.QueryRowxContext(ctx, q, imageID, userID)
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
