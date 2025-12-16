package repository

import (
	"context"
	"time"
)

type TempUploadRepository interface {
	Track(ctx context.Context, key string, userID uint, expiresAt time.Time) error
	VerifyOwnership(ctx context.Context, key string, userID uint) (bool, error)
	Untrack(ctx context.Context, key string) error
	ListExpired(ctx context.Context, now time.Time, limit int) ([]string, error)
	RemoveExpiredIndex(ctx context.Context, key string) error
}
