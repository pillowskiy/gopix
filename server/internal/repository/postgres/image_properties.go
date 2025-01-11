package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
)

type imagePropsRepository struct {
	PostgresRepository
}

func NewImagePropsRepository(db *sqlx.DB) *imagePropsRepository {
	return &imagePropsRepository{
		PostgresRepository: PostgresRepository{db},
	}
}

func (repo *imagePropsRepository) Create(
	ctx context.Context, imageID domain.ID, props *domain.ImageProperties,
) error {
	const q = `INSERT INTO image_properties (image_id, mime, ext, height, width) VALUES ($1, $2, $3, $4, $5)`

	_, err := repo.ext(ctx).ExecContext(ctx, q, imageID, props.Mime, props.Ext, props.Height, props.Width)
	if err != nil {
		return errors.Wrap(err, "ImagePropertiesRepository.Create.StructScan")
	}

	return nil
}

func (repo *imagePropsRepository) Delete(ctx context.Context, imageID domain.ID) error {
	const q = `DELETE FROM image_properties WHERE image_id = $1`

	if _, err := repo.ext(ctx).ExecContext(ctx, q, imageID); err != nil {
		return err
	}

	return nil
}

func (repo *imagePropsRepository) Properties(ctx context.Context, imageID domain.ID) (*domain.ImageProperties, error) {
	const q = `SELECT * FROM image_properties WHERE image_id = $1`

	props := new(domain.ImageProperties)
	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, imageID)

	if err := rowx.StructScan(props); err != nil {
		return nil, errors.Wrap(err, "ImageRepository.Create.StructScan")
	}

	return props, nil
}
