package postgres

import (
	"context"
	"encoding/json"
	"time"

	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type NotificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepoImpl{db: db}
}

func (r *NotificationRepoImpl) Create(ctx context.Context, n *model.Notification) error {
	var dataBytes []byte
	if n.Data != nil {
		b, err := json.Marshal(n.Data)
		if err != nil {
			return err
		}
		dataBytes = b
	}

	now := time.Now().UTC()

	e := entity.NotificationEntity{
		UserID:    n.UserID,
		Type:      string(n.Type),
		TitleKey:  n.TitleKey,
		BodyKey:   n.BodyKey,
		Data:      dataBytes,
		CreatedAt: now,
		ReadAt:    n.ReadAt,
	}

	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return err
	}

	n.ID = e.ID
	n.CreatedAt = e.CreatedAt
	return nil
}

func (r *NotificationRepoImpl) ListForUser(
	ctx context.Context,
	userID uint,
	since time.Time,
	limit int,
	cursor *time.Time,
) ([]model.Notification, *time.Time, error) {

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
			_ = json.Unmarshal(row.Data, &data)
		}

		out = append(out, model.Notification{
			ID:        row.ID,
			UserID:    row.UserID,
			Type:      model.NotificationType(row.Type),
			TitleKey:  row.TitleKey,
			BodyKey:   row.BodyKey,
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

func (r *NotificationRepoImpl) MarkRead(ctx context.Context, userID uint, notificationID uint) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", notificationID, userID).
		Update("read_at", &now).Error
}

func (r *NotificationRepoImpl) MarkAllRead(ctx context.Context, userID uint) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", &now).Error
}

func (r *NotificationRepoImpl) UnreadCount(ctx context.Context, userID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.NotificationEntity{}).
		Where("user_id = ? AND created_at >= ? AND read_at IS NULL", userID, since).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepoImpl) DeleteOlderThan(ctx context.Context, t time.Time) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("created_at < ?", t).
		Delete(&entity.NotificationEntity{})
	return res.RowsAffected, res.Error
}
