package service

import (
	"context"
	"time"

	appsvc "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
)

type NotificationServiceImpl struct {
	repo     repository.NotificationRepository
	realtime appsvc.Notifier
}

func NewNotificationService(repo repository.NotificationRepository, realtime appsvc.Notifier) *NotificationServiceImpl {
	return &NotificationServiceImpl{repo: repo, realtime: realtime}
}

func (s *NotificationServiceImpl) CreateAndPush(ctx context.Context, n model.Notification) error {
	if err := s.repo.Create(ctx, &n); err != nil {
		return err
	}
	_ = s.realtime.Push(ctx, n.UserID, n) // best effort
	return nil
}

func (s *NotificationServiceImpl) ListRecent(ctx context.Context, userID uint, limit int, cursor *time.Time) ([]model.Notification, *time.Time, error) {
	since := time.Now().AddDate(0, 0, -30).UTC()
	return s.repo.ListForUser(ctx, userID, since, limit, cursor)
}

func (s *NotificationServiceImpl) MarkRead(ctx context.Context, userID uint, notificationID uint) error {
	return s.repo.MarkRead(ctx, userID, notificationID)
}

func (s *NotificationServiceImpl) MarkAllRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *NotificationServiceImpl) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	since := time.Now().AddDate(0, 0, -30).UTC()
	return s.repo.UnreadCount(ctx, userID, since)
}
