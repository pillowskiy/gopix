package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/internal/repository/postgres/pgutils"
	"github.com/pkg/errors"
)

var commentSortQuery = pgutils.NewSortQueryBuilder().
	AddField(string(domain.CommentNewestSort), pgutils.SortField{Field: "created_at", Order: pgutils.SortOrderDESC}).
	AddField(string(domain.CommentOldestSort), pgutils.SortField{Field: "created_at", Order: pgutils.SortOrderASC})

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
	imageID domain.ID,
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
  WHERE image_id = $1 AND parent_id IS NULL ORDER BY %s LIMIT $2 OFFSET $3
  `, sortQuery)

	limit := pagInput.PerPage
	rowx, err := repo.db.QueryxContext(ctx, q, imageID, limit, (pagInput.Page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetByImageID.QueryContext: %v", err)
	}
	defer rowx.Close()

	cmts, err := pgutils.ScanToStructSliceOf[domain.DetailedComment](rowx)
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
	commentID domain.ID,
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

func (repo *commentRepository) Delete(ctx context.Context, commentID domain.ID) error {
	q := `DELETE FROM comments WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, q, commentID)
	if err != nil {
		return fmt.Errorf("CommentRepository.Delete.ExecContext: %v", err)
	}

	return nil
}

func (repo *commentRepository) Update(
	ctx context.Context,
	commnetID domain.ID,
	comment *domain.Comment,
) (*domain.Comment, error) {
	q := `UPDATE comments SET comment = COALESCE(NULLIF($1, ''), comment) WHERE id = $2 RETURNING *`

	rowx := repo.db.QueryRowxContext(ctx, q, comment.Text, commnetID)

	cmt := new(domain.Comment)
	if err := rowx.StructScan(cmt); err != nil {
		return nil, fmt.Errorf("CommentRepository.Update.StructScan: %v", err)
	}

	return cmt, nil
}

func (repo *commentRepository) HasUserCommented(
	ctx context.Context,
	commentID domain.ID,
	userID domain.ID,
) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM comments WHERE image_id = $1 AND author_id = $2)`

	rowx := repo.db.QueryRowxContext(ctx, q, commentID, userID)

	var exists bool
	if err := rowx.Scan(&exists); err != nil {
		return false, fmt.Errorf("CommentRepository.IsCommentedByUser.Scan: %v", err)
	}

	return exists, nil
}

func (repo *commentRepository) GetReplies(ctx context.Context, commentID domain.ID, userID *domain.ID) ([]domain.DetailedComment, error) {
	q := `
  SELECT
    c.*,
    u.id AS "author.id",
    u.username AS "author.username",
    u.avatar_url AS "author.avatar_url",
    EXISTS(SELECT * FROM comments_to_likes WHERE comment_id = c.id AND user_id = $2) AS "stats.liked",
    COUNT(DISTINCT cl.user_id) AS "stats.likes"
  FROM comments c
  JOIN users u ON c.author_id = u.id
  LEFT JOIN comments_to_likes cl ON c.id = cl.comment_id
  WHERE parent_id = $1 GROUP BY c.id, u.id
  `

	rows, err := repo.db.QueryxContext(ctx, q, commentID, userID)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetReplies.QueryxContext: %v", err)
	}

	cmts, err := pgutils.ScanToStructSliceOf[domain.DetailedComment](rows)
	if err != nil {
		return nil, fmt.Errorf("CommentRepository.GetReplies.scanToStructSliceOf: %v", err)
	}

	return cmts, nil
}

func (repo *commentRepository) LikeComment(ctx context.Context, commentID domain.ID, userID domain.ID) error {
	q := `INSERT INTO comments_to_likes (comment_id, user_id) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, q, commentID, userID)
	return err
}

func (repo *commentRepository) UnlikeComment(ctx context.Context, commentID domain.ID, userID domain.ID) error {
	q := `DELETE FROM comments_to_likes WHERE comment_id = $1 AND user_id = $2`

	_, err := repo.db.ExecContext(ctx, q, commentID, userID)
	return err
}

func (repo *commentRepository) HasUserLikedComment(
	ctx context.Context,
	commentID domain.ID,
	userID domain.ID,
) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM comments_to_likes WHERE comment_id = $1 AND user_id = $2)`

	rowx := repo.db.QueryRowxContext(ctx, q, commentID, userID)

	var exists bool
	if err := rowx.Scan(&exists); err != nil {
		return false, fmt.Errorf("CommentRepository.HasUserLikedComment.Scan: %v", err)
	}

	return exists, nil
}
