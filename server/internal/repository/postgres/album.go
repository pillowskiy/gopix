package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pkg/errors"
)

type albumRepository struct {
	db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) *albumRepository {
	return &albumRepository{db: db}
}

func (repo *albumRepository) Create(ctx context.Context, album *domain.Album) (*domain.Album, error) {
	q := `INSERT INTO albums (name,description,author_id) VALUES ($1, $2, $3) RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, album.Name, album.Description, album.AuthorID)

	createdAlbum := new(domain.Album)
	if err := rowx.StructScan(createdAlbum); err != nil {
		return nil, errors.Wrap(err, "AlbumRepository.Create.StructScan")
	}

	return createdAlbum, nil
}

func (repo *albumRepository) GetByID(ctx context.Context, albumID domain.ID) (*domain.Album, error) {
	q := `SELECT * FROM albums WHERE id = $1`

	rowx := repo.db.QueryRowxContext(ctx, q, albumID)

	album := new(domain.Album)
	if err := rowx.StructScan(album); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "AlbumRepository.GetByID.StructScan")
	}

	return album, nil
}

func (repo *albumRepository) GetAlbumImages(
	ctx context.Context, albumID domain.ID, pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.ImageWithAuthor], error) {
	q := `
  SELECT
    i.*,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url"
  FROM images_to_albums ia
  JOIN images i ON i.id = ia.image_id AND i.access_level = 'public'::access_level
  JOIN users u ON u.id = i.author_id
  WHERE ia.album_id = $1
  LIMIT $2 OFFSET $3
  `

	rowx, err := repo.db.QueryxContext(ctx, q, albumID, pagInput.PerPage, (pagInput.Page-1)*pagInput.PerPage)
	if err != nil {
		return nil, errors.Wrap(err, "AlbumRepository.GetAlbumImages.QueryxContext")
	}

	images, err := scanToStructSliceOf[domain.ImageWithAuthor](rowx)
	if err != nil {
		return nil, errors.Wrap(err, "AlbumRepository.GetAlbumImages.scanToStructSliceOf")
	}

	pag := &domain.Pagination[domain.ImageWithAuthor]{
		PaginationInput: *pagInput,
		Items:           images,
	}

	countQuery := `SELECT COUNT(1) FROM images_to_albums WHERE album_id = $1`
	_ = repo.db.QueryRowxContext(ctx, countQuery, albumID).Scan(&pag.Total)

	return pag, nil
}

func (repo *albumRepository) GetByAuthorID(ctx context.Context, authorID domain.ID) ([]domain.DetailedAlbum, error) {
	q := `
  SELECT
    a.*,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url",
    (
      SELECT TO_JSON(
        COALESCE(ARRAY_AGG(to_jsonb(img)) FILTER (WHERE img.id IS NOT NULL), '{}')
      )
      FROM (
        SELECT i.*
        FROM images i
        INNER JOIN images_to_albums ita ON i.id = ita.image_id
        WHERE ita.album_id = a.id LIMIT 3
      ) AS img
    ) AS "cover"
  FROM albums a
  INNER JOIN users u ON a.author_id = u.id
  WHERE a.author_id = $1 GROUP BY a.id, u.id
  `

	rows, err := repo.db.QueryxContext(ctx, q, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "AlbumRepository.GetByAuthorID.QueryxContext")
	}

	defer rows.Close()
	var albums []domain.DetailedAlbum
	for rows.Next() {
		var row domain.DetailedAlbum
		var rowCoverJSON []byte

		if err := rows.Scan(
			&row.ID,
			&row.AuthorID,
			&row.Name,
			&row.Description,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Author.ID,
			&row.Author.Username,
			&row.Author.AvatarURL,
			&rowCoverJSON,
		); err != nil {
			return nil, errors.Wrap(err, "AlbumRepository.GetByAuthorID.Scan")
		}

		row.Cover = []domain.Image{}
		if len(rowCoverJSON) > 0 {
			if err := json.Unmarshal(rowCoverJSON, &row.Cover); err != nil {
				return nil, errors.Wrap(err, "imageRepository.GetDetailed.Unmarshal")
			}
		}

		albums = append(albums, row)
	}

	return albums, nil
}

func (repo *albumRepository) Delete(ctx context.Context, albumID domain.ID) error {
	q := `DELETE FROM albums WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, q, albumID)
	return errors.Wrap(err, "AlbumRepository.Delete.ExecContext")
}

func (repo *albumRepository) Update(
	ctx context.Context,
	albumID domain.ID,
	album *domain.Album,
) (*domain.Album, error) {
	q := `
  UPDATE albums SET 
    name = COALESCE(NULLIF($1, ''), name),
    description = COALESCE(NULLIF($2, ''), description)
  WHERE id = $3 RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, album.Name, album.Description, albumID)

	updatedAlbum := new(domain.Album)
	if err := rowx.StructScan(updatedAlbum); err != nil {
		return nil, errors.Wrap(err, "AlbumRepository.Update.StructScan")
	}

	return updatedAlbum, nil
}

func (repo *albumRepository) PutImage(
	ctx context.Context,
	albumID domain.ID,
	imageID domain.ID,
) error {
	q := `INSERT INTO images_to_albums (album_id, image_id) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, q, albumID, imageID)
	if err != nil {
		return errors.Wrap(err, "AlbumRepository.PutImage.ExecContext")
	}

	return nil
}

func (repo *albumRepository) DeleteImage(
	ctx context.Context,
	albumID domain.ID,
	imageID domain.ID,
) error {
	q := `DELETE FROM album_to_images WHERE album_id = $1 AND image_id = $2`

	_, err := repo.db.ExecContext(ctx, q, albumID, imageID)
	if err != nil {
		return errors.Wrap(err, "AlbumRepository.DeleteImage.ExecContext")
	}

	return nil
}
