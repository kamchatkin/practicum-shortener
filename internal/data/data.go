package data

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"go.uber.org/zap"
)

func Set(ctx context.Context, db *storage.Storage, key string, value string, userID int64) error {
	logger := logs.NewLogger()
	logger.Info("data.Set", zap.String("key", key), zap.String("value", value))
	return (*db).Set(ctx, key, value, userID)
}

func SetBatch(ctx context.Context, db *storage.Storage, item map[string]string, userID int64) error {
	logger := logs.NewLogger()
	logger.Info("data.SetBatch", zap.Any("item", item))
	return (*db).SetBatch(ctx, item, userID)
}

func Get(ctx context.Context, db *storage.Storage, key string) (models.Alias, error) {
	logger := logs.NewLogger()
	logger.Info("data.Get", zap.String("key", key))
	return (*db).Get(ctx, key)
}

func GetBySource(ctx context.Context, db *storage.Storage, sourceKey string) (models.Alias, error) {
	logger := logs.NewLogger()
	logger.Info("data.GetBySource", zap.String("sourceKey", sourceKey))
	return (*db).GetBySource(ctx, sourceKey)
}

// RegisterUser Регистрирует нового пользователя и возвращает его ID
func RegisterUser(ctx context.Context, db *storage.Storage) (int64, error) {
	logger := logs.NewLogger()
	logger.Info("data.RegisterUser")
	return (*db).RegisterUser(ctx)
}

func UserAliases(ctx context.Context, db *storage.Storage, userID int64) ([]*models.Alias, error) {
	logger := logs.NewLogger()
	logger.Info("data.UserAliases", zap.Int64("user_id", userID))
	return (*db).UserAliases(ctx, userID)
}
