package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pillowskiy/gopix/pkg/batch"
	"github.com/pkg/errors"
)

var imageBatchConfig = batch.BatchConfig{Retries: 3, MaxSize: 10000}

type imageRepository struct {
	db          *sqlx.DB
	viewBatcher batch.Batcher[domain.ImageView]
}

func NewImageRepository(db *sqlx.DB) *imageRepository {
	repo := &imageRepository{db: db}

	repo.viewBatcher = batch.NewWithConfig(&imageBatchConfig, repo.batchViews)
	go repo.viewBatcher.Ticker(time.Minute * 5)

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

func (r *imageRepository) AddView(ctx context.Context, view *domain.ImageView) error {
	r.viewBatcher.Add(*view)
	return nil
}

func (r *imageRepository) batchViews(views []domain.ImageView) error {
	if len(views) == 0 {
		return nil
	}

	start := time.Now()
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.BeginTx")
	}
	defer tx.Rollback()

	query := `INSERT INTO images_to_views (image_id, user_id) VALUES (:image_id, :user_id) ON CONFLICT DO NOTHING`
	if _, err := tx.NamedExecContext(ctx, query, views); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.ExecContext")
	}

	analyticsQuery, analyticsParams := r.namedAnalyticsQuery("views_count", views)
	if _, err := tx.ExecContext(ctx, analyticsQuery, analyticsParams...); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.AnalyticsExecContext")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.Commit")
	}

	log.Println("imageRepository.batchViews", time.Since(start))
	return nil
}

func (r *imageRepository) namedAnalyticsQuery(
	column string,
	items []domain.ImageView,
) (string, []any) {
	counter := make(map[int]int)
	for _, item := range items {
		counter[item.ImageID]++
	}

	p := make([]string, 0, len(items)*2)
	var v []any
	for imageID, count := range counter {
		p = append(p, fmt.Sprintf("($%d::int, $%d::int)", len(p)+1, len(p)+2))
		v = append(v, imageID, count)
	}

	q := fmt.Sprintf(`
    UPDATE images_analytics AS ia SET views_count = ia.views_count + p.views_count
    FROM (VALUES %s) AS p(image_id, views_count) WHERE ia.image_id = p.image_id
  `, strings.Join(p, ","))

	return q, v
}
