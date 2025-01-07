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
) (*domain.ImageProperties, error) {
	const q = `
  INSERT INTO images_properties (image_id, mime, ext, height, width)
  VALUES ($1, $2, $3, $4, $5) RETURNING *`

	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, imageID, props.Mime, props.Ext, props.Height, props.Width)

	imgProps := new(domain.ImageProperties)
	if err := rowx.StructScan(imgProps); err != nil {
		return nil, errors.Wrap(err, "ImagePropertiesRepository.Create.StructScan")
	}

	return imgProps, nil
}
