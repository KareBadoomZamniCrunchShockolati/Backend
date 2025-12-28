package service

import (
	"challenge-app/internal/domain/repository"
	"context"
	"log"
	"time"
)

type NotificationCleanup struct {
	repo repository.NotificationRepository
}

func NewNotificationCleanup(repo repository.NotificationRepository) *NotificationCleanup {
	return &NotificationCleanup{repo: repo}
}

func (c *NotificationCleanup) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// run once at startup (optional)
	c.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *NotificationCleanup) cleanup(ctx context.Context) {
	cutoff := time.Now().AddDate(0, 0, -30).UTC()
	n, err := c.repo.DeleteOlderThan(ctx, cutoff)
	if err != nil {
		log.Println("notification cleanup error:", err)
		return
	}
	log.Println("notification cleanup deleted:", n)
}
