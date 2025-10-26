package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type VerificationRepository struct {
	client *redis.Client
}

func NewVerificationRepository(client *redis.Client) *VerificationRepository {
	return &VerificationRepository{client: client}
}

func (r *VerificationRepository) StoreVerificationCode(ctx context.Context, email, code string, expiresInMinutes int) error {
	key := fmt.Sprintf("verify:%s", email)
	return r.client.Set(ctx, key, code, time.Duration(expiresInMinutes)*time.Minute).Err()
}

func (r *VerificationRepository) GetVerificationCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("verify:%s", email)
	return r.client.Get(ctx, key).Result()
}

func (r *VerificationRepository) DeleteVerificationCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("verify:%s", email)
	return r.client.Del(ctx, key).Err()
}
