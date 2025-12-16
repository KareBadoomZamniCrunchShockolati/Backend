package storage

import (
	"context"
	"io"
	"time"
)

type PresignedUpload struct {
	Key       string
	UploadURL string
	Headers   map[string]string
	ExpiresAt time.Time
}

type ObjectStorage interface {
	// Server upload (used for profile/cover uploads)
	Upload(ctx context.Context, key string, contentType string, body io.Reader) (publicURL string, err error)

	// Client direct upload
	PresignPut(ctx context.Context, key string, contentType string, expiresIn time.Duration) (*PresignedUpload, error)
	Copy(ctx context.Context, srcKey, dstKey string) (publicURL string, err error)

	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}
