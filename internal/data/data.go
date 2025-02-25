package data

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
)

func Set(ctx context.Context, db *storage.Storage, key string, value string, userID int64) error {
	return (*db).Set(ctx, key, value, userID)
}

func SetBatch(ctx context.Context, db *storage.Storage, item map[string]string, userID int64) error {
	return (*db).SetBatch(ctx, item, userID)
}

func Get(ctx context.Context, db *storage.Storage, key string) (models.Alias, error) {
	return (*db).Get(ctx, key)
}

func GetBySource(ctx context.Context, db *storage.Storage, sourceKey string) (models.Alias, error) {
	return (*db).GetBySource(ctx, sourceKey)
}

// RegisterUser Регистрирует нового пользователя и возвращает его ID
func RegisterUser(ctx context.Context, db *storage.Storage) (int64, error) {
	return (*db).RegisterUser(ctx)
}

func UserAliases(ctx context.Context, db *storage.Storage, userID int64) ([]*models.Alias, error) {
	return (*db).UserAliases(ctx, userID)
}
