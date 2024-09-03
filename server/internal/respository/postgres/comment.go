package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
)

var commentSortQuery = NewSortQueryBuilder().
	AddField(string(domain.CommentNewestSort), SortField{Field: "created_at", Order: sortOrderDESC}).
	AddField(string(domain.CommentOldestSort), SortField{Field: "created_at", Order: sortOrderASC})

type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (repo *commentRepository) Create(
	ctx context.Context,
	comment *domain.Comment,
) (*domain.Comment, error) {
	q := `INSERT INTO comments (image_id, author_id, comment) VALUES($1, $2, $3) RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, comment.ImageID, comment.AuthorID, comment.Text)

	cmt := new(domain.Comment)
	if err := rowx.StructScan(cmt); err != nil {
		return nil, fmt.Errorf("CommentRepository.Create.StructScan: %v", err)
	}

	return cmt, nil
}

func (repo *commentRepository) GetByImageID(
	ctx context.Context,
	imageID int,
	pagInput *domain.PaginationInput,
	sort domain.CommentSortMethod,
) (*domain.Pagination[domain.DetailedComment], error) {
	sortQuery, ok := commentSortQuery.SortQuery(string(sort))
	if !ok {
		return nil, repository.ErrIncorrectInput
	}

	q := fmt.Sprintf(`
  SELECT
    c.*,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url"
  FROM comments c
  JOIN users u ON c.author_id = u.id
  WHERE image_id = $1 ORDER BY %s LIMIT $2 OFFSET $3
  `, sortQuery)

	limit := pagInput.PerPage
	rowx, err := repo.db.QueryxContext(ctx, q, imageID, limit, (pagInput.Page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetByImageID.QueryContext: %v", err)
	}
	defer rowx.Close()

	cmts, err := scanToStructSliceOf[domain.DetailedComment](rowx)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetByImageID.scanToStructSliceOf: %v", err)
	}

	pagination := &domain.Pagination[domain.DetailedComment]{
		PaginationInput: *pagInput,
		Items:           cmts,
	}

	countQuery := `SELECT COUNT(1) FROM comments WHERE image_id = $1`
	_ = repo.db.QueryRowxContext(ctx, countQuery, imageID).Scan(&pagination.Total)

	return pagination, nil
}

func (repo *commentRepository) GetByID(
	ctx context.Context,
	commentID int,
) (*domain.Comment, error) {
	q := `SELECT * FROM comments WHERE id = $1`

	rowx := repo.db.QueryRowxContext(ctx, q, commentID)

	cmt := new(domain.Comment)
	if err := rowx.StructScan(cmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("CommentRepository.GetByID.StructScan: %v", err)
	}

	return cmt, nil
}

func (repo *commentRepository) Delete(ctx context.Context, commentID int) error {
	q := `DELETE FROM comments WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, q, commentID)
	if err != nil {
		return fmt.Errorf("CommentRepository.Delete.ExecContext: %v", err)
	}

	return nil
}

func (repo *commentRepository) Update(
	ctx context.Context,
	commnetID int,
	comment *domain.Comment,
) (*domain.Comment, error) {
	q := `UPDATE comments SET comment = $1 WHERE id = $2 RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, comment.Text, commnetID)

	cmt := new(domain.Comment)
	if err := rowx.StructScan(cmt); err != nil {
		return nil, fmt.Errorf("CommentRepository.Update.StructScan: %v", err)
	}

	return cmt, nil
}

func (repo *commentRepository) HasUserCommented(
	ctx context.Context,
	commentID int,
	userID int,
) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM comments WHERE image_id = $1 AND author_id = $2)`

	rowx := repo.db.QueryRowxContext(ctx, q, commentID, userID)

	var exists bool
	if err := rowx.Scan(&exists); err != nil {
		return false, fmt.Errorf("CommentRepository.IsCommentedByUser.Scan: %v", err)
	}

	return exists, nil
}
