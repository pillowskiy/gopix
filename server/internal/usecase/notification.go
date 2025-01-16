package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pkg/errors"
)

type NotificationRepository interface {
	Push(ctx context.Context, recieverID domain.ID, notif *domain.Notification) (*domain.Notification, error)
	Stats(ctx context.Context, userID domain.ID) (*domain.NotificationStats, error)
	GetNotifications(
		ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
	) (*domain.Pagination[domain.Notification], error)
	MarkAsRead(ctx context.Context, notifs []domain.ID) error
}

type NotificationSignal interface {
	Subscribe(id string) (<-chan *domain.Notification, func())
	Publish(id string, item *domain.Notification) error
}

type notifactionUseCase struct {
	repo    NotificationRepository
	signal  NotificationSignal
	logger  logger.Logger
	waitDur time.Duration
}

func NewNotificationUseCase(repo NotificationRepository, signal NotificationSignal, logger logger.Logger) *notifactionUseCase {
	return &notifactionUseCase{repo: repo, signal: signal, logger: logger, waitDur: 2 * time.Minute}
}

func (uc *notifactionUseCase) Notify(ctx context.Context, userID domain.ID, notif *domain.Notification) error {
	createdNotif, err := uc.repo.Push(ctx, userID, notif)
	if err != nil {
		return fmt.Errorf("NotificationUseCase.Notify.Push: %w", err)
	}

	subKey := uc.subKey(userID)
	if err := uc.signal.Publish(subKey, createdNotif); err != nil {
		uc.logger.Errorf("Failed to push notification for sub %s: %v", subKey, err)
	}
	return nil
}

func (uc *notifactionUseCase) WaitForNotification(ctx context.Context, userID domain.ID) (*domain.Notification, error) {
	subKey := uc.subKey(userID)
	ch, cancel := uc.signal.Subscribe(subKey)
	defer cancel()

	ctx, ctxCancel := context.WithTimeout(ctx, uc.waitDur)
	defer ctxCancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case notif := <-ch:
		return notif, nil
	}
}

func (uc *notifactionUseCase) GetNotifications(
	ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.Notification], error) {
	return uc.getAndReadNotifications(ctx, userID, pagInput)
}

func (uc *notifactionUseCase) GetStats(ctx context.Context, userID domain.ID) (*domain.NotificationStats, error) {
	return uc.getStats(ctx, userID)
}

func (uc *notifactionUseCase) getStats(ctx context.Context, userID domain.ID) (*domain.NotificationStats, error) {
	stats, err := uc.repo.Stats(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return stats, nil
}

func (uc *notifactionUseCase) getAndReadNotifications(
	ctx context.Context, userID domain.ID, pagInput *domain.PaginationInput,
) (*domain.Pagination[domain.Notification], error) {
	notifs, err := uc.repo.GetNotifications(ctx, userID, pagInput)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	unreadNotifs := make([]domain.ID, 0, len(notifs.Items))
	for _, notif := range notifs.Items {
		if notif.Read == false {
			unreadNotifs = append(unreadNotifs, notif.ID)
		}
	}

	if err := uc.repo.MarkAsRead(ctx, unreadNotifs); err != nil {
		uc.logger.Errorf("Failed to mark as read notifications %+v: %v", unreadNotifs, err)
	}

	return notifs, nil
}

func (uc *notifactionUseCase) subKey(userID domain.ID) string {
	return fmt.Sprintf("user-%s", userID)
}
