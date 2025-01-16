package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/internal/repository/postgres/pgutils"
	"github.com/pkg/errors"
)

type notificationRepository struct {
	PostgresRepository
}

func NewNotificationRepository(db *sqlx.DB) *notificationRepository {
	return &notificationRepository{
		PostgresRepository: PostgresRepository{db},
	}
}

func (repo *notificationRepository) Push(ctx context.Context, recieverID domain.ID, notif *domain.Notification) (*domain.Notification, error) {
	const q = `INSERT INTO notifications (user_id, title, message, hidden) VALUES($1, $2, $3, $4) RETURNING *`

	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, recieverID, notif.UserID, notif.Message, notif.Hidden)

	createdNotif := new(domain.Notification)
	if err := rowx.StructScan(createdNotif); err != nil {
		return nil, fmt.Errorf("NotificationRepository.Push.ExecContext: %w", err)
	}
	return createdNotif, nil
}

func (repo *notificationRepository) Stats(ctx context.Context, userID domain.ID) (*domain.NotificationStats, error) {
	const q = `
        SELECT COUNT(*) FILTER (WHERE read = false) AS unread
        FROM notifications
        WHERE user_id = $1;
    `

	stats := new(domain.NotificationStats)
	rowx := repo.ext(ctx).QueryRowxContext(ctx, q, userID)
	if err := rowx.StructScan(stats); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("NotificationRepository.Stats.StructScan: %w", err)
	}

	return stats, nil
}

func (repo *notificationRepository) GetNotifications(
	ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.Notification], error) {
	const q = `SELECT * FROM notifications WHERE user_id = $1`

	rows, err := repo.ext(ctx).QueryxContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("NotificationRepository.Stats.QueryxContext: %w", err)
	}

	notifs, err := pgutils.ScanToStructSliceOf[domain.Notification](rows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("NotificationRepository.Stats.StructScan: %w", err)
	}

	pag := &domain.Pagination[domain.Notification]{
		Items: notifs,
	}

	return pag, nil
}

func (repo *notificationRepository) MarkAsRead(ctx context.Context, notifs []domain.ID) error {
	const baseQuery = `UPDATE notifications SET read = true WHERE id IN (%s)`

	placeholders := make([]string, len(notifs))
	args := make([]interface{}, len(notifs))
	for i, id := range notifs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	q := fmt.Sprintf(baseQuery, strings.Join(placeholders, ","))
	_, err := repo.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("NotificationRepository.MarkAsRead.ExecContext: %w", err)
	}

	return nil
}
