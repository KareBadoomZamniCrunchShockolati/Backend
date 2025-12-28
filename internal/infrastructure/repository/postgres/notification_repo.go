package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, n *model.Notification) error {
	var dataBytes []byte
	if n.Data != nil {
		b, err := json.Marshal(n.Data)
		if err != nil {
			return err
		}
		dataBytes = b
	}

	e := entity.NotificationEntity{
		UserID:    n.UserID,
		Type:      string(n.Type),
		Title:     n.Title,
		Body:      n.Body,
		Data:      dataBytes,
		CreatedAt: time.Now().UTC(),
		ReadAt:    n.ReadAt,
	}

	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return err
	}

	n.ID = e.ID
	n.CreatedAt = e.CreatedAt
	return nil
}

func (r *NotificationRepo) ListForUser(ctx context.Context, userID uint, since time.Time, limit int, cursor *time.Time) ([]model.Notification, *time.Time, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	q := r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Order("created_at DESC").
		Limit(limit)

	if cursor != nil {
		q = q.Where("created_at < ?", *cursor)
	}

	var rows []entity.NotificationEntity
	if err := q.Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	out := make([]model.Notification, 0, len(rows))
	for _, row := range rows {
		var data map[string]any
		if len(row.Data) > 0 {
			_ = json.Unmarshal(row.Data, &data) // best-effort
		}
		out = append(out, model.Notification{
			ID:        row.ID,
			UserID:    row.UserID,
			Type:      model.NotificationType(row.Type),
			Title:     row.Title,
			Body:      row.Body,
			Data:      data,
			CreatedAt: row.CreatedAt,
			ReadAt:    row.ReadAt,
		})
	}

	var nextCursor *time.Time
	if len(rows) == limit {
		t := rows[len(rows)-1].CreatedAt
		nextCursor = &t
	}

	return out, nextCursor, nil
}

func (r *NotificationRepo) MarkRead(ctx context.Context, userID uint, notificationID uint) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", notificationID, userID).
		Update("read_at", &now)

	return res.Error
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID uint) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", &now)

	return res.Error
}

func (r *NotificationRepo) UnreadCount(ctx context.Context, userID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("user_id = ? AND created_at >= ? AND read_at IS NULL", userID, since).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepo) DeleteOlderThan(ctx context.Context, t time.Time) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("created_at < ?", t).
		Delete(&entity.NotificationEntity{})
	return res.RowsAffected, res.Error
}
