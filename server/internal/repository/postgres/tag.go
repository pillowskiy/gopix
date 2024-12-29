package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/internal/repository/postgres/pgutils"
	"github.com/pkg/errors"
)

type tagRepository struct {
	PostgresRepository
}

func NewTagRepository(db *sqlx.DB) *tagRepository {
	return &tagRepository{
		PostgresRepository: PostgresRepository{db},
	}
}

func (repo *tagRepository) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	q := `INSERT INTO tags (name) VALUES($1) ON CONFLICT (name) DO NOTHING RETURNING *`

	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, tag.Name)

	createdTag := new(domain.Tag)
	if err := rowx.StructScan(createdTag); err != nil {
		return nil, errors.Wrap(err, "TagRepository.Upsert.StructScan")
	}

	return createdTag, nil
}

func (repo *tagRepository) UpsertImageTags(
	ctx context.Context, tag *domain.Tag, imageID domain.ID,
) error {
	relationQuery := `INSERT INTO images_to_tags(image_id, tag_id) VALUES($1, $2) ON CONFLICT DO NOTHING`

	err := repo.DoInTransaction(ctx, func(ctx context.Context) error {
		upTag, err := repo.GetByName(ctx, tag.Name)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				upTag, err = repo.Create(ctx, tag)
				if err != nil {
					return err
				}
			default:
				return err
			}
		}

		if _, err := repo.ext(ctx).ExecContext(ctx, relationQuery, imageID, upTag.ID); err != nil {
			return errors.Wrap(err, "TagRepository.UpsertImageTags.StructScan")
		}

		return nil
	})

	return err
}

func (repo *tagRepository) DeleteImageTag(ctx context.Context, tagID, imageID domain.ID) error {
	q := `DELETE FROM images_to_tags WHERE image_id = $1 AND tag_id = $2`

	_, err := repo.ext(ctx).ExecContext(ctx, q, imageID, tagID)
	return err
}

func (repo *tagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	q := `SELECT * FROM tags WHERE name = $1`

	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, name)

	tag := new(domain.Tag)
	if err := rowx.StructScan(tag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "TagRepository.GetByName.StructScan")
	}

	return tag, nil
}

func (repo *tagRepository) Search(ctx context.Context, name string) ([]domain.Tag, error) {
	q := `SELECT * FROM tags WHERE name LIKE $1 LIMIT 10`

	rows, err := repo.ext(ctx).QueryxContext(ctx, q, name)
	if err != nil {
		return nil, errors.Wrap(err, "TagRepository.Search.QueryxContext")
	}

	tags, err := pgutils.ScanToStructSliceOf[domain.Tag](rows)
	if err != nil {
		return nil, errors.Wrap(err, "TagRepository.Search.scanToStructSliceOf")
	}

	return tags, nil
}

func (repo *tagRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Tag, error) {
	q := `SELECT * FROM tags WHERE id = $1`

	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, id)

	tag := new(domain.Tag)
	if err := rowx.StructScan(tag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "TagRepository.GetByID.StructScan")
	}

	return tag, nil
}

func (repo *tagRepository) Delete(ctx context.Context, id domain.ID) error {
	q := `DELETE FROM tags WHERE id = $1`

	_, err := repo.ext(ctx).ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "TagRepository.Delete.ExecContext")
	}

	return nil
}
