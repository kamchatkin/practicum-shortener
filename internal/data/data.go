package data

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
)

func Set(ctx context.Context, db *storage.Storage, key string, value string) error {
	return (*db).Set(ctx, key, value)
}

func SetBatch(ctx context.Context, db *storage.Storage, item map[string]string) error {
	return (*db).SetBatch(ctx, item)
}

func Get(ctx context.Context, db *storage.Storage, key string) (models.Alias, error) {
	return (*db).Get(ctx, key)
}
