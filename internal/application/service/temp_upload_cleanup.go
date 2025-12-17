package service

import (
	"context"
	"time"

	repository_interface "challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/storage"
)

type TempUploadCleaner struct {
	repo    repository_interface.TempUploadRepository
	storage storage.ObjectStorage
}

func NewTempUploadCleaner(repo repository_interface.TempUploadRepository, storage storage.ObjectStorage) *TempUploadCleaner {
	return &TempUploadCleaner{repo: repo, storage: storage}
}

func (c *TempUploadCleaner) RunOnce(ctx context.Context, limit int) error {
	keys, err := c.repo.ListExpired(ctx, time.Now(), limit)
	if err != nil {
		return err
	}
	for _, key := range keys {
		_ = c.storage.Delete(ctx, key)       
		_ = c.repo.RemoveExpiredIndex(ctx, key)
		_ = c.repo.Untrack(ctx, key)
	}
	return nil
}
