package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	repository_interface "challenge-app/internal/domain/repository"

	"github.com/go-redis/redis/v8"
)

type TempUploadRepository struct {
	rdb *redis.Client
}

var _ repository_interface.TempUploadRepository = (*TempUploadRepository)(nil)

func NewTempUploadRepository(rdb *redis.Client) *TempUploadRepository {
	return &TempUploadRepository{rdb: rdb}
}

func tempKey(key string) string { return "tempimg:" + key }
func indexKey() string         { return "tempimg:index" }

func (t *TempUploadRepository) Track(ctx context.Context, key string, userID uint, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	if err := t.rdb.Set(ctx, tempKey(key), fmt.Sprintf("%d", userID), ttl).Err(); err != nil {
		return err
	}
	return t.rdb.ZAdd(ctx, indexKey(), &redis.Z{
		Score:  float64(expiresAt.Unix()),
		Member: key,
	}).Err()
}

func (t *TempUploadRepository) VerifyOwnership(ctx context.Context, key string, userID uint) (bool, error) {
	val, err := t.rdb.Get(ctx, tempKey(key)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	stored, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, nil
	}

	return uint(stored) == userID, nil
}

func (t *TempUploadRepository) Untrack(ctx context.Context, key string) error {
	_ = t.rdb.Del(ctx, tempKey(key)).Err()
	return t.rdb.ZRem(ctx, indexKey(), key).Err()
}

func (t *TempUploadRepository) ListExpired(ctx context.Context, now time.Time, limit int) ([]string, error) {
	res, err := t.rdb.ZRangeByScore(ctx, indexKey(), &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", now.Unix()),
		Offset: 0,
		Count:  int64(limit),
	}).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *TempUploadRepository) RemoveExpiredIndex(ctx context.Context, key string) error {
	return t.rdb.ZRem(ctx, indexKey(), key).Err()
}
