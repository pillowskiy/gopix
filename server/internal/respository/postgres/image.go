package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
)

type imageRepository struct {
	db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) *imageRepository {
	return &imageRepository{db: db}
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
    COALESCE(l.likes_count, 0) AS likes,
    COALESCE(v.views_count, 0) AS views,
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
      (SELECT image_id, COUNT(*) AS likes_count FROM images_to_likes GROUP BY image_id) l ON i.id = l.image_id
    LEFT JOIN
      (SELECT image_id, COUNT(*) AS views_count FROM images_to_views GROUP BY image_id) v ON i.id = v.image_id
    LEFT JOIN
      images_to_tags it ON i.id = it.image_id
    LEFT JOIN
      tags t ON it.tag_id = t.id
    WHERE
      i.id = $1
    GROUP BY
      i.id, u.id, l.likes_count, v.views_count;
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

	// TEMP: We shouldn't depend on JSON, probably we can use annonymous struct
	err = json.Unmarshal(tagsJSON, &detailedImage.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.GetDetailed.Unmarshal")
	}

	return &detailedImage, nil
}
