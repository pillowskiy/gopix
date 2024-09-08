package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
)

type tagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (repo *tagRepository) Upsert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	q := `INSERT INTO tags (name) VALUES($1) ON CONFLICT (name) DO NOTHING RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, tag.Name)

	createdTag := new(domain.Tag)
	if err := rowx.StructScan(createdTag); err != nil {
		return nil, errors.Wrap(err, "TagRepository.Upsert.StructScan")
	}

	return createdTag, nil
}

func (repo *tagRepository) UpsertImageTags(ctx context.Context, tag *domain.Tag, imageID int) error {
	upsertQuery := `INSERT INTO tags(name) VALUES($1) ON CONFLICT (name) DO NOTHING RETURNING *`
	getByNameQuery := `SELECT * FROM tags WHERE name = $1`
	relationQuery := `INSERT INTO images_to_tags(image_id, tag_id) VALUES($1, $2)`

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "TagRepository.UpsertImageTags.Begin")
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, upsertQuery, tag.Name); err != nil {
		return errors.Wrap(err, "TagRepository.UpsertImageTags.StructScan")
	}

	rowx := tx.QueryRowxContext(ctx, getByNameQuery, tag.Name)
	upTag := new(domain.Tag)
	if err := rowx.StructScan(upTag); err != nil {
		return errors.Wrap(err, "TagRepository.UpsertImageTags.StructScan")
	}

	if _, err := tx.ExecContext(ctx, relationQuery, imageID, upTag.ID); err != nil {
		return errors.Wrap(err, "TagRepository.UpsertImageTags.StructScan")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "TagRepository.UpsertImageTags.Commit")
	}

	return nil
}

func (repo *tagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	q := `SELECT * FROM tags WHERE name = $1`

	rowx := repo.db.QueryRowxContext(ctx, q, name)

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

	rows, err := repo.db.QueryxContext(ctx, q, name)
	if err != nil {
		return nil, errors.Wrap(err, "TagRepository.Search.QueryxContext")
	}

	tags, err := scanToStructSliceOf[domain.Tag](rows)
	if err != nil {
		return nil, errors.Wrap(err, "TagRepository.Search.scanToStructSliceOf")
	}

	return tags, nil
}

func (repo *tagRepository) GetByID(ctx context.Context, id int) (*domain.Tag, error) {
	q := `SELECT * FROM tags WHERE id = $1`

	rowx := repo.db.QueryRowxContext(ctx, q, id)

	tag := new(domain.Tag)
	if err := rowx.StructScan(tag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "TagRepository.GetByID.StructScan")
	}

	return tag, nil
}

func (repo *tagRepository) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM tags WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "TagRepository.Delete.ExecContext")
	}

	return nil
}
