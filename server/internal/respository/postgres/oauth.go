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

type oauthRepository struct {
	PostgresRepository
}

func NewOAuthRepository(db *sqlx.DB) *oauthRepository {
	return &oauthRepository{
		PostgresRepository: PostgresRepository{db},
	}
}

func (repo *oauthRepository) GetByOAuthID(ctx context.Context, oauthID string) (*domain.OAuth, error) {
	q := `SELECT * FROM oauth WHERE oauth_id = $1`

	var oauth domain.OAuth
	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, oauthID)
	if err := rowx.StructScan(&oauth); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("OAuthRepository.GetByUserID.StructScan: %v", err)
	}

	return &oauth, nil
}

func (repo *oauthRepository) Create(ctx context.Context, oauth *domain.OAuth) error {
	q := `INSERT INTO oauth (user_id, oauth_id, service) VALUES ($1, $2, $3)`

	_, err := repo.ext(ctx).ExecContext(ctx, q, oauth.UserID, oauth.OAuthID, oauth.Service)
	return err
}
