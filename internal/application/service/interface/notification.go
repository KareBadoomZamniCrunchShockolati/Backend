package serviceinterface

import (
	"context"
	"time"

	"challenge-app/internal/domain/model"
)

type NotificationService interface {
	CreateAndPush(ctx context.Context, n model.Notification) error
	ListRecent(ctx context.Context, userID uint, limit int, cursor *time.Time) ([]model.Notification, *time.Time, error)
	MarkRead(ctx context.Context, userID uint, notificationID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint) (int64, error)
}
