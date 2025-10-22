package storage

import "context"

// File represents a file stored or to be uploaded.
type File struct {
	ID   string
	Name string
	Data []byte
}

type Storage interface {
	Upload(ctx context.Context, file File) (string, error)
	Download(ctx context.Context, fileID string) (File, error)
	Delete(ctx context.Context, fileID string) error
}
