package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
)

type followingRepository struct {
	db *sqlx.DB
}

func NewFollowingRepository(db *sqlx.DB) *followingRepository {
	return &followingRepository{db: db}
}

func (repo *followingRepository) Follow(ctx context.Context, userID domain.ID, executorID domain.ID) error {
	q := `INSERT INTO following (follower_id, followed_id) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, q, executorID, userID)
	return err
}

func (repo *followingRepository) Unfollow(ctx context.Context, userID domain.ID, executorID domain.ID) error {
	q := `DELETE FROM following WHERE follower_id = $1 AND followed_id = $2`

	_, err := repo.db.ExecContext(ctx, q, executorID, userID)
	return err
}

func (repo *followingRepository) IsFollowing(ctx context.Context, followerID, followingID domain.ID) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM following WHERE follower_id = $1 AND followed_id = $2)`

	isFollowing := false
	if err := repo.db.QueryRowxContext(ctx, q, followerID, followingID).Scan(&isFollowing); err != nil {
		return isFollowing, errors.Wrap(err, "FollowingRepository.IsFollowing")
	}

	return isFollowing, nil
}

func (repo *followingRepository) Stats(
	ctx context.Context, userID domain.ID, executorID *domain.ID,
) (*domain.FollowingStats, error) {
	stats := new(domain.FollowingStats)
	if executorID != nil {
		followingQuery := `SELECT EXISTS(SELECT 1 FROM following WHERE follower_id = $1 AND followed_id = $2)`
		_ = repo.db.QueryRowxContext(ctx, followingQuery, executorID, userID).Scan(&stats.IsFollowing)
	}

	countQuery := `
    SELECT
      COUNT(*) FILTER (WHERE followed_id = $1) AS following,
      COUNT(*) FILTER (WHERE follower_id = $1) AS followers
    FROM following
  `

	rowx := repo.db.QueryRowxContext(ctx, countQuery, userID)
	if err := rowx.StructScan(stats); err != nil {
		return nil, errors.Wrap(err, "FollowingRepository.Stats.StructScan")
	}

	return stats, nil
}
