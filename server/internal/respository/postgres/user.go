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

type userRepository struct {
	PostgresRepository
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		PostgresRepository: PostgresRepository{db},
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	q := `INSERT INTO users (username, email, external, password_hash, avatar_url) VALUES($1, $2, $3, $4, $5) RETURNING *`

	u := new(domain.User)
	rowx := r.ext(ctx).QueryRowxContext(
		ctx, q, user.Username, user.Email, user.External, user.PasswordHash, user.AvatarURL,
	)
	if err := rowx.StructScan(u); err != nil {
		return nil, fmt.Errorf("UserRepository.Register.StructScan: %v", err)
	}

	return u, nil
}

func (r *userRepository) GetUnique(ctx context.Context, user *domain.User) (*domain.User, error) {
	q := `SELECT * FROM users WHERE username = $1 OR email = $2`

	u := new(domain.User)
	rowx := r.ext(ctx).QueryRowxContext(ctx, q, user.Username, user.Email)
	if err := rowx.StructScan(u); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepository.GetUnique.StructScan: %v", err)
	}

	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	q := `SELECT * FROM users WHERE id = $1`

	u := new(domain.User)
	rowx := r.ext(ctx).QueryRowxContext(ctx, q, id)
	if err := rowx.StructScan(u); err != nil {
		return nil, fmt.Errorf("UserRepository.GetById.StructScan: %v", err)
	}

	return u, nil
}

func (r *userRepository) Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error) {
	q := `UPDATE users SET 
    username = COALESCE(NULLIF($2, ''), username),
    email = COALESCE(NULLIF($3, ''), email),
    avatar_url = COALESCE(NULLIF($4, ''), avatar_url)
  WHERE id = $1 RETURNING *`

	u := new(domain.User)
	rowx := r.ext(ctx).QueryRowxContext(
		ctx, q, id,
		user.Username,
		user.Email,
		user.AvatarURL,
	)
	if err := rowx.StructScan(u); err != nil {
		return nil, fmt.Errorf("UserRepository.Update.StructScan: %v", err)
	}

	return u, nil
}

func (r *userRepository) SetPermissions(ctx context.Context, id domain.ID, permissions int) error {
	q := `UPDATE users SET permissions = $1 WHERE id = $2`

	_, err := r.ext(ctx).ExecContext(ctx, q, permissions, id)
	if err != nil {
		return err
	}

	return nil
}
