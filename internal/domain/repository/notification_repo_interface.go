package repository

import (
	"challenge-app/internal/domain/model"
	"context"
	"time"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	ListForUser(ctx context.Context, userID uint, since time.Time, limit int, cursor *time.Time) (items []model.Notification, nextCursor *time.Time, err error)
	MarkRead(ctx context.Context, userID uint, notificationID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint, since time.Time) (int64, error)
	DeleteOlderThan(ctx context.Context, t time.Time) (int64, error)
}
